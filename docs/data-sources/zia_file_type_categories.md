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

Use the **zia_file_type_categories** data source to retrieve the list of all file types, including predefined and custom file types, available for configuring rule conditions in different ZIA policies. You can retrieve predefined file types for specific file categories of policies. This data source can be referenced within the `zia_dlp_web_rules` in the attribute `file_type_categories`.

The data source supports two modes:
- **Single Result Mode**: Retrieve a specific file type by `id` or `name` (returns `id`, `name`, `parent` fields)
- **List Mode**: Retrieve all file types matching filters like `enums` (returns results in the `categories` list)

## Example Usage - Retrieve a specific file type by name

```hcl
data "zia_file_type_categories" "javascript" {
    name = "FTCATEGORY_JAVASCRIPT"
}

output "file_type_id" {
    value = data.zia_file_type_categories.javascript.id
}
```

## Example Usage - Retrieve a specific file type by ID

```hcl
data "zia_file_type_categories" "by_id" {
    id = 10
}

output "file_type_name" {
    value = data.zia_file_type_categories.by_id.name
}
```

## Example Usage - Retrieve a specific file type with enum filter

```hcl
data "zia_file_type_categories" "javascript_in_file_control" {
    enums = "FILETYPECATEGORYFORFILETYPECONTROL"
    name  = "FTCATEGORY_JAVASCRIPT"
}
```

## Example Usage - Retrieve all file types for a policy category

```hcl
# Get all file types for File Type Control policy
data "zia_file_type_categories" "all_file_control" {
    enums = "FILETYPECATEGORYFORFILETYPECONTROL"
}

# Access all categories
output "all_file_types" {
    value = data.zia_file_type_categories.all_file_control.categories
}

# Get just the names
output "file_type_names" {
    value = data.zia_file_type_categories.all_file_control.categories[*].name
}
```

## Example Usage - Retrieve file types excluding custom types

```hcl
data "zia_file_type_categories" "predefined_only" {
    enums = "FILETYPECATEGORYFORFILETYPECONTROL"
    exclude_custom_file_types = true
}

output "predefined_file_types" {
    value = data.zia_file_type_categories.predefined_only.categories
}
```

## Example Usage - Get all DLP file types

```hcl
data "zia_file_type_categories" "dlp_types" {
    enums = "ZSCALERDLP"
}

# Use in a DLP rule
resource "zia_dlp_web_rules" "example" {
    name = "Example DLP Rule"
    file_types = data.zia_file_type_categories.dlp_types.categories[*].name
    # ... other configuration
}
```

## Argument Reference

The following arguments are supported:

### Optional

All arguments are optional, but you must provide at least one of: `id`, `name`, or `enums`.

* `id` - (Integer) File type ID. When specified, returns a single file type category matching this ID.
* `name` - (String) File type name. When specified, returns a single file type category matching this name (case-insensitive).
* `enums` - (String) Enum value to filter file types for specific policy categories. Valid values:
  * `ZSCALERDLP` - Web DLP rules with content inspection
  * `EXTERNALDLP` - Web DLP rules without content inspection
  * `FILETYPECATEGORYFORFILETYPECONTROL` - File Type Control policy

  When used alone (without `id` or `name`), returns all matching file types in the `categories` attribute.
  When combined with `name`, filters the search to categories within this enum.

* `exclude_custom_file_types` - (Boolean) Whether to exclude custom file types from the results. Defaults to `false`. Set to `true` to return only predefined file types.

## Attributes Reference

The following attributes are exported:

### Single Result Mode (when `id` or `name` is provided)

* `id` - (Integer) File type category ID
* `name` - (String) File type category name
* `parent` - (String) Parent category of the file type
* `categories` - (List) Empty list when in single result mode

### List Mode (when only `enums` is provided)

* `id` - (Integer) Generated hash ID for this query
* `categories` - (List of Objects) List of file type categories matching the filter. Each object contains:
  * `id` - (Integer) File type category ID
  * `name` - (String) File type category name
  * `parent` - (String) Parent category of the file type
