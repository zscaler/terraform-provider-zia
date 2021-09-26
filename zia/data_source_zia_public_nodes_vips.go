package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePublicNodeVirtualAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePublicNodeVirtualAddressRead,
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
			"data_center": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_ips": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vpn_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gre_ips": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gre_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pac_ips": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pac_domain_name": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePublicNodeVirtualAddressRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.publicnodevips.GetPublicNodeVipAddresses()
	if err != nil {
		return nil
	}
	_ = d.Set("cloud_name", resp.CloudName)
	_ = d.Set("region", resp.Region)
	_ = d.Set("city", resp.City)
	_ = d.Set("dataCenter", resp.DataCenter)
	_ = d.Set("location", resp.Location)
	_ = d.Set("vpn_domain_name", resp.VpnDomainName)
	_ = d.Set("gre_domain_name", resp.GreDomainName)
	_ = d.Set("vpn_ips", resp.VpnIps)
	_ = d.Set("gre_ips", resp.GreIps)
	_ = d.Set("pac_ips", resp.PacIps)
	_ = d.Set("pac_domain_name", resp.PacDomainName)

	return nil
}
