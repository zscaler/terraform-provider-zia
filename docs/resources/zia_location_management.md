---
subcategory: "Location Management"
layout: "zscaler"
page_title: "ZIA: location_management"
description: |-
  Official documentation https://help.zscaler.com/zia/about-locations
  API documentation https://help.zscaler.com/zia/location-management#/locations-get
  Creates and manages ZIA locations and sub-locations.
---

# zia_location_management (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-locations)
* [API documentation](https://help.zscaler.com/zia/location-management#/locations-get)

The **zia_location_management** resource allows the creation and management of ZIA locations in the Zscaler Internet Access. This resource can then be associated with a:

* Static IP resource
* GRE Tunnel resource
* VPN credentials resource
* URL filtering, firewall filtering annd several other types of rule based resources

## Example Usage - Location Management with UFQDN VPN Credential

```hcl
resource "zia_location_management" "usa_sjc37"{
    name                        = "USA_SJC_37"
    description                 = "Created with Terraform"
    country                     = "UNITED_STATES"
    tz                          = "UNITED_STATES_AMERICA_LOS_ANGELES"
    auth_required               = true
    idle_time_in_minutes        = 720
    display_time_unit           = "HOUR"
    surrogate_ip                = true
    xff_forward_enabled         = true
    ofw_enabled                 = true
    ips_control                 = true
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
    }
    depends_on = [zia_traffic_forwarding_vpn_credentials.usa_sjc37 ]
}

resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type            = "UFQDN"
    fqdn            = "usa_sjc37@acme.com"
    comments        = "USA - San Jose IPSec Tunnel"
    pre_shared_key  = "***************"
}
```

## Example Usage - Location Management with IP VPN Credential

```hcl
# ZIA Location Management with IP VPN Credential
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
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
       ip_address = zia_traffic_forwarding_static_ip.usa_sjc37.ip_address
    }
}

resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type        = "IP"
    ip_address  =  zia_traffic_forwarding_static_ip.usa_sjc37.ip_address
    depends_on = [ zia_traffic_forwarding_static_ip.usa_sjc37 ]
    comments    = "Created via Terraform"
    pre_shared_key = "******************"
}

resource "zia_traffic_forwarding_static_ip" "usa_sjc37"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "SJC37 - Static IP"
    geo_override = false
}
```

## Example Usage - Location Management with Manual and Dynamic Location Groups

```hcl
# Retrieve ZIA Manual Location Groups
data "zia_location_groups" "this"{
    name = "SDWAN_CAN"
}

# ZIA Location Management with UFQDN VPN Credential
resource "zia_location_management" "usa_sjc37"{
    name                        = "USA_SJC_37"
    description                 = "Created with Terraform"
    country                     = "UNITED_STATES"
    tz                          = "UNITED_STATES_AMERICA_LOS_ANGELES"
    state                       = "California"
    auth_required               = true
    idle_time_in_minutes        = 720
    display_time_unit           = "HOUR"
    surrogate_ip                = true
    xff_forward_enabled         = true
    ofw_enabled                 = true
    ips_control                 = true
    profile                     = "CORPORATE"
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
    }
    static_location_groups {
      id = [data.zia_location_groups.this.id]
    }
    depends_on = [zia_traffic_forwarding_vpn_credentials.usa_sjc37 ]
}

resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type            = "UFQDN"
    fqdn            = "usa_sjc37@acme.com"
    comments        = "USA - San Jose IPSec Tunnel"
    pre_shared_key  = "***************"
}
```

## Example Usage - Location Management with Excluded Manual and Dynamic Location Groups

```hcl
# Retrieve ZIA Manual Location Groups
data "zia_location_groups" "this"{
    name = "SDWAN_CAN"
}

# ZIA Location Management with UFQDN VPN Credential
resource "zia_location_management" "usa_sjc37"{
    name                        = "USA_SJC_37"
    description                 = "Created with Terraform"
    country                     = "UNITED_STATES"
    tz                          = "UNITED_STATES_AMERICA_LOS_ANGELES"
    state                       = "California"
    auth_required               = true
    idle_time_in_minutes        = 720
    display_time_unit           = "HOUR"
    surrogate_ip                = true
    xff_forward_enabled         = true
    ofw_enabled                 = true
    ips_control                 = true
    exclude_from_dynamic_groups = true
    exclude_from_manual_groups  = true
    profile                     = "CORPORATE"
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
    }
    depends_on = [zia_traffic_forwarding_vpn_credentials.usa_sjc37 ]
}

resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type            = "UFQDN"
    fqdn            = "usa_sjc37@acme.com"
    comments        = "USA - San Jose IPSec Tunnel"
    pre_shared_key  = "***************"
}
```

```hcl
resource "zia_location_management" "usa_sjc37_office_branch01"{
    name = "USA_SJC37_Office-Branch01"
    description = "Created with Terraform"
    country = "UNITED_STATES"
    tz = "UNITED_STATES_AMERICA_LOS_ANGELES"
    profile = "CORPORATE"
    parent_id = zia_location_management.usa_sjc37.id
    depends_on = [ zia_traffic_forwarding_static_ip.usa_sjc37, zia_traffic_forwarding_vpn_credentials.usa_sjc37, zia_location_management.usa_sjc37 ]
    auth_required = true
    idle_time_in_minutes = 720
    display_time_unit = "HOUR"
    surrogate_ip = true
    ofw_enabled = true
    ip_addresses = [ "10.5.0.0-10.5.255.255" ]
    up_bandwidth = 10000
    dn_bandwidth = 10000
}
```

## Example Usage - SubLocation Management with UFQDN VPN Credential

```hcl

resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type            = "UFQDN"
    fqdn            = "usa_sjc37@acme.com"
    comments        = "USA - San Jose IPSec Tunnel"
    pre_shared_key  = "***************"
}

resource "zia_location_management" "usa_sjc37"{
    name                        = "USA_SJC_37"
    description                 = "Created with Terraform"
    country                     = "UNITED_STATES"
    tz                          = "UNITED_STATES_AMERICA_LOS_ANGELES"
    auth_required               = true
    idle_time_in_minutes        = 720
    display_time_unit           = "HOUR"
    surrogate_ip                = true
    xff_forward_enabled         = true
    ofw_enabled                 = true
    ips_control                 = true
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
    }
}

resource "zia_location_management" "usa_sjc37_office_branch01"{
    name = "USA_SJC37_Office-Branch01"
    description = "Created with Terraform"
    country = "UNITED_STATES"
    tz = "UNITED_STATES_AMERICA_LOS_ANGELES"
    profile = "CORPORATE"
    parent_id = zia_location_management.usa_sjc37.id
    zia_traffic_forwarding_vpn_credentials.usa_sjc37, zia_location_management.usa_sjc37
    auth_required = true
    idle_time_in_minutes = 720
    display_time_unit = "HOUR"
    surrogate_ip = true
    ofw_enabled = true
    ip_addresses = [ "10.5.0.0-10.5.255.255" ]
    up_bandwidth = 10000
    dn_bandwidth = 10000
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) - Location Name.
* `ip_addresses` - (Required) For locations: IP addresses of the egress points that are provisioned in the Zscaler Cloud. Each entry is a single IP address (e.g., `238.10.33.9`). For sub-locations: Egress, internal, or GRE tunnel IP addresses. Each entry is either a single IP address or range (e.g., `10.10.33.1-10.10.33.10`). The value is required if `vpn_credentials` are not defined.
* `vpn_credentials`
  * `id` - (Optional) VPN credential resource id. The value is required if `ip_addresses` are not defined.

### Optional

* `description` - (String) Additional notes or information regarding the location or sub-location. The description cannot exceed 1024 characters.
* `country` - (Optional) Country - See list of supported country names [here](https://help.zscaler.com/zia/location-management#/locations-post)
* `state` - (String) Country - State Name i.e `California`
* `tz` - (String) Timezone of the location. If not specified, it defaults to GMT. See list of supported country names [here](https://help.zscaler.com/zia/location-management#/locations-post)
* `profile` - (Optional) Profile tag that specifies the location traffic type. If not specified, this tag defaults to `Unassigned`. The supported options are: `NONE`, `CORPORATE`, `SERVER`, `GUESTWIFI`, `IOT`, `WORKLOAD`.

* `aup_block_internet_until_accepted` - (Boolean) For First Time AUP Behavior, Block Internet Access. When set, all internet access (including non-HTTP traffic) is disabled until the user accepts the AUP.
* `aup_enabled` - (Boolean) Enable AUP. When set to true, AUP is enabled for the location.
* `aup_force_ssl_inspection` - (Boolean) For First Time AUP Behavior, Force SSL Inspection. When set, Zscaler will force SSL Inspection in order to enforce AUP for HTTPS traffic.
* `aup_timeout_in_days` - (Number) Custom AUP Frequency. Refresh time (in days) to re-validate the AUP.
* `cookies_and_proxy` - (Boolean) Enable Cookies and proxy feature
* `digest_auth_enabled` - (Boolean) Enable Digest Auth feature
* `kerberos_auth` - (Boolean) Enable Kerberos Auth feature
* `auth_required` - (Boolean) Enforce Authentication. Required when ports are enabled, IP Surrogate is enabled, or Kerberos Authentication is enabled.
* `caution_enabled` - (Boolean) Enable Caution. When set to true, a caution notifcation is enabled for the location.
* `display_time_unit` - (String) Display Time Unit. The time unit to display for IP Surrogate idle time to disassociation.
* `dn_bandwidth` - (Number) Download bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `up_bandwidth` - (Number) Upload bandwidth in bytes. The value `0` implies no Bandwidth Control enforcement.
* `idle_time_in_minutes` - (Number) Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled.
* `ips_control` - (Boolean) Enable IPS Control. When set to true, IPS Control is enabled for the location if Firewall is enabled.
* `ofw_enabled` - (Boolean) Enable Firewall. When set to true, Firewall is enabled for the location.
* `parent_id` - (Number) - Parent Location ID. If this ID does not exist or is `0`, it is implied that it is a parent location. Otherwise, it is a sub-location whose parent has this ID. x-applicableTo: `SUB`
* `ports` - (List of Numbers) IP ports that are associated with the location.
* `ssl_scan_enabled` - (Boolean) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.
* `surrogate_ip` - (Boolean) Enable Surrogate IP. When set to true, users are mapped to internal device IP addresses.
* `surrogate_ip_enforced_for_known_browsers` - (Boolean) Enforce Surrogate IP for Known Browsers. When set to true, IP Surrogate is enforced for all known browsers.
* `surrogate_refresh_time_in_minutes` - (Number) Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates.
* `surrogate_refresh_time_unit` - (String) Display Refresh Time Unit. The time unit to display for refresh time for re-validation of surrogacy.
* `xff_forward_enabled` - (Boolean) Enable XFF Forwarding. When set to true, traffic is passed to Zscaler Cloud via the X-Forwarded-For (XFF) header.
* `zapp_ssl_scan_enabled` - (Boolean) This parameter was deprecated and no longer has an effect on SSL policy. It remains supported in the API payload in order to maintain backwards compatibility with existing scripts, but it will be removed in future.

* `iot_discovery_enabled` - (Boolean) Enable IOT Discovery at the location

* `iot_enforce_policy_set` - (Boolean) Enable IOT Policy at the location

* `other_sub_location` - (Boolean) If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv4 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other and it can be renamed, if required.

* `other6_sub_location` - (Boolean) If set to true, indicates that this is a default sub-location created by the Zscaler service to accommodate IPv6 addresses that are not part of any user-defined sub-locations. The default sub-location is created with the name Other6 and it can be renamed, if required. This field is applicable only if ipv6Enabled is set is true.

* `sub_loc_scope_values` - (List of Strings) Specifies values for the selected sublocation scope type

* `sub_loc_scope` - (Strings) Defines a scope for the sublocation from the available types to segregate workload traffic from a single sublocation to apply different Cloud Connector and ZIA security policies. This field is only available for the Workload traffic type sublocations whose parent locations are associated with Amazon Web Services (AWS) Cloud Connector groups. The supported options are: `VPC_ENDPOINT`, `VPC`, `NAMESPACE`, `ACCOUNT`

* `sub_loc_acc_ids` - (List of Strings) Specifies one or more Amazon Web Services (AWS) account IDs. These AWS accounts are associated with the parent location of this sublocation created in the Zscaler Cloud & Branch Connector Admin Portal.

* `ipv6_enabled` - (Boolean) If set to true, IPv6 is enabled for the location and IPv6 traffic from the location can be forwarded to the Zscaler service to enforce security policies.

* `default_extranet_ts_pool` - (Boolean) Indicates that the traffic selector specified in the extranet is the designated default traffic selector

* `default_extranet_dns` - (Boolean) Indicates that the DNS server configuration used in the extranet is the designated default DNS server

* `ipv6_dns_64prefix` - (Optional) Name-ID pair of the NAT64 prefix configured as the DNS64 prefix for the location. If specified, the DNS64 prefix is used for the IP addresses that reside in this location. If not specified, a prefix is selected from the set of supported prefixes. This field is applicable only if ipv6Enabled is set is true.

* `dynamic_location_groups` - (List of Object) Dynamic location groups the location belongs to
  * `id` - (Optional) The Identifier that uniquely identifies an entity

* `static_location_groups` - (List of Object) Manual location groups the location belongs to
  * `id` - (Optional) The Identifier that uniquely identifies an entity

* `exclude_from_dynamic_groups` - (Boolean) Enable to prevent the location from being assigned to any dynamic groups and to remove it from any dynamic groups it's already assigned to

* `exclude_from_manual_groups` - (Boolean) Enable to prevent the location from being added to manual groups and to remove it from any manual groups it's already assigned to

  **NOTE** The attributes, ``dynamic_location_groups``, and ``static_location_groups`` CANNOT be configured if the attributes `exclude_from_dynamic_groups` and/or `exclude_from_manual_groups` are set to `true`

* `extranet` - (Block, Max: 1) The ID of the extranet resource that must be assigned to the location
  * `id` - (int) The Identifier that uniquely identifies an entity

* `extranet_ip_pool` - (Block, Max: 1) The ID of the traffic selector specified in the extranet
  * `id` - (int) The Identifier that uniquely identifies an entity

* `extranet_dns` - (Block, Max: 1) The ID of the DNS server configuration used in the extranet
  * `id` - (int) The Identifier that uniquely identifies an entity

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_location_management** can be imported by using `<LOCATION_ID>` or `<LOCATION_NAME>` as the import ID.

For example:

```shell
terraform import zia_location_management.example <location_id>
```

or

```shell
terraform import zia_location_management.example <location_name>
```
