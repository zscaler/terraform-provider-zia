# Terraform Provider ZIA — Claude Code Guidelines

This file provides project-specific guidance for the `terraform-provider-zia` Terraform provider. Follow these conventions when creating, modifying, or reviewing code.

## Project Overview

This is a Terraform provider for Zscaler Internet Access (ZIA). It wraps the ZIA REST API via the Go SDK (`zscaler-sdk-go/v3`). The provider implements resources and data sources following the `terraform-plugin-sdk/v2` patterns.

## Project Structure

```
zia/                           # All Go source files (resources, data sources, tests, helpers)
  provider.go                  # Provider registration
  resource_zia_<name>.go       # Resource implementations
  data_source_zia_<name>.go    # Data source implementations
  resource_zia_<name>_test.go  # Acceptance tests
  common.go                    # Shared schema builders & expand/flatten helpers
  utils.go                     # Shared utility functions
  common/version.go            # Provider version string
docs/resources/                # Resource documentation (published to Terraform Registry)
docs/data-sources/             # Data source documentation
examples/                      # Example .tf files
```

## SDK Client Pattern

```go
zClient := meta.(*Client)
service := zClient.Service
```

## Resource ID Convention (string `id` + typed `<resource>_id`)

The ZIA API identifies objects with an integer ID, but Terraform requires the
resource's internal `id` to be a string. Every resource therefore exposes two
attributes:

- `id` — `schema.TypeString`, `Computed: true`. This is Terraform's internal
  resource ID. It is set with `d.SetId(strconv.Itoa(resp.ID))` and is what users
  reference and import against.
- `<resource>_id` (e.g. `rule_label_id`, `header_profile_id`) —
  `schema.TypeInt`, `Computed: true`. This holds the API's native integer ID and
  is what the Create/Read/Update/Delete paths pass to the SDK.

Read the integer ID back out with the shared helper
`getIntFromResourceData(d, "<resource>_id")` (in `utils.go`) — never parse
`d.Id()` ad hoc in the CRUD body. The string↔int conversion lives only at the
boundaries: `d.SetId(strconv.Itoa(resp.ID))` writes the string form, and the
importer parses `d.Id()` with `strconv.ParseInt` (falling back to a name lookup
when the import key is not numeric).

```go
"id": {
    Type:     schema.TypeString,
    Computed: true,
},
"rule_label_id": {
    Type:     schema.TypeInt,
    Computed: true,
},
```

```go
// CRUD bodies read the typed integer ID, not d.Id():
id, ok := getIntFromResourceData(d, "rule_label_id")

// Create/Read set both forms from the API response:
d.SetId(strconv.Itoa(resp.ID))
_ = d.Set("rule_label_id", resp.ID)
```

`resource_zia_rule_labels.go` is the canonical reference implementation. Do not
reinvent the conversion or expose the integer ID as the Terraform `id`.

Avoid `Computed: true` on input attributes that the user fully controls. Use it
only where the API genuinely assigns or defaults the value (the two ID fields
above, API-defaulted fields, and API-assigned nested `id`s) — see "Schema
Conventions".

## Rule-Based Resources — Critical Conventions

Rule-based resources include: `zia_ssl_inspection_rules`, `zia_firewall_filtering_rule`, `zia_firewall_dns_rules`, `zia_firewall_ips_rules`, `zia_cloud_app_control_rules`, `zia_url_filtering_rules`, `zia_dlp_web_rules`, `zia_forwarding_control_rule`, `zia_nat_control_rules`, `zia_sandbox_rules`, `zia_bandwidth_control_rules`, `zia_traffic_capture_rules`, `zia_file_type_control_rules`, `zia_casb_dlp_rules`, `zia_casb_malware_rules`.

### Order Field Validation

The `order` field MUST include `ValidateFunc: validation.IntAtLeast(1)`:

```go
"order": {
    Type:         schema.TypeInt,
    Required:     true,
    ValidateFunc: validation.IntAtLeast(1),
},
```

Negative or zero order values are not supported. Rules with negative orders are internal default/predefined rules managed by the API.

### Stripping Read-Only Fields in updateOrder Callbacks

When the provider reorders rules, the `updateOrder` callback fetches each rule via GET and sends it back via PUT with the new order. For predefined rules, the API returns read-only fields that MUST NOT be included in the PUT body:

- `Predefined` (bool)
- `DefaultRule` (bool)
- `AccessControl` (string)

The `updateOrder` callback MUST strip these before the update:

