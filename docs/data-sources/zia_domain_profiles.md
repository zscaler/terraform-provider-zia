---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA: domain_profiles"
description: |-
  Official documentation https://help.zscaler.com/zia/about-email-profiles
  API documentation https://help.zscaler.com/zia/saas-security-api#/domainProfiles/lite-get
  Retrieves the domain profile summary.
---

# zia_domain_profiles (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-email-profiles)
* [API documentation](https://help.zscaler.com/zia/saas-security-api#/domainProfiles/lite-get)

Use the **zia_domain_profiles** data source to get information about a ZIA Domain Profiles in the Zscaler Internet Access cloud or via the API. The resource can then be utilized when configuring a Web DLP Rule resource `zia_dlp_web_rules`

## Example Usage - By Name

```hcl
data "zia_domain_profiles" "this"{
    profile_name = "Example"
}
```

## Example Usage - By ID

```hcl
data "zia_domain_profiles" "this"{
    profile_id = "Example"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `profile_name` - (Optional) Domain profile name
* `profile_id` - (String) Domain profile ID

### Optional

* `description` - (String) Additional notes or information about the domain profile
* `include_company_domains` - (Boolean) Determine if the organizational domains have to be included in the domain profile
* `include_subdomains` - (String) determine whether or not to include subdomains
* `custom_domains` - (List) List of custom domains for the domain profile. There can be one or more custom domains.
* `predefined_email_domains` - (List) List of predefined email service provider domains for the domain profile
