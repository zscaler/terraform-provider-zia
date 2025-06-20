---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_network_application_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get
  API documentation https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get
  Creates and manages ZIA Cloud firewall Network Application Groups.
---

# zia_firewall_filtering_network_application_groups (Resource)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/networkApplicationGroups/{groupId}-get)

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

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_firewall_filtering_network_application_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import zia_firewall_filtering_network_application_groups.example <group_id>
```

or

```shell
terraform import zia_firewall_filtering_network_application_groups.example <group_name>
```
