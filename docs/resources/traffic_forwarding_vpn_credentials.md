---
subcategory: "Traffic Forwarding VPN Credentials"
layout: "zia"
page_title: "ZIA: traffic_forwarding_vpn_credentials"
description: |-
        Adds VPN credentials that can be associated to locations.
---


# zia_traffic_forwarding_vpn_credentials (Resource)

The **zia_traffic_forwarding_vpn_credentials** - Adds VPN credentials that can be associated to locations. VPN is one way to route traffic from customer locations to the cloud. Site-to-site IPSec VPN credentials can be identified by the cloud through one of the following methods:

* Common Name (CN) of IPSec Certificate
* VPN User FQDN - requires VPN_SITE_TO_SITE subscription
* VPN IP Address - requires VPN_SITE_TO_SITE subscription
* Extended Authentication (XAUTH) or hosted mobile UserID - requires VPN_MOBILE subscription

## Example Usage

```hcl
# ZIA Traffic Forwarding - VPN Credentials
resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type = "UFQDN"
    fqdn = "sjc-1-37@acme.com"
    comments = "created automatically"
    pre_shared_key = "newPassword123!"
}

output "zia_traffic_forwarding_vpn_credentials"{
    value = zia_traffic_forwarding_vpn_credentials.example
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created.
* `fqdn` - (Required) Fully Qualified Domain Name. Applicable only to `UFQDN` or `XAUTH` (or `HOSTED_MOBILE_USERS`) auth type.
* `pre_shared_key` - (Required) Pre-shared key. This is a required field for UFQDN and IP auth type.
* `comments` - (Optional) Additional information about this VPN credential.

:warning: The `pre_shared_key` parameter is ommitted from the output for security reasons.
