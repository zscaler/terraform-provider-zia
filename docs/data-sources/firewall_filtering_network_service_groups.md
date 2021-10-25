---
subcategory: "Firewall Filtering Network Service Groups"
layout: "zia"
page_title: "ZIA: firewall_filtering_network_service_groups"
description: |-
  Retrieve ZIA firewall rule network service groups.
  
---

# zia_firewall_filtering_network_service_groups (Data Source)

The **zia_firewall_filtering_network_service_groups** data source provides details about a specific network service groups available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering network service rule.

## Example Usage

```hcl
# ZIA Network Service Groups
data "zia_firewall_filtering_network_service_groups" "example"{
    name = "Corporate Custom SSH TCP_10022"
}

output "zia_firewall_filtering_network_service_groups" {
  value = data.zia_firewall_filtering_network_service_groups.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ip source group to be exported.
* `services` - (List of service IDs) (see [below for nested schema](#nestedatt--services))

### Read-Only

* `description` - (String)
* `id` - (Number) The ID of this resource.

<a id="nestedatt--services"></a>

### Nested Schema for `services`

Read-Only:

* `id` - (List of service IDs)
