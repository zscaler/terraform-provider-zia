---
subcategory: "Traffic Forwarding - Static IP"
layout: "zia"
page_title: "ZIA: traffic_forwarding_static_ip"
description: |-
        Gets all provisioned static IP addresses.
---

# zia_traffic_forwarding_static_ip (Data Source)

The **zia_traffic_forwarding_static_ip** - data source retrieves all provisioned static IP addresses.

## Example Usage

```hcl
# ZIA Traffic Forwarding - Virtual IP Addresses (VIPs)
data "zia_traffic_forwarding_public_node_vips" "yvr1"{
    datacenter = "YVR1"
}

output "zia_traffic_forwarding_public_node_vips_yvr1"{
    value = data.zia_traffic_forwarding_public_node_vips.yvr1
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Number) The unique identifier for the static IP address
* `ip_address` - (String) The static IP address
* `geo_override` - (Boolean) If not set, geographic coordinates and city are automatically determined from the IP address. Otherwise, the latitude and longitude coordinates must be provided.
* `latitude` - (Number) Required only if the geoOverride attribute is set. Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.
* `longitude` - (Number) Required only if the geoOverride attribute is set. Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.
* `routable_ip` - (Boolean) Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.
* `last_modification_time` - (Number) When the static IP address was last modified
* `comment` - (String) Additional information about this static IP address

`managed_by` (Set of Object)

* `id` - (Number)
* `name` (String)
* `extensions` (String)

`last_modified_by` (Set of Object)

* `id` - (Number)
* `name` (String)
* `extensions` (String)
