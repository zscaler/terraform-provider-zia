---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_gre_tunnel_info"
description: |-
  Official documentation https://help.zscaler.com/zia/traffic-forwarding-0#/orgProvisioning/ipGreTunnelInfo-get
  API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/orgProvisioning/ipGreTunnelInfo-get
  Gets a list of IP addresses with GRE tunnel details.
---

# zia_traffic_forwarding_gre_tunnel_info (Data Source)

* [Official documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/orgProvisioning/ipGreTunnelInfo-get)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/orgProvisioning/ipGreTunnelInfo-get)

The **zia_traffic_forwarding_gre_tunnel_info** data source to get information about provisioned GRE tunnel information created in the Zscaler Internet Access portal.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE tunnel
data "zia_traffic_forwarding_gre_tunnel_info" "example" {
  ip_address = "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) - Filter based on an IP address range.
* `gre_enabled` - (Optional) - Displays only ip addresses with GRE tunnel enabled

-> **NOTE** `ip_address` is the public IP address (Static IP) associated with the GRE Tunnel

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `gre_tunnel_ip` - (String) The start of the internal IP address in /29 CIDR range
* `primary_gw` - (String)
* `secondary_gw` - (String)
* `tun_id` - (Number)
* `gre_range_primary` - (String)
* `gre_range_secondary` - (String)
