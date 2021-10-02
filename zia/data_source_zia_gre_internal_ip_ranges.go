package zia

/*
import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/greinternalipranges"
)

func dataSourceGreInternalIPRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGreInternalIPRangesRead,
		Schema: map[string]*schema.Schema{
			"start_ip_address": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGreInternalIPRangesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *greinternalipranges.GREInternalIPRanges
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for gre tunnel id: %d\n", id)
		res, err := zClient.greinternalipranges.GetGREInternalIPRanges(id)
		if err != nil {
			return err
		}
		resp = res
	}

	_ = d.Set("start_ip_address", resp.StartIPAddress)
	_ = d.Set("end_ip_address", resp.EndIPAddress)

	return nil
}
*/
