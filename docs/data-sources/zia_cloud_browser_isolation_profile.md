---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZIA: cloud_browser_isolation_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-url-filtering-policy#Action
  API documentation https://help.zscaler.com/zia/browser-isolation#/browserIsolation/profiles-get
  Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, File Type Control and SSL Inspection rules.
  Get information about an Cloud Browser Isolation Profile in Zscaler Internet Access cloud.
---

# zia_cloud_browser_isolation_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-url-filtering-policy#Action)
* [API documentation](https://help.zscaler.com/zia/browser-isolation#/browserIsolation/profiles-get)

Use the **zia_cloud_browser_isolation_profile** data source to get information about an isolation profile in the Zscaler Internet Access cloud. This data source is required when configuring URL filtering rule where the action is set to `ISOLATE`

## Example Usage

```hcl
data "zia_cloud_browser_isolation_profile" "this" {
    name = "ZS_CBI_Profile1"
}
```

## Argument Reference

* `name` - (Required) This field defines the name of the isolation profile.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (string) The universally unique identifier (UUID) for the browser isolation profile.
* `url` - (string) The browser isolation profile URL
* `default_profile` - (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field
