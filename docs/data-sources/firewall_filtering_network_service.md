---
subcategory: "Firewall Policies"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): firewall_filtering_network_service"
sidebar_current: "docs-datasource-zia-firewall-filtering-network-service"
description: |-
    Get information about firewall rule network services.

---

# Data Source: zia_firewall_filtering_network_service

The **zia_firewall_filtering_network_service** data source to get information about a network service available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network service rule.

## Example Usage

```hcl
# ZIA Network Service
data "zia_firewall_filtering_network_service" "example" {
  name = "ICMP_ANY"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the application layer service that you want to control. It can include any character and spaces.
* `id` - (Optional) The ID of the application layer service to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) (Optional) Enter additional notes or information. The description cannot exceed 10240 characters.
* `type` - (String) - Supported values are: `STANDARD`, `PREDEFINED` and `CUSTOM`
* `is_name_l10n_tag` - (Bool) - Default: false
* `src_tcp_ports` - (Optional) The TCP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service
  * `start` - (Number)
  * `end` - (Number)
* `dest_tcp_ports` - (Required) The TCP destination port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
  * `start` - (Number)
  * `end` - (Number)
* `src_udp_ports` - The UDP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
  * `start` - (Number)
  * `end` - (Number)
* `dest_udp_ports` - The UDP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
  * `start` - (Number)
  * `end` - (Number)
