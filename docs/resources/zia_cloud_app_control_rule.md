---
subcategory: "Cloud App Control Policy"
layout: "zscaler"
page_title: "ZIA: cloud_app_control_rule"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-rules-cloud-app-control-policy
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/webApplicationRules/{rule_type}-get
  Creates and manages ZIA Cloud Application Control rule.
---

# Resource: zia_cloud_app_control_rule

* [Official documentation](https://help.zscaler.com/zia/adding-rules-cloud-app-control-policy)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/webApplicationRules/{rule_type}-get)

The **zia_cloud_app_control_rule** resource allows the creation and management of ZIA Cloud Application Control rules in the Zscaler Internet Access.

**NOTE** Resources or DataSources to retrieve Tenant Profile or Cloud Application Risk Profile ID information are not currently available.

## Example Usage - Using Data Source for Actions (Recommended)

```hcl
# Get valid actions for the applications
data "zia_cloud_app_control_rule_actions" "webmail_actions" {
  type       = "WEBMAIL"
  cloud_apps = ["GOOGLE_WEBMAIL", "YAHOO_WEBMAIL"]
}

resource "zia_cloud_app_control_rule" "webmail_rule" {
  name                = "WebMail Control Rule"
  description         = "Control webmail access"
  order               = 1
  rank                = 7
  state               = "ENABLED"
  type                = "WEBMAIL"

  # Use data source to get valid actions
  actions             = data.zia_cloud_app_control_rule_actions.webmail_actions.available_actions_without_isolate

  applications        = ["GOOGLE_WEBMAIL", "YAHOO_WEBMAIL"]
  device_trust_levels = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
  user_agent_types    = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
}
```

## Example Usage - AI/ML Application Control

```hcl
data "zia_cloud_app_control_rule_actions" "ai_actions" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI"]
}

resource "zia_cloud_app_control_rule" "ai_control" {
  name         = "ChatGPT Controls"
  description  = "Control ChatGPT usage"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  type         = "AI_ML"

  # Automatically gets all valid actions except ISOLATE
  actions      = data.zia_cloud_app_control_rule_actions.ai_actions.available_actions_without_isolate

  applications = ["CHATGPT_AI"]
}
```

## Example Usage - File Sharing Controls

```hcl
data "zia_cloud_app_control_rule_actions" "file_share_actions" {
  type       = "FILE_SHARE"
  cloud_apps = ["DROPBOX", "ONEDRIVE"]
}

resource "zia_cloud_app_control_rule" "file_sharing" {
  name         = "File Sharing Controls"
  description  = "Control file sharing operations"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  type         = "FILE_SHARE"

  # Returns only actions supported by both Dropbox and OneDrive
  actions      = data.zia_cloud_app_control_rule_actions.file_share_actions.available_actions_without_isolate

  applications = ["DROPBOX", "ONEDRIVE"]
}
```

## Example Usage - Cloud Browser Isolation (ISOLATE Actions)

ISOLATE actions require Cloud Browser Isolation subscription and must be used alone (cannot mix with other actions):

```hcl
data "zia_cloud_app_control_rule_actions" "chatgpt_isolate" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI"]
}

data "zia_cloud_browser_isolation_profile" "cbi_profile" {
  name = "My-CBI-Profile"
}

resource "zia_cloud_app_control_rule" "isolate_chatgpt" {
  name         = "ChatGPT Isolation"
  description  = "Isolate ChatGPT using Cloud Browser Isolation"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  type         = "AI_ML"

  # Use isolate_actions for CBI rules
  actions      = data.zia_cloud_app_control_rule_actions.chatgpt_isolate.isolate_actions

  applications = ["CHATGPT_AI"]

  # Required for ISOLATE actions
  cbi_profile {
    id   = data.zia_cloud_browser_isolation_profile.cbi_profile.id
    name = data.zia_cloud_browser_isolation_profile.cbi_profile.name
    url  = data.zia_cloud_browser_isolation_profile.cbi_profile.url
  }
}
```

## Example Usage - Filtered Actions (ALLOW Only)

```hcl
data "zia_cloud_app_control_rule_actions" "slack_allow" {
  type            = "ENTERPRISE_COLLABORATION"
  cloud_apps      = ["SLACK"]
  action_prefixes = ["ALLOW"]  # Only permissive actions
}

resource "zia_cloud_app_control_rule" "slack_allow_only" {
  name         = "Slack Allow Only"
  description  = "Allow specific Slack operations"
  order        = 1
  rank         = 7
  state        = "ENABLED"
  type         = "ENTERPRISE_COLLABORATION"

  # Only ALLOW_ actions
  actions      = data.zia_cloud_app_control_rule_actions.slack_allow.filtered_actions

  applications = ["SLACK"]
}
```

## Example Usage - With Time Validity

```hcl
data "zia_cloud_app_control_rule_actions" "social_media_actions" {
  type       = "SOCIAL_NETWORKING"
  cloud_apps = ["FACEBOOK"]
}

resource "zia_cloud_app_control_rule" "social_media_time_restricted" {
  name                  = "Social Media Time Restricted"
  description           = "Allow social media only during specified hours"
  order                 = 1
  rank                  = 7
  state                 = "ENABLED"
  type                  = "SOCIAL_NETWORKING"
  actions               = data.zia_cloud_app_control_rule_actions.social_media_actions.available_actions_without_isolate
  applications          = ["FACEBOOK"]

  enforce_time_validity = true
  validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
  validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
  validity_time_zone_id = "US/Pacific"

  time_quota            = 15
  size_quota            = 10
  device_trust_levels   = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The Cloud App Control rule name.
* `type` - (String) The Cloud App Control rule type.

### Optional

* `description` - (String) The description of the Cloud App Control rule.

* `actions` - (List of String) Refer to the Cloud Application Control. To retrieve the list of supported actions based on `type` use the resource: [zia_cloud_app_control_rule_actions](https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_app_control_rule_actions)

* `order` - (Number) The rule order of execution for the Cloud App Control rule with respect to other
* `rank` - (Number) Admin rank of the admin who creates this rule

* `state` - (String) Enables or disables the Cloud App Control rule.. The supported values are:
  * `DISABLED`
  * `ENABLED`

* `devices` (list) - Specifies devices that are managed using Zscaler Client Connector.
      - `id` - (Integer) Identifier that uniquely identifies an entity
* `device_groups` (list) - This field is applicable for devices that are managed using Zscaler Client Connector.
      - `id` - (Integer) Identifier that uniquely identifies an entity

* `device_trust_levels` - (Optional) List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation. Supported values: `ANY`, `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`

* `user_risk_score_levels` (List of String) - Indicates the user risk score level selectedd for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`
* `user_agent_types` (List of String) - User Agent types on which this rule will be applied: Returned values are: `CHROME`, `FIREFOX`, `MSIE`, `MSEDGE`,   `MSCHREDGE`, `OPERA`, `OTHER`, `SAFARI`
* `time_quota` - (Number) Time quota in minutes, after which the Cloud App Control Rules rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `size_quota` - (Number) Size quota in MB beyond which the Cloud App Control Rules rule is applied. If not set, no quota is enforced. If a policy rule action is set to `BLOCK`, this field is not applicable.
* `validity_start_time` - (String) If enforce_time_validity is set to true, the Cloud App Control Rules rule will be valid starting on this date and time. The date and time must be provided in `RFC1123` format i.e `Sun, 16 Jun 2024 15:04:05 UTC`
* `validity_end_time` - (String) If `enforce_time_validity` is set to true, the Cloud App Control Rules rule will cease to be valid on this end date and time. The date and time must be provided in `RFC1123` format i.e `Sun, 16 Jun 2024 15:04:05 UTC`

  **NOTE** Notice that according to RFC1123 the day must be provided as a double digit value for `validity_start_time` and `validity_end_time` i.e `01`, `02` etc.

* `validity_time_zone_id` - (String) If `enforce_time_validity` is set to true, the Cloud App Control Rules rule date and time will be valid based on this time zone ID. The attribute is validated against the official [IANA List](https://nodatime.org/TimeZones)

* `enforce_time_validity` - (Optional) Enforce a set a validity time period for the Cloud App Control Rules rule.

* `applications` - (List of Strings) The list of cloud applications to which the Cloud App Control rule must be applied. To retrieve the list of cloud applications, use the data source: `zia_cloud_applications`

* `eun_enabled` - (Boolean) A Boolean value that indicates whether Enhanced User Notification (EUN) is enabled for the rule.
* `eun_template_id` - (Integer) The ID of the Enhanced User Notification (EUN) template associated with the rule.
* `browser_eun_template_id` - (Integer) The ID of the Browser Enhanced User Notification (EUN) template associated with the rule.
* `tenancy_profile_ids` - (List of Numbers) This is an immutable reference to an entity. which mainly consists of id and name.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `cloud_app_risk_profile` - (List of Numbers) Name-ID pair of cloud Application Risk Profile for which rule will be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `cloud_app_instances` - (List of Numbers) Name-ID pair of cloud application instances for which rule will be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `cbi_profile` - (List) The cloud browser isolation profile to which the ISOLATE action is applied in the Cloud App Control Rules Policy rules. This block is required when the attribute `action` is set to `ISOLATE`
  * `id` - (String) The universally unique identifier (UUID) for the browser isolation profile
  * `name` - (String) Name of the browser isolation profile
  * `url` - (String) The browser isolation profile URL

* `locations` - (List of Numbers) The Name-ID pairs of locations to which the Cloud App Control rule must be applied. Maximum of up to `8` locations. When not used it implies `Any` to apply the rule to all locations.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `location_groups` - (List of Numbers) The Name-ID pairs of locations groups to which the Cloud App Control rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `users` - (List of Numbers) The Name-ID pairs of users to which the Cloud App Control rule must be applied. Maximum of up to `4` users. When not used it implies `Any` to apply the rule to all users.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `groups` - (List of Numbers) The Name-ID pairs of groups to which the Cloud App Control rule must be applied. Maximum of up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `departments` - (List of Numbers) The name-ID pairs of the departments that are excluded from the Cloud App Control rule.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `time_windows` - (List of Numbers) The Name-ID pairs of time windows to which the Cloud App Control rule must be applied. Maximum of up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `labels` - (List of Numbers) The Name-ID pairs of rule labels associated to the Cloud App Control rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.

## Important Notes

### Using the Data Source for Actions

**Best Practice**: Always use the `zia_cloud_app_control_rule_actions` data source to retrieve valid actions for your applications. The data source automatically handles:

* Application-specific action support
* Action intersections when multiple applications are configured
* Separation of ISOLATE actions from standard actions

```hcl
data "zia_cloud_app_control_rule_actions" "my_actions" {
  type       = "AI_ML"
  cloud_apps = ["CHATGPT_AI"]
}

resource "zia_cloud_app_control_rule" "example" {
  actions = data.zia_cloud_app_control_rule_actions.my_actions.available_actions_without_isolate
}
```

### ISOLATE Actions Requirements

When using ISOLATE actions:

* ISOLATE actions **cannot be mixed** with other action types (ALLOW, DENY, BLOCK, CAUTION)
* ISOLATE actions **require** `cbi_profile` block with a valid Cloud Browser Isolation profile
* ISOLATE actions **cannot** have `browser_eun_template_id` set
* Create separate rules for ISOLATE vs non-ISOLATE actions

### Multiple Applications

When configuring multiple applications in a single rule, only actions supported by ALL applications are valid. The data source automatically computes this intersection when you specify multiple cloud_apps.

### Action Validation

The resource validates actions during `terraform plan`. If invalid actions are detected, an error message will show:

* Which actions are invalid

* List of valid actions for your configuration
* Suggestion to use the data source

For more information, see the [zia_cloud_app_control_rule_actions](https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_app_control_rule_actions) data source documentation.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Policy access rule can be imported by using `<RULE_TYPE:RULE_ID>` or `<RULE_TYPE:RULE_NAME>` as the import ID.

For example:

```shell
terraform import zia_cloud_app_control_rule.this <rule_type:rule_id>
```

```shell
terraform import zia_cloud_app_control_rule.this <"rule_type:rule_name">
```
