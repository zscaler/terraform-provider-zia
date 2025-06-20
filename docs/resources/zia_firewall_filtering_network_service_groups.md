---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_service_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/networkServiceGroups-get
  API documentation https://help.zscaler.com/zia/firewall-policies#/networkServiceGroups-get
  Creates and manages ZIA Cloud firewall Network Service Groups.
---

# zia_firewall_filtering_network_service_groups (Resource)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/networkServiceGroups-get)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/networkServiceGroups-get)

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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**firewall_filtering_network_service_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import firewall_filtering_network_service_groups.example <group_id>
```

or

```shell
terraform import firewall_filtering_network_service_groups.example <group_name>
```
