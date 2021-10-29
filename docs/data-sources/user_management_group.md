---
subcategory: "User Management Group"
layout: "zia"
page_title: "ZIA: group_management"
description: |-
  Retrieve ZIA user group.
  
---

# zia_group_management (Data Source)

The **zia_group_management** -- data source provides details about a specific user group that may have been created in the Zscaler Internet Access portal. This data source can then be associated with a ZIA cloud firewall filtering rule, and URL filtering rules.

## Example Usage

```hcl

data "zia_group_management" "devops" {
 name = "DevOps"
}

output "zia_group_management" {
  value = data.zia_group_management.devops
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) User name. This appears when choosing users for policies.
* `id` - (Optional) Unique identfier for the group.
* `idp_id` - (Optional) Unique identfier for the identity provider (IdP)
* `comments` - (Optional) Additional information about the group
