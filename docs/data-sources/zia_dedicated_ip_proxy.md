---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA: dedicated_ip_proxy"
description: |-
  Get information about Dedicated IP Gateways from the Zscaler Internet Access forwarding control policy.
---

# zia_dedicated_ip_proxy (Data Source)

Use the **zia_dedicated_ip_proxy** data source to get information about a Dedicated IP Gateway from the Zscaler Internet Access forwarding control policy. This data source uses the lite API endpoint and returns a simplified set of attributes.

## Example Usage - Retrieve By Name

```hcl
data "zia_dedicated_ip_proxy" "this" {
  name = "GW01"
}
```

## Example Usage - Retrieve By ID

```hcl
data "zia_dedicated_ip_proxy" "this" {
  id = 22394221
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The unique identifier of the Dedicated IP Gateway.
* `name` - (Optional) The name of the Dedicated IP Gateway.

One of `id` or `name` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (Integer) The unique identifier of the Dedicated IP Gateway.
* `name` - (String) The name of the Dedicated IP Gateway.
* `create_time` - (Integer) Timestamp when the Dedicated IP Gateway was created.
* `last_modified_time` - (Integer) Timestamp when the Dedicated IP Gateway was last modified.
* `default` - (Boolean) Indicates whether this is the default Dedicated IP Gateway.
