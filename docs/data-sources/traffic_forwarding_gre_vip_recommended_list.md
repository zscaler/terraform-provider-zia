---
subcategory: "Traffic Forwarding GRE VIP Recommended List"
layout: "zia"
page_title: "ZIA: gre_vip_recommended_list"
description: |-
        Gets a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.
  
---

# zia_gre_vip_recommended_list (Data Source)

The **zia_gre_vip_recommended_list** - data source retrieves details about a a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE VIP Recommended List
data "zia_gre_vip_recommended_list" "example"{
    source_ip = "1.1.1.1"
}

output "zia_gre_virtual_ip_address_list_example"{
    value = data.zia_gre_virtual_ip_address_list.example
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Number) Unique identifer of the GRE virtual IP address (VIP)
* `source_ip` - (String) The public source IP address.
* `virtual_ip` - (String) GRE cluster virtual IP address (VIP)
* `private_service_edge` - (Boolean) Set to true if the virtual IP address (VIP) is a ZIA Private Service Edge
* `datacenter` - (String) Data center information
