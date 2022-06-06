---
subcategory: "User Authentication Settings"
layout: "zscaler"
page_title: "ZIA: auth_settings_urls"
description: |-
  Adds a URL to or removes a URL from the cookie authentication exempt list
---

# Resource: zia_auth_settings_urls

The **zia_auth_settings_urls** resource alows you to add or remove a URL from the cookie authentication exempt list in the Zscaler Internet Access cloud or via the API. To learn more see [URL Format Guidelines](https://help.zscaler.com/zia/url-format-guidelines)

## Example Usage

```hcl
# ZIA User Auth Settings Data Source
resource "zia_auth_settings_urls" "example" {
  urls = [
    ".okta.com",
    ".oktacdn.com",
    ".mtls.oktapreview.com",
    ".mtls.okta.com",
    "d3l44rcogcb7iv.cloudfront.net",
    "pac.zdxcloud.net",
    ".windowsazure.com",
    ".fedoraproject.org",
    "login.windows.net",
    "d32a6ru7mhaq0c.cloudfront.net",
    ".kerberos.oktapreview.com",
    ".oktapreview.com",
    "login.zdxcloud.net",
    "login.microsoftonline.com",
    "smres.zdxcloud.net",
    ".kerberos.okta.com"
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `urls` - (Required) The email address of the admin user to be exported.

### Optional

There are no optional parameters supported by this resource.
