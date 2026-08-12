---
subcategory: "PAC Files"
layout: "zscaler"
page_title: "ZIA: pac_files"
description: |-
  Official documentation https://help.zscaler.com/zia/about-hosted-pac-files
  API documentation https://help.zscaler.com/zia/hosted-pac-files#/pacFiles-get
  Get information about ZIA hosted PAC files.
---

# zia_pac_files (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-hosted-pac-files)
* [API documentation](https://help.zscaler.com/zia/hosted-pac-files#/pacFiles-get)

Use the **zia_pac_files** data source to retrieve the list of PAC files in deployed state, including default and custom PAC files. A single PAC file can be retrieved by `id` or `name`, and the PAC file content can be omitted from the results via the `filter` attribute.

## Example Usage - Retrieve All PAC Files

```hcl
data "zia_pac_files" "all" {}
```

## Example Usage - Retrieve All PAC Files Without Content

```hcl
data "zia_pac_files" "no_content" {
  filter = "pac_content"
}
```

## Example Usage - Retrieve a PAC File By Name

```hcl
data "zia_pac_files" "this" {
  name = "Default PAC"
}
```

## Example Usage - Retrieve a PAC File By ID

```hcl
data "zia_pac_files" "this" {
  id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The unique identifier of the PAC file to retrieve. When set, only the matching PAC file is returned.
* `name` - (Optional) The name of the PAC file to retrieve. When set, only the matching PAC file is returned.
* `filter` - (Optional) When set to `pac_content`, the PAC file content is omitted from the results. Supported value: `pac_content`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `pac_files` - The list of PAC files. Each entry exports the following attributes:
  * `id` - (Number) The unique identifier for the PAC file
  * `name` - (String) The name of the PAC file
  * `description` - (String) The description of the PAC file
  * `domain` - (String) The domain of your organization to which the PAC file applies
  * `pac_url` - (String) The URL location of the PAC file, auto-generated when the PAC file is first added
  * `pac_content` - (String) The content of the PAC file. Empty when the `filter` attribute is set to `pac_content`
  * `editable` - (Boolean) Indicates whether the PAC file is editable
  * `pac_sub_url` - (String) The obfuscated URL of the PAC file. Returned when `pac_url_obfuscated` is `true`
  * `pac_url_obfuscated` - (Boolean) Indicates whether the PAC file URL is obfuscated
  * `pac_verification_status` - (String) The verification status of the PAC file, indicating whether any syntax errors were identified. Supported values: `VERIFY_NOERR`, `VERIFY_ERR`, `NOVERIFY`
  * `pac_version_status` - (String) The status of the PAC file version. Supported values: `DEPLOYED`, `STAGE`, `LKG`
  * `pac_version` - (Number) The version number of the PAC file
  * `pac_commit_message` - (String) The commit message entered when the PAC file version was saved
  * `total_hits` - (Number) The number of times the PAC file was used during the last 30 days
  * `last_modified_time` - (Number) Timestamp when the PAC file was last modified, in Unix time
  * `create_time` - (Number) Timestamp when the PAC file was created, in Unix time
  * `last_modified_by` - (List) The admin that last modified the PAC file
    * `id` - (Number) Identifier that uniquely identifies an entity
    * `name` - (String) The configured name of the entity
