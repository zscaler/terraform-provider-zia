---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA: casb_dlp_rules"
description: |-
  Official documentation https://help.zscaler.com/zia/saas-security-api#/casbDlpRules-post
  API documentation https://help.zscaler.com/zia/configuring-data-rest-scanning-dlp-policy
  Adds a new SaaS Security Data at Rest Scanning DLP rule
---

# zia_casb_dlp_rules (Resource)

* [Official documentation](https://help.zscaler.com/zia/configuring-data-rest-scanning-dlp-policy)
* [API documentation](https://help.zscaler.com/zia/saas-security-api#/casbDlpRules-post)

The **zia_casb_dlp_rules** resource Adds a new SaaS Security Data at Rest Scanning DLP rule in the Zscaler Internet Access.

## Example Usage

```hcl
data "zia_casb_tenant" "this" {
  tenant_name = "Jira_Tenant01"
}

data "zia_dlp_incident_receiver_servers" "this" {
  name = "ZS_Incident_Receiver"
}

data "zia_rule_labels" "this" {
    name = "RuleLabel01
}

data "zia_dlp_engines" "this" {
    name = "PCI"
}

data "zia_admin_users" "this" {
    username = auditor01
}

resource "zia_casb_dlp_rules" "this" {
  name = "SaaS_ITSM_App_Rule"
  description = "SaaS_ITSM_App_Rule"
  order = 1
  rank = 7
  type = "OFLCASB_DLP_ITSM"
  action = "OFLCASB_DLP_REPORT_INCIDENT"
  severity = "RULE_SEVERITY_HIGH"
  without_content_inspection = false
  external_auditor_email = "jdoe@acme.com"
  file_types = [
        "FTCATEGORY_APPX",
        "FTCATEGORY_SQL",
  ]
  collaboration_scope = [
        "ANY",
  ]
  components = [
        "COMPONENT_ITSM_OBJECTS",
        "COMPONENT_ITSM_ATTACHMENTS",
  ]
 cloud_app_tenants {
    id = [data.zia_casb_tenant.this.tenant_id]
  }
 dlp_engines {
    id = [data.zia_dlp_engines.this.id]
  }
  object_types {
    id = [32, 33, 34]
  }
 labels {
    id = [data.zia_rule_labels.this.id]
  }
  zscaler_incident_receiver {
    id = data.zia_dlp_incident_receiver_servers.this.id
  }
  auditor_notification {
    id = data.zia_admin_users.this.id
  }
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Optional) Rule name.
* `order` - (Optional) Order of rule execution with respect to other SaaS Security Data at Rest Scanning DLP rules.

* `type` - (Optional) The type of SaaS Security Data at Rest Scanning DLP rule. Supported values:
  * `OFLCASB_DLP_FILE`
  * `OFLCASB_DLP_EMAIL`
  * `OFLCASB_DLP_CRM`
  * `OFLCASB_DLP_ITSM`
  * `OFLCASB_DLP_COLLAB`
  * `OFLCASB_DLP_REPO`
  * `OFLCASB_DLP_STORAGE`
  * `OFLCASB_DLP_GENAI`

* `action` - (Optional) The configured action for the policy rule. Supported values include but are not limited to:
  * `OFLCASB_DLP_REPORT_INCIDENT`
  * `OFLCASB_DLP_SHARE_READ_ONLY`
  * `OFLCASB_DLP_EXTERNAL_SHARE_READ_ONLY`
  * `OFLCASB_DLP_INTERNAL_SHARE_READ_ONLY`
  * `OFLCASB_DLP_REMOVE_PUBLIC_LINK_SHARE`
  * `OFLCASB_DLP_REVOKE_SHARE`
  * `OFLCASB_DLP_REMOVE_EXTERNAL_SHARE`
  * `OFLCASB_DLP_REMOVE_INTERNAL_SHARE`
  * `OFLCASB_DLP_REMOVE_COLLABORATORS`
  * `OFLCASB_DLP_REMOVE_INTERNAL_LINK_SHARE`
  * `OFLCASB_DLP_REMOVE_DISCOVERABLE`
  * `OFLCASB_DLP_NOTIFY_END_USER`
  * `OFLCASB_DLP_APPLY_MIP_TAG`
  * `OFLCASB_DLP_APPLY_BOX_TAG`
  * `OFLCASB_DLP_MOVE_TO_RESTRICTED_FOLDER`
  * `OFLCASB_DLP_REMOVE`
  * `OFLCASB_DLP_QUARANTINE`
  * `OFLCASB_DLP_APPLY_EMAIL_TAG`
  * `OFLCASB_DLP_APPLY_GOOGLEDRIVE_LABEL`
  * `OFLCASB_DLP_REMOVE_EXT_COLLABORATORS`
  * `OFLCASB_DLP_QUARANTINE_TO_USER_ROOT_FOLDER`
  * `OFLCASB_DLP_APPLY_WATERMARK`
  * `OFLCASB_DLP_REMOVE_WATERMARK`
  * `OFLCASB_DLP_APPLY_HEADER`
  * `OFLCASB_DLP_APPLY_FOOTER`
  * `OFLCASB_DLP_APPLY_HEADER_FOOTER`
  * `OFLCASB_DLP_REMOVE_HEADER`
  * `OFLCASB_DLP_REMOVE_FOOTER`
  * `OFLCASB_DLP_REMOVE_HEADER_FOOTER`
  * `OFLCASB_DLP_BLOCK`
  * `OFLCASB_DLP_APPLY_ATLASSIAN_CLASSIFICATION_LABEL`
  * `OFLCASB_DLP_ALLOW`
  * `OFLCASB_DLP_REDACT`

* `cloud_app_tenants` - (Block List) Name-ID pairs of the cloud application tenants for which the rule is applied.
  * `id` - (int) A unique identifier for an entity.

### Optional

* `description` - (Optional) An admin editable text-based description of the rule.
* `rank` - (Optional) Admin rank that is assigned to this rule. Mandatory when admin rank-based access restriction is enabled.
* `state` - (Optional) Administrative state of the rule. Supported values: `ENABLED`, `DISABLED`. Defaults to `ENABLED`.
* `severity` - (Optional) The severity level of the incidents that match the policy rule. Supported values:
  * `RULE_SEVERITY_HIGH`
  * `RULE_SEVERITY_MEDIUM`
  * `RULE_SEVERITY_LOW`
  * `RULE_SEVERITY_INFO`

* `bucket_owner` - (Optional) A user who inspect their buckets for sensitive data. When you choose a user, their buckets are available in the Buckets field.
* `external_auditor_email` - (Optional) Email address of the external auditor to whom the DLP email alerts are sent.

* `content_location` - (Optional) The location for the content that the Zscaler service inspects for sensitive data. Supported values:
  * `CONTENT_LOCATION_PRIVATE_CHANNEL`
  * `CONTENT_LOCATION_PUBLIC_CHANNEL`
  * `CONTENT_LOCATION_SHARED_CHANNEL`
  * `CONTENT_LOCATION_DIRECT_MESSAGE`
  * `CONTENT_LOCATION_MULTI_PERSON_DIRECT_MESSAGE`

* `collaboration_scope` - (Optional, Set of String) Collaboration scope for the rule. Supported values:
  * `ANY`
  * `COLLABORATION_SCOPE_EXTERNAL_COLLAB_VIEW`
  * `COLLABORATION_SCOPE_EXTERNAL_COLLAB_EDIT`
  * `COLLABORATION_SCOPE_EXTERNAL_LINK_VIEW`
  * `COLLABORATION_SCOPE_EXTERNAL_LINK_EDIT`
  * `COLLABORATION_SCOPE_INTERNAL_COLLAB_VIEW`
  * `COLLABORATION_SCOPE_INTERNAL_COLLAB_EDIT`
  * `COLLABORATION_SCOPE_INTERNAL_LINK_VIEW`
  * `COLLABORATION_SCOPE_INTERNAL_LINK_EDIT`
  * `COLLABORATION_SCOPE_PRIVATE_EDIT`
  * `COLLABORATION_SCOPE_PRIVATE`
  * `COLLABORATION_SCOPE_PUBLIC`

* `components` - (Optional, Set of String) Collaboration scope for the rule. Supported values:
  * `ANY`
  * ``COMPONENT_EMAIL_BODY`
  * ``COMPONENT_EMAIL_ATTACHMENT`
  * ``COMPONENT_EMAIL_SUBJECT`
  * ``COMPONENT_ITSM_OBJECTS`
  * ``COMPONENT_ITSM_ATTACHMENTS`
  * ``COMPONENT_CRM_CHATTER_MESSAGES`
  * ``COMPONENT_CRM_ATTACHMENTS_IN_OBJECTS`
  * ``COMPONENT_COLLAB_MESSAGES`
  * ``COMPONENT_COLLAB_ATTACHMENTS`
  * ``COMPONENT_CRM_CASES`
  * ``COMPONENT_GENAI_MESSAGES`
  * ``COMPONENT_GENAI_ATTACHMENTS`
  * ``COMPONENT_FILE_ATTACHMENTS`

* `file_types` - (Optional) File types for which the rule is applied. If not set, the rule is applied across all file types. See supported list [Configuring the Data at Rest Scanning DLP Policy](https://help.zscaler.com/zia/saas-security-api#/casbDlpRules-post)

* `recipient` - (Optional) Specifies if the email recipient is internal or external.
* `quarantine_location` - (Optional) Location where all the quarantined files are moved and necessary actions are taken by either deleting or restoring the data.
* `watermark_delete_old_version` - (Optional, Boolean) Specifies whether to delete an old version of the watermarked file.
* `include_criteria_domain_profile` - (Optional, Boolean) If true, criteriaDomainProfiles is included as part of the criteria.
* `include_email_recipient_profile` - (Optional, Boolean) If true, emailRecipientProfiles is included as part of the criteria.
* `without_content_inspection` - (Optional, Boolean) If true, Content Matching is set to None.
* `include_entity_groups` - (Optional, Boolean) If true, entityGroups is included as part of the criteria.
* `domains` - (Optional, Set of String) The domain for the external organization sharing the channel.

* `entity_groups` - (Block List) Name-ID pairs of entity groups that are part of the rule criteria.
  * `id` - (int) A unique identifier for an entity.

* `included_domain_profiles` - (Block List) Name-ID pairs of domain profiles included in the criteria for the rule.
  * `id` - (int) A unique identifier for an entity.

* `excluded_domain_profiles` - (Block List) Name-ID pairs of domain profiles excluded in the criteria for the rule.
  * `id` - (int) A unique identifier for an entity.

* `criteria_domain_profiles` - (Block List) Name-ID pairs of domain profiles that are mandatory in the criteria for the rule.
  * `id` - (int) A unique identifier for an entity.

* `email_recipient_profiles` - (Block List) Name-ID pairs of recipient profiles for which the rule is applied.
  * `id` - (int) A unique identifier for an entity.

* `object_types` - (Block List) List of object types for which the rule is applied.
  * `id` - (int) A unique identifier for an entity.

* `labels` - (Block List) Name-ID pairs of rule labels associated with the rule.
  * `id` - (int) A unique identifier for an entity.

* `buckets` - (Block List) The buckets for the Zscaler service to inspect for sensitive data.
  * `id` - (int) A unique identifier for an entity.

* `groups` - (Block List) Name-ID pairs of groups for which the rule is applied.
  * `id` - (int) A unique identifier for an entity.

* `departments` - (Block List) Name-ID pairs of departments for which the rule is applied.
  * `id` - (int) A unique identifier for an entity.

* `users` - (Block List) Name-ID pairs of users for which rule is applied.
  * `id` - (int) A unique identifier for an entity.

* `dlp_engines` - (Block List) The list of DLP engines to which the DLP policy rule must be applied.
  * `id` - (int) Identifier that uniquely identifies an entity.

* `zscaler_incident_receiver` - (Block, Max: 1) The Zscaler Incident Receiver details.
  * `id` - (int) A unique identifier for an entity.

* `auditor_notification` - (Block, Max: 1) Notification template used for DLP email alerts sent to the auditor.
  * `id` - (int) A unique identifier for an entity.

* `tag` - (Block, Max: 1) Tag applied to the rule.
  * `id` - (int) A unique identifier for an entity.

* `watermark_profile` - (Block, Max: 1) Watermark profile applied to the rule.
  * `id` - (int) A unique identifier for an entity.

* `redaction_profile` - (Block, Max: 1) Name-ID of the redaction profile in the criteria.
  * `id` - (int) A unique identifier for an entity.

* `casb_email_label` - (Block, Max: 1) Name-ID of the email label associated with the rule.
  * `id` - (int) A unique identifier for an entity.

* `casb_tombstone_template` - (Block, Max: 1) Name-ID of the quarantine tombstone template associated with the rule.
  * `id` - (int) A unique identifier for an entity.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_casb_dlp_rules** can be imported by using `<RULE_TYPE:RULE_ID>` or `<RULE_TYPE:RULE_NAME>` as the import ID.

For example:

```shell
terraform import zia_casb_dlp_rules.this <rule_type:rule_id>
```

```shell
terraform import zia_casb_dlp_rules.this <"rule_type:rule_name">
```
