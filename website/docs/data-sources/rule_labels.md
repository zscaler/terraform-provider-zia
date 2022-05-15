---
subcategory: "Rule Labels"
layout: "zia"
page_title: "ZIA: rule labels"
description: |-
  Retrieve ZIA rule labels details.
  
---
# zia_rule_labels (Data Source)

The **zia_rule_labels** - data source provides details about a specific rule label resource in the Zscaler Internet Access cloud or via the API. This data source can then be associated with resources such as: Firewall Rules and URL filtering rules

## Example Usage

```hcl
# ZIA Rule Labels Data Source
data "zia_rule_labels" "example" {
    name = "Example"
}

output "zia_rule_labels" {
    value = data.zia_rule_labels.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the rule label to be exported.

### Read-Only

* `id` - (String) The unique identifer for the device group.
* `description` - (String) The rule label description.
* `last_modified_time` - (String) Timestamp when the rule lable was last modified. This is a read-only field. Ignored by PUT and DELETE requests.
* `last_modified_by` - (String) The admin that modified the rule label last. This is a read-only field. Ignored by PUT requests.
* `created_by` - (String) The admin that created the rule label. This is a read-only field. Ignored by PUT requests.
* `referenced_rule_count` - (int) The number of rules that reference the label.