```go
func(id int, order OrderRule) error {
    rule, err := <sdk_package>.Get(ctx, service, id)
    if err != nil {
        return err
    }
    rule.LastModifiedTime = 0
    rule.LastModifiedBy = nil
    rule.Predefined = false
    rule.DefaultRule = false
    rule.AccessControl = ""
    rule.Order = order.Order
    rule.Rank = order.Rank
    _, err = <sdk_package>.Update(ctx, service, id, rule)
    return err
},
```

Only strip fields that exist in the SDK struct for each rule type:

| Rule Resource | Fields to Strip |
|---|---|
| ssl_inspection, firewall_filtering, firewall_dns, firewall_ips, nat_control, traffic_capture | `Predefined`, `DefaultRule`, `AccessControl` |
| cloud_app_control | `Predefined`, `AccessControl` |
| sandbox_rules, bandwidth_control | `DefaultRule`, `AccessControl` |
| dlp_web, file_type_control, casb_dlp, casb_malware | `AccessControl` |
| url_filtering, forwarding_control | No read-only fields to strip |

### Reorder Loop Architecture (`reorderAll` in `common.go`)

The ZIA API does not provide a native bulk-reorder endpoint, so the provider runs a goroutine-driven reorder loop that converges the API state to the HCL-declared order using `GET` + targeted `PUT` calls. PR #567 reshaped this loop. **Do not regress these invariants.**

1. **Diff-based passes.** Each `reorderAll` tick calls the resource's `getCurrent()` callback exactly once; that callback is implemented as a single `GetAll` returning `map[int]OrderRule` (id → {Order, Rank}). The loop then iterates the registered rules and only issues `updateOrder(id, …)` when the API's current `Order/Rank` differs from the desired values. **Never** re-PUT a rule whose order already matches.

   When implementing `getCurrent` for a new rule resource, follow this shape (special-case `casb_malware_rules`, which lacks `Rank`):

   ```go
   return func() (map[int]OrderRule, error) {
       all, err := <sdk_pkg>.GetAll(ctx, service)
       if err != nil {
           return nil, err
       }
       m := make(map[int]OrderRule, len(all))
       for _, r := range all {
           m[r.ID] = OrderRule{Order: r.Order, Rank: r.Rank}
       }
       return m, nil
   }
   ```

2. **`countOrderable` + skip-out-of-range.** Predefined / unmanaged rules can have `Order < 1` or `Order = 0`; counting them inflates the ceiling and produces `INVALID_INPUT_ARGUMENT: Rule is not allowed at order N` from the API. `reorderAll` calls `countOrderable(current)` to get only the rules with `Order >= 1`, and defers any `PUT` whose desired `Order > apiOrderable`. Those deferred rules are picked up on a later tick once new POSTs extend the orderable range.

3. **Deadlock-breaker (`maxStuckOnSkippedTicks ≈ 60s`).** Terraform's parallelism creates this state mid-apply: every in-range rule is at its target, but `K` registered rules have declared orders > `apiOrderable` and remain skipped. Without an early exit, the loop sat at the slower `maxNoProgressTicks` (≈3 min) every batch. The new fast-exit triggers when:
   - `skipped > 0`
   - all in-range rules at target (`alreadyAtTarget == size - skipped`)
   - no progress this pass (`putsIssued == 0` and `alreadyAtTarget` did not grow)
   - no PUT errors

   Returning here unblocks `waitForReorder` so the next `Create` batch can extend the range. The skipped rules remain registered and are reconciled on the next reorder cycle triggered by the new registration.

4. **Convergence requires two clean passes.** A pass with `putsIssued == 0`, `skipped == 0`, `puterrs == 0`, and a stable `size` counts as one clean pass. The loop exits only after **two** consecutive clean passes, so in-flight PUTs from the previous pass have time to settle in the SDK's `GetAll` view (the OneAPI cache invalidates parent collections on non-GET — see SDK ≥ v3.8.32 — but the API still has its own propagation lag).

5. **`updateOrder` reset on PUT error.** If any `PUT` in a pass fails (transient 429/5xx), the loop resets its stability counters so the next tick re-evaluates from scratch. This is what lets the deadlock-breaker stay aggressive without prematurely exiting on transient failures.

When adding a new rule resource, register `reorderWithBeforeReorder(resourceType, getCurrent, updateOrder, beforeReorder)` from the resource's `Create`/`Update` paths. The shared loop handles everything above — do not write a per-resource reorder loop.

### Predefined Rules — User-Facing Guidance

- Predefined rules CAN be managed via Terraform for reordering purposes
- `destroy` operations are NOT supported for predefined rules
- Not all attributes available on custom rules apply to predefined rules
- Rule orders must always be contiguous (no gaps)
- When deleting custom rules, use `terraform apply -target=<resource>` then re-adjust remaining order numbers

