---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/firewallFilteringRules-post
  API documentation https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zia/firewall-policies/firewall-filtering-rules-resource-create-firewall-filtering-rule
  Get information about firewall filtering rules — look up a single rule, or list multiple rules with server-side and JMESPath filtering.
---

# zia_firewall_filtering_rule (Data Source)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/firewallFilteringRules-post)
* [API documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zia/firewall-policies/firewall-filtering-rules-resource-create-firewall-filtering-rule)

Use the **zia_firewall_filtering_rule** data source to read firewall filtering rules from the Zscaler Internet Access cloud firewall.

The data source supports two ways of selecting rules:

1. **Single rule** — supply `rule_id` or `name` to look up exactly one rule. The matched rule's attributes are populated at the top level.
2. **List of rules** — omit `rule_id` and `name` to return every rule in the `rules` block. Narrow the list with the filter arguments below and/or the `search` (JMESPath) argument.

## Example Usage

### 1. Single-rule lookup by name

```hcl
data "zia_firewall_filtering_rule" "example" {
  name = "Office 365 One Click Rule"
}

output "rule_id" {
  value = data.zia_firewall_filtering_rule.example.id
}
```

### 2. Single-rule lookup by id

```hcl
data "zia_firewall_filtering_rule" "by_id" {
  rule_id = 186659
}
```

### 3. List every firewall filtering rule (collection mode)

```hcl
data "zia_firewall_filtering_rule" "all" {}

output "all_rule_names" {
  value = [for r in data.zia_firewall_filtering_rule.all.rules : r.name]
}
```

### 4. List rules with a server-side filter and iterate with `for_each`

```hcl
data "zia_firewall_filtering_rule" "block_drops" {
  rule_action = "BLOCK_DROP"
}

# Build a map keyed by the rule's real API id so for_each is stable.
output "block_drop_rules" {
  value = { for r in data.zia_firewall_filtering_rule.block_drops.rules : r.id => {
    name  = r.name
    order = r.order
    state = r.state
  } }
}
```

### 5. Combine multiple server-side filters

```hcl
data "zia_firewall_filtering_rule" "mac_engineering_block" {
  rule_action  = "BLOCK_DROP"
  department   = "Engineering"
  device_group = "Mac Endpoints"
}
```

### 6. Narrow further with a JMESPath expression

JMESPath runs **after** the server-side filters, against the raw API JSON. Use API `camelCase` field names (`enableFullLogging`, `defaultRule`, `destIpCategories`, …), not the Terraform `snake_case` attribute names.

```hcl
# Only BLOCK_DROP rules that are also ENABLED and have full logging on.
data "zia_firewall_filtering_rule" "active_blocks_with_logging" {
  rule_action = "BLOCK_DROP"
  search      = "[?state=='ENABLED' && enableFullLogging]"
}
```

```hcl
# Predefined rules only.
data "zia_firewall_filtering_rule" "predefined" {
  search = "[?predefined]"
}
```

```hcl
# Custom rules only (not predefined, not the default rule).
data "zia_firewall_filtering_rule" "custom" {
  search = "[?!predefined && !defaultRule]"
}
```

```hcl
# Rules whose destination IP categories include MALWARE_SITE.
data "zia_firewall_filtering_rule" "malware" {
  search = "[?contains(destIpCategories, 'MALWARE_SITE')]"
}
```

```hcl
# Rules whose name starts with "Block".
data "zia_firewall_filtering_rule" "block_prefix" {
  search = "[?starts_with(name, 'Block')]"
}
```

### 7. Use the result with `for_each` to drive another resource

```hcl
data "zia_firewall_filtering_rule" "block_drops" {
  rule_action = "BLOCK_DROP"
  search      = "[?!defaultRule && state=='ENABLED']"
}

# Example: tag each matching rule via a hypothetical downstream resource.
resource "some_module_thing" "per_block_rule" {
  for_each = { for r in data.zia_firewall_filtering_rule.block_drops.rules : r.id => r }

  rule_id   = each.value.id
  rule_name = each.value.name
  rule_order = each.value.order
}
```

## Argument Reference

### Lookup arguments

Use one of these to fetch a single rule. If both are omitted, the data source returns every matching rule in the `rules` block instead.

* `rule_id` - (Optional, Number) Numeric id of the rule to look up.
* `name` - (Optional, String) Name of the rule to look up (case-insensitive exact match).

### Filter arguments

Use any combination of these to narrow the result set. They are passed to the API so the filtering happens server-side.

