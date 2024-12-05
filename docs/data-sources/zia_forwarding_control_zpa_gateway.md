---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA): forwarding_control_zpa_gateway"
description: |-
  Get information about forwarding control zpa gateway used in IP Source Anchoring.

---
# Data Source: zia_forwarding_control_zpa_gateway

Use the **zia_forwarding_control_zpa_gateway** data source to get information about a forwarding control zpa gateway used in IP Source Anchoring integration between Zscaler Internet Access and Zscaler Private Access. This data source can then be associated with a ZIA Forwarding Control Rule.

## Example Usage

```hcl
# ZIA Forwarding Control - ZPA Gateway
data "zia_forwarding_control_zpa_gateway" "this" {
  name = "ZPA_GW01"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding control ZPA Gateway to be exported.
* `id` - (Optional) The ID of the forwarding control ZPA Gateway resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional details about the ZPA gateway
* `last_modified_by` - (list) -  Information about the admin user that last modified the ZPA gateway
  * `id` - (int) - Identifier that uniquely identifies an entity
  * `name` - (string) - The configured name of the entity
* `last_modified_time` - (int) - Timestamp when the ZPA gateway was last modified
* `type` - (string) - Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA)
* `zpa_server_group` - () - The ZPA Server Group that is configured for Source IP Anchoring
  * `external_id` - (string) An external identifier used for an entity that is managed outside of ZIA. Examples include zpaServerGroup and zpaAppSegments. This field is not applicable to ZIA-managed entities.
  * `name` - (string) The configured name of the entity
