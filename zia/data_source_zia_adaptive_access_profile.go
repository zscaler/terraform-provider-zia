package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adaptive_access"
)

func dataSourceAdaptiveAccessProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdaptiveAccessProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Adaptive Access profile ID",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Adaptive Access profile name",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Adaptive Access profile type",
			},
			"aap_index": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Adaptive Access profile index",
			},
			"iam_aap_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Adaptive Access profile ID that is used by the API for policy configuration. This field allows you to specify which Adaptive Access profiles are applied in the access policy criteria.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates whether the Adaptive Access profile is deleted",
			},
			"iam_aap_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filters the profile rules by one or more Adaptive Access profile IDs. Setting this attribute returns the matching profile rules in the `profiles` block.",
			},
			"org_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filters the profile rules by organization ID. Setting this attribute returns the matching profile rules in the `profiles` block.",
			},
			"profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of Adaptive Access profile rules returned when `iam_aap_ids` or `org_id` is set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Adaptive Access profile ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Adaptive Access profile name",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Adaptive Access profile type",
						},
						"aap_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Adaptive Access profile index",
						},
						"iam_aap_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Adaptive Access profile ID that is used by the API for policy configuration. This field allows you to specify which Adaptive Access profiles are applied in the access policy criteria.",
						},
						"deleted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "A Boolean value that indicates whether the Adaptive Access profile is deleted",
						},
					},
				},
			},
		},
	}
}

func dataSourceAdaptiveAccessProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	nameObj, nameProvided := d.GetOk("name")
	iamIDsRaw, iamIDsProvided := d.GetOk("iam_aap_ids")
	orgIDRaw, orgIDProvided := d.GetOk("org_id")

	if !nameProvided && !iamIDsProvided && !orgIDProvided {
		return diag.FromErr(fmt.Errorf("one of 'name', 'iam_aap_ids', or 'org_id' must be provided"))
	}

	// The profile-rules endpoint is only queried when the user explicitly
	// supplies one of its filter arguments (iam_aap_ids / org_id). Otherwise
	// the data source resolves a single profile by name.
	if iamIDsProvided || orgIDProvided {
		opts := &adaptive_access.GetFilterOptions{}
		if iamIDsProvided {
			for _, v := range iamIDsRaw.([]interface{}) {
				opts.IAMAapIDs = append(opts.IAMAapIDs, v.(string))
			}
		}
		if orgIDProvided {
			orgID := orgIDRaw.(int)
			opts.OrgID = &orgID
		}

		log.Printf("[INFO] Fetching adaptive access profile rules (iam_aap_ids=%v, org_id=%v)\n", opts.IAMAapIDs, opts.OrgID)
		profiles, err := adaptive_access.GetProfileRules(ctx, service, opts)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting adaptive access profile rules: %w", err))
		}
		if len(profiles) == 0 {
			return diag.FromErr(fmt.Errorf("no adaptive access profile rules found for the provided filters"))
		}

		if err := d.Set("profiles", flattenAdaptiveAccessProfiles(profiles)); err != nil {
			return diag.FromErr(err)
		}

		// Surface the first match's scalar attributes for single-filter
		// convenience and to seed a stable resource ID.
		first := profiles[0]
		d.SetId(adaptiveAccessProfileID(first, opts))
		_ = d.Set("name", first.Name)
		_ = d.Set("type", first.Type)
		_ = d.Set("aap_index", first.AapIndex)
		_ = d.Set("iam_aap_id", first.IamAapID)
		_ = d.Set("deleted", first.Deleted)
		return nil
	}

	name := nameObj.(string)
	log.Printf("[INFO] Looking up adaptive access profile by name: %s\n", name)
	profile, err := adaptive_access.GetByName(ctx, service, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(profile.ID))
	_ = d.Set("name", profile.Name)
	_ = d.Set("type", profile.Type)
	_ = d.Set("aap_index", profile.AapIndex)
	_ = d.Set("iam_aap_id", profile.IamAapID)
	_ = d.Set("deleted", profile.Deleted)
	return nil
}

func flattenAdaptiveAccessProfiles(profiles []adaptive_access.AdaptiveAccess) []interface{} {
	if len(profiles) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"type":       p.Type,
			"aap_index":  p.AapIndex,
			"iam_aap_id": p.IamAapID,
			"deleted":    p.Deleted,
		})
	}
	return out
}

// adaptiveAccessProfileID builds a stable resource ID for the filtered
// profile-rules result. When the first record carries an API profile ID it is
// preferred; otherwise the supplied filters are used so repeated reads with the
// same arguments produce the same ID.
func adaptiveAccessProfileID(p adaptive_access.AdaptiveAccess, opts *adaptive_access.GetFilterOptions) string {
	if p.IamAapID != "" {
		return p.IamAapID
	}
	if p.ID != 0 {
		return strconv.Itoa(p.ID)
	}
	id := "adaptive-access"
	for _, v := range opts.IAMAapIDs {
		id += "-" + v
	}
	if opts.OrgID != nil {
		id += "-org" + strconv.Itoa(*opts.OrgID)
	}
	return id
}
