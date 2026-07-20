---
subcategory: "Outbound Email DLP Policy"
layout: "zscaler"
page_title: "ZIA: outbound_email_dlp"
description: |-
  Official documentation https://help.zscaler.com/zia/about-email-dlp-rules
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/emailDlpRules-get
  Creates a new Outbound Email DLP policy rule
---

# zia_outbound_email_dlp (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-email-dlp-rules)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/emailDlpRules-get)

The **zia_outbound_email_dlp** resource allows the creation and management of ZIA Outbound Email DLP policy rules in the Zscaler Internet Access cloud or via the API.

~> **NOTE:** Rule orders must always be contiguous (no gaps). Deleting a rule must be followed by order number re-adjustment of the remaining rules to ensure the API honours the required order.

~> **NOTE:** The `order` attribute must always be a positive whole number starting at 1. Negative numbers and zero are **not supported** and will result in an error.

## Example Usage

```hcl
resource "zia_outbound_email_dlp" "this" {
  name                       = "Email_DLP_Rule_01"
  description                = "Email_DLP_Rule_01"
  action                     = "BLOCK"
  state                      = "ENABLED"
  order                      = 1
  severity                   = "RULE_SEVERITY_HIGH"
  without_content_inspection = false
  min_size                   = 10
  file_types                 = ["FTCATEGORY_SCZIP", "FTCATEGORY_ACCDB"]
  user_risk_score_levels     = ["HIGH", "CRITICAL"]

  dlp_engines {
    id = [3, 4]
  }

  users {
    id = [165214882]
  }

  groups {
    id = [165437858]
  }

  departments {
    id = [68759309]
  }

  labels {
    id = [4900918]
  }

  auditor {
    id = 43859078
  }

  notification_template {
    id = 5855
  }

  receiver {
    id = "1664"
  }
}
```

## Example Usage - "Configuring an External Auditor"

```hcl
resource "zia_outbound_email_dlp" "this" {
  name                   = "Email_DLP_Rule_02"
  description            = "Email_DLP_Rule_02"
  action                 = "BLOCK"
  state                  = "ENABLED"
  order                  = 1
  severity               = "RULE_SEVERITY_HIGH"
  external_auditor_email = "auditor@acme.com"

  dlp_engines {
    id = [3]
  }

  notification_template {
    id = 5855
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The DLP policy rule name.
* `order` - (Number) The rule order of execution for the DLP policy rule with respect to other rules. Must be a positive whole number starting at `1`.

### Optional

* `description` - (String) The description of the DLP policy rule.
* `state` - (String) Enables or disables the DLP policy rule. Supported values: `ENABLED`, `DISABLED`.
* `action` - (String) The action taken when traffic matches the DLP policy rule criteria. Supported values include: `BLOCK`, `ALLOW`, `CUSTOMHEADERINSERTION`.
* `severity` - (String) Indicates the severity selected for the DLP rule violation. Supported values: `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`.
* `min_size` - (Number) The minimum file size (in KB) used for evaluation of the DLP policy rule.
* `without_content_inspection` - (Boolean) Indicates a DLP policy rule without content inspection, when the value is set to `true`.
* `external_auditor_email` - (String) The email address of an external auditor to whom DLP email notifications are sent. Cannot be combined with `auditor`; when set, `notification_template` must also be set.
* `custom_header` - (String) The custom header value inserted when the rule action includes custom header insertion.
* `parent_rule` - (Number) The unique identifier of the parent rule under which an exception (sub-) rule is added. Use the parent resource's `rule_id` (integer).

* `file_types` - (List of String) The list of file type categories to which the DLP policy rule must be applied. For the complete list of supported file type categories you can use the data source [zia_file_type_categories](https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_file_type_categories).

* `content_locations` - (List of String) The list of content locations to which the DLP policy rule must be applied.

* `user_risk_score_levels` - (List of String) The user risk score levels to which the DLP policy rule must be applied. Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.

* `sub_rules` - (List of String) Set of sub-rule IDs (strings), populated from the API for a **parent** rule after read. Sub-rules are managed as their own `zia_outbound_email_dlp` resources with `parent_rule` set; do not model sub-rules as nested blocks in the parent.

* `dlp_engines` - (Block) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `users` - (Block) The Name-ID pairs of users to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `groups` - (Block) The Name-ID pairs of groups to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `departments` - (Block) The Name-ID pairs of departments to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `excluded_users` - (Block) The Name-ID pairs of users that are excluded from the DLP policy rule.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `excluded_groups` - (Block) The Name-ID pairs of groups that are excluded from the DLP policy rule.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `excluded_departments` - (Block) The Name-ID pairs of departments that are excluded from the DLP policy rule.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `time_windows` - (Block) The Name-ID pairs of time windows to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `labels` - (Block) The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `included_domain_profiles` - (Block) The Name-ID pairs of domain profiles included in the DLP policy rule.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `email_tenants` - (Block) The Name-ID pairs of email tenants to which the DLP policy rule must be applied.
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `email_recipient_profiles` - (Block) The Name-ID pairs of email recipient profiles to which the DLP policy rule must be applied. To retrieve the list of Email Recipient Profiles use the data source [zia_email_profile](https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_email_profile).
  * `id` - (List of Number) Identifier that uniquely identifies an entity.

* `auditor` - (Block) The auditor to which the DLP policy rule must be applied. Cannot be combined with `external_auditor_email`; when set, `notification_template` must also be set.
  * `id` - (Number) Identifier that uniquely identifies an entity.

* `notification_template` - (Block) The template used for DLP notification emails. To retrieve the list of Notification Templates use the resource or data source [zia_dlp_notification_templates](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_dlp_notification_templates).
  * `id` - (Number) Identifier that uniquely identifies an entity.

* `receiver` - (Block, Max 1) The Zscaler Incident Receiver associated with the DLP policy rule. Only the `id` must be set; the receiver `type` is always sent as `ZIR` automatically.
  * `id` - (String) The unique identifier of the Zscaler Incident Receiver.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_outbound_email_dlp** can be imported using either:

* `<RULE ID>` — numeric ID of the rule, or
* `<RULE NAME>` — exact rule name.

For example:

```shell
terraform import zia_outbound_email_dlp.example <rule_id>
```

or

```shell
terraform import zia_outbound_email_dlp.example <rule_name>
```
