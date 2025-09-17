---
subcategory: "Virtual Service Edges"
layout: "zscaler"
page_title: "ZIA: virtual_service_edge_node"
description: |-
    Official documentation https://help.zscaler.com/zia/about-virtual-service-edges
    API documentation https://help.zscaler.com/zia/service-edges#/virtualZenNodes-post
   Retrieves a list of ZIA Virtual Service Edge nodes.
---

# zia_virtual_service_edge_node (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-virtual-service-edges)
* [API documentation](https://help.zscaler.com/zia/service-edges#/virtualZenNodes-post)

Use the **zia_virtual_service_edge_node** data source to get information about a Virtual Service Edge Node for the specified `Name` or `ID`

## Example Usage

```hcl
data "zia_virtual_service_edge_node" "this"{
    name = "VSENode01"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) System-generated Virtual Service Edge cluster ID
* `name` - (Optional) Name of the Virtual Service Edge cluster

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - System-generated Virtual Service Edge cluster ID
* `name` - Name of the Virtual Service Edge cluster
* `status` - Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default
* `type` - The Virtual Service Edge cluster type
* `ip_sec_enabled` - A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
* `ip_address` - The Virtual Service Edge cluster IP address
* `subnet_mask` - The Virtual Service Edge cluster subnet mask
* `default_gateway` - The IP address of the default gateway to the internet
* `zgateway_id` - The Zscaler service gateway ID
* `in_production` - Represents the Virtual Service Edge instances deployed for production purposes
* `on_demand_support_tunnel_enabled` - A Boolean value that indicates whether or not the On-Demand Support Tunnel is enabled
* `establish_support_tunnel_enabled` - A Boolean value that indicates whether or not a support tunnel for Zscaler Support is enabled
* `load_balancer_ip_address` - The IP address of the load balancer. This field is applicable only when the 'deploymentMode' field is set to CLUSTER
* `deployment_mode` - Specifies the deployment mode. Select either STANDALONE or CLUSTER if you have the VMware ESXi platform. Otherwise, select only STANDALONE
* `cluster_name` - Virtual Service Edge cluster name
* `vzen_sku_type` - The Virtual Service Edge SKU type. Supported Values: SMALL, MEDIUM, LARGE
