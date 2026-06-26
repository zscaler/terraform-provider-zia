package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/http_header_control/http_header_profile"
)

func dataSourceHttpHeaderProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHttpHeaderProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for the HTTP header profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The HTTP header profile name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the HTTP header profile.",
			},
			"slot_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The slot ID assigned to the HTTP header profile.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the HTTP header profile is deleted.",
			},
			"profile_ready_for_use": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the HTTP header profile is ready for use.",
			},
			"http_header_profile_criteria": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of matching criteria evaluated by the HTTP header profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies the criteria entry.",
						},
						"header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The header evaluated by the criteria.",
						},
						"operator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operator applied to the header criteria.",
						},
						"user_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user agent evaluated by the criteria.",
						},
						"user_agent_bitmap": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user agent bitmap evaluated by the criteria.",
						},
						"user_agent_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user agent version evaluated by the criteria.",
						},
						"category_bitmap": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The URL category bitmap evaluated by the criteria.",
						},
						"cloud_app_bitmap": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The cloud application bitmap evaluated by the criteria.",
						},
					},
				},
			},
		},
	}
}

func dataSourceHttpHeaderProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *http_header_profile.HttpHeaderProfile
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	// The API does not expose a per-ID lookup endpoint, so an ID lookup fetches
	// the full list and matches locally. A name lookup uses the SDK's dedicated
	// name search (case-insensitive), matching the other data sources.
	if idProvided {
		log.Printf("[INFO] Searching for HTTP header profile by ID: %d\n", id)
		allProfiles, err := http_header_profile.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting all HTTP header profiles: %s", err))
		}
		for i := range allProfiles {
			if allProfiles[i].ID == id {
				resp = &allProfiles[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting HTTP header profile by ID %d: profile not found", id))
		}
	}

	if resp == nil && nameStr != "" {
		log.Printf("[INFO] Searching for HTTP header profile by name: %s\n", nameStr)
		profile, err := http_header_profile.GetByName(ctx, service, nameStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting HTTP header profile by name %s: %s", nameStr, err))
		}
		resp = profile
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("slot_id", resp.SlotId)
	_ = d.Set("deleted", resp.Deleted)
	_ = d.Set("profile_ready_for_use", resp.ProfileReadyForUse)

	if err := d.Set("http_header_profile_criteria", flattenHttpHeaderProfileCriteria(resp.HttpHeaderProfileCriteria)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] HTTP header profile found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}

func flattenHttpHeaderProfileCriteria(criteria []http_header_profile.HttpHeaderProfileCriteria) []interface{} {
	if len(criteria) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(criteria))
	for _, c := range criteria {
		out = append(out, map[string]interface{}{
			"id":                 c.Id,
			"header":             c.Header,
			"operator":           c.Operator,
			"user_agent":         c.UserAgent,
			"user_agent_bitmap":  c.UserAgentBitmap,
			"user_agent_version": c.UserAgentVersion,
			"category_bitmap":    c.CategoryBitmap,
			"cloud_app_bitmap":   c.CloudAppBitmap,
		})
	}
	return out
}
