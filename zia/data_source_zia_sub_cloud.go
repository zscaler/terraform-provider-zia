package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/sub_clouds"
)

func dataSourceSubCloud() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubCloudRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Unique identifier for the subcloud. Used to look up a single subcloud when provided.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Subcloud name. Used to search for a subcloud by name when provided.",
			},
			"dcs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of data centers associated with the subcloud.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier for the data center.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data center name.",
						},
						"country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country where the data center is located.",
						},
					},
				},
			},
			"exclusions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of data centers excluded from the subcloud.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The excluded datacenter reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unique identifier for the entity.",
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
										Description: "Extension attributes.",
									},
								},
							},
						},
						"last_modified_user": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The admin that last modified the exclusion.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unique identifier for the entity.",
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
										Description: "Extension attributes.",
									},
								},
							},
						},
						"country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Country where the excluded data center is located.",
						},
						"expired": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the exclusion has expired.",
						},
						"disabled_by_ops": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the exclusion was disabled by operations.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp when the exclusion was created.",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Exclusion start time (Unix timestamp).",
						},
						"start_time_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exclusion start time (UTC). Format: MM/DD/YYYY HH:MM am/pm.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Exclusion end time (Unix timestamp).",
						},
						"end_time_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exclusion end time (UTC). Format: MM/DD/YYYY HH:MM am/pm.",
						},
						"last_modified_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp when the exclusion was last modified.",
						},
					},
				},
			},
		},
	}
}

func dataSourceSubCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all subclouds")
	allSubClouds, err := sub_clouds.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all subclouds: %w", err))
	}

	log.Printf("[DEBUG] Retrieved %d subclouds", len(allSubClouds))

	var resp *sub_clouds.SubClouds
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	if idProvided {
		log.Printf("[INFO] Searching for subcloud by ID: %d", id)
		for i := range allSubClouds {
			if allSubClouds[i].ID == id {
				resp = &allSubClouds[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("subcloud with ID %d not found", id))
		}
	}

	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for subcloud by name: %s", nameStr)
		for i := range allSubClouds {
			if allSubClouds[i].Name == nameStr {
				resp = &allSubClouds[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("subcloud with name %q not found", nameStr))
		}
	}

	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("subcloud not found with name %q or id %d", nameStr, id))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	if err := d.Set("id", resp.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %w", err))
	}
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("dcs", flattenSubCloudDCs(resp.Dcs)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting dcs: %w", err))
	}
	if err := d.Set("exclusions", flattenSubCloudExclusions(resp.Exclusions)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting exclusions: %w", err))
	}

	log.Printf("[DEBUG] Subcloud found: ID=%d, Name=%s", resp.ID, resp.Name)
	return nil
}

// flattenSubCloudDCs flattens the Dcs slice from the SDK into the schema list.
func flattenSubCloudDCs(dcs []sub_clouds.DCs) []interface{} {
	if len(dcs) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(dcs))
	for _, dc := range dcs {
		out = append(out, map[string]interface{}{
			"id":      dc.ID,
			"name":    dc.Name,
			"country": dc.Country,
		})
	}
	return out
}

// flattenSubCloudExclusions flattens the Exclusions slice from the SDK into the schema list.
func flattenSubCloudExclusions(exclusions []sub_clouds.Exclusions) []interface{} {
	if len(exclusions) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(exclusions))
	for _, e := range exclusions {
		m := map[string]interface{}{
			"country":            e.Country,
			"expired":            e.Expired,
			"disabled_by_ops":    e.DisabledByOps,
			"create_time":        e.CreateTime,
			"start_time":         e.StartTime,
			"end_time":           e.EndTime,
			"start_time_utc":     FormatExclusionTimeUTC(e.StartTime),
			"end_time_utc":       FormatExclusionTimeUTC(e.EndTime),
			"last_modified_time": e.LastModifiedTime,
		}
		m["datacenter"] = flattenIDExtensionsList(e.Datacenter)
		m["last_modified_user"] = flattenLastModifiedBy(e.LastModifiedUser)
		out = append(out, m)
	}
	return out
}
