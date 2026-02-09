package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dc_exclusions"
)

func dataSourceDatacenters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatacentersRead,
		Schema: map[string]*schema.Schema{
			"datacenter_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Filter datacenters by ID. When exactly one result is returned, set to that datacenter's ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Filter datacenters by name (case-insensitive partial match). When exactly one result is returned, set to that datacenter's name.",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Filter datacenters by city (case-insensitive partial match). When exactly one result is returned, set to that datacenter's city.",
			},
			"datacenters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all datacenters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier for the datacenter.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zscaler data center name.",
						},
						"provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provider of the datacenter.",
						},
						"city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "City where the datacenter is located.",
						},
						"timezone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timezone of the datacenter.",
						},
						"lat": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Latitude coordinate (legacy field).",
						},
						"longi": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Longitude coordinate (legacy field).",
						},
						"latitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Latitude coordinate.",
						},
						"longitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Longitude coordinate.",
						},
						"gov_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is a government-only datacenter.",
						},
						"third_party_cloud": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is a third-party cloud datacenter.",
						},
						"upload_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upload bandwidth in bytes per second.",
						},
						"download_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Download bandwidth in bytes per second.",
						},
						"owned_by_customer": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the datacenter is owned by the customer.",
						},
						"managed_bcp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the datacenter is managed by BCP.",
						},
						"dont_publish": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the datacenter should not be published.",
						},
						"dont_provision": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the datacenter should not be provisioned.",
						},
						"not_ready_for_use": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the datacenter is not ready for use.",
						},
						"for_future_use": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the datacenter is reserved for future use.",
						},
						"regional_surcharge": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether there is a regional surcharge for this datacenter.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp when the datacenter was created.",
						},
						"last_modified_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp when the datacenter was last modified.",
						},
						"virtual": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is a virtual datacenter.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDatacentersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all datacenters")
	allDatacenters, err := dc_exclusions.GetDatacenters(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting datacenters: %w", err))
	}

	log.Printf("[DEBUG] Retrieved %d datacenters from API", len(allDatacenters))

	// Apply simple filtering
	filtered := allDatacenters
	filterID, hasID := d.GetOk("datacenter_id")
	filterName, hasName := d.GetOk("name")
	filterCity, hasCity := d.GetOk("city")

	if hasID || hasName || hasCity {
		var result []dc_exclusions.Datacenter
		for _, dc := range allDatacenters {
			matched := true

			if hasID && dc.ID != filterID.(int) {
				matched = false
			}

			if matched && hasName {
				nameStr := filterName.(string)
				if nameStr != "" && !strings.Contains(strings.ToLower(dc.Name), strings.ToLower(nameStr)) {
					matched = false
				}
			}

			if matched && hasCity {
				cityStr := filterCity.(string)
				if cityStr != "" && !strings.Contains(strings.ToLower(dc.City), strings.ToLower(cityStr)) {
					matched = false
				}
			}

			if matched {
				result = append(result, dc)
			}
		}
		filtered = result
		log.Printf("[INFO] Found %d datacenters matching filter criteria", len(filtered))
	}

	// Use numeric datacenter ID when exactly one result; otherwise filter-based id for state identity
	if len(filtered) == 1 {
		d.SetId(fmt.Sprintf("%d", filtered[0].ID))
		_ = d.Set("datacenter_id", filtered[0].ID)
		_ = d.Set("name", filtered[0].Name)
		_ = d.Set("city", filtered[0].City)
	} else {
		if hasID {
			d.SetId(fmt.Sprintf("id-%d", filterID.(int)))
			_ = d.Set("datacenter_id", filterID.(int))
		} else if hasName {
			d.SetId(fmt.Sprintf("name-%s", filterName.(string)))
			_ = d.Set("name", filterName.(string))
		} else if hasCity {
			d.SetId(fmt.Sprintf("city-%s", filterCity.(string)))
			_ = d.Set("city", filterCity.(string))
		} else {
			d.SetId("all")
		}
	}

	// Flatten and set the datacenters
	if err := d.Set("datacenters", flattenDatacenters(filtered)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting datacenters: %w", err))
	}

	return nil
}

// flattenDatacenters converts the SDK datacenter structs to Terraform schema format
func flattenDatacenters(datacenters []dc_exclusions.Datacenter) []interface{} {
	if len(datacenters) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(datacenters))
	for _, dc := range datacenters {
		m := map[string]interface{}{
			"id":                 dc.ID,
			"name":               dc.Name,
			"provider":           dc.Provider,
			"city":               dc.City,
			"timezone":           dc.Timezone,
			"lat":                dc.Lat,
			"longi":              dc.Longi,
			"latitude":           dc.Latitude,
			"longitude":          dc.Longitude,
			"gov_only":           dc.GovOnly,
			"third_party_cloud":  dc.ThirdPartyCloud,
			"upload_bandwidth":   dc.UploadBandwidth,
			"download_bandwidth": dc.DownloadBandwidth,
			"owned_by_customer":  dc.OwnedByCustomer,
			"managed_bcp":        dc.ManagedBcp,
			"dont_publish":       dc.DontPublish,
			"dont_provision":     dc.DontProvision,
			"not_ready_for_use":  dc.NotReadyForUse,
			"for_future_use":     dc.ForFutureUse,
			"regional_surcharge": dc.RegionalSurcharge,
			"create_time":        dc.CreateTime,
			"last_modified_time": dc.LastModifiedTime,
			"virtual":            dc.Virtual,
		}
		out = append(out, m)
	}

	return out
}
