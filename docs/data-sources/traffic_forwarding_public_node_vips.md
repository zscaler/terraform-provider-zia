---
subcategory: "Traffic Forwarding"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): traffic_forwarding_public_node_vips"
sidebar_current: "docs-datasource-zia-traffic_forwarding_public_node_vips"
description: |-
    Gets a paginated list of the virtual IP addresses (VIPs) available in the Zscaler cloud
---

# Data Source: zia_traffic_forwarding_public_node_vips

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
