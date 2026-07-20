---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_resource_group"
description: |-
  Official documentation https://help.zscaler.com/zia/about-endpoint-dlp
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpResourceGroups-post
  Creates and manages a ZIA DLP endpoint resource group (tag)
---

# zia_endpoint_dlp_resource_group (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-endpoint-dlp)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointDlpResourceGroups-post)

The **zia_endpoint_dlp_resource_group** resource allows the creation and management of ZIA DLP endpoint resource groups (tags) in the Zscaler Internet Access cloud or via the API. A resource group is channel-specific and is used to group one or more DLP endpoint resources under a single tag, which can then be referenced by Endpoint DLP rules.

~> **NOTE:** A resource group belongs to a single `channel` and cannot be moved between channels after creation. Changing the `channel` on an existing group is not supported by the API.

~> **NOTE:** The DLP endpoint resources referenced by the `resources` argument must already exist. Manage them with the [`zia_endpoint_dlp_resource`](./zia_endpoint_dlp_resource.md) resource.

## Example Usage - Empty Group

```hcl
resource "zia_endpoint_dlp_resource_group" "printers" {
  channel     = "PRINTING"
  name        = "PrinterTag01"
  description = "Finance department printers"
}
```

## Example Usage - Group With Member Resources

```hcl
resource "zia_endpoint_dlp_resource" "printer01" {
  channel     = "PRINTING"
  name        = "Printer01"
  description = "Finance department printer"

  printer {
    ip_address = "10.10.10.20"
    domain     = "acme.local"
  }
}

resource "zia_endpoint_dlp_resource_group" "printers" {
  channel     = "PRINTING"
  name        = "PrinterTag01"
  description = "Finance department printers"

  resources = [
    zia_endpoint_dlp_resource.printer01.resource_id,
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `channel` - (String) The channel the DLP endpoint resource group belongs to. Supported values: `PRINTING`, `REMOVABLE_DRIVE_TRANSFER`, `NETWORK_DRIVE_TRANSFER`. A group cannot be moved between channels after creation.
* `name` - (String) The name of the DLP endpoint resource group (tag).

### Optional

* `description` - (String) The description of the DLP endpoint resource group.
* `resources` - (Set of Number) The set of DLP endpoint resource IDs associated with this group. Membership is reconciled by adding and removing only the resources that changed, so managing this set never re-sends the full list. Leave unset (or set to an empty set) to manage a group with no members.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The Terraform resource identifier (the numeric group ID as a string).
* `group_id` - (Number) The identifier assigned to the DLP endpoint resource group by the API.

## Import

Because the read path is channel-scoped, the import ID must include the channel.

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_endpoint_dlp_resource_group** can be imported using either:

* `<CHANNEL>:<id>` — the channel and numeric ID of the group, or
* `<CHANNEL>:<name>` — the channel and exact group name.

For example:

```shell
terraform import zia_endpoint_dlp_resource_group.example "PRINTING:367"
```

or

```shell
terraform import zia_endpoint_dlp_resource_group.example "PRINTING:PrinterTag01"
```
