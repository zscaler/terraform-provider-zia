---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZIA: location_management"
description: |-
  Official documentation https://help.zscaler.com/zia/about-locations
  API documentation https://help.zscaler.com/zia/location-management#/locations-get
  Get information about Location.

---

# zia_location_management (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-locations)
* [API documentation](https://help.zscaler.com/zia/location-management#/locations-get)

Use the **zia_location_management** data source to get information about a location resource available in the Zscaler Internet Access Location Management. This resource can then be referenced in multiple other resources, such as URL Filtering Rules, Firewall rules etc.

## Example Usage

```hcl
# ZIA Location Managemeent
data "zia_location_management" "example"{
    name = "San Jose"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) - The name of the location to be exported.
* `id` - (Required) - The ID of the location to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `aup_block_internet_until_accepted` - (Boolean) For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.
* `aup_enabled` - (Boolean) Enable AUP. When set to true, AUP is enabled for the location.
* `aup_force_ssl_inspection` - (Boolean) For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler will force SSL Inspection in order to enforce AUP for HTTPS traffic.
* `aup_timeout_in_days` - (Number) Custom AUP Frequency. Refresh time (in days) to re-validate the AUP.
* `auth_required` - (Boolean) Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
* `caution_enabled` - (Boolean) Enable Caution. When set to true, a caution notifcation is enabled for the location.
* `country` - (String) Country
* `description` - (String) Additional notes or information regarding the location or sub-location. The description cannot exceed 1024 characters.
* `display_time_unit` - (String) Display Time Unit. The time unit to display for IP Surrogate idle time to disassociation.
* `dn_bandwidth` - (Number) Download bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `up_bandwidth` - (Number) Upload bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `id` - (Number) Location ID.
* `idle_time_in_minutes` - (Number) Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled.
* `ip_addresses` - (List of String) For locations: IP addresses of the egress points that are provisioned in the Zscaler Cloud. Each entry is a single IP address (e.g., `238.10.33.9`). For sub-locations: Egress, internal, or GRE tunnel IP addresses. Each entry is either a single IP address, CIDR (e.g., `10.10.33.0/24`), or range (e.g., `10.10.33.1-10.10.33.10`)).
* `ips_control` - (Boolean) Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.
* `ofw_enabled` - (Boolean) Enable Firewall. When set to true, Firewall is enabled for the location.
* `parent_id` - (Number) - Parent Location ID. If this ID does not exist or is `0`, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: `SUB`
* `ports` - (List of String) IP ports that are associated with the location.
* `profile` - (String) Profile tag that specifies the location traffic type. If not specified, this tag defaults to `Unassigned`.
* `ssl_scan_enabled` - (Boolean) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
* `surrogate_ip` - (Boolean) Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses.
* `surrogate_ip_enforced_for_known_browsers` - (Boolean) Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers.
* `surrogate_refresh_time_in_minutes` - (Number) Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates.
* `surrogate_refresh_time_unit` - (String) Display Refresh Time Unit. The time unit to display for refresh time for re-validation of surrogacy.
* `tz` - (String) Timezone of the location. If not specified, it defaults to GMT.
* `xff_forward_enabled` - (Boolean) Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
* `zapp_ssl_scan_enabled` - (Boolean) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.

* `vpn_credentials`
  * `comments` - (String) Additional information about this VPN credential.
    Additional information about this VPN credential.
  * `fqdn` - (String) Fully Qualified Domain Name. Applicable only to `UFQDN` or `XAUTH` (or `HOSTED_MOBILE_USERS`) auth type.
  * `id` - (Number) VPN credential id
  * `pre_shared_key` - (String) Pre-shared key. This is a required field for `UFQDN` and IP auth type.
  * `type` - (String) VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created.

* `location` - (List of Object)
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String)

* `managed_by` - (List of Object)
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map of String)
