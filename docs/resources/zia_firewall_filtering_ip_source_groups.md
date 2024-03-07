---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_ip_source_groups"
description: |-
  Creates and manages ZIA Cloud firewall IP source groups.
---

# Resource: zia_firewall_filtering_ip_source_groups

The **zia_firewall_filtering_ip_source_groups** resource allows the creation and management of ZIA Cloud Firewall IP source groups in the Zscaler Internet Access. This resource can then be associated with a ZIA cloud firewall filtering rule.

## Example Usage

```hcl
# Add an IP address or addresses to a new IP Source Group
resource "zia_firewall_filtering_ip_source_groups" "example" {
  name        = "Example"
  description = "Example"
  ip_addresses = [ "192.168.100.1", "192.168.100.2", "192.168.100.3"]
}
```

```hcl
# Add an IP address range(s) to a new IP Source Group
resource "zia_firewall_filtering_ip_source_groups" "example" {
  name        = "Example"
  description = "Example"
  ip_addresses = [ "192.0.2.1-192.0.2.10" ]
}
```

```hcl
# Add subnet to a new IP Source Group
resource "zia_firewall_filtering_ip_source_groups" "example" {
  name        = "Example"
  description = "Example"
  ip_addresses = [ "203.0.113.0/24" ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Source IP group name
* `ip_addresses` (Required) Source IP addresses to be added to the group. Enter any number of IP addresses. You can enter:
  * An IP address (198.51.100.100)
  * A range of IP addresses 192.0.2.1-192.0.2.10
  * An IP address with a netmask 203.0.113.0/24

### Optional

* `description` (Optional) - Description of the source IP group

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_firewall_filtering_ip_source_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import zia_firewall_filtering_ip_source_groups.example <group_id>
```

or

```shell
terraform import zia_firewall_filtering_ip_source_groups.example <group_name>
```
