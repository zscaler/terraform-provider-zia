---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_resource"
description: |-
  Official documentation https://help.zscaler.com/zia/about-endpoint-dlp
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/dlpEndpointResource-post
  Creates and manages a ZIA DLP endpoint resource
---

# zia_endpoint_dlp_resource (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-endpoint-dlp)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/dlpEndpointResource-post)

The **zia_endpoint_dlp_resource** resource allows the creation and management of ZIA DLP endpoint resources in the Zscaler Internet Access cloud or via the API. DLP endpoint resources are channel-specific and describe the printers, removable storage devices, network shares, or personal cloud storage applications evaluated by Endpoint DLP.

~> **NOTE:** The `channel` attribute is required and determines which channel-specific block applies: `printer` (PRINTING), `removable_storage` (REMOVABLE_DRIVE_TRANSFER), `network_drives` / `network_drive_type` / `server_name` (NETWORK_DRIVE_TRANSFER), and `app_id` (PERSONAL_CLOUD_STORAGE).

~> **NOTE:** Custom endpoint applications are managed with the dedicated [`zia_endpoint_dlp_custom_apps`](./zia_endpoint_dlp_custom_apps.md) resource, which uses a separate endpoint.

## Example Usage - Printing Channel

```hcl
resource "zia_endpoint_dlp_resource" "printer01" {
  channel     = "PRINTING"
  name        = "Printer01"
  description = "Finance department printer"

  printer {
    unc        = "\\\\printserver\\finance"
    ip_address = "10.10.10.20"
    domain     = "acme.local"
  }
}
```

## Example Usage - Removable Storage Channel

```hcl
resource "zia_endpoint_dlp_resource" "usb01" {
  channel     = "REMOVABLE_DRIVE_TRANSFER"
  name        = "ApprovedUSB01"
  description = "Approved encrypted USB drive"

  removable_storage {
    vendor_id     = "0x1234"
    product_id    = "0x5678"
    serial_number = "SN-0001"
  }
}
```

## Example Usage - Network Drive Channel

```hcl
resource "zia_endpoint_dlp_resource" "share01" {
  channel            = "NETWORK_DRIVE_TRANSFER"
  name               = "NetworkShare01"
  description        = "Finance network share"
  network_drive_type = "MAPPED"
  server_name        = "fileserver01"

  network_drives {
    network_path = "\\\\fileserver01\\finance"
  }
}
```

## Example Usage - Personal Cloud Storage Channel

```hcl
resource "zia_endpoint_dlp_resource" "cloud01" {
  channel = "PERSONAL_CLOUD_STORAGE"
  name    = "PersonalDropbox"
  app_id  = 1
}
```

## Argument Reference

The following arguments are supported:

### Required

* `channel` - (String) The channel of the DLP endpoint resource. Supported values: `PRINTING`, `REMOVABLE_DRIVE_TRANSFER`, `NETWORK_DRIVE_TRANSFER`, `PERSONAL_CLOUD_STORAGE`.
* `name` - (String) The name of the DLP endpoint resource.

### Optional

* `description` - (String) The description of the DLP endpoint resource.
* `network_drive_type` - (String) The network drive type of the DLP endpoint resource. Applies to the `NETWORK_DRIVE_TRANSFER` channel.
* `server_name` - (String) The server name of the DLP endpoint resource. Applies to the `NETWORK_DRIVE_TRANSFER` channel.
* `app_id` - (Number) The application identifier of the DLP endpoint resource. Applies to the `PERSONAL_CLOUD_STORAGE` channel.

* `network_drives` - (Block List) The list of network drives associated with the DLP endpoint resource. Applies to the `NETWORK_DRIVE_TRANSFER` channel.
  * `network_path` - (String) The network path of the network drive.

* `printer` - (Block List, Max 1) The printer details of the DLP endpoint resource. Applies to the `PRINTING` channel.
  * `unc` - (String) The Universal Naming Convention (UNC) path of the printer.
  * `ip_address` - (String) The IP address of the printer.
  * `domain` - (String) The domain of the printer.

* `removable_storage` - (Block List, Max 1) The removable storage details of the DLP endpoint resource. Applies to the `REMOVABLE_DRIVE_TRANSFER` channel.
  * `vendor_id` - (String) The vendor identifier of the removable storage device.
  * `product_id` - (String) The product identifier of the removable storage device.
  * `serial_number` - (String) The serial number of the removable storage device.


## Import

Because the read path is channel-scoped, the import ID must include the channel.

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_endpoint_dlp_resource** can be imported using either:

* `<CHANNEL>:<id>` — the channel and numeric ID of the resource, or
* `<CHANNEL>:<name>` — the channel and exact resource name.

For example:

```shell
terraform import zia_endpoint_dlp_resource.example "PRINTING:12345"
```

or

```shell
terraform import zia_endpoint_dlp_resource.example "PRINTING:Printer01"
```
