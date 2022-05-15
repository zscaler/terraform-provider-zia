---
subcategory: "Department"
layout: "zia"
page_title: "ZIA: department_management"
description: |-
  Retrieve ZIA user department details.
  
---
# zia_department_management (Data Source)

The **zia_department_management** -data source provides details about a specific user department created in the Zscaler Internet Access cloud or via the API. This data source can then be associated with several ZIA resources such as: URL filtering rules, Cloud Firewall rules, and locations.

## Example Usage

```hcl
# ZIA User Department Data Source
data "zia_department_management" "engineering" {
 name = "Engineering"
}

output "zia_department_management_engineering" {
  value = data.zia_department_management.engineering
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the App Connector Group to be exported.
* `id` - (Optional) The ID of this resource.
* `idp_id` - (Optional) Unique identfier for the identity provider (IdP)
* `comments` - (Optional) Additional information about the group
* `deleted` - (Boolean)
