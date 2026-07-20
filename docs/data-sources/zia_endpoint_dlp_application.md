---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_application"
description: |-
  Official documentation https://help.zscaler.com/zia/about-endpoint-dlp
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplications-get
  Retrieves a ZIA Endpoint DLP application by name or ID
---

# zia_endpoint_dlp_application (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-endpoint-dlp)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplications-get)

Use the **zia_endpoint_dlp_application** data source to retrieve a single endpoint application by `application_name` or `id`. This is typically used to resolve an application's `zapp_id` so it can be referenced by the [`zia_endpoint_dlp_application_group`](../resources/zia_endpoint_dlp_application_group.md) resource.

~> **NOTE:** The endpoint application catalogue can be extensive. The `search`, `os_type`, and `application_type` arguments are pushed to the service as native filters to narrow the result set before the lookup, which is more efficient than filtering the full list on the client.

## Example Usage - By Name

```hcl
data "zia_endpoint_dlp_application" "notepad" {
  application_name = "Notepad++ (Windows)"
}
```

## Example Usage - By ID

```hcl
data "zia_endpoint_dlp_application" "notepad" {
  id = 500000085
}
```

## Example Usage - Narrowing With Filters

```hcl
data "zia_endpoint_dlp_application" "dropbox" {
  application_name = "Dropbox"
  os_type          = "MAC_OS"
  application_type = "WELLKNOWN"
}
```

## Argument Reference

The following arguments are supported. One of `id` or `application_name` must be set.

* `id` - (Number) The unique identifier (`resourceId`) of the endpoint application.
* `application_name` - (String) The name of the endpoint application.
* `search` - (String) Search string applied server-side against application names to narrow the result set before matching. When omitted, `application_name` is used for server-side narrowing.
* `os_type` - (String) Filters the results server-side by operating system (e.g. `WINDOWS_OS`, `MAC_OS`).
* `application_type` - (String) Filters the results server-side by application type (e.g. `WELLKNOWN`, `CUSTOM`, `DISCOVERED`).

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `zapp_id` - (String) The Zscaler Client Connector application identifier. This is the value referenced by `zia_endpoint_dlp_application_group.resources`.
* `description` - (String) The description of the endpoint application.
* `file_name` - (String) The file name of the endpoint application.
* `original_file_name` - (String) The original file name of the endpoint application.
* `bundle_id` - (String) The bundle identifier of the endpoint application.
* `digitally_signed` - (Boolean) Indicates whether the endpoint application is digitally signed.
* `mod_uid` - (Number) The modification unique identifier of the endpoint application.
* `last_modified_time` - (Number) Timestamp when the endpoint application was last modified.
* `deleted` - (Boolean) Indicates whether the endpoint application is deleted.
* `version` - (List) Version metadata for the endpoint application.
  * `version`, `z_ver_id_md32`, `threat_type`, `threat_level`, `bundle_id`, `code_signing_certificate_status`, `threat_level_updated`
* `versions` - (List) Additional version metadata for the endpoint application. Same fields as `version`.
