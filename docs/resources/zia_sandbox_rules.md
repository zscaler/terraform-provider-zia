---
subcategory: "Sandbox Policy & Settings"
layout: "zscaler"
page_title: "ZIA: sandbox_rules"
description: |-
  Creates and manages Sandbox Rules.
---

# Resource: zia_sandbox_rules

The **zia_sandbox_rules** resource allows the creation and management of SAndbox rules in the Zscaler Internet Access.

## Example Usage

```hcl
data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_group_management" "normal_internet" {
    name = "Normal_Internet"
}

resource "zia_sandbox_rules" "this" {
    name                 = "SandboxRule01"
    description          = "SandboxRule01"
    rank                 = 7
    order                = 1
    first_time_enable    = true
    ml_action_enabled    = true
    first_time_operation = "ALLOW_SCAN"
    ba_rule_action       = "BLOCK"
    state                = "ENABLED"
    ba_policy_categories = ["ADWARE_BLOCK", "BOTMAL_BLOCK", "ANONYP2P_BLOCK", "RANSOMWARE_BLOCK"]
    file_types           = ["FTCATEGORY_P7Z",
        "FTCATEGORY_MS_WORD",
        "FTCATEGORY_PDF_DOCUMENT",
        "FTCATEGORY_TAR",
        "FTCATEGORY_SCZIP",
        "FTCATEGORY_WINDOWS_EXECUTABLES",
        "FTCATEGORY_HTA",
        "FTCATEGORY_FLASH",
        "FTCATEGORY_RAR",
        "FTCATEGORY_MS_EXCEL",
        "FTCATEGORY_VISUAL_BASIC_SCRIPT",
        "FTCATEGORY_MS_POWERPOINT",
        "FTCATEGORY_WINDOWS_LIBRARY",
        "FTCATEGORY_POWERSHELL",
        "FTCATEGORY_APK",
        "FTCATEGORY_ZIP",
        "FTCATEGORY_BZIP2",
        "FTCATEGORY_JAVA_APPLET",
        "FTCATEGORY_MS_RTF"]
    protocols            = [
        "FOHTTP_RULE",
        "FTP_RULE",
        "HTTPS_RULE",
        "HTTP_RULE",
    ]
    departments {
        id = [ data.zia_department_management.engineering.id ]
    }
    groups {
        id = [ data.zia_group_management.normal_internet.id ]
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the Firewall IPS policy rule
* `order` - (Required) Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and the Rule Order reflects this rule's place in the order.

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

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

`Devices`

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
