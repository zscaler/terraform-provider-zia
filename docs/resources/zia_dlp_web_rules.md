---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_web_rules"
description: |-
  Creates and manages ZIA DLP Web Rules.
---

# Resource: zia_dlp_web_rules

The **zia_dlp_web_rules** resource allows the creation and management of ZIA DLP Web Rules in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
resource "zia_dlp_web_rules" "test" {
    name                        = "Test"
    description                 = "Test"
    action                      = "ALLOW"
    state                       = "ENABLED"
    order                       = 1
    rank                        = 7
    protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
    cloud_applications          = ["ZENDESK", "LUCKY_ORANGE", "MICROSOFT_POWERAPPS", "MICROSOFTLIVEMEETING"]
    without_content_inspection  = false
    match_only                  = false
    ocr_enabled                 = true
    file_types                = [ "WINDOWS_META_FORMAT", "BITMAP", "JPEG", "PNG", "TIFF"]
    min_size                    = 20
    zscaler_incident_receiver   = true
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP policy rule name.
* `order` - (Required) The rule order of execution for the DLP policy rule with respect to other rules.

### Optional

* `description` - (Optional) The description of the DLP policy rule.
* `external_auditor_email` - (Optional) The email address of an external auditor to whom DLP email notifications are sent.
* `match_only` - (Optional) The match only criteria for DLP engines.
* `without_content_inspection` - (Optional) Indicates a DLP policy rule without content inspection, when the value is set to true.
* `ocr_enabled` - (Optional) Enables or disables image file scanning. When OCR is enabled only the following ``file_types`` are supported: ``WINDOWS_META_FORMAT``, ``BITMAP``, ``JPEG``, ``PNG``, ``TIFF``
* `zscaler_incident_receiver` - (Optional) Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.

* `action` - (Optional) The action taken when traffic matches the DLP policy rule criteria. The supported values are:
  * `ANY`
  * `NONE`
  * `BLOCK`
  * `ALLOW`
  * `ICAP_RESPONSE`

* `state` - (Optional) Enables or disables the DLP policy rule.. The supported values are:
  * `DISABLED`
  * `ENABLED`

* `file_types` - (Optional) The list of file types to which the DLP policy rule must be applied. For the complete list of supported file types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)
* `cloud_applications` - (Optional) The list of cloud applications to which the DLP policy rule must be applied. For the complete list of supported cloud applications refer to the  [ZIA API documentation](https://help.zscaler.com/zia/data-loss-prevention#/webDlpRules-post)

* `notification_template` - (Optional) The template used for DLP notification emails.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `auditor` - (Optional) The auditor to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `url_categories` - (Optional) The list of URL categories to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `dlp_engines` - (Optional) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `locations` - (Optional) The Name-ID pairs of locations to which the DLP policy rule must be applied. Maximum of up to `8` locations. When not used it implies `Any` to apply the rule to all locations.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `location_groups` - (Optional) The Name-ID pairs of locations groups to which the DLP policy rule must be applied. Maximum of up to `32` location groups. When not used it implies `Any` to apply the rule to all location groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `users` - (Optional) The Name-ID pairs of users to which the DLP policy rule must be applied. Maximum of up to `4` users. When not used it implies `Any` to apply the rule to all users.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `groups` - (Optional) The Name-ID pairs of groups to which the DLP policy rule must be applied. Maximum of up to `8` groups. When not used it implies `Any` to apply the rule to all groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `departments` - (Optional) The name-ID pairs of the departments that are excluded from the DLP policy rule.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `excluded_users` - (Optional) The name-ID pairs of the users that are excluded from the DLP policy rule. Maximum of up to `256` users.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `excluded_departments` - (Optional) The name-ID pairs of the groups that are excluded from the DLP policy rule. Maximum of up to `256` departments.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `excluded_groups` - (Optional) The name-ID pairs of the groups that are excluded from the DLP policy rule. Maximum of up to `256` groups.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `time_windows` - (Optional) The Name-ID pairs of time windows to which the DLP policy rule must be applied. Maximum of up to `2` time intervals. When not used it implies `always` to apply the rule to all time intervals.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `labels` The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `icap_server` The DLP server, using ICAP, to which the transaction content is forwarded.
  * `id` - (Optional) Identifier that uniquely identifies an entity
