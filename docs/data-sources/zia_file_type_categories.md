---
subcategory: "File Type Control Policy"
layout: "zscaler"
page_title: "ZIA: zia_file_type_categories"
description: |-
  Official documentation https://help.zscaler.com/zia/about-file-type-control
  API documentation https://help.zscaler.com/zia/file-type-control-policy#/fileTypeCategories-get
  Retrieves the list of all file type
---

# zia_file_type_categories (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-file-type-control)
* [API documentation](https://help.zscaler.com/zia/file-type-control-policy#/fileTypeCategories-get)

Use the **zia_file_type_categories** Retrieves the list of all file types, including predefined and custom file types, available for configuring rule conditions in different ZIA policies. You can retrieve predefined file types for specific file categories of policies. This datasource can be referenced within the `zia_dlp_web_rules` in the attribute `file_type_categories`

## Example Usage - Retrieve a File Type Category by name

```hcl
data "zia_file_type_categories" "this1" {
    name = "FileType01"
}
```

## Example Usage - Retrieve a File Type Category by ID

```hcl
data "zia_file_type_categories" "this" {
    id = 12134558
}
```

## Example Usage - Retrieve a File Type Category with enum filter

```hcl
data "zia_file_type_categories" "this2" {
    name = "FileType01"
    enums = ["ZSCALERDLP"]
}
```

## Example Usage - Retrieve a File Type Category with multiple enum filters

```hcl
data "zia_file_type_categories" "this3" {
    name = "FileType01"
    enums = ["ZSCALERDLP", "EXTERNALDLP"]
}
```

## Example Usage - Retrieve a File Type Category excluding custom file types

```hcl
data "zia_file_type_categories" "this4" {
    name = "FileType01"
    exclude_custom_file_types = true
}
```

## Example Usage - Retrieve a File Type Category with all optional parameters

```hcl
data "zia_file_type_categories" "this5" {
    name = "FileType01"
    enums = ["FILETYPECATEGORYFORFILETYPECONTROL"]
    exclude_custom_file_types = true
}
```

## Argument Reference

The following arguments are supported:

### Required

At least one of the following must be provided:

* `id` - (Integer) File type ID. This ID is assigned and maintained exclusively for custom file types, and this value is different from the file type ID (i.e., fileTypeId field).
* `name` - (String) File type name. Used to search for a file type category by name.

### Optional

* `enums` - (List of Strings) Enum values to filter file types for specific policy categories. Valid values:
  * `ZSCALERDLP` - Filter for Zscaler DLP policy categories
  * `EXTERNALDLP` - Filter for External DLP policy categories
  * `FILETYPECATEGORYFORFILETYPECONTROL` - Filter for File Type Control policy categories

  Multiple enum values can be specified to filter across different policy categories.

* `exclude_custom_file_types` - (Boolean) A Boolean value specifying whether custom file types must be excluded from the list or not. Defaults to `false`. Set to `true` to exclude custom file types and only return predefined file types.

## Attributes Reference

The following attributes are exported:

* `id` - (Integer) File type ID
* `name` - (String) File type name
* `parent` - (String) Parent category of the file type
