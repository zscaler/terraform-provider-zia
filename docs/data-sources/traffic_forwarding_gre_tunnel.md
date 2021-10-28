---
subcategory: "Traffic Forwarding GRE Tunnel"
layout: "zia"
page_title: "ZIA: traffic_forwarding_gre_tunnel"
description: |-
  Gets the provisioned GRE Tunnel information
  
---

# zia_traffic_forwarding_gre_tunnel (Data Source)

The **zia_traffic_forwarding_gre_tunnel** - data source provides details about a specific provisioned GRE tunnel information created in the Zscaler Internet Access portal.

## Example Usage

```hcl
# ZIA Traffic Forwarding - GRE tunnel
data "zia_traffic_forwarding_gre_tunnel" "example"{
    source_ip = "100.1.1.1"
}

output "zia_traffic_forwarding_gre_tunnel_example" {
  value = data.zia_traffic_forwarding_gre_tunnel.example
}
```

## Argument Reference

The following arguments are supported:

* `source_ip` - (Required) - The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN. This IP address must be provisioned within the Zscaler service using the /staticIP endpoint.

### Read-Only

* `internal_ip_range` - (String) The start of the internal IP address in /29 CIDR range
* `comment` - (String) Additional information about this GRE tunnel
* `ip_unnumbered` - (Boolean) This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null
* `last_modification_time` - (Number) When the GRE tunnel information was last modified
* `within_country` - (Boolean) Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address

`primary_dest_vip` - (Set of Object) The primary destination data center and virtual IP address (VIP) of the GRE tunnel

* `city` -(String)
* `country_code` -(String)
* `datacenter` -(String)
* `id` -(Number)
* `latitude` -(Number)
* `longitude` -(Number)
* `private_service_edge` -(Boolean)
* `region` -(String)
* `virtual_ip` -(String)

`secondary_dest_vip` -(Set of Object) The secondary destination data center and virtual IP address (VIP) of the GRE tunnel

* `city` -(String)
* `country_code` -(String)
* `datacenter` -(String)
* `id` -(Number)
* `latitude` -(Number)
* `longitude` -(Number)
* `private_service_edge` -(Boolean)
* `region` -(String)
* `virtual_ip` -(String)

`managed_by` - (Set of Object) SD-WAN Partner that manages the location. If a partner does not manage the locaton, this is set to Self.

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` - (String) The configured name of the entity
* `extensions` - (Map of String)

`last_modified_by` - (Set of Object) Who modified the GRE tunnel information last

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` - (String) The configured name of the entity
* `extensions` - (Map of String)
