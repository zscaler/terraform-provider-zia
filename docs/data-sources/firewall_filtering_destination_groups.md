---
subcategory: "Firewall Filtering Destination Groups"
layout: "zia"
page_title: "ZIA: firewall_filtering_destination_groups"
description: |-
  Retrieve ZIA firewall rule destination groups.
  
---

# zia_firewall_filtering_destination_groups (Data Source)

The **zia_firewall_filtering_destination_groups** data source provides details about a specific destination groups option available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

## Example Usage

```hcl
# ZIA Destination Groups
data "zia_firewall_filtering_destination_groups" "example" {
    name = "example"
}

output "zia_firewall_filtering_destination_groups_example" {
    value = zia_firewall_filtering_destination_groups.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the time window to be exported.
* `id` - (Optional) The ID of the time window resource.

### Read-Only

* `addresses` - (List of String)
* `countries` - (List of String)
* `description` - (String)
* `ip_categories` - (List of String)
* `type` - (String)
