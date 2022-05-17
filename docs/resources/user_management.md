---
subcategory: "User Management"
layout: "zia"
page_title: "ZIA: user_management"
page_title: "Zscaler Internet Access (ZIA):  user_management"
sidebar_current: "docs-resource-zia-user-management"
description: |-
  Creates and manages ZIA local user accounts.

---
# zia_user_management (Resource)

The **zia_user_management** resource allows the creation and management of local user account in the Zscaler Internet Access cloud. The user account resource can then be associated with several different types of resource within the ZIA tenant.

## Example Usage

```hcl

data "zia_group_management" "normal_internet" {
 name = "Normal_Internet"
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}

# ZIA Local User Account
resource "zia_user_management" "john_ashcroft" {
 name         = "John Ashcroft"
 email        = "john.ashcroft@acme.com"
 password     = "P@ssw0rd123*"
 groups {
  id = data.zia_group_management.normal_internet.id
  }
 department {
  id = data.zia_department_management.engineering.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) User name. This appears when choosing users for policies.
* `email` - (Required) User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization.
* `password` - (Required) User's password. Applicable only when authentication type is Hosted DB. Password strength must follow what is defined in the auth settings.

* `groups` - (Required) List of Groups a user belongs to. Groups are used in policies.
  * `id` - (Required) Unique identfier for the group

* `department` - (Required) Department a user belongs to
  * `id` - (Required) Department ID

!> **WARNING:** The password parameter is considered sensitive information and is omitted in case terraform output is configured.

## Optional

The following attributes are supported:

* `comments` - (Optional) Additional information about this user.
* `temp_auth_email` - (Optional) Temporary Authentication Email. If you enabled one-time tokens or links, enter the email address to which the Zscaler service sends the tokens or links. If this is empty, the service will send the email to the User email.
