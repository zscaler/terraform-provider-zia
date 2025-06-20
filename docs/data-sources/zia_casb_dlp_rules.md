---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA): casb_dlp_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/saas-security-api#/casbDlpRules-post
  API documentation https://help.zscaler.com/zia/configuring-data-rest-scanning-dlp-policy
  Retrieves the SaaS Security Data at Rest Scanning Data Loss Prevention (DLP) rules based on the specified rule type.

---
# zia_casb_dlp_rules (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-data-rest-scanning-dlp-policy)
* [API documentation](https://help.zscaler.com/zia/saas-security-api#/casbDlpRules-post)

Use the **zia_casb_dlp_rules** data source to get information about SaaS Security Data at Rest Scanning Data Loss Prevention (DLP) rules based on the specified rule type.

## Example Usage - By Name

```hcl

data "zia_casb_dlp_rules" "this" {
  name = "SaaS_ITSM_App_Rule"
  type = "OFLCASB_DLP_ITSM"
}
```

## Example Usage - By ID

```hcl

data "zia_casb_dlp_rules" "this" {
  id = 154658
  type = "OFLCASB_DLP_ITSM"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) System-generated identifier for the SaaS Security Data at Rest Scanning DLP rule.
* `name` - (Optional) Rule name.
* `type` - (Optional) The type of SaaS Security Data at Rest Scanning DLP rule.
  * `OFLCASB_DLP_FILE`
  * `OFLCASB_DLP_EMAIL`
  * `OFLCASB_DLP_CRM`
  * `OFLCASB_DLP_ITSM`
  * `OFLCASB_DLP_COLLAB`
  * `OFLCASB_DLP_REPO`
  * `OFLCASB_DLP_STORAGE`
  * `OFLCASB_DLP_GENAI`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `order` - (int) Order of rule execution with respect to other SaaS Security Data at Rest Scanning DLP rules.
* `rank` - (int) Rank of the rule.
* `last_modified_time` - (int) Last modification time of the rule.
* `state` - (string) Administrative state of the rule.
* `action` - (string) The configured action for the policy rule.
* `severity` - (string) The severity level of the incidents that match the policy rule.
* `description` - (string) An admin editable text-based description of the rule.
* `bucket_owner` - (string) A user who inspects their buckets for sensitive data.
* `external_auditor_email` - (string) Email address of the external auditor to whom the DLP email alerts are sent.
* `content_location` - (string) The location for the content that the Zscaler service inspects for sensitive data.
* `number_of_internal_collaborators` - (string) Number of internal collaborators for files shared within an organization.
* `number_of_external_collaborators` - (string) Number of external collaborators for files shared outside of an organization.
* `recipient` - (string) Specifies if the email recipient is internal or external.
* `quarantine_location` - (string) Location where all quarantined files are moved for action.
* `access_control` - (string) Access privilege of this rule based on the admin's RBA state.
* `watermark_delete_old_version` - (bool) Specifies whether to delete an old version of the watermarked file.
* `include_criteria_domain_profile` - (bool) If true, `criteriaDomainProfiles` is included in the criteria.
* `include_email_recipient_profile` - (bool) If true, `emailRecipientProfiles` is included in the criteria.
* `without_content_inspection` - (bool) If true, Content Matching is set to None.
* `include_entity_groups` - (bool) If true, `entityGroups` is included in the criteria.
* `components` - (List of String) List of components for which the rule is applied.
* `collaboration_scope` - (List of String) Collaboration scope for the rule.
* `domains` - (List of String) Domain for the external organization sharing the channel.
* `file_types` - (List of String) File types to which the rule is applied.

### Block Attributes

Each of the following blocks supports nested attributes:

#### `cloud_app_tenants`

* `id` - (int) A unique identifier for a tenant.
* `name` - (string) The configured name of the tenant.

#### `entity_groups`

* `id` - (int) A unique identifier for the entity group.
* `name` - (string) The configured name of the entity group.

#### `included_domain_profiles`, `excluded_domain_profiles`, `criteria_domain_profiles`

* `id` - (int) A unique identifier for the domain profile.
* `name` - (string) The configured name of the domain profile.

#### `email_recipient_profiles`

* `id` - (int) A unique identifier for the email recipient profile.
* `name` - (string) The configured name of the email recipient profile.

#### `buckets`

* `id` - (int) A unique identifier for the bucket.
* `name` - (string) The configured name of the bucket.

#### `object_types`

* `id` - (int) A unique identifier for the object type.
* `name` - (string) The configured name of the object type.

#### `dlp_engines`

* `id` - (int) Identifier that uniquely identifies the DLP engine.
* `name` - (string) Name of the DLP engine.
* `extensions` - (Map of String) Optional metadata for the DLP engine.

#### `labels`

* `id` - (int) A unique identifier for the label.
* `name` - (string) The configured name of the label.

#### `groups`

* `id` - (int) A unique identifier for the group.
* `name` - (string) The configured name of the group.

#### `departments`

* `id` - (int) A unique identifier for the department.
* `name` - (string) The configured name of the department.

#### `users`

* `id` - (int) A unique identifier for the user.
* `name` - (string) The configured name of the user.

#### `zscaler_incident_receiver`

* `id` - (int) A unique identifier for the incident receiver.
* `name` - (string) The configured name of the incident receiver.

#### `auditor_notification`

* `id` - (int) A unique identifier for the notification.
* `name` - (string) The configured name of the notification.

#### `tag`

* `id` - (int) A unique identifier for the tag.
* `name` - (string) The configured name of the tag.

#### `watermark_profile`

* `id` - (int) A unique identifier for the watermark profile.
* `name` - (string) The configured name of the watermark profile.

#### `redaction_profile`

* `id` - (int) A unique identifier for the redaction profile.
* `name` - (string) The configured name of the redaction profile.

#### `casb_email_label`

* `id` - (int) A unique identifier for the email label.
* `name` - (string) The configured name of the email label.

#### `casb_tombstone_template`

* `id` - (int) A unique identifier for the tombstone template.
* `name` - (string) The configured name of the tombstone template.
