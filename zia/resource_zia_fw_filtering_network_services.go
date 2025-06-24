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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
)

func resourceFWNetworkServices() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkServicesCreate,
		ReadContext:   resourceNetworkServicesRead,
		UpdateContext: resourceNetworkServicesUpdate,
		DeleteContext: resourceNetworkServicesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("network_service_id", idInt)
				} else {
					resp, err := networkservices.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("network_service_id", resp.ID)
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
			"network_service_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"tag":            getCloudFirewallNwServicesTag(),
			"src_tcp_ports":  resourceNetworkPortsSchema("src tcp ports"),
			"dest_tcp_ports": resourceNetworkPortsSchema("dest tcp ports"),
			"src_udp_ports":  resourceNetworkPortsSchema("src udp ports"),
			"dest_udp_ports": resourceNetworkPortsSchema("dest udp ports"),
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD",
					"PREDEFINED",
					"CUSTOM",
				}, false),
			},
			"is_name_l10n_tag": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNetworkServicesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandNetworkServices(d)
	log.Printf("[INFO] Creating network services\n%+v\n", req)

	resp, err := networkservices.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia network services request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("network_service_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceNetworkServicesRead(ctx, d, meta)
}

func resourceNetworkServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no network services id is set"))
	}
	resp, err := networkservices.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia network services %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting network services :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("network_service_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("tag", resp.Tag)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

	if err := d.Set("src_tcp_ports", flattenNetwordPorts(resp.SrcTCPPorts)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dest_tcp_ports", flattenNetwordPorts(resp.DestTCPPorts)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("src_udp_ports", flattenNetwordPorts(resp.SrcUDPPorts)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_udp_ports", flattenNetwordPorts(resp.DestUDPPorts)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkServicesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		log.Printf("[ERROR] network service ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network service ID: %v\n", id)
	req := expandNetworkServices(d)
	if _, err := networkservices.Get(ctx, service, req.ID); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := networkservices.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceNetworkServicesRead(ctx, d, meta)
}

func resourceNetworkServicesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		log.Printf("[ERROR] network service id ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting network service ID: %v\n", (d.Id()))
	err := DetachRuleIDNameExtensions(
		zClient,
		id,
		"NwServices",
		func(r *filteringrules.FirewallFilteringRules) []common.IDNameExtensions {
			return r.NwServices
		},
		func(r *filteringrules.FirewallFilteringRules, ids []common.IDNameExtensions) {
			r.NwServices = ids
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := networkservices.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] network service deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandNetworkServices(d *schema.ResourceData) networkservices.NetworkServices {
	id, _ := getIntFromResourceData(d, "network_service_id")
	result := networkservices.NetworkServices{
		ID:            id,
		Name:          d.Get("name").(string),
		Tag:           d.Get("tag").(string),
		Description:   d.Get("description").(string),
		Type:          d.Get("type").(string),
		IsNameL10nTag: d.Get("is_name_l10n_tag").(bool),
	}
	srcTcpPorts := expandNetworkPorts(d, "src_tcp_ports")
	if srcTcpPorts != nil {
		result.SrcTCPPorts = srcTcpPorts
	}

	destTcpPorts := expandNetworkPorts(d, "dest_tcp_ports")
	if destTcpPorts != nil {
		result.DestTCPPorts = destTcpPorts
	}

	SrcUdpPorts := expandNetworkPorts(d, "src_udp_ports")
	if SrcUdpPorts != nil {
		result.SrcUDPPorts = SrcUdpPorts
	}

	DestUdpPorts := expandNetworkPorts(d, "dest_udp_ports")
	if DestUdpPorts != nil {
		result.DestUDPPorts = DestUdpPorts
	}

	return result
}
