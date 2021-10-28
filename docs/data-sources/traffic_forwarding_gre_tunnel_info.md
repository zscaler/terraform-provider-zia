---
subcategory: "Traffic Forwarding GRE Tunnel Info"
layout: "zia"
page_title: "ZIA: traffic_forwarding_gre_tunnel"
description: |-
  Gets the provisioned GRE Tunnel information
---

# zia_traffic_forwarding_gre_tunnel_info (Data Source)

The **zia_traffic_forwarding_gre_tunnel_info** - data source provides details about a specific provisioned GRE tunnel information created in the Zscaler Internet Access portal.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE tunnel
data "zia_traffic_forwarding_gre_tunnel_info" "example" {
  ip_address = "1.1.1.1"
}

output "zia_traffic_forwarding_gre_tunnel_info_example2" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.example2
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) - Filter based on an IP address range.
* `gre_enabled` - (Optional) - Displays only ip addresses with GRE tunnel enabled

### Read-Only

* `gre_tunnel_ip` - (String) The start of the internal IP address in /29 CIDR range
* `primary_gw` - (String)
* `secondary_gw` - (String)
* `tun_id` - (Number)
* `gre_range_primary` - (String)
* `gre_range_secondary` - (String)
