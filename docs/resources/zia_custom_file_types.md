---
subcategory: "File Type Control Policy"
layout: "zscaler"
page_title: "ZIA: zia_custom_file_types"
description: |-
  Official documentation https://help.zscaler.com/zia/about-file-type-control
  API documentation https://help.zscaler.com/zia/file-type-control-policy#/customFileTypes-get
  Adds a new custom file type
---

# zia_custom_file_types (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-file-type-control)
* [API documentation](https://help.zscaler.com/zia/file-type-control-policy#/customFileTypes-get)

The **zia_custom_file_types** resource allows the creation and management of ZIA custom file type in the Zscaler Internet Access.

## Example Usage

```hcl
resource "zia_custom_file_types" "this" {
  name = "FileType02"
  description = "FileType02"
  extension = "tf"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) Custom file type name
* `extension` - (Required) Custom file type name

### Optional

* `description` - (Optional) Additional information about the custom file type, if any.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_custom_file_types** can be imported by using `<FILE ID>` or `<FILE NAME>` as the import ID.

For example:

```shell
terraform import zia_custom_file_types.example <file_id>
```

or

```shell
terraform import zia_custom_file_types.example <file_name>
```
