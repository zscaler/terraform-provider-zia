---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_static_ip"
description: |-
    Official documentation https://help.zscaler.com/zia/about-static-ip
    API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/staticIP-get
    Gets static IP address for the specified ID
---

# zia_traffic_forwarding_static_ip (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-static-ip)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/staticIP-get)

Use the **zia_traffic_forwarding_static_ip** data source to get information about all provisioned static IP addresses. This resource can then be utilized when creating a GRE Tunnel or VPN Credential resource of Type `IP`

## Example Usage

```hcl
# ZIA Traffic Forwarding - Static IPs
data "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The static IP address
* `id` - (Optional) The unique identifier for the static IP address

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Number) The unique identifier for the static IP address
* `ip_address` - (String) The static IP address
* `geo_override` - (Boolean) If not set, geographic coordinates and city are automatically determined from the IP address. Otherwise, the latitude and longitude coordinates must be provided.
* `latitude` - (Number) Required only if the geoOverride attribute is set. Latitude with 7 digit precision after decimal point, ranges between `-90` and `90` degrees.
* `longitude` - (Number) Required only if the geoOverride attribute is set. Longitude with 7 digit precision after decimal point, ranges between `-180` and `180` degrees.
* `routable_ip` - (Boolean) Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.
* `last_modification_time` - (Number) When the static IP address was last modified
* `comment` - (String) Additional information about this static IP address

* `managed_by` (Set of Object)
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)

* `last_modified_by` (Set of Object)
      - `id` - (String) Identifier that uniquely identifies an entity
      - `name` - (String) The configured name of the entity
      - `extensions` - (Map of String)