### Rule Documentation Requirements

All rule-based resource docs (`docs/resources/zia_*_rule*.md`) MUST include these three notes before "Example Usage":

```markdown
~> **NOTE:** Predefined rules can be managed via the Terraform provider for reordering purposes; however, `destroy` operations are not supported for predefined rules, and not all attributes available on custom rules apply to them. When deleting existing custom rules, use the Terraform `-target` flag to target the specific rule to be removed.

~> **NOTE:** Rule orders must always be contiguous (no gaps). Deleting a rule must be followed by order number re-adjustment of the remaining rules to ensure the API honours the required order.

~> **NOTE:** The `order` attribute must always be a positive whole number starting at 1. Negative numbers and zero are **not supported** and will result in an error.
```

## Detach-Before-Delete for Shared/Referenced Objects

The ZIA API refuses to delete an object that is still referenced by a rule
(it returns `RESOURCE_IN_USE` / "associated with N rule(s)"). Resources whose
objects can be attached to rules must therefore detach themselves from every
referencing rule in their `Delete` before calling the SDK delete.

All detach helpers live in `utils.go`. Do not create a separate file for them,
and do not reimplement the loop per resource.

There are three shapes, depending on how the rule stores the reference:

1. **`[]common.IDNameExtensions` on URL filtering rules** — use
   `DetachURLFilteringRuleRef(ctx, client, id, label, getResources, setResources)`.
   The `get`/`set` closures select the field (e.g. `r.HTTPHeaderProfiles`,
   `r.HTTPHeaderActionProfiles`). Used by `zia_http_header_profile` and
   `zia_http_header_action_profile`.

2. **`[]common.IDNameExtensions` on firewall filtering / DLP web rules** — use
   the existing `DetachRuleIDNameExtensions` / `DetachDLPEngineIDNameExtensions`
   helpers (same closure shape, different rule type).

3. **URL categories (attached across many rule types)** — call
   `detachURLCategoryFromAllResources(ctx, client, categoryID)` from the
   `zia_url_categories` `Delete`. It fans out **concurrently** to every rule
   type that can reference a category and only issues a write (`Update`) for a
   rule that actually contains the category — rule types where the category is
   not attached incur a single read and no writes. A single settle delay runs
   at the end only if at least one rule was updated.

   The per-resource detachers it dispatches to handle two storage shapes:
   - **`[]string` of category IDs** (url_filtering, ssl_inspection, sandbox,
     file_type_control, bandwidth_classes, ftp_control_policy, and the mirrored
     `dest_ip_categories`/`res_categories` pair on firewall_dns / firewall_ips).
     Use the shared `removeStringFromSlice` primitive; for the mirrored pair,
     strip the ID from **both** fields so the API's request/response match
     requirement stays satisfied.
   - **`[]common.IDNameExtensions`** (dlp_web_rules), matched on the nested
     numeric `id` (or `name`).

When adding a new rule resource that can reference a URL category, add a
`detachCategoryFrom<Resource>` function in `utils.go` and register it in the
`detachURLCategoryFromAllResources` dispatcher slice. When adding a new
non-category object that can be attached to URL filtering rules, reuse
`DetachURLFilteringRuleRef` with field closures — do not write a new loop.

Performance rule: never issue an `Update` for a rule that does not actually
reference the object being deleted. Discovery reads (`GetAll` per rule type) are
acceptable and should run concurrently; writes must be gated on an actual match.

## Common Troubleshooting

### "Request body is invalid" on predefined rule reorder

The API rejects PUT requests for predefined rules when read-only fields (`Predefined`, `DefaultRule`, `AccessControl`) are included. The fix is to strip these fields in the `updateOrder` callback (see above).

### Negative order creates corrupted rule

Setting `order = -1` may be accepted by the API and stored as `order = 0`, creating an undeletable rule. Prevention: `ValidateFunc: validation.IntAtLeast(1)` on all rule `order` fields.

### Non-contiguous rule orders cause drift

Rule orders must be sequential with no gaps. If a rule at order 5 is deleted, rules at orders 6+ must be re-adjusted to fill the gap, or the API will not honour the requested positions.

### `INVALID_INPUT_ARGUMENT: Rule is not allowed at order N`

