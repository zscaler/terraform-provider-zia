---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA: casb_tenant"
description: |-
  Official documentation https://help.zscaler.com/zia/about-saas-application-tenants
  API documentation https://help.zscaler.com/zia/saas-security-api#/casbTenant/lite-get
  Retrieves information about the SaaS application tenant
---

# zia_casb_tenant (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-saas-application-tenants)
* [API documentation](https://help.zscaler.com/zia/saas-security-api#/casbTenant/lite-get)

Use the **zia_casb_tenant** data source to get information about a ZIA SaaS Application Tenants in the Zscaler Internet Access cloud or via the API.

## Example Usage - By Name

```hcl
data "zia_casb_tenant" "this"{
    tenant_name = "Bitbucket"
}
```

## Example Usage - By ID

```hcl
data "zia_casb_tenant" "this"{
    tenant_id = "11743520"
}
```

## Example Usage - Use Optional Parameters

```hcl
data "zia_casb_tenant" "this"{
    tenant_name = "Bitbucket"
    active_only = true
    app = "BITBUCKET"
    filter_by_feature = ["CASB", "SSPM"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `tenant_name` - (Optional) Tenant Name
* `tenant_id` - (Optional) Tenant ID

### Optional Filtering Options

* `active_only` - (Boolean) Indicates that the tenant is in use. Policies are enforced for this SaaS application.
* `include_deleted` - (Boolean) Include Deleted tenants
* `app_type` - (String) Specifies the SaaS application type [Available values](https://help.zscaler.com/zia/saas-security-api#/casbTenant/lite-get)
* `app` - (String) Specifies the sanctioned SaaS application [Available values](https://help.zscaler.com/zia/saas-security-api#/casbTenant/lite-get)
* `scan_config_tenants_only` - (Boolean) Specifies the tenant for which the scan is already configured
* `include_bucket_ready_s3_tenants` - (Boolean) For the AWS S3 SaaS application, this parameter indicates that the buckets have been read and are ready for use in policies and scan configurations.
* `filter_by_feature` - (List) Filters the SaaS application tenant by feature [Available values](https://help.zscaler.com/zia/saas-security-api#/casbTenant/lite-get)

### Optional

* `saas_application` - (String) SaaS tenant application i.e `BITBUCKET`
* `re_auth` - (Boolean) Enable tenant re-authentication
* `last_modified_time` - (Number) The date and time the tenant was last modified.
* `last_tenant_validation_time` - (Number) Last time the tenant was validated
* `features_supported` - (List) List of supported features
* `status` - (List) Status of the Saas Tenant
* `tenant_deleted` - (Boolean) If the tenant was deleted
* `tenant_webhook_enabled` - (Boolean) If the tenant webhook feature is enabled
* `zscaler_app_tenant_id` - (Block Set) Zscaler Application tenant ID and Name
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
