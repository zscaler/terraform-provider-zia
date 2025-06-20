---
subcategory: "SaaS Security API"
layout: "zscaler"
page_title: "ZIA): casb_tombstone_template"
description: |-
  Official documentation https://help.zscaler.com/zia/about-quarantine-tombstone-file-templates
  API documentation https://help.zscaler.com/zia/saas-security-api#/quarantineTombstoneTemplate/lite-get
  Retrieves the templates for the tombstone file created when a file is quarantined

---
# zia_casb_tombstone_template (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-quarantine-tombstone-file-templates)
* [API documentation](https://help.zscaler.com/zia/saas-security-api#/quarantineTombstoneTemplate/lite-get)

Use the **zia_casb_tombstone_template** data source to get information about templates for the tombstone file created when a file is quarantined.

## Example Usage - By Name

```hcl

data "zia_casb_tombstone_template" "this" {
  name = "TombstoneTemplate01"
}
```

## Example Usage - By ID

```hcl

data "zia_casb_tombstone_template" "this" {
  id = 154658
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Int) Tombstone file template ID
* `name` - (String) Tombstone file template name

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

N/A
