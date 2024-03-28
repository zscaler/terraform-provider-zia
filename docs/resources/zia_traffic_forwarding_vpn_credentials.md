---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_vpn_credentials"
description: |-
    Creates and manages VPN credentials that can be associated to locations.
---

# Resource: zia_traffic_forwarding_vpn_credentials

The **zia_traffic_forwarding_vpn_credentials** creates and manages VPN credentials that can be associated to locations. VPN is one way to route traffic from customer locations to the cloud. Site-to-site IPSec VPN credentials can be identified by the cloud through one of the following methods:

* Common Name (CN) of IPSec Certificate
* VPN User FQDN - requires VPN_SITE_TO_SITE subscription
* VPN IP Address - requires VPN_SITE_TO_SITE subscription
* Extended Authentication (XAUTH) or hosted mobile UserID - requires VPN_MOBILE subscription

## Example Usage

```hcl
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
# ZIA Traffic Forwarding - VPN Credentials (UFQDN)
resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type            = "UFQDN"
    fqdn            = "sjc-1-37@acme.com"
    comments        = "Example"
    pre_shared_key = "*********************"
}
```

```hcl
# ZIA Traffic Forwarding - VPN Credentials (IP)
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_traffic_forwarding_vpn_credentials" "example"{
    type            = "IP"
    ip_address      = zia_traffic_forwarding_static_ip.example.ip_address
    comments        = "Example"
    pre_shared_key  = "*********************"
    depends_on = [ zia_traffic_forwarding_static_ip.example ]
}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address      =  "1.1.1.1"
    routable_ip     = true
    comment         = "Example"
    geo_override    = true
    latitude        = -36.848461
    longitude       = 174.763336
}
```

~> **NOTE** For VPN Credentials of Type `IP` a static IP resource must be created first.

## Argument Reference

The following arguments are supported:

* `type` - (Required) VPN authentication type (i.e., how the VPN credential is sent to the server). It is not modifiable after VpnCredential is created. The supported values are: `UFQDN` and `IP`
* `fqdn` - (Required) Fully Qualified Domain Name. Applicable only to `UFQDN` or `XAUTH` (or `HOSTED_MOBILE_USERS`) auth type.
* `pre_shared_key` - (Required) Pre-shared key. This is a required field for UFQDN and IP auth type.

### Optional

* `comments` - (Optional) Additional information about this VPN credential.
* `ip_address` - (Optional) IP Address for the VON credentials. The parameter becomes required if `type = IP`

!> **WARNING:** The `pre_shared_key` parameter is ommitted from the output for security reasons.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_traffic_forwarding_vpn_credentials** can be imported by using one of the following prefixes as the import ID:

* `'IP'` - Imports all VPN Credentials of type IP

```shell
$ terraform import zia_traffic_forwarding_vpn_credentials.example 'IP'
```

* `'UFQDN'` - Imports all VPN Credentials of type UFQDN

```shell
$ terraform import zia_traffic_forwarding_vpn_credentials.this 'UFQDN'
```

* `UFQDN'` - Imports a VPN Credentials of type UFQDN containing a specific UFQDN address

```shell
$ terraform import zia_traffic_forwarding_vpn_credentials.example 'testvpn@example.com'
```

* `IP Address'` - Imports a VPN Credentials of type IP containing a specific IP address

```shell
$ terraform import zia_traffic_forwarding_vpn_credentials.example '1.1.1.1'
```
