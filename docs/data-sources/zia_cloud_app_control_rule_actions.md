---
subcategory: "Cloud App Control Policy"
layout: "zscaler"
page_title: "ZIA: cloud_app_control_rule_actions"
description: |-
  Official documentation https://help.zscaler.com/zia/adding-rules-cloud-app-control-policy
  API documentation https://help.zscaler.com/zia/cloud-app-control-policy#/webApplicationRules/{rule_type}-get
  Get information about ZIA Cloud App Control Rules.
---

# zia_cloud_app_control_rule_actions (Data Source)

* [Official documentation](https://help.zscaler.com/zia/adding-rules-cloud-app-control-policy)
* [API documentation](https://help.zscaler.com/zia/cloud-app-control-policy#/webApplicationRules/)

Use the **zia_cloud_app_control_rule_actions** data source to fetch the granular actions supported for the applications.

## Example Usage

```hcl

data "zia_cloud_app_control_rule_actions" "dropbox_actions" {
  type        = "ENTERPRISE_COLLABORATION"
  cloud_apps  = ["SLACK"]
}
```

## Argument Reference

The following arguments are supported:

* `type` - (String) The rule type selected from the available options
* `cloud_apps` - (String) The cloud application name. To retrieve the available list of DNS tunnelapplications use the data source: `zia_cloud_applications`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

N/A
