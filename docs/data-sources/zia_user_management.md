---
subcategory: "User Management"
layout: "zscaler"
page_title: "ZIA: user_management"
description: |-
    Official documentation https://help.zscaler.com/zia/about-url-filteringhttps://help.zscaler.com/zia/user-management#/users-get
    API documentation https://help.zscaler.com/zia/user-management#/users-get
    Gets a list of all users and allows user filtering by name, department, or group
---

# zia_user_management (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-url-filteringhttps://help.zscaler.com/zia/user-management#/users-get)
* [API documentation](https://help.zscaler.com/zia/about-url-filteringhttps://help.zscaler.com/zia/user-management#/users-get)

Use the **zia_user_management** data source to get information about a user account that may have been created in the Zscaler Internet Access portal or via API. This data source can then be associated with a ZIA cloud firewall filtering rule, and URL filtering rules.

## Example Usage

```hcl
# ZIA Local User Account
data "zia_user_management" "adam_ashcroft" {
 name = "Adam Ashcroft"
}
```

### Example Usage - With JMESPath Search

```hcl
# Use JMESPath to pre-filter users by department before matching by name
data "zia_user_management" "adam_ashcroft" {
 name   = "Adam Ashcroft"
 search = "[?department.name == 'Engineering']"
}
```

```hcl
# Filter to admin users only
data "zia_user_management" "admin_user" {
 name   = "Jane Smith"
 search = "[?adminUser == `true`]"
}
```

```hcl
# Filter users whose name contains a specific string
data "zia_user_management" "user" {
 name   = "Adam Ashcroft"
 search = "[?contains(name, 'Adam')]"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) User name. This appears when choosing users for policies.
* `id` - (Optional) The ID of the time window resource.
* `search` - (Optional) A [JMESPath](https://jmespath.org/) expression to filter results client-side after all pages have been retrieved from the API. The expression is applied to the list of users before name or ID matching. This is useful in large environments to narrow down the candidate set. Field names in expressions must use the API's camelCase names (e.g., `name`, `email`, `department`, `adminUser`, `type`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `email` - (Required) User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization
* `admin_user` - (String) True if this user is an Admin user. readOnly: `true` default: `false`
* `comments` - (String) Additional information about this user.
* `password` -(String, Sensitive)
* `temp_auth_email` - (String) Temporary Authentication Email. If you enabled one-time tokens or links, enter the email address to which the Zscaler service sends the tokens or links. If this is empty, the service will send the email to the User email.
* `auth_methods` - (String) Type of authentication method to be enabled. Supported values are: ``BASIC`` and ``DIGEST``
* `type` - (String) User type. Provided only if this user is not an end user. The supported types are:
  * `SUPERADMIN`
  * `ADMIN`
  * `AUDITOR`
  * `GUEST`
  * `REPORT_USER`
  * `UNAUTH_TRAFFIC_DEFAULT`

* `department` - (String) Department a user belongs to
  * `id` - (Number) Department ID
  * `name` - (String) Department name
  * `idp_id` - (Number) Identity provider (IdP) ID
  * `comments` - (String) Additional information about this department
  * `deleted` - (Boolean) default: `false`

* `groups` - (String) List of Groups a user belongs to. Groups are used in policies.
  * `id` - (Number) Unique identfier for the group
  * `name` - (String) Group name
  * `idp_id` - (Number) Unique identfier for the identity provider (IdP)
  * `comments` - (String) Additional information about the group
