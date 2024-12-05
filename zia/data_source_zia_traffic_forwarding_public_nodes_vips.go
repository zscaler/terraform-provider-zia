package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/virtualipaddress"
)

func dataSourceTrafficForwardingPublicNodeVIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingPublicNodeVIPsRead,
		Schema: map[string]*schema.Schema{
			"cloud_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"city": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gre_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"gre_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pac_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"pac_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTrafficForwardingPublicNodeVIPsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *virtualipaddress.ZscalerVIPs
	datacenter, _ := d.Get("datacenter").(string)
	if datacenter != "" {
		log.Printf("[INFO] Getting data for datacenter name: %s\n", datacenter)
		res, err := virtualipaddress.GetZscalerVIPs(ctx, service, datacenter)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(datacenter)
		_ = d.Set("cloud_name", resp.CloudName)
		_ = d.Set("region", resp.Region)
		_ = d.Set("city", resp.City)
		_ = d.Set("datacenter", resp.DataCenter)
		_ = d.Set("location", resp.Location)
		_ = d.Set("vpn_domain_name", resp.VPNDomainName)
		_ = d.Set("gre_domain_name", resp.GREDomainName)
		_ = d.Set("vpn_ips", resp.VPNIPs)
		_ = d.Set("gre_ips", resp.GREIPs)
		_ = d.Set("pac_ips", resp.PACIPs)
		_ = d.Set("pac_domain_name", resp.PACDomainName)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any datacenter with name '%s'", datacenter))
	}

	return nil
}
