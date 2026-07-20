---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_custom_apps"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-dlp-and-endpoint-resources
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplications/customApp-post
  Creates and manages a ZIA DLP custom endpoint application
---

# zia_endpoint_dlp_custom_apps (Resource)

* [Official documentation](https://help.zscaler.com/zia/adding-dlp-and-endpoint-resources)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplications/customApp-post)

The **zia_endpoint_dlp_custom_apps** resource allows the creation and management of custom endpoint applications used by ZIA Endpoint DLP. Custom applications describe an application binary (by name, operating system, and file details) that Endpoint DLP can evaluate on a managed device.

## Example Usage - Windows OS

```hcl
resource "zia_endpoint_dlp_custom_apps" "app01" {
  name        = "CustomApp01"
  description = "Internal transfer utility"
  channel     = "APPLICATION_FILE_ACCESS"

  application {
    os_type            = "WINDOWS_OS"
    file_name          = "transfer.exe"
    original_file_name = "transfer.exe"
    digitally_signed   = true
  }
}
```

## Example Usage - Windows OS

```hcl
resource "zia_endpoint_dlp_custom_apps" "app02" {
  channel            = "APPLICATION_FILE_ACCESS"
  name               = "app02"
  description        = "app02"

  application {
    os_type = "MAC_OS"
    file_name = "app02"
    bundle_id = "12334567890"
    digitally_signed = true
  }
}

```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The name of the custom endpoint application.
* `application` - (Block List, Max 1) The application details of the custom endpoint application.
  * `os_type` - (String) The operating system type of the custom endpoint application. Supported values: `WINDOWS_OS`, `MAC_OS`.
  * `file_name` - (String) The file name of the custom endpoint application.
  * `original_file_name` - (String) The original file name of the custom endpoint application.
  * `bundle_id` - (String) The bundle identifier of the custom endpoint application (typically used on macOS).
  * `digitally_signed` - (Boolean) Indicates whether the custom endpoint application is digitally signed.

### Optional

* `channel` - (String) The channel of the custom endpoint application. Custom applications only support `APPLICATION_FILE_ACCESS` (the default).
* `description` - (String) The description of the custom endpoint application.
* `app_id` - (Number) The application identifier of the custom endpoint application. This value is assigned by the service; leave it unset to let the service allocate one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The Terraform resource identifier.
* `custom_app_id` - (Number) The unique identifier assigned to the custom endpoint application by the service.
* `application_type` - (String) The application type of the endpoint application.
* `zapp_id` - (String) The Zscaler Client Connector application identifier.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_endpoint_dlp_custom_apps** can be imported by using `<ID>` or `<NAME>` as the import ID.

For example:

```shell
terraform import zia_endpoint_dlp_custom_apps.example "8914"
```

or

```shell
terraform import zia_endpoint_dlp_custom_apps.example "CustomApp01"
```
