package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/nss_servers"
)

func dataSourceNSSServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNSSServersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the nss server",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The NSS server name",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the status of the NSS server",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of the NSS server",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the NSS Server",
			},
			"icap_svr_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ICAP server ID",
			},
		},
	}
}

func dataSourceNSSServersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *nss_servers.NSSServers
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for NSS Server id: %d\n", id)
		res, err := nss_servers.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for NSS Server name: %s\n", name)
		res, err := nss_servers.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("status", resp.Status)
		_ = d.Set("state", resp.State)
		_ = d.Set("type", resp.Type)
		_ = d.Set("icap_svr_id", resp.IcapSvrId)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any NSS Server name '%s' or id '%d'", name, id))
	}

	return nil
}
