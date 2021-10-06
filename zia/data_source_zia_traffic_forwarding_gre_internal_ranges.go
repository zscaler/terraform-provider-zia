package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTrafficForwardingGreInternalIPRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingGreInternalIPRangesRead,
		Schema: map[string]*schema.Schema{
			"start_ip_address": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceTrafficForwardingGreInternalIPRangesRead(d *schema.ResourceData, m interface{}) error {

	return nil
}
