package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api"
)

func dataSourceDomainProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainProfilesRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Domain profile ID",
			},
			"profile_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Domain profile name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional notes or information about the domain profile",
			},
			"include_company_domains": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean flag to determine if the organizational domains have to be included in the domain profile",
			},
			"include_subdomains": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean flag to determine whether or not to include subdomains",
			},
			"custom_domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of custom domains for the domain profile. There can be one or more custom domains.",
			},
			"predefined_email_domains": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of predefined email service provider domains for the domain profile",
			},
		},
	}
}

func dataSourceDomainProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var matched *saas_security_api.DomainProfiles

	// Get all domain profiles
	profiles, err := saas_security_api.GetDomainProfiles(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve domain profiles: %w", err))
	}

	id, idOk := getIntFromResourceData(d, "id")
	name, _ := d.Get("profile_name").(string)

	for _, profile := range profiles {
		// Match by ID (if provided)
		if idOk && profile.ProfileID == id {
			matched = &profile
			break
		}

		// Match by name (if ID not matched or not provided)
		if name != "" && profile.ProfileName == name {
			matched = &profile
			break
		}
	}

	if matched == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any domain profile with name '%s' or id '%d'", name, id))
	}

	// Populate the schema fields
	d.SetId(fmt.Sprintf("%d", matched.ProfileID))
	_ = d.Set("profile_name", matched.ProfileName)
	_ = d.Set("description", matched.Description)
	_ = d.Set("include_company_domains", matched.IncludeCompanyDomains)
	_ = d.Set("include_subdomains", matched.IncludeSubdomains)
	_ = d.Set("custom_domains", matched.CustomDomains)
	_ = d.Set("predefined_email_domains", matched.PredefinedEmailDomains)

	return nil
}
