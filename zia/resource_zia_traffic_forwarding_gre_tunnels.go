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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/virtualipaddress"
)

func resourceTrafficForwardingGRETunnel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficForwardingGRETunnelCreate,
		ReadContext:   resourceTrafficForwardingGRETunnelRead,
		UpdateContext: resourceTrafficForwardingGRETunnelUpdate,
		DeleteContext: resourceTrafficForwardingGRETunnelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("tunnel_id", idInt)
				} else {
					resp, err := gretunnels.GetByIPAddress(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("tunnel_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"tunnel_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the GRE tunnel.",
			},
			"source_ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The source IP address of the GRE tunnel. This is typically a static IP address in the organization or SD-WAN.",
				ValidateFunc: validation.IsIPAddress,
			},
			"within_country": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Restrict the data center virtual IP addresses (VIPs) only to those within the same country as the source IP address",
			},
			"primary_dest_vip": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The primary destination data center and virtual IP address (VIP) of the GRE tunnel",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "GRE cluster virtual IP ID",
						},
						"virtual_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "GRE cluster virtual IP address (VIP)",
							ValidateFunc: validation.IsIPAddress,
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Data center information",
						},
					},
				},
			},
			"secondary_dest_vip": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The secondary destination data center and virtual IP address (VIP) of the GRE tunnel",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "GRE cluster virtual IP ID",
						},
						"virtual_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "GRE cluster virtual IP address (VIP)",
							ValidateFunc: validation.IsIPAddress,
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Data center information",
						},
					},
				},
			},
			"internal_ip_range": {
				Type: schema.TypeString,
				// Computed:     true,
				Optional:     true,
				Description:  "The start of the internal IP address in /29 CIDR range",
				ValidateFunc: validation.IsIPv4Address,
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "When within_country is enabled, you must set this to the country code.",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional information about this GRE tunnel",
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"ip_unnumbered": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "This is required to support the automated SD-WAN provisioning of GRE tunnels, when set to true gre_tun_ip and gre_tun_id are set to null",
			},
		},
	}
}

func resourceTrafficForwardingGRETunnelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandGRETunnel(d)
	log.Printf("[INFO] Creating zia gre tunnel\n%+v\n", req)

	// Handle asssignVipsIfNotSet error
	if err := asssignVipsIfNotSet(ctx, d, zClient, &req); err != nil {
		return diag.Errorf("error assigning VIPs: %v", err)
	}

	// Create GRE Tunnel
	resp, _, createErr := gretunnels.CreateGreTunnels(ctx, service, &req)
	if createErr != nil {
		return diag.Errorf("error creating GRE tunnel: %v", createErr)
	}

	log.Printf("[INFO] Created zia gre tunnel request. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("tunnel_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.Errorf("error triggering activation: %v", activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceTrafficForwardingGRETunnelRead(ctx, d, meta)
}

func asssignVipsIfNotSet(ctx context.Context, d *schema.ResourceData, zClient *Client, req *gretunnels.GreTunnels) diag.Diagnostics {
	service := zClient.Service

	if (req.PrimaryDestVip == nil || (req.PrimaryDestVip.VirtualIP == "" && req.PrimaryDestVip.ID == 0)) ||
		(req.SecondaryDestVip == nil || (req.SecondaryDestVip.VirtualIP == "" && req.SecondaryDestVip.ID == 0)) {
		// one of the vips not set, pick 2 from the recommandedVips
		countryCode, ok := getStringFromResourceData(d, "country_code")
		var pair []virtualipaddress.GREVirtualIPList
		if ok {
			vips, err := virtualipaddress.GetPairZSGREVirtualIPsWithinCountry(ctx, service, req.SourceIP, countryCode)
			if err != nil {
				log.Printf("[ERROR] Got: %v\n", err)
				vips, err = virtualipaddress.GetZSGREVirtualIPList(ctx, service, req.SourceIP, 2)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			pair = *vips
		} else {
			vips, err := virtualipaddress.GetZSGREVirtualIPList(ctx, service, req.SourceIP, 2)
			if err != nil {
				return diag.FromErr(err)
			}
			pair = *vips
		}
		req.PrimaryDestVip = &gretunnels.PrimaryDestVip{ID: pair[0].ID, VirtualIP: pair[0].VirtualIp, Datacenter: pair[0].DataCenter}
		req.SecondaryDestVip = &gretunnels.SecondaryDestVip{ID: pair[1].ID, VirtualIP: pair[1].VirtualIp, Datacenter: pair[0].DataCenter}
	}
	return nil
}

func resourceTrafficForwardingGRETunnelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "tunnel_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no Traffic Forwarding GRE Tunnel id is set"))
	}
	resp, err := gretunnels.GetGreTunnels(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing gre tunnel %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting gre tunnel:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("tunnel_id", resp.ID)
	_ = d.Set("source_ip", resp.SourceIP)
	_ = d.Set("internal_ip_range", resp.InternalIpRange)
	if resp.WithinCountry != nil {
		_ = d.Set("within_country", *resp.WithinCountry)
	}
	_ = d.Set("comment", resp.Comment)
	_ = d.Set("ip_unnumbered", resp.IPUnnumbered)
	if err := d.Set("primary_dest_vip", flattenGrePrimaryDestVipSimple(resp.PrimaryDestVip)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("secondary_dest_vip", flattenGreSecondaryDestVipSimple(resp.SecondaryDestVip)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenGrePrimaryDestVipSimple(primaryDestVip *gretunnels.PrimaryDestVip) interface{} {
	return []map[string]interface{}{
		{
			"id":         primaryDestVip.ID,
			"virtual_ip": primaryDestVip.VirtualIP,
			"datacenter": primaryDestVip.Datacenter,
		},
	}
}

func flattenGreSecondaryDestVipSimple(secondaryDestVip *gretunnels.SecondaryDestVip) interface{} {
	return []map[string]interface{}{
		{
			"id":         secondaryDestVip.ID,
			"virtual_ip": secondaryDestVip.VirtualIP,
			"datacenter": secondaryDestVip.Datacenter,
		},
	}
}

func resourceTrafficForwardingGRETunnelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "tunnel_id")
	if !ok {
		log.Printf("[ERROR] gre tunnel ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating gre tunnel ID: %v\n", id)
	req := expandGRETunnel(d)

	err := asssignVipsIfNotSet(ctx, d, zClient, &req)
	if err != nil {
		return diag.Errorf("error assigning VIPs: %v", err)
	}
	if _, err := gretunnels.GetGreTunnels(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := gretunnels.UpdateGreTunnels(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}
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

	return resourceTrafficForwardingGRETunnelRead(ctx, d, meta)
}

func resourceTrafficForwardingGRETunnelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "tunnel_id")
	if !ok {
		log.Printf("[ERROR] gre tunnel ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting gre tunnel ID: %v\n", id)

	if _, err := gretunnels.DeleteGreTunnels(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] gre tunnel deleted")
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

	return nil
}

func expandGRETunnel(d *schema.ResourceData) gretunnels.GreTunnels {
	id, _ := getIntFromResourceData(d, "tunnel_id")
	withinCountry := d.Get("within_country").(bool)
	result := gretunnels.GreTunnels{
		ID:              id,
		SourceIP:        d.Get("source_ip").(string),
		InternalIpRange: d.Get("internal_ip_range").(string),
		WithinCountry:   &withinCountry,
		Comment:         d.Get("comment").(string),
		IPUnnumbered:    d.Get("ip_unnumbered").(bool),
	}
	primaryDestVip := expandPrimaryDestVip(d)
	if primaryDestVip != nil {
		result.PrimaryDestVip = primaryDestVip
	}
	secondaryDestVip := expandSecondaryDestVip(d)
	if secondaryDestVip != nil {
		result.SecondaryDestVip = secondaryDestVip
	}
	return result
}

func expandPrimaryDestVip(d *schema.ResourceData) *gretunnels.PrimaryDestVip {
	vipsObj, ok := d.GetOk("primary_dest_vip")
	if !ok {
		return nil
	}
	vips, ok := vipsObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(vips.List()) > 0 {
		vipObj := vips.List()[0]
		vip, ok := vipObj.(map[string]interface{})
		if !ok {
			return nil
		}
		r := &gretunnels.PrimaryDestVip{
			VirtualIP:  vip["virtual_ip"].(string),
			Datacenter: vip["datacenter"].(string),
		}

		if id, ok := vip["id"].(int); ok && id != 0 {
			r.ID = id
		}
		return r
	}
	return nil
}

func expandSecondaryDestVip(d *schema.ResourceData) *gretunnels.SecondaryDestVip {
	vipsObj, ok := d.GetOk("secondary_dest_vip")
	if !ok {
		return nil
	}
	vips, ok := vipsObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(vips.List()) > 0 {
		vipObj := vips.List()[0]
		vip, ok := vipObj.(map[string]interface{})
		if !ok {
			return nil
		}
		r := &gretunnels.SecondaryDestVip{
			VirtualIP:  vip["virtual_ip"].(string),
			Datacenter: vip["datacenter"].(string),
		}
		if id, ok := vip["id"].(int); ok && id != 0 {
			r.ID = id
		}

		return r
	}
	return nil
}
