---
subcategory: "Firewall Filtering Network Application Groups"
layout: "zia"
page_title: "ZIA: firewall_filtering_network_application_groups"
description: |-
  Retrieve ZIA firewall rule network application groups.
  
---

# zia_firewall_filtering_network_application_groups (Data Source)

The **zia_firewall_filtering_network_application_groups** data source provides details about a specific network application group available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network application rule.

## Example Usage

```hcl
# ZIA Network Application Groups
data "zia_firewall_filtering_network_application_groups" "example" {
    name = "example"
}

output "zia_firewall_filtering_ip_source_groups_example" {
    value = data.zia_firewall_filtering_network_application_groups.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.

### Read-Only

* `description` - (String)
* `network_applications` - (List of String)
* `id` - The ID of this resource.
