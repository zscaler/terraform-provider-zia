---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_endpoint_resource_group_tag"
description: |-
  Official documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpResourceGroups/{channel}-get
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpResourceGroups/{channel}-get
  Get the list of DLP endpoint resource tags defined for a channel
---

# zia_dlp_endpoint_resource_group_tag (Data Source)

* [Official documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpResourceGroups/{channel}-get)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpResourceGroups/{channel}-get)

Use the **zia_dlp_endpoint_resource_group_tag** data source to get the list of DLP endpoint resource tags, in name-ID pairs, defined for a channel in the Zscaler Internet Access cloud or via the API. These tag IDs are commonly referenced by the `resource_groups` argument of DLP rule resources.

Tags are channel-specific. Provide the `channel` to return all tags for that channel, or additionally filter the result to a single tag with `name` or `id`.

## Example Usage

```hcl
# Return all tags defined for a channel
data "zia_dlp_endpoint_resource_group_tag" "all" {
  channel = "PRINTING"
}
```

```hcl
# Filter to a single tag by name
data "zia_dlp_endpoint_resource_group_tag" "printer" {
  channel = "PRINTING"
  name    = "PrinterTag01"
}
```

```hcl
# Filter to a single tag by ID
data "zia_dlp_endpoint_resource_group_tag" "printer" {
  channel = "PRINTING"
  id      = 367
}
```

## Argument Reference

The following arguments are supported:

* `channel` - (Required, String) The DLP endpoint resource channel whose tags are returned. Supported values: `PRINTING`, `REMOVABLE_DRIVE_TRANSFER`, `NETWORK_DRIVE_TRANSFER`, `PERSONAL_CLOUD_STORAGE`.
* `name` - (Optional, String) The name of a specific tag to return. When set, the result is filtered to this tag.
* `id` - (Optional, Number) The unique identifier of a specific tag to return. When set, the result is filtered to this tag.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `tags` - (List of Object) The list of tags, in name-ID pairs, defined for the specified channel.
  * `id` - (Number) Identifier that uniquely identifies the tag.
  * `name` - (String) The configured name of the tag.
  * `description` - (String) The description of the tag.

* `tags_by_name` - (Map of Number) A convenience map of tag name to tag ID, allowing a single tag ID to be looked up directly by name.

## Referencing tag IDs

This data source returns a **list** of tags. Two exported attributes make the result easy to consume:

* `tags` — the full list of `{ id, name, description }` objects.
* `tags_by_name` — a map of tag name to tag ID, for looking up a **single** tag ID by name.

### Get a single tag ID by name

Index the `tags_by_name` map with the tag name. This returns a single number:

```hcl
# Scalar single value → 367
data.zia_dlp_endpoint_resource_group_tag.printer.tags_by_name["PrinterTag01"]
```

When the target argument expects a list (such as `resource_groups.id`), wrap the single value in `[ ]`:

```hcl
resource_groups {
  id = [data.zia_dlp_endpoint_resource_group_tag.printer.tags_by_name["PrinterTag01"]]
}
```

Equivalent lookup against the raw `tags` list (no map), if preferred:

```hcl
one([for t in data.zia_dlp_endpoint_resource_group_tag.printer.tags : t.id if t.name == "PrinterTag01"])
```

### Get all tag IDs

`tags[*].id` is a splat expression that is **already a list** of numbers, so assign it directly — do **not** wrap it in `[ ]` (that would produce a list of lists):

```hcl
resource_groups {
  id = data.zia_dlp_endpoint_resource_group_tag.all.tags[*].id
}
```

### Get a specific subset of tag IDs by name

```hcl
resource_groups {
  id = [for n in ["PrinterTag01", "PrinterTag02"] : data.zia_dlp_endpoint_resource_group_tag.all.tags_by_name[n]]
}
```

~> **NOTE:** Use the bare splat `tags[*].id` (no brackets) when passing **all** tag IDs, and wrap `tags_by_name["<name>"]` in `[ ]` when passing a **single** tag ID into a list argument. A single value needs the brackets; a list (splat / comprehension) does not.
