---
subcategory: "Rule Labels"
layout: "zia"
page_title: "ZIA: rule labels"
description: |-
  Creates ZIA rule labels.
  
---
# zia_rule_labels (Resource)

The **zia_rule_labels** - creates a rule label resource in the Zscaler Internet Access cloud or via the API. This resource can then be associated with resources such as: Firewall Rules and URL filtering rules

## Example Usage

```hcl
# ZIA Rule Labels Resource
resource "zia_rule_labels" "example" {
    name = "Example"
    description = "Example"
}

output "zia_rule_labels" {
    value = data.zia_rule_labels.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the devices to be created.
* `description` - (String) The rule label description.

### Read-Only

* `id` - (String) The unique identifer for the device group.
* `last_modified_time` - (String) Timestamp when the rule lable was last modified. This is a read-only field. Ignored by PUT and DELETE requests.
* `last_modified_by` - (String) The admin that modified the rule label last. This is a read-only field. Ignored by PUT requests.
* `created_by` - (String) The admin that created the rule label. This is a read-only field. Ignored by PUT requests.
* `referenced_rule_count` - (int) The number of rules that reference the label.
