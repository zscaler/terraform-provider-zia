---
subcategory: "User Management"
layout: "zscaler"
page_title: "ZIA: group_management"
description: |-
  Gets a list of user user group.
---

# zia_group_management (Data Source)

Use the **zia_group_management** data source to get information about a user group that may have been created in the Zscaler Internet Access portal. This data source can then be associated with a ZIA cloud firewall filtering rule, and URL filtering rules.

## Example Usage

```hcl
data "zia_group_management" "devops" {
 name = "DevOps"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the user group
* `id` - (Optional) Unique identfier for the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `idp_id` - (Optional) Unique identfier for the identity provider (IdP)
* `comments` - (Optional) Additional information about the group
