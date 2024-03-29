---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_web_rules"
description: |-
  Creates and manages ZIA DLP Web Rules.
---

# Resource: zia_dlp_web_rules

The **zia_dlp_web_rules** resource allows the creation and management of ZIA DLP Web Rules in the Zscaler Internet Access cloud or via the API.

⚠️ **WARNING:** Zscaler Internet Access DLP supports a maximum of 127 Web DLP Rules to be created via API.

## Example Usage - OCR ENABLED

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
    file_types                = [ "BITMAP", "JPEG", "PNG", "TIFF"]
    min_size                    = 20
    zscaler_incident_receiver   = true
}
```

## Example Usage - "ALL_OUTBOUND" File Type

```hcl
data "zia_dlp_engines" "this" {
  predefined_engine_name = "EXTERNAL"
}

resource "zia_dlp_web_rules" "this" {
  name                       = "Example"
  description                = "Example"
  action                     = "BLOCK"
  order                      = 1
  rank                       = 7
  state                      = "ENABLED"
  # ocr_enabled              = true
  protocols                  = [ "FTP_RULE", "HTTPS_RULE", "HTTP_RULE" ]
  file_types                 = [ "ALL_OUTBOUND" ]
  zscaler_incident_receiver  = true
  without_content_inspection = false
  user_risk_score_levels     = [ "LOW", "MEDIUM", "HIGH", "CRITICAL" ]
  severity                   = "RULE_SEVERITY_HIGH"
  dlp_engines {
    id = [ data.zia_dlp_engines.this.id ]
  }
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
  * `without_content_inspection` must be set to false if `file_types` is not defined.

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

  * ~> Note: `BITMAP`, `JPEG`, `PNG`, and `TIFF` file types are exclusively supported when optical character recognition `ocr_enabled` is set to `true` for DLP rules with content inspection.

  * ~> Note: `ALL_OUTBOUND` file type is applicable only when the predefined DLP engine called `EXTERNAL` is used and when the attribute `without_content_inspection` is set to `false`.

  * ~> Note: `ALL_OUTBOUND` file type cannot be used alongside any any other file type.

* `severity` - (String) Indicates the severity selected for the DLP rule violation: Returned values are:  `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`

* `user_risk_score_levels` (Optional) - Indicates the user risk score level selectedd for the DLP rule violation: Returned values are: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

* `parent_rule`(Optional) - The unique identifier of the parent rule under which an exception rule is added.
 ~> Note: Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings.

* `sub_rules`(List) - The list of exception rules added to a parent rule.
 ~> Note: All attributes within the WebDlpRule model are applicable to the sub-rules. Values for each rule are specified by using the WebDlpRule object Exception rules can be configured only when the inline DLP rule evaluation type is set to evaluate all DLP rules in the DLP Advanced Settings.

* `notification_template` - (Optional) The template used for DLP notification emails.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `auditor` - (Optional) The auditor to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `url_categories` - (Optional) The list of URL categories to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity

* `dlp_engines` - (Optional) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Optional) Identifier that uniquely identifies an entity. Maximum of up to `4` dlp engines. When not used it implies `Any` to apply the rule to all locations.

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

* `workload_groups` (Optional) The list of preconfigured workload groups to which the policy must be applied
  * `id` - (Optional) A unique identifier assigned to the workload group
  * `name` - (Optional) The name of the workload group

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_dlp_web_rules** can be imported by using `<RULE ID>` or `<RULE NAME>` as the import ID.

For example:

```shell
terraform import zia_dlp_web_rules.example <rule_id>
```

or

```shell
terraform import zia_dlp_web_rules.example <rule_name>
```
