---
subcategory: "Admin Roles"
layout: "zia"
page_title: "ZIA: admin_roles"
description: |-
  Retrieve ZIA administrator role details.
  
---
# zia_admin_roles (Data Source)

The **zia_admin_roles** data source provides details about a specific admin role created in the Zscaler Private Access cloud or via the API. This data source can then be associated with a ZIA administrator account.

## Example Usage

```hcl
# ZIA Admin Roles Data Source
data "zia_admin_roles" "foo" {
  name = "Super Admin"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the Admin role to be exported.
* `id` - (Optional) The ID of this resource.

### Read-Only Attributes

**admin_acct_access** (String)
**analysis_access** (String)
**dashboard_access** (String)
**is_auditor** (Boolean)
**is_non_editable** (Boolean)
**logs_limit** (String)
**permissions** (List of String)
**policy_access** (String)
**rank** (Number)
**report_access** (String)
**role_type** (String)
**username_access** (String)
