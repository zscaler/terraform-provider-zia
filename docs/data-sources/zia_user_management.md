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
# Look up a user by display name (exact match)
data "zia_user_management" "adam_ashcroft" {
 name = "Adam Ashcroft"
}
```

```hcl
# Look up a user by email address (exact match, case-insensitive)
data "zia_user_management" "adam_ashcroft_by_email" {
 email = "adam.ashcroft@acme.com"
}
```

```hcl
# Look up a user by numeric ID
data "zia_user_management" "adam_ashcroft_by_id" {
 id = 29309058
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

```hcl
# Combine email lookup with a department guard — fails at plan time if the
# user is moved out of Engineering, even though the email still resolves
data "zia_user_management" "engineering_only" {
 email  = "adam.ashcroft@acme.com"
 search = "[?department.name == 'Engineering']"
}
```

## Argument Reference

The following arguments are supported. Exactly one of `id`, `name`, or `email` must be provided. When more than one is provided, the lookup is performed in the order: `id`, then `email`, then `name`.

* `id` - (Optional) Numeric ID of the user. When provided, this is the most specific lookup and short-circuits any other criteria.
* `name` - (Optional) User display name. The match is **exact** and case-sensitive against the value returned by the API.
* `email` - (Optional) User email address. The match is **exact** and case-insensitive. Use this when the display name is not unique or unknown — the provider issues a partial-match query against the `/users` endpoint using the email value to narrow the candidate pool, then matches the `email` field exactly client-side.
* `search` - (Optional) A [JMESPath](https://jmespath.org/) expression to filter results client-side after all pages have been retrieved from the API. The expression is applied to the full list of users before `id`/`email`/`name` matching, so it acts as a true pre-filter against the entire population (when `search` is set, the provider deliberately bypasses the API-side `name=<lookup>` query parameter so the JMESPath can evaluate against the full set rather than a name/email-narrowed slice). Field names in expressions must use the API's camelCase names (e.g., `name`, `email`, `department.name`, `adminUser`, `type`). If the expression excludes the target user, the subsequent lookup will fail with a "user not found" error that explicitly references the `search` expression — verify the expression references valid fields (e.g. `department.name`, not `department.email`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:
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
