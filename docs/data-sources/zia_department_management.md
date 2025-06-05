---
subcategory: "User Management"
layout: "zscaler"
page_title: "ZIA: department_management"
description: |-
  Gets a list of user departments details.

---
# Data Source: zia_department_management

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

## Argument Reference

The following arguments are supported:

* `name` - (String) Name of the user department
* `id` - (String) ID of the user department

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `idp_id` - (Optional) Unique identfier for the identity provider (IdP)
* `comments` - (Optional) Additional information about this department
* `deleted` - (Boolean) default: false