The provider attempted a `PUT` whose `Order` exceeded the API's currently orderable count. Root cause is almost always either (a) `countOrderable` was bypassed in a new resource's reorder wiring, or (b) `getCurrent()` returned a stale `GetAll` view. Verify:
- The resource's `getCurrent` callback returns `map[int]OrderRule` and is invoked from `reorderWithBeforeReorder` (do not reimplement the loop).
- The vendored Zscaler SDK is ≥ v3.8.32 so parent-collection cache invalidation fires on non-GET requests (otherwise `GetAll` returns a pre-PUT snapshot and the diff calculation is wrong).

### Excessive duplicate `PUT`s during apply (one rule receiving 10+ updates)

Symptom of regressing the diff-based reorder. Confirm `reorderAll` in `common.go` still calls `getCurrent()` once per pass and skips rules whose `Order/Rank` already match. Re-PUTting every registered rule per cycle was the SUP-3988 bug; never restore that behaviour.

### Apply hangs ~3 minutes per batch with no errors

The `maxStuckOnSkippedTicks` deadlock-breaker may have been removed or its predicate broken. The fast-exit must trigger when `skipped > 0`, all in-range rules are at target, and no PUTs/progress occurred this pass — otherwise the slower `maxNoProgressTicks` safety net runs and each `Create` batch waits ~3 min for it to time out.

## JMESPath Client-Side Filtering

