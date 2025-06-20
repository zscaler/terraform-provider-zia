---
subcategory: "User Authentication Settings"
layout: "zscaler"
page_title: "ZIA: auth_settings_urls"
description: |-
  Official documentation https://help.zscaler.com/zia/url-format-guidelines
  API documentation https://help.zscaler.com/zia/user-authentication-settings#/authSettings/exemptedUrls-get
  Gets a list of URLs that were exempted from cookie authentication and SSL Inspection.
---

# zia_auth_settings_urls (Data Source)

* [Official documentation](https://help.zscaler.com/zia/url-format-guidelines)
* [API documentation](https://help.zscaler.com/zia/user-authentication-settings#/authSettings/exemptedUrls-get)

Use the **zia_auth_settings_urls** data source to get a list of URLs that were exempted from cookie authentiation and SSL Inspection in the Zscaler Internet Access cloud or via the API. To learn more see [URL Format Guidelines](https://help.zscaler.com/zia/url-format-guidelines)

## Example Usage

```hcl
# ZIA User Auth Settings Data Source
data "zia_auth_settings_urls" "foo" {}
```

## Argument Reference

This data source can be executed without the need of additional parameters.

## Attribute Reference

N/A
