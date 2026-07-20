---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-endpoint-dlp-policy-rules#Rules
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpRules-get
  Get information about ZIA Endpoint DLP Rules.
---

# zia_endpoint_dlp_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-endpoint-dlp-policy-rules#Rules)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpRules-get)

Use the **zia_endpoint_dlp_rules** data source to get information about a ZIA Endpoint DLP Rules in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP Web Rule by name
data "zia_endpoint_dlp_rules" "example"{
    name = "Rule01"
}
```

```hcl
# Retrieve a DLP Web Rule by ID
data "zia_endpoint_dlp_rules" "example"{
    name = "4522856946"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String) The DLP policy rule name. Used to look up the rule.
* `id` - (Number) The unique identifier of the DLP policy rule. Used to look up the rule.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `order` - (Number) The rule order of execution for the DLP policy rule with respect to other rules.
* `rank` - (Number) Admin rank of the admin who created this rule.
* `state` - (String) Enables or disables the DLP policy rule. Returned values: `ENABLED`, `DISABLED`.
* `description` - (String) The description of the DLP policy rule.
* `action` - (String) The action taken when traffic matches the DLP policy rule criteria. Returned values: `ANY`, `NONE`, `BLOCK`, `CONFIRM`, `ALLOW`.
* `severity` - (String) Indicates the severity selected for the DLP rule violation. Returned values: `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`.
* `data_transfer_method` - (String) The data transfer method to which the DLP policy rule must be applied.
* `network_type` - (String) The network type to which the DLP policy rule must be applied.
* `min_size` - (Number) The minimum file size (in KB) used for evaluation of the DLP policy rule.
* `without_content_inspection` - (Boolean) Indicates a DLP policy rule without content inspection, when the value is set to `true`.
* `external_auditor_email` - (String) The email address of an external auditor to whom DLP email notifications are sent.
* `eun_enabled` - (Boolean) Indicates whether the End User Notification (EUN) is enabled for the DLP policy rule.
* `eun_template_id` - (Number) The unique identifier of the End User Notification (EUN) template.
* `uc_template_id` - (Number) The unique identifier of the user confirmation template.
* `last_modified_time` - (Number) Timestamp when the DLP policy rule was last modified.

* `file_types` - (List of String) The list of file types to which the DLP policy rule must be applied.

* `device_trust_levels` - (List of String) The list of device trust levels for which the rule must be applied. Returned values: `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`.

* `user_risk_score_levels` - (List of String) The user risk score levels to which the DLP policy rule must be applied. Returned values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.

* `parent_rule` - (Number) The unique identifier of the parent rule under which an exception rule is added.

* `sub_rules` - (List of String) The list of exception (sub-) rule IDs added to a parent rule.

* `last_modified_by` - (Block) The admin that modified the DLP policy rule last.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `notification_template` - (Block) The template used for DLP notification emails.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `auditor` - (Block) The auditor to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `dlp_engines` - (Block) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `users` - (Block) The Name-ID pairs of users to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `groups` - (Block) The Name-ID pairs of groups to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `departments` - (Block) The Name-ID pairs of departments to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `devices` - (Block) The Name-ID pairs of devices to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `device_groups` - (Block) The Name-ID pairs of device groups to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `resources` - (Block) The Name-ID pairs of resources to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `resource_groups` - (Block) The Name-ID pairs of resource groups to which the DLP policy rule must be applied.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `labels` - (Block) The Name-ID pairs of rule labels associated to the DLP policy rule.
  * `id` - (Number) Identifier that uniquely identifies an entity.
  * `name` - (String) The configured name of the entity.
  * `extensions` - (Map of String) Additional metadata for the entity.

* `receiver` - (Block) The Zscaler Incident Receiver associated with the DLP policy rule.
  * `id` - (Number) Unique identifier for the receiver.
  * `name` - (String) Name of the receiver.
  * `type` - (String) Type of the receiver.
  * `tenant` - (Block) Tenant information for the receiver.
    * `id` - (Number) Unique identifier for the tenant.
    * `name` - (String) Name of the tenant.
    * `external_id` - (String) External identifier for the tenant.
    * `extensions` - (Map of String) Additional properties for the tenant.

* `end_point_applications` - (Block) The list of endpoint applications to which the DLP policy rule must be applied.
  * `resource_id` - (Number) The unique identifier of the endpoint application.
  * `application_name` - (String) The name of the endpoint application.
  * `application_type` - (String) The type of the endpoint application.
  * `zapp_id` - (String) The Zscaler Client Connector identifier of the endpoint application.
  * `os_type` - (String) The operating system type associated with the endpoint application.
  * `description` - (String) The description of the endpoint application.
  * `bundle_id` - (String) The bundle identifier of the endpoint application.
  * `filename` - (String) The file name of the endpoint application.
  * `original_file_name` - (String) The original file name of the endpoint application.
  * `digitally_signed` - (Boolean) Indicates whether the endpoint application is digitally signed.
  * `deleted` - (Boolean) Indicates whether the endpoint application has been deleted.
  * `mod_uid` - (Number) The modification identifier of the endpoint application.
  * `last_modified_time` - (Number) Timestamp when the endpoint application was last modified.
  * `version` - (Block) The current version details of the endpoint application (see `version` block below).
  * `versions` - (Block) The list of version details of the endpoint application (see `version` block below).

* `end_point_application_groups` - (Block) The list of endpoint application groups to which the DLP policy rule must be applied.
  * `group_id` - (Number) The unique identifier of the endpoint application group.
  * `name` - (String) The name of the endpoint application group.
  * `description` - (String) The description of the endpoint application group.
  * `mod_uid` - (Number) The modification identifier of the endpoint application group.
  * `last_modified_time` - (Number) Timestamp when the endpoint application group was last modified.
  * `end_point_applications` - (Block) The endpoint applications belonging to the group (same structure as `end_point_applications` above).

The `version`/`versions` blocks expose the following attributes:

* `version` - (String) The version string of the endpoint application.
* `z_ver_id_md32` - (Number) The internal version identifier.
* `threat_type` - (Number) The threat type associated with the version.
* `threat_level` - (String) The threat level associated with the version.
* `bundle_id` - (String) The bundle identifier associated with the version.
* `code_signing_certificate_status` - (Number) The code signing certificate status.
* `threat_level_updated` - (Boolean) Indicates whether the threat level was updated.
