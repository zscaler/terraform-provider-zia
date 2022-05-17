---
subcategory: "Rule Labels"
layout: "zia"
page_title: "Zscaler Internet Access (ZIA): rule_labels"
sidebar_current: "docs-resource-zia-rule-labels"
description: |-
  Creates and manages ZIA rule labels.
---


# Resource: zia_rule_labels

The **zia_rule_labels** resource allows the creation and management of rule labels in the Zscaler Internet Access cloud or via the API. This resource can then be associated with resources such as: Firewall Rules and URL filtering rules

## Example Usage

```hcl
# ZIA Rule Labels Resource
resource "zia_rule_labels" "example" {
    name        = "Example"
    description = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the devices to be created.

### Optional

* `description` - (String) The rule label description.
