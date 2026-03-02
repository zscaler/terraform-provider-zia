package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/email_profiles"
)

func dataSourceEmailProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the email recipient profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the email recipient profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the email recipient profile.",
			},
			"emails": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The list of recipient email addresses.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceEmailProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all email profiles")
	allProfiles, err := email_profiles.GetAll(ctx, service, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all email profiles: %w", err))
	}

	log.Printf("[DEBUG] Retrieved %d email profiles", len(allProfiles))

	var resp *email_profiles.EmailProfiles
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	if idProvided {
		log.Printf("[INFO] Searching for email profile by ID: %d", id)
		for i := range allProfiles {
			if allProfiles[i].ID == id {
				resp = &allProfiles[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("email profile with ID %d not found", id))
		}
	}

	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for email profile by name: %s", nameStr)
		for i := range allProfiles {
			if strings.EqualFold(allProfiles[i].Name, nameStr) {
				resp = &allProfiles[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("email profile with name %q not found", nameStr))
		}
	}

	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("no email profile found with name %q or id %d", nameStr, id))
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
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %w", err))
	}
	if err := d.Set("emails", resp.Emails); err != nil {
		return diag.FromErr(fmt.Errorf("error setting emails: %w", err))
	}

	log.Printf("[DEBUG] Email profile found: ID=%d, Name=%s", resp.ID, resp.Name)
	return nil
}
