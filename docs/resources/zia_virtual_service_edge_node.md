---
subcategory: "Virtual Service Edges"
layout: "zscaler"
page_title: "ZIA: virtual_service_edge_node"
description: |-
    Official documentation https://help.zscaler.com/zia/about-virtual-service-edges
    API documentation https://help.zscaler.com/zia/service-edges#/virtualZenNodes-post
   Retrieves a list of ZIA Virtual Service Edge nodes.
---

# zia_virtual_service_edge_node (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-virtual-service-edges)
* [API documentation](https://help.zscaler.com/zia/service-edges#/virtualZenNodes-post)

Use the **zia_virtual_service_edge_node** resource allows the creation and management of Service Edge Node objects in the Zscaler Internet Access.
This resource can then be referenced within a `zia_virtual_service_edge_cluster` resource to create a cluster of Virtual Service Edge nodes.

## Example Usage

```hcl
resource "zia_virtual_service_edge_node" "this" {
  name                              = "VSENode01"
  status                            = "ENABLED"
  type                              = "VZEN"
  ip_address                        = "10.0.0.10"
  subnet_mask                       = "255.255.255.0"
  default_gateway                   = "10.0.0.1"
  zgateway_id                       = 12345
  in_production                     = true
  load_balancer_ip_address          = "10.0.0.20"
  deployment_mode                   = "STANDALONE"
  vzen_sku_type                     = "MEDIUM"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the Virtual Service Edge node
* `status` - (Optional) Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default

* `type` - (Optional) The Virtual Service Edge cluster type. For the complete list of supported types refer to the  [ZIA API documentation](https://help.zscaler.com/zia/service-edges#/virtualZenNodes-post).

* `ip_sec_enabled` - (Optional) A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
* `ip_address` - (Optional) The Virtual Service Edge cluster IP address. **Note**: Only IPv4 addresses are supported
* `subnet_mask` - (Optional) The Virtual Service Edge cluster subnet mask i.e `255.255.255.0`
* `default_gateway` - (Optional) The IP address of the default gateway to the internet. **Note**: Only IPv4 addresses are supported
* `zgateway_id` - (Optional) The Zscaler service gateway ID
* `in_production` - (Optional) Represents the Virtual Service Edge instances deployed for production purposes
* `on_demand_support_tunnel_enabled` - (Optional) A Boolean value that indicates whether or not the On-Demand Support Tunnel is enabled
* `establish_support_tunnel_enabled` - (Optional) A Boolean value that indicates whether or not a support tunnel for Zscaler Support is enabled
* `load_balancer_ip_address` - (Optional) The IP address of the load balancer. This field is applicable only when the 'deploymentMode' field is set to `CLUSTER`. **Note**: Only IPv4 addresses are supported

* `deployment_mode` - (Optional) Specifies the deployment mode. Select either `STANDALONE` or `CLUSTER` if you have the VMware ESXi platform. Otherwise, select only STANDALONE

* `vzen_sku_type` - (Optional) The Virtual Service Edge SKU type. Supported Values: `SMALL`, `MEDIUM`, `LARGE`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_virtual_service_edge_node** can be imported by using `<NODE_ID>` or `<NODE_NAME>` as the import ID.

For example:

```shell
terraform import zia_virtual_service_edge_node.example <node_id>
```

or

```shell
terraform import zia_virtual_service_edge_node.example <node_name>
```
