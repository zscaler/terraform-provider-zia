---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_web_rules"
description: |-
  Get information about ZIA DLP Web Rules.
---

# Data Source: zia_dlp_web_rules

Use the **zia_dlp_web_rules** data source to get information about a ZIA DLP Web Rules in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP Web Rule by name
data "zia_dlp_web_rules" "example"{
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP policy rule name.
rules.

### Optional

* `description` - (String) The description of the DLP policy rule.
* `order` - (Number) The rule order of execution for the DLP policy rule with respect to other
* `external_auditor_email` - (String) The email address of an external auditor to whom DLP email notifications are sent.
* `match_only` - (Bool) The match only criteria for DLP engines.
* `without_content_inspection` - (Bool) Indicates a DLP policy rule without content inspection, when the value is set to true.
* `ocr_enabled` - (Bool) Enables or disables image file scanning.
* `zscaler_incident_receiver` - (Bool) Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.
* `last_modified_time` - (Number) Timestamp when the DLP policy rule was last modified.

* `access_control` - (String) The access privilege for this DLP policy rule based on the admin's state. The supported values are:
  * `NONE`
  * `READ_ONLY`
  * `READ_WRITE`

* `action` - (String) The action taken when traffic matches the DLP policy rule criteria. The supported values are:
  * `ANY`
  * `NONE`
  * `BLOCK`
  * `ALLOW`
  * `ICAP_RESPONSE`

* `state` - (String) Enables or disables the DLP policy rule.. The supported values are:
  * `DISABLED`
  * `ENABLED`

* `file_types` - (String) The list of file types to which the DLP policy rule must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)
* `cloud_applications` - (Optional) The list of cloud applications to which the DLP policy rule must be applied. For the complete list of supported cloud applications refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)

* `last_modified_by` - (Number)  The admin that modified the DLP policy rule last.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `notification_template` - (Optional) The template used for DLP notification emails.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `auditor` - (Optional) The auditor to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `url_categories` - (Optional) The list of URL categories to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `dlp_engines` - (Optional) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `locations` - (Optional) The Name-ID pairs of locations to which the DLP policy rule must be applied. Maximum of up to `8` locations. When not used it implies `Any` to apply the rule to all locations.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `location_groups` - (Optional) The Name-ID pairs of locations groups to which the DLP policy rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `users` - (Optional) The Name-ID pairs of users to which the DLP policy rule must be applied. Maximum of up to `4` users. When not used it implies `Any` to apply the rule to all users.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `excluded_users` - (Optional) The name-ID pairs of the users that are excluded from the DLP policy rule. Maximum of up to `256` users.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `groups` - (Optional) The Name-ID pairs of groups to which the DLP policy rule must be applied. Maximum of up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `excluded_groups` - (Optional) The name-ID pairs of the groups that are excluded from the DLP policy rule. Maximum of up to `256` groups.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `departments` - (Optional) The name-ID pairs of the departments that are excluded from the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `excluded_departments` - (Optional) The name-ID pairs of the groups that are excluded from the DLP policy rule. Maximum of up to `256` departments.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `time_windows` - (Optional) The Name-ID pairs of time windows to which the DLP policy rule must be applied. Maximum of up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `labels` - (Optional) The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `icap_server` - (Optional) The DLP server, using ICAP, to which the transaction content is forwarded.
  * `id` - (Number) Identifier that uniquely identifies an entity
