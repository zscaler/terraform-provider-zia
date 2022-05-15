---
subcategory: "Firewall Filtering Network Service"
layout: "zia"
page_title: "ZIA: firewall_filtering_network_service"
description: |-
  Retrieve ZIA firewall rule network service.
  
---

# zia_firewall_filtering_network_service (Data Source)

The **zia_firewall_filtering_network_service** data source provides details about a specific network service available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network service rule.

## Example Usage

```hcl
# ZIA Network Service
data "zia_firewall_filtering_network_service_groups" "example"{
    name = "Corporate Custom SSH TCP_10022"
}

output "zia_firewall_filtering_network_service_groups" {
  value = data.zia_firewall_filtering_network_service_groups.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Enter a name for the application layer service that you want to control. It can include any character and spaces.
* `src_tcp_ports` - (Optional) The TCP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
* `dest_tcp_ports` - (Required) The TCP destination port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
* `src_udp_ports` - The UDP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.
* `dest_udp_ports` - The UDP source port number (example: 50) or port number range (example: 1000-1050), if any, that is used by the network service.

### Optional

* `description` - (String) (Optional) Enter additional notes or information. The description cannot exceed 10240 characters.
