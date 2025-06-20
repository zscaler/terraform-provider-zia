---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_application_services"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/networkServices-get
  API documentation https://help.zscaler.com/zia/firewall-policies#/networkServices-get
  Get information about firewall rule network application services.
---

# zia_firewall_filtering_application_services (Data Source)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/networkServices-get)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/networkServices-get)

The **zia_firewall_filtering_application_services** data source to get information about a network application services available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network application services rule.

## Example Usage

```hcl
# ZIA Network Application Service
data "zia_firewall_filtering_application_services" "example" {
  name = "SKYPEFORBUSINESS"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "FILE_SHAREPT_ONEDRIVE"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "EXCHANGEONLINE"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "M365COMMON"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "ZOOMMEETING"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "WEBEXMEETING"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "WEBEXCALLING"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "RINGCENTRALMEETING"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "GOTOMEETING"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "GOTOMEETING_INROOM"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "LOGMEINMEETING"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "LOGMEINRESCUE"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "AWS"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "GCP"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "ZSCALER_CLOUD_ENDPOINTS"
}

data "zia_firewall_filtering_application_services" "example" {
  name = "TALK_DESK"
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
