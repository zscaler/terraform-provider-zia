---
subcategory: "Admin Users"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): admin_users"
sidebar_current: "docs-resource-zia-admin-users"
description: |-
  Creates and manages ZIA administrator users.
---

# Resource: zia_admin_users

The **zia_admin_users** resource allows the creation and management of ZIA admin user account created in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_department_management" "engineering" {
  name = "Engineering"
}

resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  user_name                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "AeQ9E5w8B$"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator User"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope {
    type = "DEPARTMENT"
    scope_entities {
      id = [data.zia_department_management.engineering.id]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `login_name` - (Required) The email address of the admin user to be exported.
* `username` - (Required) The username of the admin user to be exported.
* `password` - (Required) The username of the admin user to be exported.
* `role` - (Required) Role of the admin. This is not required for an auditor.
  * `id` - (Required) Identifier that uniquely identifies an entity

!> **WARNING:** The password parameter is considered sensitive information and is omitted in case terraform output is configured.

### Optional

* `email` - (Optional) Admin or auditor's email address.
* `comments` - (Optional) Additional information about the admin or auditor.
* `disabled` - (Optional) Indicates whether or not the admin account is disabled.
* `is_auditor` - (Optional) Indicates whether the user is an auditor. This attribute is subject to change.
* `is_exec_mobile_app_enabled` - (Optional) Indicates whether or not Executive Insights App access is enabled for the admin.
* `is_non_editable` - (Optional) Indicates whether or not the admin can be edited or deleted.
* `is_password_expired` - (Optional) Indicates whether or not an admin's password has expired.
* `is_password_login_allowed` - (Optional) The default is true when SAML Authentication is disabled. When SAML Authentication is enabled, this can be set to false in order to force the admin to login via SSO only.
* `is_product_update_comm_enabled` - (Optional) Communication setting for Product Update.
* `is_security_report_comm_enabled` - (Optional) Communication for Security Report is enabled.
* `is_service_update_comm_enabled` - (Optional) Communication setting for Service Update.

* `admin_scope` - (Optional) The admin's scope. A scope is required for admins, but not applicable to auditors. This attribute is subject to change.
  * `scope_group_member_entities` - (Optional) Only applicable for the LOCATION_GROUP admin scope type, in which case this attribute gives the list of ID/name pairs of locations within the location group.
    * `id` - (Optional) Identifier that uniquely identifies an entity
    * `name` - (Optional) The configured name of the entity
  * `type` - (Optional) The admin scope type. The attribute name is subject to change.
  * `scope_entities` - (Optional) Based on the admin scope type, the entities can be the ID/name pair of departments, locations, or location groups.
    * `id` - (Optional) Identifier that uniquely identifies an entity
    * `name` - (Optional) The configured name of the entity
