---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA: tenant_restriction_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/about-tenant-profiles
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-get
  Retrieves the domain profile summary.
---

# zia_tenant_restriction_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-tenant-profiles)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-get)

Use the **zia_tenant_restriction_profile** data source to get information about a ZIA Domain Profiles in the Zscaler Internet Access cloud or via the API. The resource can then be utilized when configuring a Web DLP Rule resource `zia_dlp_web_rules`

## Example Usage - By Name

```hcl
data "zia_tenant_restriction_profile" "this"{
    name = "MiicrosoftO365Login"
}
```

## Example Usage - By ID

```hcl
data "zia_tenant_restriction_profile" "this"{
    id = "5421656"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Optional) The unique identifier for the tenant restriction profile
* `id` - (String) The tenant restriction profile name

### Optional

* `description` - (String) Additional information about the profile
* `app_type` - (String) Restricted tenant profile application type
* `item_type_primary` - (String) Tenant profile primary item type
* `item_type_secondary` - (String) Tenant profile secondary item type
* `item_data_primary` - (List) Tenant profile primary item data
* `item_data_secondary` - (List) Tenant profile secondary item data
* `item_value` - (List) Tenant profile item value for YouTube category
* `restrict_personal_o365_domains` - (Boolean) Flag to restrict personal domains for Office 365
* `allow_google_consumers` - (Boolean) Flag to allow Google consumers
* `ms_login_services_tr_v2` - (Boolean) Flag to decide between v1 and v2 for tenant restriction on MSLOGINSERVICES
* `allow_google_visitors` - (Boolean) Flag to allow Google visitors
* `allow_gcp_cloud_storage_read` - (Boolean) Flag to allow or disallow cloud storage resources for GCP
