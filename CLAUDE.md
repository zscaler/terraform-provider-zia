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

## Common Troubleshooting

### "Request body is invalid" on predefined rule reorder

The API rejects PUT requests for predefined rules when read-only fields (`Predefined`, `DefaultRule`, `AccessControl`) are included. The fix is to strip these fields in the `updateOrder` callback (see above).

### Negative order creates corrupted rule

Setting `order = -1` may be accepted by the API and stored as `order = 0`, creating an undeletable rule. Prevention: `ValidateFunc: validation.IntAtLeast(1)` on all rule `order` fields.

### Non-contiguous rule orders cause drift

Rule orders must be sequential with no gaps. If a rule at order 5 is deleted, rules at orders 6+ must be re-adjusted to fill the gap, or the API will not honour the requested positions.

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
