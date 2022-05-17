---
subcategory: "Firewall Policies"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): firewall_filtering_network_service_groups"
sidebar_current: "docs-datasource-zia-firewall-filtering-network-application-groups"
description: |-
  Get information about firewall rule network service groups.

---

# Data Source: zia_firewall_filtering_network_service_groups

Use the **zia_firewall_filtering_network_service_groups** data source to get information about a network service groups available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network service rule.

## Example Usage

```hcl
# ZIA Network Service Groups
data "zia_firewall_filtering_network_service_groups" "example"{
    name = "Corporate Custom SSH TCP_10022"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.
* `id` - (Optional) The ID of the ip source group to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String)
* `services` - (Number) The ID of this resource.
  * `id` - (Number)
  * `name` - (String)
  * `tag` - (String)
  * `src_tcp_ports` (Optional) The TCP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service
         - `start` - (Number)
         - `end` - (Number)
  * `dest_tcp_ports` - (Required) The TCP destination port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
         - `start` - (Number)
         - `end` - (Number)
  * `src_udp_ports`
         - `start` - (Number)
         - `end` - (Number)
  * `dest_udp_ports`
         - `start` - (Number)
         - `end` - (Number)
  * `type` - (String) - Supported values are: `STANDARD`, `PREDEFINED` and `CUSTOM`
  * `is_name_l10n_tag` - (Bool) - Default: false
