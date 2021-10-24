---
subcategory: "Group"
layout: "zia"
page_title: "ZIA: group_management"
description: |-
  Retrieve ZIA user group details.
  
---

# zia_group_management (Data Source)

The **zia_group_management** data source provides details about a specific user group created in the Zscaler Internet Access cloud or via the API. This data source can then be associated with several ZIA resources such as: URL filtering rules, Cloud Firewall rules, and locations.

## Example Usage

```hcl
# ZIA User Group Data Source
data "zia_group_management" "finance" {
 name = "Finance"
}

output "zia_group_management_finance" {
  value = data.zia_group_management.finance
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the App Connector Group to be exported.
* `id` - (Optional) The ID of this resource.

### Read-Only

* `comments` - (String)
* `deleted` - (Boolean)
* `idp_id` - (Number)
