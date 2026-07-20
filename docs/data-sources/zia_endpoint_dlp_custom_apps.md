---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_custom_apps"
description: |-
  Official documentation https://help.zscaler.com/zia/about-endpoint-dlp
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplications-get
  Gets information about a ZIA DLP custom endpoint application
---

# zia_endpoint_dlp_custom_apps (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-endpoint-dlp)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplications-get)

Use the **zia_endpoint_dlp_custom_apps** data source to get information about a custom endpoint application used by ZIA Endpoint DLP. The application can be looked up by `id` or by `name`.

## Example Usage - By Name

```hcl
data "zia_endpoint_dlp_custom_apps" "this" {
  name = "CustomApp01"
}
```

## Example Usage - By ID

```hcl
data "zia_endpoint_dlp_custom_apps" "this" {
  id = 8914
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Number) The unique identifier of the custom endpoint application.
* `name` - (String) The name of the custom endpoint application.
* `search` - (String) Optional search string used to narrow the results when looking up by name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) The description of the custom endpoint application.
* `os_type` - (String) The operating system type of the custom endpoint application.
* `file_name` - (String) The file name of the custom endpoint application.
* `original_file_name` - (String) The original file name of the custom endpoint application.
* `bundle_id` - (String) The bundle identifier of the custom endpoint application.
* `digitally_signed` - (Boolean) Indicates whether the custom endpoint application is digitally signed.
* `application_type` - (String) The application type of the endpoint application.
* `zapp_id` - (String) The Zscaler Client Connector application identifier.
* `mod_uid` - (Number) The modification unique identifier of the custom endpoint application.
* `last_modified_time` - (Number) Timestamp when the custom endpoint application was last modified.
* `deleted` - (Boolean) Indicates whether the custom endpoint application is deleted.
* `version` - (List) Version details of the custom endpoint application.
  * `version` - (String) The version string.
  * `z_ver_id_md32` - (Number) The internal version identifier.
  * `threat_type` - (Number) The threat type classification.
  * `threat_level` - (String) The threat level classification.
  * `bundle_id` - (String) The bundle identifier for this version.
  * `code_signing_certificate_status` - (Number) The code signing certificate status.
  * `threat_level_updated` - (Boolean) Indicates whether the threat level was updated.
* `versions` - (List) The list of known versions of the custom endpoint application. Each element exposes the same attributes as `version`.
