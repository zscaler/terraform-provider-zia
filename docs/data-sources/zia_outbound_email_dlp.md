---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: outbound_email_dlp"
description: |-
  Official documentation https://help.zscaler.com/zia/about-email-dlp-rules
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/emailDlpRules-get
  Get information about ZIA Outbound Email DLP Rules.
---

# zia_outbound_email_dlp (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-email-dlp-rules)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/emailDlpRules-get)

Use the **zia_outbound_email_dlp** data source to get information about a ZIA Outbound Email DLP Rule in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve an Outbound Email DLP Rule by name
data "zia_outbound_email_dlp" "example" {
  name = "Email_DLP_Rule_01"
}
```

```hcl
# Retrieve an Outbound Email DLP Rule by ID
data "zia_outbound_email_dlp" "example" {
  id = 4522856946
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String) The DLP policy rule name. Used to look up the rule.
* `id` - (Number) The unique identifier of the DLP policy rule. Used to look up the rule.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `order` - (Number) The rule order of execution for the DLP policy rule with respect to other rules.
* `state` - (String) Enables or disables the DLP policy rule. Returned values: `ENABLED`, `DISABLED`.
* `description` - (String) The description of the DLP policy rule.
* `action` - (String) The action taken when traffic matches the DLP policy rule criteria.
* `severity` - (String) Indicates the severity selected for the DLP rule violation.
* `min_size` - (Number) The minimum file size (in KB) used for evaluation of the DLP policy rule.
* `without_content_inspection` - (Boolean) Indicates a DLP policy rule without content inspection, when the value is set to `true`.
* `external_auditor_email` - (String) The email address of an external auditor to whom DLP email notifications are sent.
* `custom_header` - (String) The custom header value inserted when the rule action includes custom header insertion.
* `parent_rule` - (Number) The unique identifier of the parent rule under which an exception rule is added.
* `last_modified_time` - (Number) Timestamp when the DLP policy rule was last modified.
* `file_types` - (List of String) The list of file type categories to which the DLP policy rule must be applied.
* `content_locations` - (List of String) The list of content locations to which the DLP policy rule must be applied.
* `user_risk_score_levels` - (List of String) The user risk score levels to which the DLP policy rule must be applied.
* `sub_rules` - (List of String) The list of exception (sub-) rule IDs added to a parent rule.

* `dlp_engines` - (List of Object) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `users` - (List of Object) The Name-ID pairs of users to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `groups` - (List of Object) The Name-ID pairs of groups to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `departments` - (List of Object) The Name-ID pairs of departments to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `excluded_users` - (List of Object) The Name-ID pairs of users that are excluded from the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `excluded_groups` - (List of Object) The Name-ID pairs of groups that are excluded from the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `excluded_departments` - (List of Object) The Name-ID pairs of departments that are excluded from the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `time_windows` - (List of Object) The Name-ID pairs of time windows to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `labels` - (List of Object) The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `included_domain_profiles` - (List of Object) The Name-ID pairs of domain profiles included in the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `email_tenants` - (List of Object) The Name-ID pairs of email tenants to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `email_recipient_profiles` - (List of Object) The Name-ID pairs of email recipient profiles to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `auditor` - (List of Object) The auditor to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `notification_template` - (List of Object) The template used for DLP notification emails.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.

* `receiver` - (List of Object) The Zscaler Incident Receiver associated with the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies the receiver.
  * `name` - (String) The configured name of the receiver.
  * `type` - (String) The type of the receiver.
  * `tenant` - (List of Object) Tenant information for the receiver.
    * `id` - (Number) Identifier that uniquely identifies the tenant.
    * `name` - (String) The configured name of the tenant.
    * `external_id` - (String) External identifier for the tenant.
    * `extensions` - (Map of String) Additional tenant attributes.

* `last_modified_by` - (List of Object) The admin that modified the DLP policy rule last.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
