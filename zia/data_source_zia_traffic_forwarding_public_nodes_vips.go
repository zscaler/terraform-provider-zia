package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/virtualipaddresslist"
)

func dataSourceTrafficForwardingPublicNodeVIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingPublicNodeVIPsRead,
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

func dataSourceTrafficForwardingPublicNodeVIPsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *virtualipaddresslist.ZscalerVIPs
	datacenter, _ := d.Get("datacenter").(string)
	if resp == nil && datacenter != "" {
		log.Printf("[INFO] Getting data for datacenter name: %s\n", datacenter)
		res, err := zClient.virtualipaddresslist.GetZscalerVIPs(datacenter)
		if err != nil {
			return err
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
		return fmt.Errorf("couldn't find any datacenter with name '%s'", datacenter)
	}

	return nil
}
