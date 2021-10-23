---
subcategory: "Admin Users"
layout: "zia"
page_title: "ZIA: admin_users"
description: |-
  Retrieve ZIA administrator user details.
  
---
# zia_admin_users (Data Source)

The **zia_admin_users** data source provides details about a specific admin user account created in the Zscaler Private Access cloud or via the API. This data source can then be associated with a ZIA administrator role.

## Example Usage

```hcl
# ZIA Admin User Data Source
data "zia_admin_users" "john_doe" {
  login_name = "john.doe@example.com"
}

output "zia_admin_users_john_doe" {
    value = data.zia_admin_users.john_doe
}
```

## Argument Reference

The following arguments are supported:

* `login_name` - (Required) The email address of the admin user to be exported.
* `id` - (Optional) The ID of this resource.

### Read-Only

- **admin_scope** (Set of Object) (see [below for nested schema](#nestedatt--admin_scope))
- **comments** (String)
- **disabled** (Boolean)
- **email** (String)
- **exec_mobile_app_tokens** (List of Object) (see [below for nested schema](#nestedatt--exec_mobile_app_tokens))
- **is_auditor** (Boolean)
- **is_exec_mobile_app_enabled** (Boolean)
- **is_non_editable** (Boolean)
- **is_password_expired** (Boolean)
- **is_password_login_allowed** (Boolean)
- **is_product_update_comm_enabled** (Boolean)
- **is_security_report_comm_enabled** (Boolean)
- **is_service_update_comm_enabled** (Boolean)
- **pwd_last_modified_time** (Number)
- **role** (Set of Object) (see [below for nested schema](#nestedatt--role))
- **user_name** (String)

<a id="nestedatt--admin_scope"></a>
### Nested Schema for `admin_scope`

Read-Only:

- **scope_entities** (List of Object) (see [below for nested schema](#nestedobjatt--admin_scope--scope_entities))
- **scope_group_member_entities** (List of Object) (see [below for nested schema](#nestedobjatt--admin_scope--scope_group_member_entities))
- **type** (String)

<a id="nestedobjatt--admin_scope--scope_entities"></a>
### Nested Schema for `admin_scope.scope_entities`

Read-Only:

- **id** (Number)
- **name** (String)


<a id="nestedobjatt--admin_scope--scope_group_member_entities"></a>
### Nested Schema for `admin_scope.scope_group_member_entities`

Read-Only:

- **id** (Number)
- **name** (String)



<a id="nestedatt--exec_mobile_app_tokens"></a>
### Nested Schema for `exec_mobile_app_tokens`

Read-Only:

- **cloud** (String)
- **create_time** (Number)
- **device_id** (String)
- **device_name** (String)
- **name** (String)
- **org_id** (Number)
- **token** (String)
- **token_expiry** (Number)
- **token_id** (String)


<a id="nestedatt--role"></a>
### Nested Schema for `role`

Read-Only:

- **extensions** (Map of String)
- **id** (Number)
- **name** (String)


