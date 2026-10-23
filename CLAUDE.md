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

Rule-based resources (18, all on the shared reorder engine): `zia_ssl_inspection_rules`, `zia_firewall_filtering_rule`, `zia_firewall_dns_rule`, `zia_firewall_ips_rule`, `zia_cloud_app_control_rule`, `zia_url_filtering_rules`, `zia_dlp_web_rules`, `zia_forwarding_control_rule`, `zia_nat_control_rules`, `zia_sandbox_rules`, `zia_bandwidth_control_rule`, `zia_traffic_capture_rules`, `zia_file_type_control_rules`, `zia_casb_dlp_rules`, `zia_casb_malware_rules`, `zia_endpoint_dlp_rules`, `zia_endpoint_dlp_sub_rules`, `zia_outbound_email_dlp`.

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

4. **Convergence requires two clean passes — unless the cycle never moved a rule.** A pass with `putsIssued == 0`, `skipped == 0`, `puterrs == 0`, and a stable `size` counts as one clean pass. The loop exits only after **two** consecutive clean passes, so in-flight PUTs from the previous pass have time to settle in the SDK's `GetAll` view (the OneAPI cache invalidates parent collections on non-GET — see SDK ≥ v3.8.32 — but the API still has its own propagation lag). Exception (PR #614): when the cycle has issued **no** PUT at all (`cyclePuts == 0`) and the pass finds every registered rule in place, nothing is in flight, so the loop exits after that single pass. This is the common case for order-only updates after a delete closed the gap, and it halves the wait (~30s per cycle). Do not extend the shortcut to cycles that issued PUTs.

5. **`updateOrder` reset on PUT error.** If any `PUT` in a pass fails (transient 429/5xx), the loop resets its stability counters so the next tick re-evaluates from scratch. This is what lets the deadlock-breaker stay aggressive without prematurely exiting on transient failures.

6. **Oscillation breaker (`maxNonImprovingPutTicks ≈ 3 min`, SUP-4231).** When the declared orders are unsatisfiable — duplicate order values in the tenant, or unmanaged rules interleaved inside the declared order range — every PUT displaces a neighbour, so each pass keeps issuing PUTs forever and `maxNoProgressTicks` never fires (`putsIssued > 0` counts as progress). The loop tracks the best `alreadyAtTarget` seen for the current registration set; after `maxNonImprovingPutTicks` consecutive PUT-issuing passes that never exceed that best, it WARN-logs the contested rules (id, API order vs declared order) and returns so the apply completes with residual drift, instead of livelocking until terraform's context deadline (observed as a 2.5h hang). New registrations (`size` change) restart the window. Every issued order update is also INFO-logged (`moving rule N from order X to order Y`) so oscillating rules are identifiable from a customer trace alone. Tests: `TestReorder_OscillationBreaker_ExitsBoundedTime` (ping-pong livelock must exit bounded) and `TestReorder_TransientDisplacement_StillConverges` (bounded displacement must still fully converge — the breaker must not fire on healthy runs).

7. **Bounded read failures.** A failed `getCurrent()` skips every exit condition of the pass, so `reorderAll` counts consecutive read failures separately (`readFailures`) and returns after `maxNoProgressTicks` of them; any successful read resets the count, so transient failures behave exactly as before and no other counter is touched. The waiting Create/Update then reads the rule back and surfaces the API error. Tests: `TestReorder_GetCurrentAlwaysFails_ExitsBounded` and `TestReorder_GetCurrentTransientFailures_StillConverges`.

   The `end` closure passed to `updateOrderedRule` (the end-of-list position, see below) must return the error from its rule list call, never just log it and index an empty result.

