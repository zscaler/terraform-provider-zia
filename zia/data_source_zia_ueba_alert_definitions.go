package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_ueba_alerts/alert_definitions"
)

func dataSourceUEBAAlertDefinitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUEBAAlertDefinitionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The system-generated identifier of the alert definition.",
			},
			"alert_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The alert name that identifies the threat or event type the alert is triggered for.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the alert rule.",
			},
			"occurrence": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the occurrence of an ongoing alert for a specific threat or event type.",
			},
			"traffic_change_percent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the percentage change in traffic.",
			},
			"interval": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time span within which an event's occurrence triggers an alert.",
			},
			"scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies if the alert is triggered for a user, location, department, or organization.",
			},
			"severity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The threat severity based on which the alert is triggered.",
			},
			"comments": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the triggered alert.",
			},
			"entity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An immutable reference to the entity the alert is scoped to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A unique identifier for the entity.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity.",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceUEBAAlertDefinitionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all UEBA alert definitions\n")
	all, err := alert_definitions.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all UEBA alert definitions: %s", err))
	}

	var resp *alert_definitions.AlertDefinitions
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("alert_name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	if idProvided {
		for i := range all {
			if all[i].ID == id {
				resp = &all[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting UEBA alert definition by ID %d: not found", id))
		}
	}

	if resp == nil && nameProvided && nameStr != "" {
		for i := range all {
			if all[i].AlertName == nameStr {
				resp = &all[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting UEBA alert definition by name %s: not found", nameStr))
		}
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("either 'id' or 'alert_name' must be provided"))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("alert_name", resp.AlertName)
	_ = d.Set("status", resp.Status)
	_ = d.Set("occurrence", resp.Occurrence)
	_ = d.Set("traffic_change_percent", resp.TrafficChangePercent)
	_ = d.Set("interval", resp.Interval)
	_ = d.Set("scope", resp.Scope)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("comments", resp.Comments)
	if err := d.Set("entity", flattenUEBAAlertEntity(resp.Entity)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] UEBA alert definition found: ID=%d, Name=%s\n", resp.ID, resp.AlertName)
	return nil
}

func flattenUEBAAlertEntity(entity *common.IDNameExtensions) []interface{} {
	if entity == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"id":         entity.ID,
			"name":       entity.Name,
			"extensions": entity.Extensions,
		},
	}
}
