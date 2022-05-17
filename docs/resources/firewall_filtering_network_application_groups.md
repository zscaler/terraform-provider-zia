---
subcategory: "Firewall Policies"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): firewall_filtering_network_application_groups"
sidebar_current: "docs-resource-zia-firewall-filtering-network-application-groups"
description: |-
  Creates and manages ZIA Cloud firewall Network Application Groups.
---


# Resource: zia_firewall_filtering_network_application_groups

The **zia_firewall_filtering_network_application_groups** resource allows the creation and management of ZIA Cloud Firewall IP source groups in the Zscaler Internet Access. This resource can then be associated with a ZIA cloud firewall filtering rule.

## Example Usage

```hcl
# Add applications to a network application group
resource "zia_firewall_filtering_network_application_groups" "example" {
  name        = "Example"
  description = "Example"
  network_applications = [ "LDAP", "LDAPS", "SRVLOC"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Network application group name
* `network_applications` - (Required) Any number of applications to be added to the group
  * Refer to the Zscaler API Swagger for the complete list of applications [ZIA API Guide](https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups-get)

### Optional

* `description` (Optional) - Description of the network application group
