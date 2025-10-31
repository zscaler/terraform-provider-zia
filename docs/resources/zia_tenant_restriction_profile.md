---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA: tenant_restriction_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-tenant-profiles
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-post
  Retrieves the domain profile summary.
---

# zia_tenant_restriction_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-tenant-profiles)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-post)

Use the **zia_tenant_restriction_profile** resource creates and manages tenant retriction profiles in the Zscaler Internet Access cloud.

## Example Usage - Create O365 Tenant Restriction Profile

```hcl
resource "zia_tenant_restriction_profile" "this" {
  name = "ACME_MSFT_CA"
  description = "ACME_MSFT_CA"
  restrict_personal_o365_domains = true
  app_type = "MSLOGINSERVICES"
  item_data_primary = ["11111111-1111-1111-1111-111111111111"]
  item_data_secondary = ["acme.com"]
  item_type_primary = "TENANT_RESTRICTION_TENANT_DIRECTORY"
  item_type_secondary = "TENANT_RESTRICTION_TENANT_NAME"
}
```

## Example Usage - Create O365 V2 Tenant Restriction Profile

```hcl
resource "zia_tenant_restriction_profile" "this2" {
  name = "ACME_MSFT_CA_v2"
  description = "ACME_MSFT_CA_v2"
  ms_login_services_tr_v2 = true
  app_type = "MSLOGINSERVICES"
  item_data_primary = ["11111111-1111-1111-1111-111111111111:quadsj"]
  item_type_primary = "TENANT_RESTRICTION_TENANT_POLICY_ID"
}
```

## Example Usage - Create YouTube Tenant Restriction Profile

```hcl
resource "zia_tenant_restriction_profile" "this3" {
  name = "YouTube01_Profile"
  description = "YouTube01_Profile"
  app_type = "YOUTUBE"
  item_value = ["TENANT_RESTRICTION_ACTION_OR_ADVENTURE"]
  item_type_primary = "TENANT_RESTRICTION_CATEGORY_ID"
}
```

## Example Usage - Create Dropbox Tenant Restriction Profile

```hcl
resource "zia_tenant_restriction_profile" "this4" {
  name = "Dropbox_Profile"
  description = "Dropbox_Profile"
  app_type = "DROPBOX"
  item_data_primary = [139732608]
  item_type_primary = "TENANT_RESTRICTION_TEAM_ID"
}
```

## Example Usage - Create Google Tenant Restriction Profile

```hcl
resource "zia_tenant_restriction_profile" "this5" {
  name = "Google_Profile01"
  description = "Google_Profile01"
  allow_google_consumers = false
  allow_google_visitors = false
  app_type = "GOOGLE"
  item_data_primary = ["acme.com"]
  item_type_primary = "TENANT_RESTRICTION_DOMAIN"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Optional) The unique identifier for the tenant restriction profile

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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_tenant_restriction_profile** can be imported by using `<PROFILE ID>` or `<PROFILE NAME>` as the import ID.

For example:

```shell
terraform import zia_tenant_restriction_profile.example <profile_id>
```

or

```shell
terraform import zia_tenant_restriction_profile.example <profile_name>
```
