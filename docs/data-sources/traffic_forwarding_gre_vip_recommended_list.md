---
subcategory: "Traffic Forwarding"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): gre_vip_recommended_list"
sidebar_current: "docs-datasource-zia-gre-vip-recommended-list"
description: |-
    Gets a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.

---

# Data Source: zia_gre_vip_recommended_list

Use the **zia_gre_vip_recommended_list** data source to get information about a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE VIP Recommended List
data "zia_gre_vip_recommended_list" "example"{
    source_ip = "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) - Filter based on an IP address range.
* `id` - (Number) Unique identifer of the GRE virtual IP address (VIP)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `source_ip` - (String) The public source IP address.
* `virtual_ip` - (String) GRE cluster virtual IP address (VIP)
* `private_service_edge` - (Boolean) Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge
* `datacenter` - (String) Data center information
