---
subcategory: "Admin & Role Management"
layout: "zscaler"
page_title: "ZIA: admin_roles"
description: |-
  Get information about ZIA administrator roles.
---

# Data Source: zia_admin_roles

Use the **zia_admin_roles** data source to get information about an admin role created in the Zscaler Internet Access cloud or via the API. This data source can then be associated with a ZIA administrator account.

## Example Usage

```hcl
# ZIA Admin Roles Data Source
data "zia_auth_settings_urls" "example" {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Admin role to be exported.
* `id` - (Optional) The ID of the admin role to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `admin_acct_access` (String)
* `analysis_access` (String)
* `dashboard_access` (String) Dashboard access permission. Supported values are: `NONE`, `READ_ONLY`
* `is_auditor` (Boolean) Indicates whether this is an auditor role.
* `is_non_editable` (Boolean) Indicates whether or not this admin user is editable/deletable.
* `logs_limit` (String) Log range limit. Returned values are: `UNRESTRICTED`, `MONTH_1`, `MONTH_2`, `MONTH_3`, `MONTH_4`, `MONTH_5`, `MONTH_6`
* `permissions` (List of String) List of functional areas to which this role has access. This attribute is subject to change.
* `policy_access` (String) Policy access permission. Returned values are: `NONE`, `READ_ONLY`,`READ_WRITE`
* `rank` (Number) Admin rank of this admin role. This is applicable only when admin rank is enabled in the advanced settings. Default value is 7 (the lowest rank). The assigned admin rank determines the roles or admin users this user can manage, and which rule orders this admin can access.
* `report_access` (String) Report access permission. Returned values are: `NONE`, `READ_ONLY`,`READ_WRITE`
* `role_type` (String) The admin role type. ()This attribute is subject to change.) Supported values are:  `ORG_ADMIN`, `EXEC_INSIGHT`, `EXEC_INSIGHT_AND_ORG_ADMIN`, `SDWAN`
* `username_access` (String) Username access permission. When set to NONE, the username will be obfuscated. Supported values are: `NONE|READ_ONLY`
