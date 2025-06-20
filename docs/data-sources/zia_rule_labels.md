---
subcategory: "Rule Labels"
layout: "zscaler"
page_title: "ZIA: rule_labels"
description: |-
  Official documentation https://help.zscaler.com/zia/about-rule-labels
  API documentation https://help.zscaler.com/zia/rule-labels#/ruleLabels-get
  Get information about rule labels details.
---

# zia_rule_labels (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-rule-labels)
* [API documentation](https://help.zscaler.com/zia/rule-labels#/ruleLabels-get)

Use the **zia_rule_labels** data source to get information about a rule label resource in the Zscaler Internet Access cloud or via the API. This data source can then be associated with resources such as: Firewall Rules and URL filtering rules

## Example Usage

```hcl
# ZIA Rule Labels Data Source
data "zia_rule_labels" "example" {
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the rule label to be exported.
* `id` - (String) The unique identifer for the rule label.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) The rule label description.
* `last_modified_time` - (String) Timestamp when the rule lable was last modified. This is a read-only field. Ignored by PUT and DELETE requests.
* `last_modified_by` - (String) The admin that modified the rule label last. This is a read-only field. Ignored by PUT requests.
* `created_by` - (String) The admin that created the rule label. This is a read-only field. Ignored by PUT requests.
* `referenced_rule_count` - (int) The number of rules that reference the label.
