---
subcategory: "User Management"
layout: "zscaler"
page_title: "ZIA: department_management"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-departments
  API documentation https://help.zscaler.com/zia/user-management#/departments-get
  Gets a list of user departments details.
---

# zia_department_management (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-departments)
* [API documentation](https://help.zscaler.com/zia/user-management#/departments-get)

Use the **zia_department_management** data source to get information about user department created in the Zscaler Internet Access cloud or via the API. This data source can then be associated with several ZIA resources such as: URL filtering rules, Cloud Firewall rules, and locations.

## Example Usage

```hcl
# ZIA User Department Data Source
data "zia_department_management" "engineering" {
 name = "Engineering"
}
```

```hcl
# ZIA User Department Data Source
data "zia_department_management" "finance" {
 name = "Finance"
}
```

### Example Usage - With JMESPath Search

```hcl
# Use JMESPath to pre-filter departments by name pattern
data "zia_department_management" "engineering" {
 name   = "Engineering"
 search = "[?contains(name, 'Eng')]"
}
```

```hcl
# Filter out deleted departments
data "zia_department_management" "active_dept" {
 name   = "Finance"
 search = "[?deleted == `false`]"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String) Name of the user department
* `id` - (String) ID of the user department
* `search` - (Optional) A [JMESPath](https://jmespath.org/) expression to filter results client-side after all pages have been retrieved from the API. The expression is applied to the list of departments before name or ID matching. This is useful in large environments to narrow down the candidate set. Field names in expressions must use the API's camelCase names (e.g., `name`, `idpId`, `comments`, `deleted`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `idp_id` - (Optional) Unique identfier for the identity provider (IdP)
* `comments` - (Optional) Additional information about this department
* `deleted` - (Boolean) default: false
