package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_nodes"
)

func dataSourceVZENNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVZENNodesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "System-generated Virtual Service Edge cluster ID",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the Virtual Service Edge cluster",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge cluster type",
			},
			"ip_sec_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge cluster IP address",
			},
			"subnet_mask": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge cluster subnet mask",
			},
			"default_gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the default gateway to the internet",
			},
			"zgateway_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Zscaler service gateway ID",
			},
			"in_production": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Represents the Virtual Service Edge instances deployed for production purposes",
			},
			"on_demand_support_tunnel_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates whether or not the On-Demand Support Tunnel is enabled",
			},
			"establish_support_tunnel_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates whether or not a support tunnel for Zscaler Support is enabled",
			},
			"load_balancer_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the load balancer. This field is applicable only when the 'deploymentMode' field is set to CLUSTER",
			},
			"deployment_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the deployment mode. Select either STANDALONE or CLUSTER if you have the VMware ESXi platform. Otherwise, select only STANDALONE",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Virtual Service Edge cluster name",
			},
			"vzen_sku_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Virtual Service Edge SKU type. Supported Values: SMALL, MEDIUM, LARGE",
			},
		},
	}
}

func dataSourceVZENNodesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *vzen_nodes.VZENNodes
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for VZEN node id: %d\n", id)
		res, err := vzen_nodes.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for VZEN node name: %s\n", name)
		res, err := vzen_nodes.GetNodeByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("id", resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("status", resp.Status)
		_ = d.Set("type", resp.Type)
		_ = d.Set("ip_address", resp.IPAddress)
		_ = d.Set("subnet_mask", resp.SubnetMask)
		_ = d.Set("default_gateway", resp.DefaultGateway)
		_ = d.Set("ip_sec_enabled", resp.IPSecEnabled)
		_ = d.Set("zgateway_id", resp.ZGatewayID)
		_ = d.Set("in_production", resp.InProduction)
		_ = d.Set("on_demand_support_tunnel_enabled", resp.OnDemandSupportTunnelEnabled)
		_ = d.Set("establish_support_tunnel_enabled", resp.EstablishSupportTunnelEnabled)
		_ = d.Set("load_balancer_ip_address", resp.LoadBalancerIPAddress)
		_ = d.Set("deployment_mode", resp.DeploymentMode)
		_ = d.Set("cluster_name", resp.ClusterName)
		_ = d.Set("vzen_sku_type", resp.VzenSkuType)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any vzen node name '%s' or id '%d'", name, id))
	}

	return nil
}