The provider supports an optional `search` attribute on select data sources that enables client-side filtering via [JMESPath](https://jmespath.org/) expressions. This feature is powered by the `zscaler-sdk-go` JMESPath integration — the SDK applies the expression after all pages have been fetched from the API, before results are returned to the provider.

### How It Works

1. The data source checks for a `search` attribute in the Terraform configuration
2. If present, the context is enriched via `zscaler.ContextWithJMESPath(ctx, expression)`
3. The SDK's pagination engine fetches all pages as usual
4. `ApplyJMESPathFromContext` applies the JMESPath filter to the aggregated results
5. The filtered results are returned to the provider for local name/ID matching

### Supported Data Sources

| Data Source | Filterable Fields (camelCase) |
|---|---|
| `zia_group_management` | `name`, `idpId`, `comments` |
| `zia_user_management` | `name`, `email`, `department`, `adminUser`, `type` |
| `zia_department_management` | `name`, `idpId`, `comments`, `deleted` |
| `zia_devices` | `name`, `osType`, `osVersion`, `deviceModel`, `ownerName` |
| `zia_cloud_applications` | `app`, `appName`, `parent`, `parentName` |
| `zia_location_groups` | `name`, `groupType`, `comments`, `predefined` |
| `zia_location_management` | `name`, `country`, `sslScanEnabled`, `ofwEnabled`, `authRequired`, `profile` |

### Implementation Pattern

When adding `search` to a new data source:

```go
import "github.com/zscaler/zscaler-sdk-go/v3/zscaler"

// 1. Add to schema
"search": {
    Type:        schema.TypeString,
    Optional:    true,
    Description: "JMESPath expression to filter results client-side.",
},

// 2. Enrich context before SDK calls
if searchExpr, ok := d.GetOk("search"); ok {
    ctx = zscaler.ContextWithJMESPath(ctx, searchExpr.(string))
    log.Printf("[INFO] JMESPath filter set: %s\n", searchExpr.(string))
}
```

### Key Rules

- The `search` attribute MUST always be `Optional` — existing behavior is unchanged when omitted
- Field names in JMESPath expressions use the API's **camelCase** names (e.g., `idpId`, not `idp_id`)
- JMESPath filtering narrows the pool BEFORE local name/ID matching — if the filter excludes the target, the lookup will fail with "not found"
- Debug logs are emitted by the SDK when JMESPath is active (visible with `TF_LOG=DEBUG`)

## User-Facing Writing Conventions

These rules apply to ALL user-facing prose in this repo:

- `docs/resources/*.md`
- `docs/data-sources/*.md`
- `docs/guides/*.md` (including `release-notes.md`)
- `CHANGELOG.md`
- The "Description" string on schema fields (it appears in `terraform-plugin-docs` output and Registry pages)
- Any `~> NOTE` callouts inside docs

### 1. Write for the user, not for the SDK

User-facing text MUST NOT reference Go SDK identifiers, internal types, or implementation primitives. The user never reads the SDK, so do not mention:

- SDK package paths (`zscaler-sdk-go`, `filteringrules`, `oneapiclient`, `oneapiconfig`)
- SDK function names (`GetAll`, `GetByName`, `Get`, `Create`, `Update`)
- SDK struct names (`GetAllFilterOptions`, `OrderRule`, `FirewallFilteringRules`)
- Provider-internal helpers (`reorderAll`, `countOrderable`, `reorderWithBeforeReorder`, `waitForReorder`)
- Internal terms like "pagination engine", "JMESPath context enrichment", "diff-based PUTs", "parent-collection invalidation"

Instead use neutral product language: *"the data source returns…"*, *"filter applied server-side"*, *"runs after the response is received"*, *"used to look up a single rule"*, *"used to narrow results"*.

These internal mechanics belong in `CLAUDE.md`, `.cursor/rules/`, and `.claude/skills/` — never in the published docs or changelog.

### 2. "Breaking change" is reserved for actual breaks

Do NOT label a change as a breaking change unless it satisfies BOTH:

- A working `terraform plan` / `apply` / `import` from a previous release will now error or change behaviour observably (drift, state mismatch, attribute removal, output value change), AND
- The user has no automatic migration path — they must edit HCL or state to recover.

Specifically NOT breaking, even if the user has to update HCL:

- Renaming an INPUT argument when the underlying capability is preserved under another argument (e.g. "look up by id" still works, the argument name just moved from `id` to `rule_id`). Document the rename in one sentence inline; do not use `~>` callouts, banners, or a `> Breaking change:` block.
- Adding a new optional argument or output attribute.
- Adding a new resource or data source.
- Tightening a `ValidateFunc` on a field that was already documented as having that constraint.

When in doubt, ask the user before adding a breaking-change banner — banners scare adopters and we only spend that scariness budget when it's earned.

### 3. Changelog and release-notes entries describe outcomes

`CHANGELOG.md` and `docs/guides/release-notes.md` entries are read by users deciding whether to upgrade. They should:

- Start with the user-visible improvement (*"Enhanced the X data source to return a list of rules…"*), not the implementation (*"Refactored X to call GetAll once and apply JMESPath context-enrichment…"*).
- Name the resource(s) / data source(s) affected so users can `Ctrl+F` for them.
- Reference the GitHub PR with the `[PR #NNN](https://github.com/zscaler/terraform-provider-zia/pull/NNN)` link convention already in use.
- Reference the issue (`[issue #NNN](…)`) when the change addresses a specific report.
- Stay to one paragraph per PR unless multiple distinct improvements ship together.
- NEVER include SDK terminology (see rule 1).

The internal mechanics — "the reorder engine now…", "parent-collection cache invalidation in OneAPI v3.8.32…", "diff-based pass with `countOrderable` deferral" — belong in `CLAUDE.md` and the architecture sections of the cursor/claude internal docs, not in the changelog.

## Schema Conventions

- Booleans with `omitempty` in the SDK: use `Optional: true, Computed: true`
- API-defaulted fields: use `Optional: true, Computed: true`
- Write-only fields (passwords, keys): preserve from prior state in Read
- Nested blocks with API-assigned IDs: use `Computed: true` on the nested `id` field

## Build and Test

```bash
# Build the provider
go build ./...

# Run a specific acceptance test
TF_ACC=1 go test ./zia/ -v -run TestAccResource<Name>Basic -timeout 120m

# Run sweepers to clean up test resources
go test ./zia/ -v -sweep=global -sweep-run=zia_<name> -timeout 30m

# Debug logging
TF_LOG=DEBUG ZSCALER_SDK_VERBOSE=true ZSCALER_SDK_LOG=true terraform apply -no-color 2>&1 | tee /tmp/tf-debug.log
```

## Release Versioning

Every release MUST update:
1. `zia/common/version.go` — version string
2. `GNUmakefile` — all three version occurrences in `build13` target
3. `CHANGELOG.md` — new entry at the top
4. `docs/guides/release-notes.md` — same entry, update `Last updated` line

## Critical Rules

1. NEVER create a resource without an accompanying acceptance test and documentation
2. ALWAYS register new resources in `provider.go` and `resource_type.go`
3. ALWAYS include activation handling in Create, Update, Delete
4. ALWAYS support import by both numeric ID and name
5. ALWAYS add a sweeper for resources with Delete
6. ALWAYS add `ValidateFunc: validation.IntAtLeast(1)` to `order` fields on rule resources
7. ALWAYS strip read-only fields in `updateOrder` callbacks for rule resources
8. ALWAYS include the three predefined-rule / order-validation notes in rule resource docs
9. Use existing helpers from `common.go` and `utils.go` — never reimplement
10. NEVER reference SDK identifiers, internal helpers, or implementation primitives in user-facing docs, changelog entries, or release notes — see "User-Facing Writing Conventions"
11. NEVER label a change as a breaking change unless a previously working configuration genuinely breaks — argument renames where the capability is preserved are a one-line inline note, not a banner
