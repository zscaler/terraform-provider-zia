package zia

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/sub_clouds"
)

func resourceSubCloud() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSubCloudRead,
		CreateContext: resourceSubCloudCreate,
		UpdateContext: resourceSubCloudUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("cloud_id", idInt)
				}
				diags := resourceSubCloudRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read subcloud import: %s", diags[0].Summary)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the subcloud as a string (Terraform resource ID).",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier for the subcloud as an integer. Used as the path parameter for the PUT API.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Subcloud name. This attribute is read-only and cannot be updated via the API.",
			},
			"exclusions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of data centers excluded from the subcloud.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The excluded datacenter reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Unique identifier for the datacenter.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Datacenter name.",
									},
									"country": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Country where the datacenter is located.",
									},
								},
							},
						},
						"country": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Country where the excluded data center is located.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Exclusion end time (Unix timestamp). Either end_time or end_time_utc must be set.",
						},
						"end_time_utc": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Exclusion end time (UTC). Format: MM/DD/YYYY HH:MM:SS am/pm. If set, overrides end_time.",
							ValidateFunc: ValidateExclusionTimeUTC,
						},
					},
				},
			},
			"dcs": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Set of data centers associated with the subcloud (read-only).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSubCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	cloudID, ok := getIntFromResourceData(d, "cloud_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("cloud_id is required for subcloud creation/update"))
	}

	// Warn if user provides name (which will be ignored as it cannot be set via API)
	if name, ok := d.GetOk("name"); ok && name.(string) != "" {
		log.Printf("[WARN] The 'name' attribute cannot be set via the API. The name from the API response will be used instead.")
	}

	req := expandSubCloud(d)
	req.ID = cloudID

	log.Printf("[INFO] Creating/updating subcloud id=%d", cloudID)
	_, _, err := sub_clouds.Update(ctx, service, cloudID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating subcloud: %w", err))
	}

	d.SetId(strconv.Itoa(cloudID))
	_ = d.Set("cloud_id", cloudID)
	return resourceSubCloudRead(ctx, d, meta)
}

func resourceSubCloudRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "cloud_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no cloud_id is set"))
	}
	allSubClouds, err := sub_clouds.GetAll(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting subclouds: %w", err))
	}

	var resp *sub_clouds.SubClouds
	for i := range allSubClouds {
		if allSubClouds[i].ID == id {
			resp = &allSubClouds[i]
			break
		}
	}
	if resp == nil {
		d.SetId("")
		log.Printf("[WARN] Subcloud id=%d not found, removing from state", id)
		return nil
	}

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("cloud_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("dcs", flattenSubCloudDCs(resp.Dcs))
	_ = d.Set("exclusions", flattenSubCloudExclusionsForResource(resp.Exclusions))
	return nil
}

func resourceSubCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	cloudID, ok := getIntFromResourceData(d, "cloud_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no cloud_id is set"))
	}

	// Warn if user tries to update name (which cannot be updated)
	if d.HasChange("name") {
		log.Printf("[WARN] The 'name' attribute cannot be updated via the API. The name from the API response will be used instead.")
	}

	req := expandSubCloud(d)
	req.ID = cloudID

	log.Printf("[INFO] Updating subcloud id=%d", cloudID)
	_, _, err := sub_clouds.Update(ctx, service, cloudID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating subcloud: %w", err))
	}

	return resourceSubCloudRead(ctx, d, meta)
}

func expandSubCloud(d *schema.ResourceData) *sub_clouds.SubClouds {
	id, _ := getIntFromResourceData(d, "cloud_id")
	req := &sub_clouds.SubClouds{
		ID: id,
		// Name is intentionally omitted - it cannot be updated via the API
		// The API will return the current name in the response
		Exclusions: expandSubCloudExclusions(d.Get("exclusions")),
	}
	return req
}

func expandSubCloudExclusions(v interface{}) []sub_clouds.Exclusions {
	if v == nil {
		return nil
	}
	set, ok := v.(*schema.Set)
	if !ok {
		return nil
	}
	list := set.List()
	if len(list) == 0 {
		return nil
	}
	out := make([]sub_clouds.Exclusions, 0, len(list))
	for _, item := range list {
		raw, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		endTime, err := resolveExclusionTimeFromMap(raw, "end_time", "end_time_utc")
		if err != nil {
			log.Printf("[WARN] sub_cloud exclusion end_time: %v", err)
			continue
		}

		e := sub_clouds.Exclusions{
			Country: getStringFromMap(raw, "country"),
			EndTime: endTime,
		}
		if dcList, ok := raw["datacenter"].([]interface{}); ok && len(dcList) > 0 {
			if dcMap, ok := dcList[0].(map[string]interface{}); ok {
				dc := &common.IDNameExtensions{
					ID:   getIntFromMap(dcMap, "id"),
					Name: getStringFromMap(dcMap, "name"),
				}
				if country := getStringFromMap(dcMap, "country"); country != "" {
					dc.Extensions = map[string]interface{}{"country": country}
				}
				e.Datacenter = dc
			}
		}
		out = append(out, e)
	}
	return out
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getIntFromMap(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch t := v.(type) {
		case int:
			return t
		case int64:
			return int(t)
		case float64:
			return int(t)
		}
	}
	return 0
}

// resolveExclusionTimeFromMap uses the shared ResolveExclusionTimeEpoch util for consistency with zia_dc_exclusions.
func resolveExclusionTimeFromMap(raw map[string]interface{}, epochKey, utcKey string) (int, error) {
	_, hasEpoch := raw[epochKey]
	epoch := getIntFromMap(raw, epochKey)
	utcStr := getStringFromMap(raw, utcKey)
	return ResolveExclusionTimeEpoch(hasEpoch, epoch, utcStr)
}

func flattenSubCloudExclusionsForResource(exclusions []sub_clouds.Exclusions) []interface{} {
	if len(exclusions) == 0 {
		return nil
	}
	// Sort by datacenter ID so Read produces a stable order and matches config order when config lists by same key.
	exc := make([]sub_clouds.Exclusions, len(exclusions))
	copy(exc, exclusions)
	sort.Slice(exc, func(i, j int) bool {
		idI, idJ := 0, 0
		if exc[i].Datacenter != nil {
			idI = exc[i].Datacenter.ID
		}
		if exc[j].Datacenter != nil {
			idJ = exc[j].Datacenter.ID
		}
		return idI < idJ
	})
	out := make([]interface{}, 0, len(exc))
	for _, e := range exc {
		m := map[string]interface{}{
			"country":      e.Country,
			"end_time_utc": FormatExclusionTimeUTC(e.EndTime),
		}
		m["datacenter"] = flattenSubCloudExclusionDatacenterForResource(e.Datacenter)
		out = append(out, m)
	}
	return out
}

// flattenSubCloudExclusionDatacenterForResource returns a single-element list with id and name only so state matches config and full-object hash does not drift.
func flattenSubCloudExclusionDatacenterForResource(dc *common.IDNameExtensions) []interface{} {
	if dc == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{"id": dc.ID, "name": dc.Name},
	}
}
