---
subcategory: "Data Loss Prevention"
layout: "zia"
page_title: "ZIA: dlp_dictionary_predefined_identifiers"
description: |-
  Get information about DLP Predefined Identifiers.
---

# Data Source: zia_dlp_dictionary_predefined_identifiers

Use the **zia_dlp_dictionary_predefined_identifiers** data source to get information about the list of predefined identifiers that are available for selection in the specified hierarchical DLP dictionary.

## Example Usage

```hcl
data "zia_dlp_dictionary_predefined_identifiers" "this" {
  name = "CRED_LEAKAGE"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the hierarchical DLP dictionary.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the hierarchical DLP dictionary.
* `name` - (String) The name of the hierarchical DLP dictionary: Supported values: `ASPP_LEAKAGE`, `CRED_LEAKAGE`, `EUIBAN_LEAKAGE`, `PPEU_LEAKAGE`, `USDL_LEAKAGE`.
* `predefined_identifiers` - (List) The list of hierarchical DLP dictionary values.
