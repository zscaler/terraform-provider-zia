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

	// Always fetch all time windows and search locally
	log.Printf("[INFO] Fetching all time windows\n")
	allTimeWindows, err := timewindow.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all time windows: %s", err))
	}

	log.Printf("[DEBUG] Retrieved %d time windows\n", len(allTimeWindows))

	var resp *timewindow.TimeWindow
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	// Search by ID first if provided
	if idProvided {
		log.Printf("[INFO] Searching for time window by ID: %d\n", id)
		for _, tw := range allTimeWindows {
			if tw.ID == id {
				resp = &tw
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting time window by ID %d: time window not found", id))
		}
	}

	// Search by name if not found by ID and name is provided
	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for time window by name: %s\n", nameStr)
		for _, tw := range allTimeWindows {
			if tw.Name == nameStr {
				resp = &tw
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting time window by name %s: time window not found", nameStr))
		}
	}

	// If neither ID nor name provided, or no match found
	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("couldn't find any time window with name '%s' or id '%d'", nameStr, id))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	// Set the resource data
	d.SetId(fmt.Sprintf("%d", resp.ID))
	err = d.Set("name", resp.Name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	err = d.Set("start_time", resp.StartTime)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting start_time: %s", err))
	}
	err = d.Set("end_time", resp.EndTime)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting end_time: %s", err))
	}
	err = d.Set("day_of_week", resp.DayOfWeek)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting day_of_week: %s", err))
	}

	log.Printf("[DEBUG] Time window found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}
