---
subcategory: "Admin Users"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): admin_users"
sidebar_current: "docs-datasource-zia-admin-users"
description: |-
  Get information about ZIA administrator users.
---

# Data Source: zia_admin_users

Use the **zia_admin_users** data source to get information about an admin user account created in the Zscaler Internet Access cloud or via the API. This data source can then be associated with a ZIA administrator role.

## Example Usage

```hcl
# ZIA Admin User Data Source by login_name
data "zia_admin_users" "john_doe" {
  login_name = "john.doe@example.com"
}
```

```hcl
# ZIA Admin User Data Source by username
data "zia_admin_users" "john_doe" {
  username = "John Doe"
}
```

## Argument Reference

The following arguments are supported:

* `login_name` - (Required) The email address of the admin user to be exported.
* `username` - (Required) The username of the admin user to be exported.
* `id` - (Optional) The ID of the admin user to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `email` - (String) Admin or auditor's email address.
* `comments` - (String) Additional information about the admin or auditor.
* `disabled` - (Boolean) Indicates whether or not the admin account is disabled.
* `is_auditor` - (Boolean) Indicates whether the user is an auditor. This attribute is subject to change.
* `is_exec_mobile_app_enabled` - (Boolean) Indicates whether or not Executive Insights App access is enabled for the admin.
* `is_non_editable` - (Boolean) Indicates whether or not the admin can be edited or deleted.
* `is_password_expired` - (Boolean) Indicates whether or not an admin's password has expired.
* `is_password_login_allowed` - (Boolean) The default is true when SAML Authentication is disabled. When SAML Authentication is enabled, this can be set to false in order to force the admin to login via SSO only.
* `is_product_update_comm_enabled` - (Boolean) Communication setting for Product Update.
* `is_security_report_comm_enabled` - (Boolean) Communication for Security Report is enabled.
* `is_service_update_comm_enabled` - (Boolean) Communication setting for Service Update.

* `role` - (Set of Object) Role of the admin. This is not required for an auditor.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity

* `admin_scope` - (Set of Object) The admin's scope. Only applicable for the LOCATION_GROUP admin scope type, in which case this attribute gives the list of ID/name pairs of locations within the location group.
  * `scope_group_member_entities` - (Number) Only applicable for the LOCATION_GROUP admin scope type, in which case this attribute gives the list of ID/name pairs of locations within the location group.
    * `id` - (Number) Identifier that uniquely identifies an entity
    * `name` - (String) The configured name of the entity
  * `type` - (String) The admin scope type. The attribute name is subject to change.
  * `scope_entities` - (String) Based on the admin scope type, the entities can be the ID/name pair of departments, locations, or location groups.
    * `id` - (Number) Identifier that uniquely identifies an entity
    * `name` - (String) The configured name of the entity

* `exec_mobile_app_tokens` - (List of Object)
  * `cloud` - (String)
  * `org_id` - (Number)
  * `name` - (String)
  * `token_id` - (String)
  * `token` - (String)
  * `token_expiry` - (Number)
  * `create_time` - (Number)
  * `device_id` - (String)
  * `device_name` - (String)
