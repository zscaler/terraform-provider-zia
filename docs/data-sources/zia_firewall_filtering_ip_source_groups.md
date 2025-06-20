---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA: firewall_filtering_destination_groups"
description: |-
  Official documentation https://help.zscaler.com/zia/firewall-policies#/ipSourceGroups-get
  API documentation https://help.zscaler.com/zia/firewall-policies#/ipSourceGroups-get
  Get information about IP Source groups.
---


# zia_firewall_filtering_ip_source_groups (Data Source)

* [Official documentation](https://help.zscaler.com/zia/firewall-policies#/ipSourceGroups-get)
* [API documentation](https://help.zscaler.com/zia/firewall-policies#/ipSourceGroups-get)

Use the **zia_firewall_filtering_ip_source_groups** data source to get information about ip source groups available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

## Example Usage

```hcl
# ZIA IP Source Groups
data "zia_firewall_filtering_ip_source_groups" "example" {
    name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.
* `id` - (Optional) The ID of the ip source group resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of this resource.
* `description` - (String)
* `ip_addresses` - (List of String)
