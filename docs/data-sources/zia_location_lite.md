---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZIA: location_lite"
description: |-
  Get information about Location Lite.

---

# Data Source: zia_location_lite

Use the **zia_location_lite** data source to get information about a location in lite mode option available in the Zscaler Internet Access. This data source can be used to retrieve the Road Warrior location to then associated with one of the following resources: ``zia_url_filtering_rules``, ``zia_firewall_filtering_rule`` and ``zia_dlp_web_rules`

```hcl
# Retrieve ZIA Location Lite
data "zia_location_lite" "this" {
 name = "Road Warrior"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Location group name
* `id` - (Optional) Unique identifier for the location group

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `kerberos_auth` - (Boolean)
* `digest_auth_enabled` - (Boolean)
* `parent_id` - (Number) - Parent Location ID. If this ID does not exist or is `0`, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: `SUB`
* `tz` - (String) Timezone of the location. If not specified, it defaults to GMT.
* `zapp_ssl_scan_enabled` - (Boolean) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
* `xff_forward_enabled` - (Boolean) Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
* `surrogate_ip` - (Boolean) Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses.
* `surrogate_ip_enforced_for_known_browsers` - (Boolean) Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers.
* `ofw_enabled` - (Boolean) Enable Firewall. When set to true, Firewall is enabled for the location.
* `ips_control` - (Boolean) Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.
* `aup_enabled` - (Boolean) Enable AUP. When set to true, AUP is enabled for the location.
* `caution_enabled` - (Boolean) Enable Caution. When set to true, a caution notifcation is enabled for the location.
* `aup_block_internet_until_accepted` - (Boolean) For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.
* `aup_force_ssl_inspection` - (Boolean) For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler will force SSL Inspection in order to enforce AUP for HTTPS traffic.
* `ec_location` - (Boolean)
* `other_sub_location` - (Boolean) If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv4 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other and it can be renamed, if required.
* `other6_sub_location` - (Boolean) If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv6 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other6 and it can be renamed, if required. This field is applicable only if ipv6Enabled is set is true
* `ipv6_enabled` - (Number) If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.