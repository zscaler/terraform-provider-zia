---
subcategory: "Admin Roles"
layout: "zscaler"
page_title: "ZIA: admin_roles"
description: |-
  Official documentation https://help.zscaler.com/zia/about-role-management
  API documentation https://help.zscaler.com/zia/admin-role-management#/adminRoles-get
  Creates and manages ZIA admin roles.
---

# zia_admin_roles (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-role-management)
* [API documentation](https://help.zscaler.com/zia/admin-role-management#/adminRoles-get)

The **zia_admin_roles** resource allows the creation and management of admin roles in the Zscaler Internet Access cloud or via the API.

## Example Usage - Create Admin Role

```hcl
# ZIA Admin Roles Resource
resource "zia_admin_roles" "this" {
  name               = "AdminRoleTerraform"
  rank               = 7
  alerting_access    = "READ_WRITE"
  dashboard_access   = "READ_WRITE"
  report_access      = "READ_WRITE"
  analysis_access    = "READ_ONLY"
  username_access    = "READ_ONLY"
  device_info_access = "READ_ONLY"
  admin_acct_access  = "READ_WRITE"
  policy_access      = "READ_WRITE"
  permissions = [
 "NSS_CONFIGURATION", "LOCATIONS", "HOSTED_PAC_FILES", "EZ_AGENT_CONFIGURATIONS",
 "SECURE_AGENT_NOTIFICATIONS", "VPN_CREDENTIALS", "AUTHENTICATION_SETTINGS", "STATIC_IPS",
 "GRE_TUNNELS", "CLIENT_CONNECTOR_PORTAL", "SECURE", "POLICY_RESOURCE_MANAGEMENT",
 "CUSTOM_URL_CAT", "OVERRIDE_EXISTING_CAT", "TENANT_PROFILE_MANAGEMENT", "COMPLY",
 "SSL_POLICY", "ADVANCED_SETTINGS", "PROXY_GATEWAY", "SUBCLOUDS", "IDENTITY_PROXY_SETTINGS",
 "USER_MANAGEMENT", "APIKEY_MANAGEMENT", "FIREWALL_DNS", "VZEN_CONFIGURATION",
 "PARTNER_INTEGRATION", "USER_ACCESS", "CUSTOMER_ACCT_INFO", "CUSTOMER_SUBSCRIPTION",
 "CUSTOMER_ORG_SETTINGS", "ZIA_TRAFFIC_CAPTURE", "REMOTE_ASSISTANCE_MANAGEMENT"
  ]
}
```

## Example Usage - Create Admin SDWAN Role

```hcl
resource "zia_admin_roles" "this" {
  name               = "SDWANAdminRoleTerraform"
  rank               = 7
  policy_access      = "READ_WRITE"
  alerting_access    = "NONE"
  role_type          = "SDWAN"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name of the admin role

### Optional

* `rank` - (String) Admin rank of this admin role. This is applicable only when admin rank is enabled in the advanced settings. Default value is 7 (the lowest rank). The assigned admin rank determines the roles or admin users this user can manage, and which rule orders this admin can access
* `policy_access` - (String) Policy access permission. Supported Values: `READ_WRITE` and `READ_ONLY`
* `alerting_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `dashboard_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `report_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `analysis_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `username_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `device_info_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `admin_acct_access` - (String) The rule label description. Supported Values: `READ_WRITE` and `READ_ONLY`
* `is_auditor` - (String) The rule label description.
* `feature_permissions` - (map) The rule label description.
* `feature_ext_feature_permissionspermissions` - (String) Supported Values: `"INCIDENT_WORKFLOW": "NONE"`, `"INCIDENT_WORKFLOW": "RESTRICTED"`, `"INCIDENT_WORKFLOW": "NONE"`

* `is_non_editable` - (bool) Indicates whether or not this admin user is editable
* `logs_limit` - (String) Log range limit. Supported Values: `UNRESTRICTED`, `MONTH_1`, `MONTH_2`, `MONTH_3`, `MONTH_4`, `MONTH_5`, `MONTH_6`
* `role_type` - (String) The admin role type. This attribute is subject to change. Supported value: `"ORG_ADMIN|EXEC_INSIGHT|EXEC_INSIGHT_AND_ORG_ADMIN|SDWAN"`
* `report_time_duration` - (String) Time duration allocated to the report dashboard. The default value of -1 indicates that no time restriction is applied to the report dashboard. Time Unit is in hours.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_rule_labels** can be imported by using `<LABEL_ID>` or `<LABEL_NAME>` as the import ID.

For example:

```shell
terraform import zia_rule_labels.example <label_id>
```

or

```shell
terraform import zia_rule_labels.example <label_name>
```
