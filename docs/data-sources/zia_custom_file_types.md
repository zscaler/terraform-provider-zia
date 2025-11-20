---
subcategory: "File Type Control Policy"
layout: "zscaler"
page_title: "ZIA: zia_custom_file_types"
description: |-
  Official documentation https://help.zscaler.com/zia/about-file-type-control
  API documentation https://help.zscaler.com/zia/file-type-control-policy#/fileTypeRules-post
  Retrieves all the rules in the File Type Control policy.
---

# zia_custom_file_types (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-file-type-control)
* [API documentation](https://help.zscaler.com/zia/file-type-control-policy#/fileTypeRules-post)

Use the **zia_custom_file_types** data source to retrieves File Type Control rules.

## Example Usage

```hcl
# Retrieve a File Type Control Rule by name
data "zia_custom_file_types" "this" {
    name = "Example"
}
```

```hcl
# Retrieve a File Type Control Rule by ID
data "zia_custom_file_types" "this" {
    id = "12134558"
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `id` - (Integer) Custom file type ID. This ID is assigned and maintained exclusively for custom file types, and this value is different from the file type ID (i.e., fileTypeId field).
* `name` - (String) Custom file type name
* `description` - (String) Additional information about the custom file type, if any.
* `extension` - (String) Specifies the file type extension. The maximum extension length is 10 characters. Existing Zscaler extensions cannot be added to custom file types.
* `file_type_id` - (Integer) File type ID. This ID is assigned and maintained for all file types including predefined and custom file types, and this value is different from the custom file type ID.
