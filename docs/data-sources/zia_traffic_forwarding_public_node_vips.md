---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_public_node_vips"
description: |-
    Official documentation https://help.zscaler.com/zia/about-gre-tunnels
    API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels-post
    Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud
---

# zia_traffic_forwarding_public_node_vips (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-gre-tunnels)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels-post)

Use the **zia_traffic_forwarding_public_node_vips** data source to retrieve a paginated list of virtual IP addresses (VIPs) available in the Zscaler cloud.

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

* `cloud_name` - (String) Cloud Name
* `region` - (String) Region
* `city` - (String) City
* `datacenter` - (String) Data-center Name
* `location` - (String) Location Coordinates
* `vpn_ips` - (List of String) VPN IPs
* `vpn_domain_name` - (String) VPN Server DN
* `gre_ips` - (List of String) GRE IPs
* `gre_domain_name` - (String) Proxy Host Name
* `pac_ips` - (List of String) Pac IPs
* `pac_domain_name` - (String) Pac Server DN
