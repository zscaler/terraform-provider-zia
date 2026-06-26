package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/http_header_control/http_header_action_profile"
)

func dataSourceHttpHeaderActionProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHttpHeaderActionProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for the HTTP header action profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The HTTP header action profile name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the HTTP header action profile.",
			},
			"slot_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The slot ID assigned to the HTTP header action profile.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the HTTP header action profile is deleted.",
			},
			"profile_ready_for_use": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the HTTP header action profile is ready for use.",
			},
			"http_header_action_profile_keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of header key/value pairs applied by the action profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies the header key/value pair.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The header key.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The header value.",
						},
					},
				},
			},
		},
	}
}

func dataSourceHttpHeaderActionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *http_header_action_profile.HttpHeaderActionProfile
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
		log.Printf("[INFO] Searching for HTTP header action profile by ID: %d\n", id)
		allProfiles, err := http_header_action_profile.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting all HTTP header action profiles: %s", err))
		}
		for i := range allProfiles {
			if allProfiles[i].ID == id {
				resp = &allProfiles[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting HTTP header action profile by ID %d: profile not found", id))
		}
	}

	if resp == nil && nameStr != "" {
		log.Printf("[INFO] Searching for HTTP header action profile by name: %s\n", nameStr)
		profile, err := http_header_action_profile.GetByName(ctx, service, nameStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting HTTP header action profile by name %s: %s", nameStr, err))
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

	if err := d.Set("http_header_action_profile_keys", flattenHttpHeaderActionProfileKeys(resp.HttpHeaderActionProfileKeys)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] HTTP header action profile found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}

func flattenHttpHeaderActionProfileKeys(keys []http_header_action_profile.HttpHeaderActionProfileKeys) []interface{} {
	if len(keys) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(keys))
	for _, k := range keys {
		out = append(out, map[string]interface{}{
			"id":    k.ID,
			"key":   k.Key,
			"value": k.Value,
		})
	}
	return out
}
