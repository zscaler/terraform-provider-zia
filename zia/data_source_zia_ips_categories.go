package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ips_control_policies/ips_signature_rules"
)

func dataSourceIPSCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPSCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the IPS category.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the IPS category. This is the value used by the `res_categories` and `dest_ip_categories` attributes on firewall rules.",
			},
			"back_end_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The descriptive name of the IPS category as displayed in the admin portal.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the IPS category.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the IPS category has been deleted.",
			},
			"predefined": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the IPS category is predefined by the service.",
			},
			"ips_signature_rules_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of IPS signature rules associated with the category.",
			},
			"categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Every available IPS category. Populated when both `id` and `name` are omitted, so that the available category names can be discovered.",
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
						"back_end_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"predefined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ips_signature_rules_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIPSCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	log.Printf("[INFO] Fetching all IPS categories")
	allCategories, err := ips_signature_rules.GetIPSCategories(ctx, service)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all IPS categories: %w", err))
	}

	log.Printf("[DEBUG] Retrieved %d IPS categories", len(allCategories))

	var resp *ips_signature_rules.IPSCategories
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	if idProvided {
		log.Printf("[INFO] Searching for IPS category by ID: %d", id)
		for i := range allCategories {
			if allCategories[i].Id == id {
				resp = &allCategories[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("IPS category with ID %d not found", id))
		}
	}

	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for IPS category by name: %s", nameStr)
		for i := range allCategories {
			if strings.EqualFold(allCategories[i].Name, nameStr) {
				resp = &allCategories[i]
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("IPS category with name %q not found", nameStr))
		}
	}

	// Neither id nor name was supplied, so return the whole list and let the
	// caller pick from it. The "id" attribute is an integer, so this synthetic
	// id has to stay numeric or the SDK panics when it writes the state. There
	// are no list-mode filters to tell one unfiltered read from another, so a
	// constant is sufficient here.
	if resp == nil {
		log.Printf("[INFO] Returning all %d IPS categories", len(allCategories))
		d.SetId("0")
		if err := d.Set("categories", flattenIPSCategories(allCategories)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting categories: %w", err))
		}
		return nil
	}

	d.SetId(fmt.Sprintf("%d", resp.Id))
	if err := d.Set("id", resp.Id); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %w", err))
	}
	if err := d.Set("name", resp.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %w", err))
	}
	if err := d.Set("back_end_name", resp.BackEndName); err != nil {
		return diag.FromErr(fmt.Errorf("error setting back_end_name: %w", err))
	}
	if err := d.Set("description", resp.Description); err != nil {
		return diag.FromErr(fmt.Errorf("error setting description: %w", err))
	}
	if err := d.Set("deleted", resp.Deleted); err != nil {
		return diag.FromErr(fmt.Errorf("error setting deleted: %w", err))
	}
	if err := d.Set("predefined", resp.Predefined); err != nil {
		return diag.FromErr(fmt.Errorf("error setting predefined: %w", err))
	}
	if err := d.Set("ips_signature_rules_count", resp.IpsSignatureRulesCount); err != nil {
		return diag.FromErr(fmt.Errorf("error setting ips_signature_rules_count: %w", err))
	}

	log.Printf("[DEBUG] IPS category found: ID=%d, Name=%s", resp.Id, resp.Name)
	return nil
}

func flattenIPSCategories(categories []ips_signature_rules.IPSCategories) []interface{} {
	result := make([]interface{}, 0, len(categories))
	for _, c := range categories {
		result = append(result, map[string]interface{}{
			"id":                        c.Id,
			"name":                      c.Name,
			"back_end_name":             c.BackEndName,
			"description":               c.Description,
			"deleted":                   c.Deleted,
			"predefined":                c.Predefined,
			"ips_signature_rules_count": c.IpsSignatureRulesCount,
		})
	}
	return result
}
