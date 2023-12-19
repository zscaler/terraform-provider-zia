---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZIA: cloud_browser_isolation_profile"
description: |-
  Get information about an Cloud Browser Isolation Profile in Zscaler Internet Access cloud.
---

# Data Source: zia_cloud_browser_isolation_profile

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
