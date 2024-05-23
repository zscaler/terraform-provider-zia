package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/timewindow"
)

func dataSourceFWTimeWindow() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWTimeWindowRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"day_of_week": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceFWTimeWindowRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *timewindow.TimeWindow
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting time window : %s\n", name)
		res, err := zClient.timewindow.GetTimeWindowByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("start_time", resp.StartTime)
		_ = d.Set("end_time", resp.EndTime)
		_ = d.Set("day_of_week", resp.DayOfWeek)

	} else {
		return fmt.Errorf("couldn't find any time window with name '%s'", name)
	}

	return nil
}
