---
subcategory: "User Management"
layout: "zscaler"
page_title: "ZIA: group_management"
description: |-
  Official documentation https://help.zscaler.com/zia/user-management#/groups-get
  API documentation https://help.zscaler.com/zia/user-management#/groups-get
  Gets a list of user user group.
---

# zia_group_management (Data Source)

* [Official documentation](https://help.zscaler.com/zia/user-management#/groups-get)
* [API documentation](https://help.zscaler.com/zia/user-management#/groups-get)

Use the **zia_group_management** data source to get information about a user group that may have been created in the Zscaler Internet Access portal. This data source can then be associated with a ZIA cloud firewall filtering rule, and URL filtering rules.

## Example Usage

```hcl
data "zia_group_management" "devops" {
 name = "DevOps"
}
```

### Example Usage - With JMESPath Search

```hcl
# Use JMESPath to pre-filter groups before matching by name
data "zia_group_management" "devops" {
 name   = "DevOps"
 search = "[?contains(name, 'Dev')]"
}
```

```hcl
# Filter groups by IdP ID
data "zia_group_management" "idp_group" {
 name   = "Engineering"
 search = "[?idpId == `0`]"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the user group
* `id` - (Optional) Unique identfier for the group.
* `search` - (Optional) A [JMESPath](https://jmespath.org/) expression to filter results client-side after all pages have been retrieved from the API. The expression is applied to the list of groups before name or ID matching. This is useful in large environments to narrow down the candidate set. Field names in expressions must use the API's camelCase names (e.g., `name`, `idpId`, `comments`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `idp_id` - (Optional) Unique identfier for the identity provider (IdP)
* `comments` - (Optional) Additional information about the group
