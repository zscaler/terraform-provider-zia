---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_edm_schema"
description: |-
  Official documentation https://help.zscaler.com/zia/about-exact-data-match
  API documentation https://help.zscaler.com/zia/data-loss-prevention#/dlpExactDataMatchSchemas-get
  Gets the list of DLP Exact Data Match (EDM) templates
---

# zia_dlp_edm_schema (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-exact-data-match)
* [API documentation](https://help.zscaler.com/zia/data-loss-prevention#/dlpExactDataMatchSchemas-get)

Use the **zia_dlp_edm_schema** data source to get information about a the list of DLP Exact Data Match (EDM) templates in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP Exact Data Match (EDM) by name
data "zia_dlp_edm_schema" "this"{
    project_name = "Example"
}
```

```hcl
# Retrieve a DLP Exact Data Match (EDM) by ID
data "zia_dlp_edm_schema" "example"{
    schema_id = 1234567890
}
```

## Argument Reference

The following arguments are supported:

### Required

* `project_name` - (Required) The EDM schema (i.e., EDM template) name.

### Optional

* `schema_id` - (Number) The identifier (1-65519) for the EDM schema (i.e., EDM template) that is unique within the organization.
* `revision` - (String) The revision number of the CSV file upload to the Index Tool.
* `file_name` - (String) The generated filename, excluding the extention.
* `original_file_name` - (String) The name of the CSV file uploaded to the Index Tool.
* `file_upload_status` - (String) The status of the EDM template's CSV file upload to the Index Tool
* `orig_col_count` - (String) The total count of actual columns selected from the CSV file.
* `created_by` - (Number)  The login name (or userid) the admin who created the EDM schema (i.e., EDM template).
  * `id` - (Number) Identifier that uniquely identifies an entity
* `last_modified_by` - (Number)  The admin that modified the DLP policy rule last.
  * `id` - (Number) Identifier that uniquely identifies an entity
* `last_modified_time` - (Number) Timestamp when the DLP policy rule was last modified.
* `cells_used` - (Number) The total number of cells used by the EDM schema (i.e., EDM template).
* `schema_active` - (Bool) Indicates the status of a specified EDM schema (i.e., EDM template)
* `token_list` - (Bool) The list of tokens (or criteria) to match against. Up to 16 tokens can be selected for data matching
  * `name` - (String) The token (i.e., criteria) name.
  * `type` - (String) The token type.
  * `primary_key` (Bool) Indicates whether the token is a primary key.
  * `original_column` (Number) The column position for the token in the original CSV file uploaded to the Index Tool, starting from 1.
  * `hash_file_column_order` (Number) The column position for the token in the hashed file, starting from 1.
  * `col_length_bitmap` (Number) The length of the column bitmap in the hashed file.
* `schedule_present` - (Bool) Indicates whether the EDM schema (i.e., the EDM template configured using the Index Tool) has a schedule.
* `schedule` - (List) The schedule details, if present for the EDM schema (i.e., EDM template).
  * `schedule_type` - (String) The schedule type for the EDM schema (i.e., EDM template), Monthly, Weekly, Daily, or None.
  * `schedule_day_of_month` - (String) The day of the month the EDM schema (i.e., EDM template) is scheduled for.
  * `schedule_day_of_week` - (List String) The day of the week the EDM schema (i.e., EDM template) is scheduled for.
  * `schedule_time` - (Number) The time of the day (in minutes) that the EDM schema (i.e., EDM template) is scheduled for.
  * `schedule_disabled` - (Number) If set to true, the schedule for the EDM schema (i.e., EDM template) is temporarily in a disabled state.
