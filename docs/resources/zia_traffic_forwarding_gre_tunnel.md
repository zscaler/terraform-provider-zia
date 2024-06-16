---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: zia_traffic_forwarding_gre_tunnel"
description: |-
  Creates and manages GRE tunnel configuration.
---

# Resource: zia_traffic_forwarding_gre_tunnel

The **zia_traffic_forwarding_gre_tunnel** resource allows the creation and management of GRE tunnel configuration in the Zscaler Internet Access (ZIA) portal.

-> **Note:** The provider automatically query the Zscaler cloud for the primary and secondary destination datacenter and virtual IP address (VIP) of the GRE tunnel. The attribute can be overriden if needed by setting the parameters: `primary_dest_vip` and `secondary_dest_vip`.

## Example Usage - Unnumbered

```hcl
# Creates a numbered GRE Tunnel
resource "zia_traffic_forwarding_gre_tunnel" "example" {
  source_ip         = zia_traffic_forwarding_static_ip.example.ip_address
  comment           = "Example"
  within_country    = true
  country_code      = "US"
  ip_unnumbered     = true
  depends_on        = [ zia_traffic_forwarding_static_ip.example ]
}

# ZIA Traffic Forwarding - Static IP
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address      =  "1.1.1.1"
    routable_ip     = true
    comment         = "Example"
    geo_override    = true
    latitude        = 37.418171
    longitude       = -121.953140
}
```

```hcl
data "zia_traffic_forwarding_gre_vip_recommended_list" "this"{
    source_ip = zia_traffic_forwarding_static_ip.this.ip_address
    required_count = 2
}

data "zia_gre_internal_ip_range_list" "this"{
    required_count = 10
}

resource "zia_traffic_forwarding_static_ip" "this"{
    ip_address =  "50.98.112.169"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = 49.0526
    longitude = -122.8291
}

resource "zia_traffic_forwarding_gre_tunnel" "this" {
  source_ip      = zia_traffic_forwarding_static_ip.this.ip_address
  comment        = "GRE Tunnel Created with Terraform"
  within_country = false
  country_code   = "US"
  ip_unnumbered  = true
  primary_dest_vip {
    datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].datacenter
    virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].virtual_ip
  }
  secondary_dest_vip {
    datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].datacenter
    virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].virtual_ip
  }
  depends_on     = [zia_traffic_forwarding_static_ip.this]
}
```

-> **Note:** Although the example shows 2 valid attributes defined (datacenter, virtual_ip) within the primary_dest_vip and secondary_dest_vip, only one attribute is required. If setting the datacenter name as the attribute i.e YVR1. The provider will automatically select the available VIP.

-> **Note:** To obtain the datacenter codes and/or virtual_ips, refer to the following [Zscaler Portal](https://config.zscaler.com/zscloud.net/cenr) and choose your cloud tenant.

-> **Note:** The provider will automatically query and set the Zscaler cloud for the next available `/29` internal IP range to be used in a numbered GRE tunnel.

## Example Usage - Numbered

```hcl
data "zia_gre_internal_ip_range_list" "this"{
    required_count = 1
}

# ZIA Traffic Forwarding - Static IP
resource "zia_traffic_forwarding_static_ip" "this"{
    ip_address      =  "1.1.1.1"
    routable_ip     = true
    comment         = "Example"
    geo_override    = true
    latitude        = 37.418171
    longitude       = -121.953140
}

# Creates a Numbered GRE Tunnel
resource "zia_traffic_forwarding_gre_tunnel" "this" {
  source_ip       = zia_traffic_forwarding_static_ip.this.ip_address
  comment         = "Example"
  internal_ip_range = data.zia_gre_internal_ip_range_list.this.list[0].start_ip_address
  within_country  = true
  country_code    = "US"
  ip_unnumbered   = true
  depends_on      = [ zia_traffic_forwarding_static_ip.this ]

  lifecycle {
    ignore_changes = [
      internal_ip_range,
    ]
  }
}
```

-> **Note:** When configuring a numbered GRE Tunnel where the attribute `internal_ip_range` is defined, we must set the lifecycle block to ignore changes to the ``internal_ip_range`` attribute unless it is explicitly changed in the Terraform configuration.

## Argument Reference

The following arguments are supported:

### Required

* `source_ip` (required) The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN. This IP address must be provisioned within the Zscaler service using the /staticIP endpoint.

### Optional

* `within_country` (Optional) Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address
* `comment` (Optional) Additional information about this GRE tunnel
* `country_code` (Optional) When within_country is enabled, you must set this to the country code.
* `ip_unnumbered` (Optional) This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null
* `primary_dest_vip`**` (Optional) The primary destination data center and virtual IP address (VIP) of the GRE tunnel.
  * `id` - (Optional) Unique identifer of the GRE virtual IP address (VIP)
  * `virtual_ip` (Optional) GRE cluster virtual IP address (VIP)

* `secondary_dest_vip` (Optional) The secondary destination data center and virtual IP address (VIP) of the GRE tunnel.
  * `id` - (Optional) Unique identifer of the GRE virtual IP address (VIP)
  * `virtual_ip` (Optional) GRE cluster virtual IP address (VIP)

* `internal_ip_range` (Optional) The start of the internal IP address in /29 CIDR range. Automatically set by the provider if `ip_unnumbered` is set to `false`.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_traffic_forwarding_gre_tunnel** can be imported by using `<TUNNEL_ID>` as the import ID.

For example:

```shell
terraform import zia_traffic_forwarding_gre_tunnel.example <tunnel_id>
```

or

```shell
terraform import zia_traffic_forwarding_gre_tunnel.example <engine_name>
```
