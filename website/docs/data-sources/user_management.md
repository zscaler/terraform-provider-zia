---
subcategory: "User Management"
layout: "zia"
page_title: "ZIA: user_management"
description: |-
  Retrieve ZIA user account.
  
---
# zia_user_management (Data Source)

The **zia_user_management` -data source provides details about a specific user account that may have been created in the Zscaler Internet Access portal or via API. This data source can then be associated with a ZIA cloud firewall filtering rule, and URL filtering rules.

## Example Usage

```hcl
# ZIA Local User Account
data "zia_user_management" "adam_ashcroft" {
 name = "Adam Ashcroft"
}

output "zia_user_management" {
  value = data.zia_user_management.adam_ashcroft
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) User name. This appears when choosing users for policies.
* `id` - (Optional) The ID of the time window resource.

## Attribute Reference

The following attributes are supported:

* `email` - (Required) User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization
* `admin_user` - (Optional) True if this user is an Admin user. readOnly: `true` default: `false`
* `comments` - (Optional) Additional information about this user.
* `password` -(String, Sensitive)
* `temp_auth_email` - (String) Temporary Authentication Email. If you enabled one-time tokens or links, enter the email address to which the Zscaler service sends the tokens or links. If this is empty, the service will send the email to the User email.
* `type` - (String) User type. Provided only if this user is not an end user.

`department` - (Required) Department a user belongs to

* `id` - (Number) Department ID
* `name` - (String) Department name
* `idp_id` - (Number) Identity provider (IdP) ID
* `comments` - (String) Additional information about this department
* `deleted` - (Boolean) default: `false`

`groups` - (Required) List of Groups a user belongs to. Groups are used in policies.

* `id` - (Number) Unique identfier for the group
* `name` - (String) Group name
* `idp_id` - (Number) Unique identfier for the identity provider (IdP)
* `comments` - (String) Additional information about the group
