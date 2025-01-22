package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_idm_profile_lite"
)

func dataSourceDLPIDMProfileLite() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPIDMProfileLiteRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the IDM template (i.e., IDM profile)",
			},
			"template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IDM template name",
			},
			"client_vm": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "This is an immutable reference to an entity. which mainly consists of id and name",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"num_documents": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of documents associated to the IDM template.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The date and time the IDM template was last modified.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the DLP policy rule last.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Filter to include only active IDM profiles",
			},
		},
	}
}

func dataSourceDLPIDMProfileLiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dlp_idm_profile_lite.DLPIDMProfileLite
	activeOnly := d.Get("active_only").(bool) // Retrieve the active_only value

	profileLiteID, ok := getIntFromResourceData(d, "profile_id")
	if ok {
		log.Printf("[INFO] Getting data for dlp idm profile id: %d\n", profileLiteID)
		res, err := dlp_idm_profile_lite.GetDLPProfileLiteID(ctx, service, profileLiteID, activeOnly) // Use activeOnly here
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	profileLiteName, _ := d.Get("template_name").(string)
	if resp == nil && profileLiteName != "" {
		log.Printf("[INFO] Getting data for dlp idm template name: %s\n", profileLiteName)
		res, err := dlp_idm_profile_lite.GetDLPProfileLiteByName(ctx, service, profileLiteName, activeOnly) // Use activeOnly here
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ProfileID))
		_ = d.Set("profile_id", resp.ProfileID)
		_ = d.Set("template_name", resp.TemplateName)
		_ = d.Set("num_documents", resp.NumDocuments)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.ModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("client_vm", flattenIDExtensionsList(resp.ClientVM)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any dlp idm profile name '%s' or id '%d'", profileLiteName, profileLiteID))
	}

	return nil
}
