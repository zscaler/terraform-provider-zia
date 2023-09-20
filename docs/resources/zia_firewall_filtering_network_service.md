---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_service"
description: |-
      Creates and manages a ZIA Cloud firewall network service.
---

# Resource: zia_firewall_filtering_network_service

The **zia_firewall_filtering_network_service** resource allows the creation and management of ZIA Cloud Firewall IP network services in the Zscaler Internet Access. This resource can then be associated with a ZIA cloud firewall filtering rule and network service group resources.

## Example Usage

```hcl
resource "zia_firewall_filtering_network_service" "example" {
  name        = "example"
  description = "example"
  src_tcp_ports {
    start = 5000
  }
  src_tcp_ports {
    start = 5001
  }
  src_tcp_ports {
    start = 5002
    end   = 5005
  }
  dest_tcp_ports {
    start = 5000
  }
    dest_tcp_ports {
    start = 5001
  }
  dest_tcp_ports {
    start = 5003
    end   = 5005
  }
  type = "CUSTOM"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the service
* `type` - (Required) Supported values: `STANDARD`, `PREDEFINED`, `CUSTOM`

!> **NOTE:** Resources of type `PREDEFINED` are built-in resources within the ZIA cloud and must be imported before the Terraform execution. Attempting to update the resource directly will return `DUPLICATE_ITEM` error message. To import a predefined built-in resource use the following command for example: `terraform import zia_firewall_filtering_network_service.this "DHCP"`

* `src_tcp_ports` - (Required) The TCP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service
  * `start` - (Required)
  * `end` - (Required)
* `dest_tcp_ports` - (Required) The TCP destination port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
  * `start` - (Number)
  * `end` - (Number)
* `src_udp_ports` - The UDP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
  * `start` - (Number)
  * `end` - (Number)
* `dest_udp_ports` - The UDP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
  * `start` - (Number)
  * `end` - (Number)

-> **NOTE** The `end` port parameter must always be greater than the value defined in the `start` port.

### Optional

* `description` - (Optional) Description of the service
* `is_name_l10n_tag` - (Optional
* `tag` - (Optional) The following values are supported: `"ICMP_ANY`, `"UDP_ANY"`, `"TCP_ANY"`, `"OTHER_NETWORK_SERVICE"`, `"DNS"`, `"NETBIOS"`, `"FTP"`, `"GNUTELLA"`, `"H_323"`, `"HTTP"`, `"HTTPS"`, `"IKE"`, `"IMAP"`, `"ILS"`, `"IKE_NAT"`, `"IRC"`, `"LDAP"`, `"QUIC"`, `"TDS"`, `"NETMEETING"`, `"NFS"`, `"NTP"`, `"SIP"`, `"SNMP"`, `"SMB"`, `"SMTP"`, `"SSH"`, `"SYSLOG"`, `"TELNET"`, `"TRACEROUTE"`, `"POP3"`, `"PPTP"`, `"RADIUS"`, `"REAL_MEDIA"`, `"RTSP"`, `"VNC"`, `"WHOIS"`, `"KERBEROS_SEC"`, `"TACACS"`, `"SNMPTRAP"`, `"NMAP"`, `"RSYNC"`, `"L2TP"`, `"HTTP_PROXY"`, `"PC_ANYWHERE"`, `"MSN"`, `"ECHO"`, `"AIM"`, `"IDENT"`, `"YMSG"`, `"SCCP"`, `"MGCP_UA"`, `"MGCP_CA"`, `"VDO_LIVE"`, `"OPENVPN"`, `"TFTP"`, `"FTPS_IMPLICIT"`, `"ZSCALER_PROXY_NW_SERVICES"`, `"GRE_PROTOCOL"`, `"ESP_PROTOCOL"`, `"DHCP"`
