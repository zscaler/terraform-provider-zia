---
subcategory: "Traffic Forwarding - GRE Tunnel"
layout: "zia"
page_title: "ZIA: zia_traffic_forwarding_gre_tunnel"
description: |-
        Adds a GRE tunnel configuration.
---

# zia_traffic_forwarding_gre_tunnel (Resource)

The **zia_traffic_forwarding_gre_tunnel** - resource provisions a GRE tunnel configuration in the Zscaler Internet Access (ZIA) portal.

:warning: The provider automatically query the Zscaler cloud for the primary and secondary destination datacenter and virtual IP address (VIP) of the GRE tunnel. The parameter can be overriden if needed by setting the parameters: `primary_dest_vip` and `secondary_dest_vip`.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE Tunnel
// Creates a numbered GRE Tunnel
resource "zia_traffic_forwarding_gre_tunnel" "telus_home_internet_01_gre01" {
  source_ip = zia_traffic_forwarding_static_ip.example.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  within_country = true
  country_code = "CA"
  ip_unnumbered = false
  depends_on = [ zia_traffic_forwarding_static_ip.example ]
}

# ZIA Traffic Forwarding - Static IP
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = -36.848461
    longitude = 174.763336
}
```

:warning: The provider will automatically query and set the Zscaler cloud for the next available `/29` internal IP range to be used in a numbered GRE tunnel.

```hcl
# ZIA Traffic Forwarding - GRE Tunnel
// Creates an unnumbered GRE Tunnel
resource "zia_traffic_forwarding_gre_tunnel" "telus_home_internet_01_gre01" {
  source_ip = zia_traffic_forwarding_static_ip.example.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  within_country = true
  country_code = "CA"
  ip_unnumbered = true
  depends_on = [ zia_traffic_forwarding_static_ip.example ]
}

# ZIA Traffic Forwarding - Static IP
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = -36.848461
    longitude = 174.763336
}
```

## Argument Reference

The following arguments are supported:

* `source_ip` (required) The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN. This IP address must be provisioned within the Zscaler service using the /staticIP endpoint.

`primary_dest_vip`**` (Block Set) The primary destination data center and virtual IP address (VIP) of the GRE tunnel.

* `id` - (Required) Unique identfier for the group
* `virtual_ip` (Optional) GRE cluster virtual IP address (VIP)

`secondary_dest_vip` (Block Set) The secondary destination data center and virtual IP address (VIP) of the GRE tunnel.

* `id` - (Required) Unique identfier for the group
* `virtual_ip` (Optional) GRE cluster virtual IP address (VIP)

* `internal_ip_range` (Optional) The start of the internal IP address in /29 CIDR range. Automatically set by the provider if `ip_unnumbered` is set to `false`.
* `within_country` (Optional) Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address
* `comment` (Optional) Additional information about this GRE tunnel
* `country_code` (Optional) When within_country is enabled, you must set this to the country code.
* `ip_unnumbered` (Optional) This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null
* `last_modification_time` (Optional) When the GRE tunnel information was last modified

`last_modified_by` (Optional) When the GRE tunnel information was last modified

* `id` (Number) Identifier that uniquely identifies an entity
* `extensions` (Map of String)
