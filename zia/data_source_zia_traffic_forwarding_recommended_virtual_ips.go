package zia

/*
import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGreVirtualIPAddressesList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGreVirtualIPAddressesListRead,
		Schema: map[string]*schema.Schema{
			"virtual_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_service_edge": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGreVirtualIPAddressesListRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.virtualipaddresslist.GetZSGREVirtualIPList()
	if err != nil {
		return nil
	}

	d.SetId(resp.ID)
	_ = d.Set("virtual_ip", resp.VirtualIp)
	_ = d.Set("private_service_edge", resp.PrivateServiceEdge)
	_ = d.Set("datacenter", resp.DataCenter)

	return nil
}
*/
