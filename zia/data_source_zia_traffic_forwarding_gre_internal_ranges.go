package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/greinternalipranges"
)

func dataSourceTrafficForwardingGreInternalIPRangeList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingGreInternalIPRangesRead,
		Schema: map[string]*schema.Schema{
			"required_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTrafficForwardingGreInternalIPRangesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	count, ok := getIntFromResourceData(d, "required_count")
	if !ok {
		count = 1
	}
	ipRanges, err := zClient.greinternalipranges.GetGREInternalIPRange(count)
	if err != nil {
		return err
	}
	d.SetId("internal-ip-range-list")
	_ = d.Set("list", flattenInternalIpRangeList(*ipRanges))

	return nil
}

func flattenInternalIpRangeList(list []greinternalipranges.GREInternalIPRange) []interface{} {
	result := make([]interface{}, len(list))
	for i, ip := range list {
		result[i] = map[string]interface{}{
			"start_ip_address": ip.StartIPAddress,
			"end_ip_address":   ip.EndIPAddress,
		}
	}
	return result
}
