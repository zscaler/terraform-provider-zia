package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
)

func dataBandwdithClasses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBandwdithClassesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the bandwidth class",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of the bandwidth class",
			},
			"urls": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The URLs included in the bandwidth class",
			},
			"web_applications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The web conferencing applications included in the bandwidth class",
			},
			"application_service_groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The application service groups included in the bandwidth class",
			},
			"network_applications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The network applications included in the bandwidth class",
			},
			"network_services": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The network services included in the bandwidth class",
			},
			"url_categories": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The URL categories to add to the bandwidth class",
			},
			"applications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The applications included in the bandwidth class",
			},
			"is_name_l10n_tag": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the bandwidth class name is localized",
			},
		},
	}
}

func dataBandwdithClassesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *bandwidth_classes.BandwidthClasses
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for bandwidth class id: %d\n", id)
		res, err := bandwidth_classes.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for bandwidth class name: %s\n", name)
		res, err := bandwidth_classes.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("getfile_size", resp.GetfileSize)
		_ = d.Set("file_size", resp.FileSize)
		_ = d.Set("type", resp.Type)
		_ = d.Set("urls", resp.Urls)
		_ = d.Set("web_applications", resp.WebApplications)
		_ = d.Set("application_service_groups", resp.ApplicationServiceGroups)
		_ = d.Set("network_applications", resp.NetworkApplications)
		_ = d.Set("network_services", resp.NetworkServices)
		_ = d.Set("url_categories", resp.UrlCategories)
		_ = d.Set("applications", resp.Applications)
		_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any bandwidth class name '%s' or id '%d'", name, id))
	}

	return nil
}
