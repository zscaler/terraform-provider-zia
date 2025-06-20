---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_destination_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-post
  API documentation https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-post
  Creates and manages ZIA Cloud firewall IP destination groups.
---

# zia_firewall_filtering_destination_groups (Resource)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-post)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-post)

The **zia_firewall_filtering_destination_groups** resource allows the creation and management of ZIA Cloud Firewall IP destination groups in the Zscaler Internet Access. This resource can then be associated with a ZIA cloud firewall filtering rule.

## Example Usage

```hcl
# IP Destination Group of Type DSTN_FQDN
resource "zia_firewall_filtering_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
}
```

```hcl
# IP Destination Group of Type DSTN_IP
resource "zia_firewall_filtering_destination_groups" "example_dstn_ip" {
  name        = "Example Destination IP"
  description = "Example Destination IP"
  type        = "DSTN_IP"
  addresses = ["3.217.228.0-3.217.231.255",
    "3.235.112.0-3.235.119.255",
    "52.23.61.0-52.23.62.25",
    "35.80.88.0-35.80.95.255"]
}
```

```hcl
# IP Destination Group of Type DSTN_DOMAIN
resource "zia_firewall_filtering_destination_groups" "example_dstn_domain" {
  name          = "Example Destination Domain"
  description   = "Example Destination Domain"
  type          = "DSTN_DOMAIN"
  addresses     = ["acme.com", "acme1.com"]
}
```

```hcl
# IP Destination Group of Type DSTN_OTHER
resource "zia_firewall_filtering_destination_groups" "example_dstn_other" {
  name          = "Example Destination Other"
  description   = "Example Destination Other"
  type          = "DSTN_OTHER"
  countries     = ["COUNTRY_CA"]
  ip_categories = ["CUSTOM_01", "CUSTOM_02"]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Destination IP group name
* `type` - (Required) Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs). The supported values are:
  * `DSTN_IP`
  * `DSTN_FQDN`
  * `DSTN_DOMAIN`
  * `DSTN_OTHER`
* `addresses` (Required) Destination IP addresses, domains or FQDNs within the group

### Optional

* `countries` (Optional) Destination IP address counties. You can identify destinations based on the location of a server.
* `description` (Optional) Additional information about the destination IP group
* `ip_categories` (Optional) Destination IP address URL categories. You can identify destinations based on the URL category of the domain. See list of all IP Categories [Here](https://help.zscaler.com/zia/firewall-policies#/ipDestinationGroups-get)
  * !> **WARNING:** The `ip_categories` attribute only accepts custom URL categories.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_firewall_filtering_destination_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import zia_firewall_filtering_destination_groups.example <group_id>
```

or

```shell
terraform import zia_firewall_filtering_destination_groups.example <group_name>
```