* `predefined_rule_count` - (Optional, Bool) Include the predefined rule count in the response.
* `rule_name` - (Optional, String) Partial match on rule name.
* `rule_label` - (Optional, String)
* `rule_label_id` - (Optional, Number)
* `rule_order` - (Optional, String)
* `rule_description` - (Optional, String)
* `rule_action` - (Optional, String) e.g. `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, `EVAL_NWAPP`.
* `location` - (Optional, String)
* `department` - (Optional, String)
* `group` - (Optional, String)
* `user` - (Optional, String)
* `device` - (Optional, String)
* `device_group` - (Optional, String)
* `device_trust_level` - (Optional, String)
* `nw_application` - (Optional, String)
* `filter_src_ips` - (Optional, String) Source IPs filter.
* `filter_dest_addresses` - (Optional, String) Destination addresses filter.
* `filter_src_ip_groups` - (Optional, String) Source IP groups filter.
* `filter_dest_ip_groups` - (Optional, String) Destination IP groups filter.
* `filter_nw_services` - (Optional, String) Network services filter.
* `filter_dest_ip_categories` - (Optional, String) Destination IP categories filter.

~> **Note:** the six `filter_*` arguments above are filter inputs that share their name with per-rule output attributes (`src_ips`, `dest_addresses`, etc.). The `filter_` prefix disambiguates them.

### Client-side `search` argument

* `search` - (Optional, String) [JMESPath](https://jmespath.org/) expression applied to the results after the server-side filters above. Field names in the expression use the API's `camelCase` form (e.g. `enableFullLogging`, `destIpCategories`), not the Terraform `snake_case` attribute names.

  ```hcl
  search = "[?action=='BLOCK_DROP' && state=='ENABLED']"
  ```

## Attribute Reference

### Top-level attributes (single-rule mode)

When `rule_id` or `name` resolves to exactly one rule, the following attributes are populated at the top level. In collection mode they are at their zero values; use `rules[*]` instead.

* `id` - (String) Terraform resource id (in single-rule mode this is the rule's numeric id; in collection mode it is a constant placeholder).
* `description` - Additional notes or information. Cannot exceed 10,240 characters.
* `order` - Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on).
* `state` - `ENABLED` rules are actively enforced; `DISABLED` rules keep their place in the rule order but are skipped.
* `action` - `ALLOW`, `BLOCK_DROP`, `BLOCK_RESET`, `BLOCK_ICMP`, or `EVAL_NWAPP`.
* `rank` - Admin rank (`0`-`7`). Defaults to `7`.
* `last_modified_time` - Unix epoch of the most recent change.
* `access_control` - API-managed access-control level.
* `enable_full_logging` - Whether the rule logs the full session.
* `default_rule` - Whether this is the per-tenant default rule.
* `predefined` - Whether the rule is a Zscaler-predefined rule (cannot be deleted).

`Who, Where and When` attributes:

* `locations` - Up to 8 locations. Implies `Any` when empty.
  * `id` - Identifier that uniquely identifies an entity
  * `name` - The configured name of the entity
  * `extensions` - (Map of String)
* `location_groups` - Up to 32 location groups. Implies `Any` when empty.
  * Same nested shape as `locations`.
* `users` - Up to 4 users. Implies `Any` when empty.
  * Same nested shape as `locations`.
* `groups` - Up to 8 groups. Implies `Any` when empty.
  * Same nested shape as `locations`.
* `departments` - Any number of departments. Implies `Any` when empty.
  * Same nested shape as `locations`.
* `time_windows` - Up to 2 time intervals. Implies `Always` when empty.
  * Same nested shape as `locations`.

`Network services / applications` attributes:

* `nw_service_groups` - Predefined or custom network service groups the rule applies to.
* `nw_services` - Predefined or custom network services. Empty means `Any`.
* `nw_application_groups` - Predefined application groups the rule controls.
* `nw_applications` - Predefined applications the rule controls.

`Source` attributes:

* `src_ip_groups` - Source IP address groups.
  * Same nested shape as `locations`.
* `src_ips` - (List of String) Individual IPs, subnets, or ranges.

`Destination` attributes:

* `dest_addresses` - (List of String) IPs, subnets, ranges, or FQDNs.
* `dest_countries` - (List of String) 2-letter ISO-3166-Alpha-2 country codes, e.g. `US`, `CA`.
* `dest_ip_categories` - (List of String) URL-category-based destinations.
* `dest_ip_groups` - Destination IP address groups.
  * Same nested shape as `locations`.

`App service` attributes:

* `app_service_groups` - Application service groups the rule applies to.
  * Same nested shape as `locations`.
* `app_services` - Application services the rule applies to.
  * Same nested shape as `locations`.

`Other` attributes:

* `labels` - Labels applicable to the rule.
  * Same nested shape as `locations`.
* `device_groups` - Device groups the rule applies to.
  * Same nested shape as `locations`.
* `devices` - Individual devices the rule applies to.
  * Same nested shape as `locations`.
* `device_trust_levels` - (List of String) Allowed device trust levels.
* `workload_groups` - (List) Preconfigured workload groups the policy applies to.
  * `id` - (Number) Unique identifier
  * `name` - (String) Workload group name
  * `description` - (String)
  * `expression` - (String)
  * `expression_json` - (Block) Structured expression tree.
* `zpa_app_segments` - (List) ZPA Application Segments the rule applies to (ZPA Gateway forwarding only).
  * `id` - (Number)
  * `name` - (String)
  * `external_id` - (String)
* `last_modified_by` - (Block) Identifier of the admin who last modified the rule.
  * Same nested shape as `locations`.

### `rules` block (collection mode)

* `rules` - (List of Object) When in collection mode, each element is one matching rule. The element schema is identical to the top-level per-rule attributes above, **including** a per-row `id` field that holds the rule's real API id. Use this id as the stable key for `for_each` or downstream resource wiring.

Example accessors:

```hcl
data.zia_firewall_filtering_rule.<x>.rules[0].id
data.zia_firewall_filtering_rule.<x>.rules[*].name
{ for r in data.zia_firewall_filtering_rule.<x>.rules : r.id => r.order }
```
