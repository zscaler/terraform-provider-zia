package zia

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/greinternalipranges"
)

func dataSourceTrafficForwardingGreInternalIPRangeList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingGreInternalIPRangesRead,
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

func dataSourceTrafficForwardingGreInternalIPRangesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	count, ok := getIntFromResourceData(d, "required_count")
	if !ok {
		count = 1
	}
	ipRanges, err := greinternalipranges.GetGREInternalIPRange(ctx, service, count)
	if err != nil {
		return diag.FromErr(err)
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