8. **Relative placement while rules are still being created (`relativeTargets`, PR #615, SUP-4231).** While a registered rule's declared order is past the orderable count, the list is incomplete (Terraform creates in batches of ~10). Placing in-range rules at their final absolute positions then is wasted work: each later batch is inserted among them and shifts them, so they were moved again after every batch — the customer's 150-rule firewall deployment made ~3,500 moves for 125 rules (~800 PUT/hr against the API's 400/hr per endpoint), took 4.5h and failed. While incomplete, the present registered rules are placed in declared order directly below the unregistered rules above them (`min(others, d-1) + rank`); once complete, absolute declared orders are used, so the final order is unchanged. Unit test: `TestFirewallFilteringRules_BatchCreatePlacesEachRuleAboutOnce` (~150 moves; ~7,000 on 4.8.14).

9. **Two-phase pass with an up-to-date snapshot.** A move shifts the rules between the old and new position, so after each PUT the pass updates its snapshot (`applyMove`) instead of comparing the next rules against stale positions (which re-PUT rules already in place). Phase 1 leaves the longest run of rules already in relative order (`inOrderRun`, longest increasing subsequence of current positions) and moves only the others; phase 2 moves anything still off target. A swap costs 2 PUTs, a new rule 1. The test fake (`fakeAPI.updateOrder`) models the same shifting — keep it realistic.

10. **The reorder runs on its own context.** The `getCurrent`/`updateOrder` closures use `reorderCtx := context.WithoutCancel(ctx)`: the cycle serves every waiting rule, so it must not die with the create timeout (default 60m) of the rule that started it — that produced SUP-4231's bursts of `context deadline exceeded` exactly 60 minutes after each cycle's first create. Test: `TestFirewallFilteringRules_ReorderSurvivesCancelledCreate`.

When adding a new rule resource, register `reorderWithBeforeReorder(resourceType, getCurrent, updateOrder, beforeReorder)` from the resource's `Create`/`Update` paths, with closures over `reorderCtx`. The shared loop handles everything above — do not write a per-resource reorder loop.

### Create/Update Placement (`updateOrderedRule` in `common.go`, PR #614)

**Invariant: never write a position computed from a rule list read earlier.** A rule of the same type deleted while the request is in flight — by the same apply or outside Terraform — shrinks the list, so an absolute position read before the delete (typically "the last order") no longer exists and the API rejects the write with `INVALID_INPUT_ARGUMENT` `Rule is not allowed at order N`. This was the Bridgewater failure: a mid-stack rule deleted and the rules below renumbered in one apply; every Update moved its rule to the last order from a list read before the 5–12s DELETE completed, and all of them failed. The delete alone had already produced the renumbered positions.

Every rule `Update` writes through `updateOrderedRule(d, known, current, end, write)` before registering with `reorderWithBeforeReorder`. It picks the mode from the planned diff (`orderedRuleUpdateMode`, which uses `d.HasChangesExcept("order")` — it reads the plan, so unit tests must build the update from a real plan, see Build and Test):

| Mode | When | Write |
|---|---|---|
| `ruleUpdateOrderOnly` | only `order` changed | **none** — the reorder engine places the rule and skips it if already in place |
| `ruleUpdateInPlace` | other attributes changed, `rank` did not | the attributes, at the rule's **current** position (`known` from the existence-check GET, re-read via `current` only on retry) |
| `ruleUpdateMoveToEnd` | `rank` changed | at `end()` |

- `end` for the 10 resources that historically parked rules (url_filtering, firewall_filtering, firewall_dns, firewall_ips, ssl_inspection, file_type_control, sandbox, cloud_app_control, forwarding_control, traffic_capture): `endOfList(orders)` from a fresh list (last order, rank 7 — allowed whatever the ranks above; the reason the PR #514 park existed). Apply the resource's default-rule filter first (e.g. sandbox `filterOutDefaultRule`).
- `end` for the 8 resources that write the declared order (bandwidth, nat, casb_dlp, casb_malware, dlp_web, endpoint_dlp, endpoint_dlp_sub, outbound_email_dlp): the **declared** order and rank. Do not move these to the end: some of their `updateOrder` callbacks write only `Order` (bandwidth), so a parked rank 7 would never be corrected. endpoint_dlp parks at the tail with its declared rank.
- Resources without `rank` (casb_malware, outbound_email_dlp) never hit the rank path; their `write` sets only `Order`.
- Capture `intendedOrder`/`intendedRank` **before** calling `updateOrderedRule` and register those with the reorder engine — `write` mutates `req.Order`/`req.Rank`.
- `writeAtFreshPosition` retries a write rejected by `isStaleOrderError` (INVALID_INPUT_ARGUMENT + `not allowed at order N`) up to `maxStaleOrderRetries` (3), re-reading the position each time. It must run **before** `failFastOnErrorCodes`, which treats every INVALID_INPUT_ARGUMENT as fatal (so the old in-loop INVALID_INPUT_ARGUMENT retries never ran).

**Create:** creates post at a starting order cached per apply (`<resource>StartingOrder`, per rule type / per parent where scoped). On `isStaleOrderError`, reset the cache and read it again (loop-based creates: reset + restore `req.Order`/`req.Rank` to the declared values + `continue`, bounded by `staleOrderRetries`; loop-free creates: a `startingOrder(refresh bool)` closure and a bounded retry after the POST; url_filtering uses `writeAtFreshPosition`). One concurrent delete is tolerated anyway (a create may use count+1); two or more are not.

**Create at the end of the list (PR #615):** after a successful create, advance the cached starting order (`<resource>StartingOrder++` under its lock) so the next rule is created after it instead of at the same position, which pushed every rule already placed down by one. Read the cached value under the lock too (`-race` flags it otherwise).

**Delete** needs no ordering logic: once nothing writes a stale position, there is nothing to race.

### Predefined Rules — User-Facing Guidance

- Predefined rules CAN be managed via Terraform for reordering purposes
- `destroy` operations are NOT supported for predefined rules
- Not all attributes available on custom rules apply to predefined rules
- Rule orders must always be contiguous (no gaps)
- When deleting a custom rule, renumber the remaining rules; from v4.8.14 the delete and the renumbering can run in the same apply (before v4.8.14: `terraform apply -target=<resource>` first, then renumber in a second apply)

### Rule Documentation Requirements

All rule-based resource docs (`docs/resources/zia_*_rule*.md`) MUST include these four notes before "Example Usage":

```markdown
~> **NOTE:** Predefined rules can be managed via the Terraform provider for reordering purposes; however, `destroy` operations are not supported for predefined rules, and not all attributes available on custom rules apply to them.

~> **NOTE:** Rule orders must always be contiguous (no gaps). When you delete a rule, renumber the remaining rules so the orders stay contiguous. From v4.8.14, the deletion and the renumbering can be applied together in one `terraform apply`.

~> **NOTE:** The `order` attribute must always be a positive whole number starting at 1. Negative numbers and zero are **not supported** and will result in an error.

~> **NOTE:** Creating or updating rules takes longer than for most other resources. The ZIA API does not provide a way to set the position of several rules in one request, so after a rule is created or updated the provider sends additional requests to move the rules into the order declared in your configuration, and waits until all rules of this type are in place. While this happens, Terraform reports the resource as still creating or modifying; the provider is working in the background. Do not interrupt the run, as cancelling it can leave the rules partially reordered.
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

The provider attempted a `PUT` whose `Order` exceeded the API's currently orderable count. Root cause is one of (a) a write at a position computed from a rule list read before a concurrent delete completed — the Bridgewater delete-and-renumber case; a resource's Create/Update not going through `updateOrderedRule` / the stale-order retry (see "Create/Update Placement"), (b) `countOrderable` was bypassed in a new resource's reorder wiring, or (c) `getCurrent()` returned a stale `GetAll` view. A customer trace for (a) shows the PUTs carrying the old last order sent before `rule deleted` is logged, and no `reorderAll` lines. Verify:
- The resource's `getCurrent` callback returns `map[int]OrderRule` and is invoked from `reorderWithBeforeReorder` (do not reimplement the loop).
- The vendored Zscaler SDK is ≥ v3.8.32 so parent-collection cache invalidation fires on non-GET requests (otherwise `GetAll` returns a pre-PUT snapshot and the diff calculation is wrong).

### Excessive duplicate `PUT`s during apply (one rule receiving 10+ updates)

Symptom of regressing the diff-based reorder. Confirm `reorderAll` in `common.go` still calls `getCurrent()` once per pass and skips rules whose `Order/Rank` already match. Re-PUTting every registered rule per cycle was the SUP-3988 bug; never restore that behaviour.

### Apply hangs ~3 minutes per batch with no errors

The `maxStuckOnSkippedTicks` deadlock-breaker may have been removed or its predicate broken. The fast-exit must trigger when `skipped > 0`, all in-range rules are at target, and no PUTs/progress occurred this pass — otherwise the slower `maxNoProgressTicks` safety net runs and each `Create` batch waits ~3 min for it to time out.

### `zia_firewall_dns_rule`: `Redirect IP address must be provided.` / `redirect_ip` drift (API defect workaround)

After a DNS rule has been updated, `GET /firewallDnsRules/{id}` stops returning `redirectIp`, while `PUT` still rejects the rule without it (and the list endpoint can return an older version of the rule that still has it). The provider works around this in `resource_zia_firewall_dns_rules.go`: Create/Update record the configured `redirect_ip` per rule ID (`rememberDNSRedirectIP` in `utils.go`), the reorder `updateOrder` callbacks restore it when the GET omits it, and Read keeps the known value when the API omits it and the action is unchanged. Do not "simplify" this away; remove it only once the API returns `redirectIp` again.

### `zia_dlp_web_rules`: apply never finishes; reorder log shows `from order N (rank 0) to order N (rank 7)` every pass

The DLP web rules API can store a different rank from the one sent (observed: rank 7 sent, rank 0 stored, with Admin Ranking disabled), while other rule types keep it. Because `reorderAll` only treats a rule as placed when order AND rank match, registering the configured rank made DLP rules never reach target, and the deadlock-breaker never released deferred rules. Create/Update therefore register the rank the API stored (`dlpWebRuleStoredRank` in `utils.go`). Do not change them back to `req.Rank`.

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

# Unit tests (offline: in-memory mock of the ZIA API, no credentials)
TF_ACC= go test ./zia -run 'TestOrderedRules_|TestURLFilteringRules_|TestIsStaleOrderError|TestReorder' -count=1 -race

# Run a specific acceptance test
TF_ACC=1 go test ./zia/ -v -run TestAccResource<Name>Basic -timeout 120m

# Run sweepers to clean up test resources
go test ./zia/ -v -sweep=global -sweep-run=zia_<name> -timeout 30m

# Debug logging
TF_LOG=DEBUG ZSCALER_SDK_VERBOSE=true ZSCALER_SDK_LOG=true terraform apply -no-color 2>&1 | tee /tmp/tf-debug.log
```

### Unit tests for rule ordering

- `zia/unit_fixture_test.go`: `newUnitFixture` (SDK OneAPI mock harness `common.NewZIATest`, reorder engine reset, 10ms tick) and `fakeRuleList` — a JSON rule list that behaves like the API: contiguous orders, a PUT moves the rule and shifts the others, an order past the end is rejected with `Rule is not allowed at order N`, writes are serialized, `deleteDelay` makes a delete slow and `onGetRule` lands a delete between a read and the following write.
- `zia/common_rule_order_unit_test.go`: `orderedRuleSpecs()` lists every rule resource (paths, ID attribute, extra state such as `type` / `parent_rule` / an SSL `action` block); `TestOrderedRules_OrderOnlyUpdateSendsNoWrite` and `TestOrderedRules_InPlaceUpdateKeepsPosition` run on all of them. **Add a spec when adding a rule resource.**
- `zia/resource_zia_url_filtering_rules_unit_test.go`: the end-to-end delete-and-renumber race, stale-position retries on update and create, rank change.
- Build updates with `f.updateData(r, id, state, config)` (a real plan via `r.Diff`). A `ResourceData` built from state plus `d.Set` has no diff, so `HasChangesExcept` sees nothing and every update looks order-only.
- Unit tests import `github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common`; the local `vendor/` is gitignored — run `go mod vendor` after `go.mod` changes so `make build13` keeps working.

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
8. ALWAYS route rule `Update` writes through `updateOrderedRule`, and retry a create's cached starting order on `isStaleOrderError` — NEVER write a position computed from a rule list read earlier
9. ALWAYS include the four predefined-rule / order-validation / reordering-duration notes in rule resource docs
10. Use existing helpers from `common.go` and `utils.go` — never reimplement
11. NEVER reference SDK identifiers, internal helpers, or implementation primitives in user-facing docs, changelog entries, or release notes — see "User-Facing Writing Conventions"
12. NEVER label a change as a breaking change unless a previously working configuration genuinely breaks — argument renames where the capability is preserved are a one-line inline note, not a banner
