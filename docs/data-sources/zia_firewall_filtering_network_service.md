---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_service"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/networkServices-get
  API documentation https://help.zscaler.com/zia/firewall-policies#/networkServices-get
  Get information about firewall rule network services.
---

# zia_firewall_filtering_network_service (Data Source)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/networkServices-get)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/networkServices-get)

The **zia_firewall_filtering_network_service** data source to get information about a network service available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network service rule.

## Example Usage

```hcl
data "zia_firewall_filtering_network_service" "example_protocol" {
  protocol = "TCP"
}

# Test locale only
data "zia_firewall_filtering_network_service" "example_locale" {
  locale = "en-US"
}

# Test protocol + locale
data "zia_firewall_filtering_network_service" "example_both" {
  protocol = "TCP"
  locale   = "en-US"
}

data "zia_firewall_filtering_network_service" "dns" {
  name = "DNS"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the application layer service that you want to control. It can include any character and spaces.
See the [Network Services API](https://help.zscaler.com/zia/firewall-policies#/networkServices-get) for the list of available services. Check the attribute `tag` in the API documentation for the updated list.

* `id` - (Optional) The ID of the application layer service to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) Enter additional notes or information. The description cannot exceed 10240 characters.
* `protocol` - (String) Filter based on the network service protocol. Supported values are: `TCP`, `UDP`, `ICMP`, `GRE`, `ESP`, `OTHER`
* `locale` - (String) When set to one of the supported locales (i.e., `en-US`, `de-DE`, `es-ES`, `fr-FR`, `ja-JP`,`fr-FR`, `ja-JP`, `zh-CN`), the network service's description is localized into the requested language.
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
* `tag` - (string) - Supported network services names returned by the API. See the [Network Services API](https://help.zscaler.com/zia/firewall-policies#/networkServices-get) for the list of available services. Check the attribute `tag` in the API documentation for the updated list.
