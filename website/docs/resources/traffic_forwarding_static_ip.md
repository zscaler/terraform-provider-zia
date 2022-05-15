---
subcategory: "Traffic Forwarding - Static IP"
layout: "zia"
page_title: "ZIA: traffic_forwarding_static_ip"
description: |-
        Adds a Static IP address.
---

# zia_traffic_forwarding_static_ip (Resource)

The **zia_traffic_forwarding_static_ip** - resource provisions a static IP address in the Zscaler Internet Access (ZIA) portal.

## Example Usage

```hcl
# ZIA Traffic Forwarding - Virtual IP Addresses (VIPs)
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = -36.848461
    longitude = 174.763336
}

output "zia_traffic_forwarding_static_ip"{
    value = zia_traffic_forwarding_static_ip.example
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The static IP address
* `latitude` - (Required) Required only if the geoOverride attribute is set. Latitude with 7 digit precision after decimal point, ranges between -90 and 90 degrees.
* `longitude` - (Required) Required only if the geoOverride attribute is set. Longitude with 7 digit precision after decimal point, ranges between -180 and 180 degrees.
* `geo_override` - (Optional) If not set, geographic coordinates and city are automatically determined from the IP address. Otherwise, the latitude and longitude coordinates must be provided.
* `routable_ip` - (Optional) Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.
* `last_modification_time` - (Optional) When the static IP address was last modified
* `comment` - (Optional) Additional information about this static IP address

`managed_by` (Set of Object)

* `id` - (Optional)

`last_modified_by` (Set of Object)

* `id` - (Optional)
