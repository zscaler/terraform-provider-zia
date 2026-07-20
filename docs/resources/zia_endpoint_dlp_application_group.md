---
subcategory: "Endpoint Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: endpoint_dlp_application_group"
description: |-
  Official documentation https://help.zscaler.com/zia/about-endpoint-dlp
  API documentation https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplicationGroups-post
  Creates and manages a ZIA Endpoint DLP application group
---

# zia_endpoint_dlp_application_group (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-endpoint-dlp)
* [API documentation](https://help.zscaler.com/legacy-apis/endpoint-data-loss-prevention-dlp-policy#/endPointApplicationGroups-post)

The **zia_endpoint_dlp_application_group** resource allows the creation and management of ZIA Endpoint DLP application groups in the Zscaler Internet Access cloud or via the API. An application group bundles one or more endpoint applications under a single tag, which can then be referenced by Endpoint DLP rules through their `end_point_application_groups` selector.

~> **NOTE:** Application groups only support the `APPLICATION_FILE_ACCESS` channel. The `channel` argument defaults to that value and cannot be set to anything else.

~> **NOTE:** Each member is identified by its application ID (the endpoint application's `zappId`). Use the [`zia_endpoint_dlp_application`](../data-sources/zia_endpoint_dlp_application.md) data source to look up the `zapp_id` of an application by name.

## Example Usage - Empty Group

```hcl
resource "zia_endpoint_dlp_application_group" "this" {
  name        = "Tag01"
  description = "test"
}
```

## Example Usage - Group With Member Applications

```hcl
resource "zia_endpoint_dlp_application_group" "this" {
  name        = "Tag01"
  description = "test"

  resources = [
    "500000085",
    "600000022",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The name of the endpoint application group.

### Optional

* `channel` - (String) The channel the endpoint application group belongs to. Only `APPLICATION_FILE_ACCESS` is supported, which is also the default.
* `description` - (String) The description of the endpoint application group.
* `resources` - (Set of String) The set of endpoint application IDs (the endpoint applications' `zappId`) that belong to this group. Leave unset to manage a group with no members.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The Terraform resource identifier (the numeric group ID as a string).
* `group_id` - (Number) The identifier assigned to the endpoint application group by the API.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_endpoint_dlp_application_group** can be imported using either the numeric group ID or the exact group name.

For example:

```shell
terraform import zia_endpoint_dlp_application_group.example 366
```

or

```shell
terraform import zia_endpoint_dlp_application_group.example Tag01
```
