---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA): forwarding_control_proxy_gateway"
description: |-
  Get information about forwarding control proxy gateway.

---
# Data Source: zia_forwarding_control_proxy_gateway

Use the **zia_forwarding_control_proxy_gateway** data source to retrieve the proxy gateway information. This data source can then be associated with the attribute `proxy_gateway` when creating a Forwarding Control Rule via the resource: `zia_forwarding_control_rule`

## Example Usage

```hcl
# ZIA Forwarding Control - Proxy Gateway
data "zia_forwarding_control_proxy_gateway" "this" {
  name = "Proxy_GW01"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forwarding control Proxy Gateway to be exported.
* `id` - (Optional) The ID of the forwarding control Proxy Gateway resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional details about the Proxy gateway

* `last_modified_by` - (list) -  Information about the admin user that last modified the Proxy gateway
  * `id` - (int) - Identifier that uniquely identifies an entity
  * `name` - (string) - The configured name of the entity

* `last_modified_time` - (int) - Timestamp when the ZPA gateway was last modified

* `type` - (string) - Indicates whether the type of Proxy gateway. Returned values are: `PROXYCHAIN`, `ZIA`, or `ECSELF`

* `fail_closed` - (Boolean) - Indicates whether fail close is enabled to drop the traffic or disabled to allow the traffic when both primary and secondary proxies defined in this gateway are unreachable.

* `primary_proxy` - (Set of String) - The primary proxy for the gateway. This field is not applicable to the Lite API.
  * `id` - (string) A unique identifier for the primary proxy gateway
  * `name` - (string) The configured name for the primary proxy gateway

* `secondary_proxy` - () - The secondary proxy for the gateway. This field is not applicable to the Lite API.
  * `id` - (string) A unique identifier for the secondary proxy gateway
  * `name` - (string) The configured name for the secondary proxy gateway
