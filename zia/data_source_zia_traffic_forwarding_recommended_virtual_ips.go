package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/virtualipaddresslist"
)

func dataSourceGreVirtualIPAddressesList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGreVirtualIPAddressesListRead,
		Schema: map[string]*schema.Schema{
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"geo_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"required_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"virtual_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_service_edge": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGreVirtualIPAddressesListRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	count, ok := getIntFromResourceData(d, "required_count")
	if !ok {
		count = 1
	}
	sourceIP, ok := getStringFromResourceData(d, "source_ip")
	if !ok {
		return fmt.Errorf("please provide a source_ip for the vips list")
	}
	resp, err := zClient.virtualipaddresslist.GetZSGREVirtualIPList(sourceIP, count)
	if err != nil {
		return err
	}
	d.SetId(sourceIP)
	_ = d.Set("list", flattenVIPList(*resp))

	return nil
}

func flattenVIPList(list []virtualipaddresslist.GREVirtualIPList) []interface{} {
	result := make([]interface{}, len(list))
	for i, vip := range list {
		result[i] = map[string]interface{}{
			"id":                   vip.ID,
			"virtual_ip":           vip.VirtualIp,
			"private_service_edge": vip.PrivateServiceEdge,
			"datacenter":           vip.DataCenter,
		}
	}
	return result
}
