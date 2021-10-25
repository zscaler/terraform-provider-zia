---
subcategory: "Firewall Filtering IP Source Groups"
layout: "zia"
page_title: "ZIA: firewall_filtering_ip_source_groups"
description: |-
  Retrieve ZIA firewall rule IP source groups.
  
---

# zia_firewall_filtering_ip_source_groups (Data Source)

The **zia_firewall_filtering_ip_source_groups** data source provides details about a specific ip source groups available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

## Example Usage

```hcl
# ZIA IP Source Groups
data "zia_firewall_filtering_ip_source_groups" "example" {
    name = "example"
}

output "zia_firewall_filtering_ip_source_groups_example" {
    value = data.zia_firewall_filtering_ip_source_groups.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.

### Read-Only

* `id` - The ID of this resource.
* `description` - (String)
* `ip_addresses` - (List of String)
