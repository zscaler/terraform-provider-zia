package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/virtualipaddresslist"
)

func dataSourceGreVirtualIPAddressesList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGreVirtualIPAddressesListRead,
		Schema: map[string]*schema.Schema{
			"source_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"geo_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"virtual_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_service_edge": {
							Type:     schema.TypeBool,
							Computed: true,
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
	sourceIP, ok := d.GetOk("source_ip")
	if !ok {
		return fmt.Errorf("please provide a source_ip for the vips list")
	}
	sIP, ok := sourceIP.(string)
	if !ok {
		return fmt.Errorf("please provide a source_ip for the vips list")
	}
	resp, err := zClient.virtualipaddresslist.GetZSGREVirtualIPList(sIP)
	if err != nil {
		return err
	}
	d.SetId(sIP)
	_ = d.Set("list", flattenVIPList(*resp))

	return nil
}

func flattenVIPList(list []virtualipaddresslist.GREVirtualIPList) []interface{} {
	log.Printf("[ERROR] Got: %v\n", list)
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
