---
subcategory: "Location Management"
layout: "zia"
page_title: "ZIA: location_management"
description: |-
    Adds new locations and sub-locations.
  
---
# zia_location_management (Resource)

The **zia_location_management** - resource provisions a location resource configuration in the Zscaler Internet Access (ZIA) portal. This resource can then be associated with a:

* Static IP resource
* GRE Tunnel resource
* VPN credentials resource
* URL filtering and firewall filtering rules

## Example Usage

```hcl
# ZIA Location Management
resource "zia_location_management" "usa_sjc37"{
    name = "USA_SJC_37"
    description = "Created with Terraform"
    country = "UNITED_STATES"
    tz = "UNITED_STATES_AMERICA_LOS_ANGELES"
    auth_required = true
    idle_time_in_minutes = 720
    display_time_unit = "HOUR"
    surrogate_ip = true
    xff_forward_enabled = true
    ofw_enabled = true
    ips_control = true
    ip_addresses = [ zia_traffic_forwarding_static_ip.usa_sjc37.ip_address ]
    depends_on = [ zia_traffic_forwarding_static_ip.usa_sjc37, zia_traffic_forwarding_vpn_credentials.usa_sjc37 ]
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.vpn_credental_id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
    }
}

resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type = "UFQDN"
    fqdn = "usa_sjc37@acme.com"
    comments = "USA - San Jose IPSec Tunnel"
    pre_shared_key = "P@ass0rd123!"
}

resource "zia_traffic_forwarding_static_ip" "usa_sjc37"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "SJC37 - Static IP"
    geo_override = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) - Location Name.
* `ip_addresses` - (Required) For locations: IP addresses of the egress points that are provisioned in the Zscaler Cloud. Each entry is a single IP address (e.g., `238.10.33.9`). For sub-locations: Egress, internal, or GRE tunnel IP addresses. Each entry is either a single IP address, CIDR (e.g., `10.10.33.0/24`), or range (e.g., `10.10.33.1-10.10.33.10`)). The value is required if `vpn_credentials` are not defined.

`vpn_credentials`

* `id` - (Number) VPN credential resource id. The value is required if `ip_addresses` are not defined.

* `country` - (String) Country
* `tz` - (String) Timezone of the location. If not specified, it defaults to GMT.
* `description` - (String) Additional notes or information regarding the location or sub-location. The description cannot exceed 1024 characters.
* `profile` - (Optional) Profile tag that specifies the location traffic type. If not specified, this tag defaults to `Unassigned`. The supported options are: `NONE`, `CORPORATE`, `SERVER`, `GUESTWIFI`, `IOT`.

* `aup_block_internet_until_accepted` - (Optional) For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.
* `aup_enabled` - (Optional) Enable AUP. When set to true, AUP is enabled for the location.
* `aup_force_ssl_inspection` - (Optional) For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler will force SSL Inspection in order to enforce AUP for HTTPS traffic.
* `aup_timeout_in_days` - (Optional) Custom AUP Frequency. Refresh time (in days) to re-validate the AUP.
* `auth_required` - (Optional) Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
* `caution_enabled` - (Optional) Enable Caution. When set to true, a caution notifcation is enabled for the location.
* `display_time_unit` - (String) Display Time Unit. The time unit to display for IP Surrogate idle time to disassociation.
* `dn_bandwidth` - (Optional) Download bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `up_bandwidth` - (Optional) Upload bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `idle_time_in_minutes` - (Number) Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled.
* `ips_control` - (Optional) Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.
* `ofw_enabled` - (Optional) Enable Firewall. When set to true, Firewall is enabled for the location.
* `parent_id` - (Optional) - Parent Location ID. If this ID does not exist or is `0`, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: `SUB`
* `ports` - (Optional) IP ports that are associated with the location.
* `ssl_scan_enabled` - (Optional) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
* `surrogate_ip` - (Optional) Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses.
* `surrogate_ip_enforced_for_known_browsers` - (Optional) Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers.
* `surrogate_refresh_time_in_minutes` - (Optional) Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates.
* `surrogate_refresh_time_unit` - (Optional) Display Refresh Time Unit. The time unit to display for refresh time for re-validation of surrogacy.
* `xff_forward_enabled` - (Optional) Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
* `zapp_ssl_scan_enabled` - (Optional) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.

`managed_by` - (List of Object)
    *`id` - (Optional) Identifier that uniquely identifies an entity
    * `name` - (Optional) The configured name of the entity
    * `extensions` - (Optional)
