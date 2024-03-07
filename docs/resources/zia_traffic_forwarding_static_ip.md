---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_static_ip"
description: |-
    Creates and manages static IP addresses.
---

# Resource: zia_traffic_forwarding_static_ip

The **zia_traffic_forwarding_static_ip** resource allows the creation and management of static ip addresses in the Zscaler Internet Access cloud. The resource, can then be associated with other resources such as:

* VPN Credentials of type `IP`
* Location Management
* GRE Tunnel

## Example Usage

```hcl
# ZIA Traffic Forwarding - Static IP
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address      =  "1.1.1.1"
    routable_ip     = true
    comment         = "Example"
    geo_override    = true
    latitude        = -36.848461
    longitude       = 174.763336
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The static IP address

### Optional

* `comment` - (Optional) Additional information about this static IP address
* `latitude` - (Optional) Required only if the geoOverride attribute is set. Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.
* `longitude` - (Optional) Required only if the geoOverride attribute is set. Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.
* `geo_override` - (Optional) If not set, geographic coordinates and city are automatically determined from the IP address. Otherwise, the latitude and longitude coordinates must be provided.
* `routable_ip` - (Optional) Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.

* `managed_by` (Set of Object)
  * `id` - (Optional)

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Static IP resources can be imported by using `<STATIC IP ID>` or `<IP ADDRESS>`as the import ID.

```shell
terraform import zia_traffic_forwarding_static_ip.example <static_ip_id>
```

or

```shell
terraform import zpa_app_connector_group.example <ip_address>
```
