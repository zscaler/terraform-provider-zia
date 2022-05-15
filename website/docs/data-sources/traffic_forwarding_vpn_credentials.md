---
subcategory: "Traffic Forwarding VPN Credentials"
layout: "zia"
page_title: "ZIA: traffic_forwarding_vpn_credentials"
description: |-
        Gets VPN credentials that can be associated to locations.
---


# zia_traffic_forwarding_vpn_credentials (Data Source)

The **zia_traffic_forwarding_vpn_credentials** - data source retrieves VPN credentials that can be associated to locations. VPN is one way to route traffic from customer locations to the cloud. Site-to-site IPSec VPN credentials can be identified by the cloud through one of the following methods:

* Common Name (CN) of IPSec Certificate
* VPN User FQDN - requires VPN_SITE_TO_SITE subscription
* VPN IP Address - requires VPN_SITE_TO_SITE subscription
* Extended Authentication (XAUTH) or hosted mobile UserID - requires VPN_MOBILE subscription

## Example Usage

```hcl
# ZIA Traffic Forwarding - VPN Credentials
data "zia_traffic_forwarding_vpn_credentials" "example"{
    fqdn = "sjc-1-37@acme.com"
}

output "zia_vpn_credentials"{
    value = data.zia_traffic_forwarding_vpn_credentials.example
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Number) VPN credential id
* `type` - (String) VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created.
* `fqdn` - (String) Fully Qualified Domain Name. Applicable only to `UFQDN` or `XAUTH` (or `HOSTED_MOBILE_USERS`) auth type.
* `pre_shared_key` - (String) Pre-shared key. This is a required field for UFQDN and IP auth type.
* `comments` - (String) Additional information about this VPN credential.

`location` - (Set of Object) Location that is associated to this VPN credential. Non-existence means not associated to any location.

* `id` - (Number) Identifier that uniquely identifies an entity
* `name` - (String) The configured name of the entity
* `extensions` - (Map of String)

* `managed_by` - (Set of Object) SD-WAN Partner that manages the location. If a partner does not manage the locaton, this is set to Self.
* `id` - (Number) Identifier that uniquely identifies an entity
* `name` - (String) The configured name of the entity
* `extensions` - (Map of String)

:warning: The `pre_shared_key` parameter is ommitted from the output for security reasons.
