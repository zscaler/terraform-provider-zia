---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_cloud_to_cloud_ir"
description: |-
  Official documentation https://help.zscaler.com/zia/dlp-cloud-cloud-incident-forwarding
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/cloudToCloudIR-get
  Retrieves Cloud-to-Cloud Incident Receiver (C2CIR) information configured in the ZIA Admin Portal
---

# zia_dlp_cloud_to_cloud_ir (Data Source)

* [Official documentation](https://help.zscaler.com/zia/dlp-cloud-cloud-incident-forwarding)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/cloudToCloudIR-get)

Use the **zia_dlp_cloud_to_cloud_ir** data source to get information about Cloud-to-Cloud Incident Receiver (C2CIR) tenants configured in the ZIA Admin Portal. This data source retrieves detailed information about C2CIR configurations including tenant authorization, onboardable entities, and validation status. The retrieved information can be used in Web DLP Rules [zia_dlp_web_rules](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_dlp_web_rules) or CASB DLP Rules [zia_casb_dlp_rules](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_casb_dlp_rules).

## Example Usage

```hcl
# Retrieve the C2CIR by name
data "zia_dlp_cloud_to_cloud_ir" "this" {
  name = "AzureTenant01"
}

# Output the retrieved information
output "zia_dlp_cloud_to_cloud_ir" {
  value = data.zia_dlp_cloud_to_cloud_ir.this
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The name of the Cloud-to-Cloud Incident Receiver tenant to retrieve.

## Attributes Reference

The following attributes are exported:

* `id` - (Number) The unique identifier for the C2CIR tenant.
* `name` - (String) The name of the C2CIR tenant.
* `status` - (List of String) The current status of the C2CIR tenant (e.g., `CASB_TENANT_ACTIVE`).
* `modified_time` - (Number) Timestamp when the C2CIR tenant was last modified.
* `last_tenant_validation_time` - (Number) Timestamp of the last tenant validation.
* `last_validation_msg` - (List) Last validation message information.
  * `error_msg` - (String) Error message from validation.
  * `error_code` - (Number) Error code from validation.
* `last_modified_by` - (List) Information about who last modified the C2CIR tenant.
  * `id` - (Number) Unique identifier for the modifier.
  * `name` - (String) Name of the modifier.
  * `external_id` - (String) External identifier for the modifier.
  * `extensions` - (Map) Additional properties for the modifier.
* `onboardable_entity` - (List) Information about the onboardable entity.
  * `id` - (Number) Unique identifier for the onboardable entity.
  * `name` - (String) Name of the onboardable entity.
  * `type` - (String) Type of the onboardable entity (e.g., `SAAS_TENANT`).
  * `enterprise_tenant_id` - (String) Enterprise tenant ID.
  * `application` - (String) Application name (e.g., `SLACK`).
  * `last_validation_msg` - (List) Last validation message for the onboardable entity.
    * `error_msg` - (String) Error message from validation.
    * `error_code` - (Number) Error code from validation.
  * `tenant_authorization_info` - (List) Tenant authorization information.
    * `access_token` - (String) Access token for authorization.
    * `bot_token` - (String) Bot token for authorization.
    * `redirect_url` - (String) Redirect URL for authorization.
    * `type` - (String) Authorization type (e.g., `SLACK_BOT`).
    * `env` - (String) Environment (e.g., `SALESFORCE_PRODUCTION`).
    * `temp_auth_code` - (String) Temporary authorization code.
    * `subdomain` - (String) Subdomain for the tenant.
    * `apicp` - (String) API CP configuration.
    * `client_id` - (String) Client ID for authorization.
    * `client_secret` - (String) Client secret for authorization.
    * `secret_token` - (String) Secret token for authorization.
    * `user_name` - (String) Username for authorization.
    * `user_pwd` - (String) User password for authorization.
    * `instance_url` - (String) Instance URL for the tenant.
    * `role_arn` - (String) Role ARN for authorization.
    * `quarantine_bucket_name` - (String) Quarantine bucket name.
    * `cloud_trail_bucket_name` - (String) Cloud trail bucket name.
    * `bot_id` - (String) Bot ID for authorization.
    * `org_api_key` - (String) Organization API key.
    * `external_id` - (String) External identifier.
    * `enterprise_id` - (String) Enterprise identifier.
    * `cred_json` - (String) Credential JSON.
    * `role` - (String) Role for authorization (e.g., `READ`).
    * `organization_id` - (String) Organization identifier.
    * `workspace_name` - (String) Workspace name.
    * `workspace_id` - (String) Workspace identifier.
    * `qtn_channel_url` - (String) Quarantine channel URL.
    * `features_supported` - (List of String) Supported features (e.g., `CASB`).
    * `mal_qtn_lib_name` - (String) Malware quarantine library name.
    * `dlp_qtn_lib_name` - (String) DLP quarantine library name.
    * `credentials` - (String) Credentials for authorization.
    * `token_endpoint` - (String) Token endpoint for authorization.
    * `rest_api_endpoint` - (String) REST API endpoint.
    * `smir_bucket_config` - (List) SMIR bucket configuration.
      * `id` - (Number) Unique identifier for the SMIR bucket.
      * `config_name` - (String) Configuration name for the bucket.
      * `metadata_bucket_name` - (String) Metadata bucket name URL.
      * `data_bucket_name` - (String) Data bucket name URL.
    * `qtn_info` - (List) Quarantine information.
      * `admin_id` - (String) Administrator identifier.
      * `qtn_folder_path` - (String) Quarantine folder path.
      * `mod_time` - (Number) Modification time.
    * `qtn_info_cleared` - (Boolean) Whether quarantine information is cleared.
  * `zscaler_app_tenant_id` - (List) Zscaler app tenant ID information.
    * `id` - (Number) Unique identifier for the Zscaler app tenant.
    * `name` - (String) Name of the Zscaler app tenant.
    * `external_id` - (String) External identifier for the Zscaler app tenant.
    * `extensions` - (Map) Additional properties for the Zscaler app tenant
