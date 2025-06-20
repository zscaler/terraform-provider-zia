---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA): sandbox_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/about-sandbox
  API documentation https://help.zscaler.com/zia/sandbox-policy-settings#/sandboxRules-get
  Retrieves the list of all Sandbox policy rules configured in the ZIA Admin Portal.
---

# zia_sandbox_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-sandbox)
* [API documentation](https://help.zscaler.com/zia/sandbox-policy-settings#/sandboxRules-get)

Use the **zia_sandbox_rules** data source to get information about a sandbox rule in the Zscaler Internet Access.

## Example Usage

```hcl
# ZIA Sandbox Rule by name
data "zia_sandbox_rules" "this" {
    name = "Default BA Rule"
}
```

```hcl
# ZIA Sandbox Rule by ID
data "zia_sandbox_rules" "this" {
    id = "12365478"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Sandbox rule
* `id` - (Optional) Unique identifier for the Sandbox rule

## Read-Only

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Enter additional notes or information. The description cannot exceed 10,240 characters.
* `order` - (Integer) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.
* `state` - (String) The state of the rule indicating whether it is enabled or disabled. Supported values: `ENABLED` or `DISABLED`
* `rank` - (Integer) The admin rank specified for the rule based on your assigned admin rank. Admin rank determines the rule order that can be specified for the rule. Admin rank can be configured if it is enabled in the Advanced Settings.
* `ba_rule_action` - (String) The action configured for the rule that must take place if the traffic matches the rule criteria. Supported Values: `ALLOW` or `BLOCK`
* `first_time_enable` - (Boolean) A Boolean value indicating whether a First-Time Action is specifically configured for the rule. The First-Time Action takes place when users download unknown files. The action to be applied is specified using the firstTimeOperation field.
* `first_time_operation` - (String) The action that must take place when users download unknown files for the first time. Supported Values: `ALLOW_SCAN`, `QUARANTINE`, `ALLOW_NOSCAN`, `QUARANTINE_ISOLATE`
* `ml_action_enabled` - (Boolean) A Boolean value indicating whether to enable or disable the AI Instant Verdict option to have the Zscaler service use AI analysis to instantly assign threat scores to unknown files. This option is available to use only with specific rule actions such as Quarantine and Allow and Scan for First-Time Action.
* `by_threat_score` - (Integer)
* `default_rule` - (Boolean) Value that indicates whether the rule is the Default Cloud IPS Rule or not

* `url_categories` - (List of Strings) The list of URL categories to which the DLP policy rule must be applied.
* `file_types` - (List of Strings) File type categories for which the policy is applied. If not set, the rule is applied across all file types.

`Who, Where and When` supports the following attributes:

* `locations` - (List of Objects) You can manually select up to `8` locations. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `location_groups` - (List of Objects)You can manually select up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `users` - (List of Objects) You can manually select up to `4` general and/or special users. When not used it implies `Any` to apply the rule to all users.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `groups` - (List of Objects) You can manually select up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `departments` - (List of Objects) Apply to any number of departments When not used it implies `Any` to apply the rule to all departments.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `labels` (List of Objects) Labels that are applicable to the rule.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `zpa_app_segments` (List of Objects) The ZPA application segments to which the rule applies
      - `id` - (Integer) Identifier that uniquely identifies an entity
