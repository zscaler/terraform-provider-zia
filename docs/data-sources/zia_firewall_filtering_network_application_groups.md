---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_application_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get
  API documentation https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get
  Get information about Network Application groups.
---


# zia_firewall_filtering_network_application_groups (Data Source)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get)

Use the **zia_firewall_filtering_network_application_groups** data source to get information about network application groups available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

## Example Usage

```hcl
# ZIA IP Source Groups
data "zia_firewall_filtering_network_application_groups" "example" {
    name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the ip source group resource.
* `name` - (Required) The name of the ip source group to be exported.
* `network_applications` - (Required) Any number of applications to be added to the group
  * Refer to the Zscaler API Swagger for the complete list of applications [ZIA API Guide](https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups-get)

### Optional

* `description` (Optional) - Description of the network application group
