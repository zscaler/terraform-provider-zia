package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/gretunnelinfo"
)

func dataSourceTrafficForwardingIPGreTunnelInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingIPGreTunnelInfoRead,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gre_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gre_tunnel_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_gw": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secondary_gw": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tun_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gre_range_primary": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gre_range_secondary": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTrafficForwardingIPGreTunnelInfoRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.gretunnelinfo

	var resp *gretunnelinfo.GRETunnelInfo
	id, ok := getStringFromResourceData(d, "ip_address")
	if ok {
		log.Printf("[INFO] Getting data for gre tunnel id: %s\n", id)
		res, err := gretunnelinfo.GetGRETunnelInfo(service, id)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.TunID))
		_ = d.Set("ip_address", resp.IPaddress)
		_ = d.Set("gre_enabled", resp.GREEnabled)
		_ = d.Set("gre_tunnel_ip", resp.GREtunnelIP)
		_ = d.Set("primary_gw", resp.PrimaryGW)
		_ = d.Set("secondary_gw", resp.SecondaryGW)
		_ = d.Set("tun_id", resp.TunID)
		_ = d.Set("gre_range_primary", resp.GRERangePrimary)
		_ = d.Set("gre_range_secondary", resp.GRERangeSecondary)

	} else {
		return fmt.Errorf("couldn't find any info for gre tunnel with id '%s'", id)
	}

	return nil
}
