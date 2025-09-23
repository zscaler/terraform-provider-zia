package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/timewindow"
)

func dataSourceFWTimeWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFWTimeWindowRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
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

func dataSourceFWTimeWindowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *timewindow.TimeWindow
	var searchCriteria string

	// Check if searching by ID
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting time window by id: %d\n", id)
		searchCriteria = fmt.Sprintf("id=%d", id)

		// Get all time windows and find the one with matching ID
		allTimeWindows, err := timewindow.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, tw := range allTimeWindows {
			if tw.ID == id {
				resp = &tw
				break
			}
		}
	}

	// Check if searching by name (only if ID search didn't find anything)
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting time window by name: %s\n", name)
		searchCriteria = fmt.Sprintf("name=%s", name)

		res, err := timewindow.GetTimeWindowByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("couldn't find any time window with %s", searchCriteria))
	}

	return nil
}
