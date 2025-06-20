---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_gre_tunnel"
description: |-
  Official documentation https://help.zscaler.com/zia/about-gre-tunnels
  API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels-post
  Gets provisioned GRE tunnel information.
---

# zia_traffic_forwarding_gre_tunnel (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-gre-tunnels)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/greTunnels-post)

The **zia_traffic_forwarding_gre_tunnel** data source to get information about provisioned GRE tunnel information created in the Zscaler Internet Access portal.

## Example Usage - Retrieve GRE Tunnel by Source IP

```hcl
data "zia_traffic_forwarding_gre_tunnel" "example" {
  source_ip = "1.1.1.1"
}
```

## Example Usage - Retrieve GRE Tunnel by ID

```hcl
data "zia_traffic_forwarding_gre_tunnel" "example" {
  id = 1419395
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) - Unique identifier of the static IP address that is associated to a GRE tunnel
* `source_ip` - (Optional) - The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN.

-> **NOTE** `source_ip` is the public IP address (Static IP) associated with the GRE Tunnel

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `within_country` (Boolean) Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address
* `comment` (String) Additional information about this GRE tunnel
* `country_code` (String) When within_country is enabled, you must set this to the country code.
* `ip_unnumbered` (Boolean) This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null
* `primary_dest_vip`**` (List) The primary destination data center and virtual IP address (VIP) of the GRE tunnel.
  * `id` - (Number) Unique identifer of the GRE virtual IP address (VIP)
  * `virtual_ip` (String) GRE cluster virtual IP address (VIP)

* `secondary_dest_vip` (List) The secondary destination data center and virtual IP address (VIP) of the GRE tunnel.
  * `id` - (Number) Unique identifer of the GRE virtual IP address (VIP)
  * `virtual_ip` (String) GRE cluster virtual IP address (VIP)

* `internal_ip_range` (String) The start of the internal IP address in /29 CIDR range. Automatically set by the provider if `ip_unnumbered` is set to `false`.
