---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_endpoint_resource_channels"
description: |-
  Official documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/dlpEndpointResource/{channel}-get
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/dlpEndpointResource/{channel}-get
  Get information about DLP resources configured for the specified channel
---

# zia_dlp_endpoint_resource_channels (Data Source)

* [Official documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/dlpEndpointResource/{channel}-get)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/dlpEndpointResource/{channel}-get)

Use the **zia_dlp_endpoint_resource_channels** data source to get information about a DLP resource configured for the specified channel in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP endpoint resource by name within a channel
data "zia_dlp_endpoint_resource_channels" "example" {
  channel = "NETWORK_DRIVE_TRANSFER"
  name    = "NetworkShare01"
}
```

```hcl
# Retrieve a DLP endpoint resource by ID within a channel
data "zia_dlp_endpoint_resource_channels" "example" {
  channel = "PRINTING"
  id      = 4522856946
}
```

## Argument Reference

The following arguments are supported:

* `channel` - (Required, String) The DLP endpoint resource channel to query. Supported values: `PRINTING`, `REMOVABLE_DRIVE_TRANSFER`, `NETWORK_DRIVE_TRANSFER`, `PERSONAL_CLOUD_STORAGE`.
* `name` - (Optional, String) The name of the DLP endpoint resource. Used to look up the resource within the specified channel.
* `id` - (Optional, Number) The unique identifier of the DLP endpoint resource. Used to look up the resource within the specified channel.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `is_predefined` - (Boolean) Indicates whether the DLP endpoint resource is predefined.
* `network_drive_type` - (String) The network drive type of the DLP endpoint resource.
* `description` - (String) The description of the DLP endpoint resource.
* `server_name` - (String) The server name of the DLP endpoint resource.
* `app_id` - (Number) The application identifier of the DLP endpoint resource.

* `network_drives` - (List of Object) The list of network drives associated with the DLP endpoint resource.
  * `network_path` - (String) The network path of the network drive.

* `printer` - (List of Object) The printer details of the DLP endpoint resource.
  * `unc` - (String) The Universal Naming Convention (UNC) path of the printer.
  * `ip_address` - (String) The IP address of the printer.
  * `domain` - (String) The domain of the printer.

* `removable_storage` - (List of Object) The removable storage details of the DLP endpoint resource.
  * `vendor_id` - (String) The vendor identifier of the removable storage device.
  * `product_id` - (String) The product identifier of the removable storage device.
  * `serial_number` - (String) The serial number of the removable storage device.

* `application` - (List of Object) The application details of the DLP endpoint resource.
  * `os_type` - (String) The operating system type of the application.
  * `file_name` - (String) The file name of the application.
  * `original_file_name` - (String) The original file name of the application.
  * `bundle_id` - (String) The bundle identifier of the application.
  * `digitally_signed` - (Boolean) Indicates whether the application is digitally signed.
