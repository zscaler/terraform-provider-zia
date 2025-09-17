package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_nodes"
)

func resourceVZENNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVZENNodeCreate,
		ReadContext:   resourceVZENNodeRead,
		UpdateContext: resourceVZENNodeUpdate,
		DeleteContext: resourceVZENNodeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("node_id", idInt)
				} else {
					resp, err := vzen_nodes.GetNodeByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("node_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the Virtual Service Edge node",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
					"DISABLED_BY_SERVICE_PROVIDER",
					"NOT_PROVISIONED_IN_SERVICE_PROVIDER",
					"IN_TRIAL",
				}, false),
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge cluster type",
			},
			"ip_sec_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge cluster IP address",
				ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
					v, ok := i.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
						return warnings, errors
					}
					// Allow empty string or valid IP address
					if v == "" {
						return warnings, errors
					}
					return validation.IsIPAddress(i, k)
				},
			},
			"subnet_mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge cluster subnet mask",
			},
			"default_gateway": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the default gateway to the internet",
				ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
					v, ok := i.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
						return warnings, errors
					}
					// Allow empty string or valid IP address
					if v == "" {
						return warnings, errors
					}
					return validation.IsIPAddress(i, k)
				},
			},
			"in_production": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Represents the Virtual Service Edge instances deployed for production purposes",
			},
			"on_demand_support_tunnel_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value that indicates whether or not the On-Demand Support Tunnel is enabled",
			},
			"establish_support_tunnel_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value that indicates whether or not a support tunnel for Zscaler Support is enabled",
			},
			"load_balancer_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the load balancer. This field is applicable only when the 'deploymentMode' field is set to CLUSTER",
				ValidateFunc: func(i interface{}, k string) (warnings []string, errors []error) {
					v, ok := i.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
						return warnings, errors
					}
					// Allow empty string or valid IP address
					if v == "" {
						return warnings, errors
					}
					return validation.IsIPAddress(i, k)
				},
			},
			"deployment_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the deployment mode. Select either STANDALONE or CLUSTER if you have the VMware ESXi platform. Otherwise, select only STANDALONE",
				ValidateFunc: validation.StringInSlice([]string{
					"STANDALONE",
					"CLUSTER",
				}, false),
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Virtual Service Edge cluster name",
			},
			"vzen_sku_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Virtual Service Edge SKU type. Supported Values: SMALL, MEDIUM, LARGE",
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL",
					"MEDIUM",
					"LARGE",
				}, false),
			},
		},
	}
}

func resourceVZENNodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandVZENNodes(d)
	log.Printf("[INFO] Creating ZIA vzen nodes\n%+v\n", req)

	resp, _, err := vzen_nodes.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA vzen nodes request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("node_id", resp.ID)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceVZENNodeRead(ctx, d, meta)
}

func resourceVZENNodeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "node_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no vzen nodes id is set"))
	}
	resp, err := vzen_nodes.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia vzen nodes %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia vzen nodes:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("node_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("status", resp.Status)
	_ = d.Set("type", resp.Type)
	_ = d.Set("ip_address", resp.IPAddress)
	_ = d.Set("subnet_mask", resp.SubnetMask)
	_ = d.Set("default_gateway", resp.DefaultGateway)
	_ = d.Set("ip_sec_enabled", resp.IPSecEnabled)
	_ = d.Set("in_production", resp.InProduction)
	_ = d.Set("on_demand_support_tunnel_enabled", resp.OnDemandSupportTunnelEnabled)
	_ = d.Set("establish_support_tunnel_enabled", resp.EstablishSupportTunnelEnabled)
	_ = d.Set("load_balancer_ip_address", resp.LoadBalancerIPAddress)
	_ = d.Set("deployment_mode", resp.DeploymentMode)
	_ = d.Set("cluster_name", resp.ClusterName)
	_ = d.Set("vzen_sku_type", resp.VzenSkuType)

	return nil
}

func resourceVZENNodeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "node_id")
	if !ok {
		log.Printf("[ERROR] vzen nodes ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia vzen nodes ID: %v\n", id)

	req := expandVZENNodes(d)
	if _, err := vzen_nodes.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := vzen_nodes.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceVZENClusterRead(ctx, d, meta)
}

func resourceVZENNodeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "node_id")
	if !ok {
		log.Printf("[ERROR] vzen nodes ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia vzen nodes ID: %v\n", (d.Id()))

	if _, err := vzen_nodes.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia vzen nodes deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandVZENNodes(d *schema.ResourceData) vzen_nodes.VZENNodes {
	id, _ := getIntFromResourceData(d, "node_id")
	result := vzen_nodes.VZENNodes{
		ID:                            id,
		Name:                          d.Get("name").(string),
		Status:                        d.Get("status").(string),
		Type:                          d.Get("type").(string),
		IPAddress:                     d.Get("ip_address").(string),
		SubnetMask:                    d.Get("subnet_mask").(string),
		DefaultGateway:                d.Get("default_gateway").(string),
		IPSecEnabled:                  d.Get("ip_sec_enabled").(bool),
		InProduction:                  d.Get("in_production").(bool),
		OnDemandSupportTunnelEnabled:  d.Get("on_demand_support_tunnel_enabled").(bool),
		EstablishSupportTunnelEnabled: d.Get("establish_support_tunnel_enabled").(bool),
		LoadBalancerIPAddress:         d.Get("load_balancer_ip_address").(string),
		DeploymentMode:                d.Get("deployment_mode").(string),
		ClusterName:                   d.Get("cluster_name").(string),
		VzenSkuType:                   d.Get("vzen_sku_type").(string),
	}
	return result
}
