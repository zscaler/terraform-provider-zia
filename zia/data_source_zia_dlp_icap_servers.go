package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_icap_servers"
)

func dataSourceDLPICAPServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPICAPServersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for a DLP server.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DLP server name.",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DLP server URL.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DLP server status",
			},
		},
	}
}

func dataSourceDLPICAPServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dlp_icap_servers.DLPICAPServers
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp icap server id: %d\n", id)
		res, err := dlp_icap_servers.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp icap server name: %s\n", name)
		res, err := dlp_icap_servers.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("url", resp.URL)
		_ = d.Set("status", resp.Status)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any dlp icap server name '%s' or id '%d'", name, id))
	}

	return nil
}
