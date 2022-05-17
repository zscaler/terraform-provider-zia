---
subcategory: "Firewall Policies"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): firewall_filtering_network_service_groups"
sidebar_current: "docs-resource-zia-firewall-filtering-network-service-groups"
description: |-
  Creates and manages ZIA Cloud firewall Network Service Groups.
---

# Resource: zia_firewall_filtering_network_service_groups

The **zia_firewall_filtering_network_service_groups** resource allows the creation and management of ZIA Cloud Firewall IP network service groups in the Zscaler Internet Access. This resource can then be associated with a ZIA cloud firewall filtering rule.

## Example Usage

```hcl
data "zia_firewall_filtering_network_service" "example1" {
  name = "FTP"
}

data "zia_firewall_filtering_network_service" "example2" {
  name = "NETBIOS"
}

data "zia_firewall_filtering_network_service" "example3" {
  name = "DNS"
}

# Add network services to a network services group
resource "zia_firewall_filtering_network_service_groups" "example"{
    name        = "example"
    description = "example"
    services {
        id = [
            data.zia_firewall_filtering_network_service.example1.id,
            data.zia_firewall_filtering_network_service.example2.id,
            data.zia_firewall_filtering_network_service.example3.id
        ]
    }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Name of the network service group
* `services` - (Required) Any number of network services ID to be added to the group

### Optional

* `description` (Optional) - Description of the network services group
