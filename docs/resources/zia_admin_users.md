---
subcategory: "Admin & Role Management"
layout: "zscaler"
page_title: "ZIA: admin_users"
description: |-
  Official documentation https://help.zscaler.com/zia/about-administrators
  API documentation https://help.zscaler.com/zia/admin-role-management#/adminUsers-get
  Creates and manages ZIA administrator users.
---

# zia_admin_users (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-administrators)
* [API documentation](https://help.zscaler.com/zia/admin-role-management#/adminUsers-get)

The **zia_admin_users** resource allows the creation and management of ZIA admin user account created in the Zscaler Internet Access cloud or via the API.

## Example Usage - Organization Scope

```hcl
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  user_name                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "*********************"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator User"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "ORGANIZATION"
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}
```

## Example Usage - Department Scope

```hcl
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  user_name                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "*********************"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator User"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "DEPARTMENT"
    admin_scope_entities {
        id = [ data.zia_department_management.engineering.id, data.zia_department_management.sales.id ]
    }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_department_management" "engineering" {
  name = "Engineering"
}
```

## Example Usage - Location Scope

```hcl
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  user_name                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
  password                        = "*********************"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator User"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "LOCATION"
    admin_scope_entities {
        id = [ data.zia_location_management.au_sydney_branch01.id ]
    }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_location_management" "au_sydney_branch01" {
  name = "AU - Sydney - Branch01"
}
```

## Example Usage - Location Group Scope

```hcl
######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_admin_users" "john_smith" {
  login_name                      = "john.smith@acme.com"
  user_name                       = "John Smith"
  email                           = "john.smith@acme.com"
  is_password_login_allowed       = true
 password                         = "*********************"
  is_security_report_comm_enabled = true
  is_service_update_comm_enabled  = true
  is_product_update_comm_enabled  = true
  comments                        = "Administrator User"
  role {
    id = data.zia_admin_roles.super_admin.id
  }
  admin_scope_type = "LOCATION_GROUP"
    admin_scope_entities {
        id = [ data.zia_location_groups.corporate_user_traffic_group.id ]
    }
}

data "zia_admin_roles" "super_admin" {
  name = "Super Admin"
}

data "zia_location_groups" "corporate_user_traffic_group" {
  name = "Corporate User Traffic Group"
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

* `admin_scope_type` - (Optional) The admin's scope. A scope is required for admins, but not applicable to auditors. This attribute is subject to change. Support values are: `ORGANIZATION`, `DEPARTMENT`, `LOCATION`, `LOCATION_GROUP`
  * `admin_scope_entities` - (Optional) Based on the admin scope type, the entities can be the ID/name pair of departments, locations, or location groups.
    * `id` - (Optional) Identifier that uniquely identifies an entity

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_admin_users** can be imported by using `<ADMIN ID>` or `<LOGIN NAME>` as the import ID.

For example:

```shell
terraform import zia_admin_users.example <admin_id>
```

or

```shell
terraform import zia_admin_users.example <login_name>
```

⚠️ **NOTE :**:  This provider do not import the password attribute value during the importing process.
