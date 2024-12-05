package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/vpncredentials"
)

func dataSourceTrafficForwardingVPNCredentials() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingVPNCredentialsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pre_shared_key": {
				Type:     schema.TypeString,
				Computed: true,
				// Sensitive: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"managed_by": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTrafficForwardingVPNCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *vpncredentials.VPNCredentials
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for vpn credential id: %d\n", id)
		res, err := vpncredentials.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	fqdn, _ := d.Get("fqdn").(string)
	if resp == nil && fqdn != "" {
		log.Printf("[INFO] Getting data for vpn credential fqdn: %s\n", fqdn)
		res, err := vpncredentials.GetByFQDN(ctx, service, fqdn)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	vpnType, _ := d.Get("type").(string)
	if resp == nil && vpnType != "" {
		log.Printf("[INFO] Getting data for vpn credential type: %s\n", vpnType)
		res, err := vpncredentials.GetVPNByType(ctx, service, vpnType, nil, nil, nil)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(res) > 0 {
			resp = &res[0] // Assuming you want the first result
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("type", resp.Type)
		_ = d.Set("fqdn", resp.FQDN)
		_ = d.Set("ip_address", resp.IPAddress)
		// _ = d.Set("pre_shared_key", resp.PreSharedKey)
		_ = d.Set("comments", resp.Comments)
		if err := d.Set("location", flattenVPNCredentialsLocation(resp.Location)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("managed_by", flattenVPNCredentialsManagedBy(resp.ManagedBy)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any vpn credentials with fqdn '%s' or id '%d'", fqdn, id))
	}

	return nil
}

// Want to simplify this. This flattening function will be used in multiple places.
func flattenVPNCredentialsLocation(location *vpncredentials.Location) interface{} {
	if location == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":         location.ID,
			"name":       location.Name,
			"extensions": location.Extensions,
		},
	}
}

// Want to simplify this. This flattening function will be used in multiple places.
func flattenVPNCredentialsManagedBy(managedBy *vpncredentials.ManagedBy) interface{} {
	if managedBy == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":         managedBy.ID,
			"name":       managedBy.Name,
			"extensions": managedBy.Extensions,
		},
	}
}
