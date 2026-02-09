package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dc_exclusions"
)

func dataSourceDCExclusions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDCExclusionsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter exclusions by the resource id (datacenter ID as string). Use this to look up a specific exclusion by its zia_dc_exclusions resource id.",
			},
			"datacenter_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter exclusions by datacenter ID (dcid).",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Filter exclusions by datacenter name (case-insensitive partial match). When filtering by datacenter_id and exactly one exclusion is returned, this is set to that exclusion's datacenter name.",
			},
			"exclusions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of DC exclusion entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The exclusion identifier (datacenter ID as string). Matches the zia_dc_exclusions resource id.",
						},
						"dc_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Datacenter ID (dcid) for the exclusion.",
						},
						"expired": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the exclusion has expired.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp when the exclusion window starts.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unix timestamp when the exclusion window ends.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the DC exclusion.",
						},
						"dc_name_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Datacenter ID from the dcName reference.",
						},
						"dc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Datacenter name from the dcName reference.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDCExclusionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all DC exclusions")
	all, err := dc_exclusions.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting DC exclusions: %w", err))
	}

	log.Printf("[DEBUG] Retrieved %d DC exclusions from API", len(all))

	filtered := all
	filterID := 0
	hasID := false
	if idStr, ok := d.GetOk("id"); ok && idStr.(string) != "" {
		if parsed, err := strconv.Atoi(idStr.(string)); err == nil {
			filterID = parsed
			hasID = true
		}
	}
	if !hasID {
		if dcID, ok := d.GetOk("datacenter_id"); ok {
			filterID = dcID.(int)
			hasID = true
		}
	}
	filterName, hasName := d.GetOk("name")

	if hasID || hasName {
		var result []dc_exclusions.DCExclusions
		for _, ex := range all {
			matched := true
			if hasID && ex.DcID != filterID {
				matched = false
			}
			if matched && hasName {
				nameStr := filterName.(string)
				dcName := ""
				if ex.DcName != nil {
					dcName = ex.DcName.Name
				}
				if nameStr != "" && !strings.Contains(strings.ToLower(dcName), strings.ToLower(nameStr)) {
					matched = false
				}
			}
			if matched {
				result = append(result, ex)
			}
		}
		filtered = result
		log.Printf("[INFO] Found %d DC exclusions matching filter criteria", len(filtered))
		if hasID && len(filtered) == 0 {
			return diag.FromErr(fmt.Errorf("no DC exclusion found with datacenter id %d", filterID))
		}
	}

	// Use the datacenter ID when we have exactly one result; otherwise use a filter-based id for state identity
	if len(filtered) == 1 {
		d.SetId(fmt.Sprintf("%d", filtered[0].DcID))
	} else if hasID {
		d.SetId(fmt.Sprintf("id-%d", filterID))
	} else if hasName {
		d.SetId(fmt.Sprintf("name-%s", filterName.(string)))
	} else {
		d.SetId("all")
	}

	// Set top-level name so it is never null: use filter value when provided, otherwise derive from single result
	if hasName {
		_ = d.Set("name", filterName.(string))
	} else if len(filtered) == 1 && filtered[0].DcName != nil {
		_ = d.Set("name", filtered[0].DcName.Name)
	} else {
		_ = d.Set("name", "")
	}

	if err := d.Set("exclusions", flattenDCExclusions(filtered)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting exclusions: %w", err))
	}

	return nil
}

func flattenDCExclusions(list []dc_exclusions.DCExclusions) []interface{} {
	if len(list) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(list))
	for _, ex := range list {
		m := map[string]interface{}{
			"id":          fmt.Sprintf("%d", ex.DcID),
			"dc_id":       ex.DcID,
			"expired":     ex.Expired,
			"start_time":  ex.StartTime,
			"end_time":    ex.EndTime,
			"description": ex.Description,
		}
		if ex.DcName != nil {
			m["dc_name_id"] = ex.DcName.ID
			m["dc_name"] = ex.DcName.Name
		} else {
			m["dc_name_id"] = 0
			m["dc_name"] = ""
		}
		out = append(out, m)
	}
	return out
}
