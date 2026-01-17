---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA): forwarding_control_dedicated_ip_gateway"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-forwarding-policy
  Get information about forwarding control dedicated IP gateway.
---

# zia_forwarding_control_dedicated_ip_gateway (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-forwarding-policy)

Use the **zia_forwarding_control_dedicated_ip_gateway** data source to get information about a dedicated IP gateway that can be associated with a ZIA Forwarding Control Rule using the `ENATDEDIP` forwarding method.

## Example Usage

```hcl
data "zia_forwarding_control_dedicated_ip_gateway" "this" {
  name = "Default"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the dedicated IP gateway.
* `id` - (Optional) The ID of the dedicated IP gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional details about the dedicated IP gateway
* `create_time` - (int) - Timestamp when the dedicated IP gateway was created
* `last_modified_by` - (list) - Information about the admin user that last modified the dedicated IP gateway
  * `id` - (int) - Identifier that uniquely identifies an entity
  * `name` - (string) - The configured name of the entity
* `last_modified_time` - (int) - Timestamp when the dedicated IP gateway was last modified
* `default` - (bool) - Indicates whether this is the default dedicated IP gateway
