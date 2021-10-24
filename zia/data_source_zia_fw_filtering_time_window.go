package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/timewindow"
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
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"day_of_week": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func dataSourceFWTimeWindowRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *timewindow.TimeWindow
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting time window id: %d\n", id)
		res, err := zClient.timewindow.GetTimeWindow(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
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
		return fmt.Errorf("couldn't find any time window with name '%s' or id '%d'", name, id)
	}

	return nil
}
