---
subcategory: "Rule Labels"
layout: "zscaler"
page_title: "ZIA: rule_labels"
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
