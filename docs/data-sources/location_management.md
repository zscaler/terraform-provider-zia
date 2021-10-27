---
subcategory: "Location Management"
layout: "zia"
page_title: "ZIA: location_management"
description: |-
  Retrieve ZIA Location.
  
---

# zia_location_management (Data Source)

The **zia_location_management** -data source provides details about a specific location resource available in the Zscaler Internet Access Location Management.

## Example Usage

```hcl
# ZIA Network Service
data "zia_location_management" "example"{
    name = "San Jose"
}

output "zia_location_management_example" {
  value = data.zia_location_management.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) - Location Name.

### Read-Only

* `aup_block_internet_until_accepted` -(Boolean)
* `aup_enabled` - (Boolean)
* `aup_force_ssl_inspection` - (Boolean)
* `aup_timeout_in_days` - (Number)
* `auth_required` - (Boolean)
* `caution_enabled` - (Boolean)
* `country` - (String) Country
* `description` - (String)
* `display_time_unit` - (String)
* `dn_bandwidth` - (Number) Download bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `up_bandwidth` - (Number) Upload bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `id` - (Number) Location ID.
* `idle_time_in_minutes` - (Number)
* `ip_addresses` - (List of String) For locations: IP addresses of the egress points that are provisioned in the Zscaler Cloud. Each entry is a single IP address (e.g., `238.10.33.9`). For sub-locations: Egress, internal, or GRE tunnel IP addresses. Each entry is either a single IP address, CIDR (e.g., `10.10.33.0/24`), or range (e.g., `10.10.33.1-10.10.33.10`)).
* `ips_control` - (Boolean)
* `ofw_enabled` - (Boolean)
* `parent_id` - (Number) - Parent Location ID. If this ID does not exist or is `0`, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: `SUB`
* `ports` - (String) IP ports that are associated with the location.
* `profile` - (String)
* `ssl_scan_enabled` - (Boolean)
* `surrogate_ip` - (Boolean)
* `surrogate_ip_enforced_for_known_browsers` - (Boolean)
* `surrogate_refresh_time_in_minutes` - (Number)
* `surrogate_refresh_time_unit` - (String)
* `tz` - (String) Timezone of the location. If not specified, it defaults to GMT.
* `xff_forward_enabled` -(Boolean)
* `zapp_ssl_scan_enabled` -(Boolean)

`vpn_credentials`

* `comments` - (String) Additional information about this VPN credential.
Additional information about this VPN credential.
* `fqdn` - (String) Fully Qualified Domain Name. Applicable only to `UFQDN` or `XAUTH` (or `HOSTED_MOBILE_USERS`) auth type.
* `id` - (Number) VPN credential id
* `pre_shared_key` - (String) Pre-shared key. This is a required field for `UFQDN` and IP auth type.
* `type` - (String) VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created.

`location` - (List of Object)
    *`id` - (Number) Identifier that uniquely identifies an entity
    * `name` - (String) The configured name of the entity
    * `extensions` -(Map of String)

`managed_by` - (List of Object)
    *`id` - (Number) Identifier that uniquely identifies an entity
    * `name` - (String) The configured name of the entity
    * `extensions` -(Map of String)
