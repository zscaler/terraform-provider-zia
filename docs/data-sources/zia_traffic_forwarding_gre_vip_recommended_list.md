---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_gre_vip_recommended_list"
description: |-
    Official documentation https://help.zscaler.com/zia/about-gre-tunnels
    API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels-post
    Gets a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.
---

# zia_traffic_forwarding_gre_vip_recommended_list (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-gre-tunnels)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels-post)

Use the **zia_traffic_forwarding_gre_vip_recommended_list** data source to get information about a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE VIP Recommended List
data "zia_traffic_forwarding_gre_vip_recommended_list" "this"{
    source_ip = "1.1.1.1"
    required_count = 2
}
```

## Example Usage - With Overridden Geo Coordinates

```hcl
# ZIA Traffic Forwarding - GRE VIP Recommended List
data "zia_traffic_forwarding_gre_vip_recommended_list" "this"{
    source_ip = "1.1.1.1"
    required_count = 2
    latitude     = 22.2914
    longitude    = 114.1445
}
```

## Argument Reference

The following arguments are supported:

* `source_ip` - (Required) - Filter based on an IP address range.
* `required_count` - (Optional)  - Number of IP address to be exported.
* `id` - (Number) Unique identifer of the GRE virtual IP address (VIP)

## Attribute Reference

In addition to all arguments above, the following optional attributes can be used to manipulate the recommended list filtering:

* `source_ip` - (String) The public source IP address.
* `routable_ip` - (Boolean) The routable IP address.
* `within_country_only` - (Boolean) Search within country only.
* `include_private_service_edge` - (Boolean) Include ZIA Private Service Edge VIPs.
* `include_current_vips` - (Boolean) Include currently assigned VIPs.
* `latitude` - (Number) The latitude coordinate of the GRE tunnel source.
* `longitude` - (Number) The longitude coordinate of the GRE tunnel source.
* `subcloud` - (String) The longitude coordinate of the GRE tunnel source.

In addition to all arguments above, the following optional attributes are exported:

* `list` - The list of all recommended returned Virtual IP Addresses (VIPs)
  * `id` - (Number) Unique identifer of the GRE virtual IP address (VIP)
  * `virtual_ip` - (String) GRE cluster virtual IP address (VIP)
  * `private_service_edge` - (Boolean) Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge
  * `datacenter` - (String) Data center information
  * `city` - (String) Data center city information
  * `country_code` - (String) Data center country code information in ISO 3166 Alpha-2
  * `region` - (String) Data center region information.
  * `latitude` - (Number) The latitude coordinate of the GRE tunnel source.
  * `longitude` - (Number) The longitude coordinate of the GRE tunnel source.
