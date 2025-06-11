package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/alerts"
)

func dataSourceSubscriptionAlerts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubscriptionAlertsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the nss server",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The NSS server name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enables or disables the status of the NSS server",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The health of the NSS server",
			},
			"pt0_severities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from cookie authentication",
			},
			"secure_severities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from cookie authentication",
			},
			"manage_severities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from cookie authentication",
			},
			"comply_severities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from cookie authentication",
			},
			"system_severities": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Cloud applications that are exempted from cookie authentication",
			},
		},
	}
}

func dataSourceSubscriptionAlertsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *alerts.AlertSubscriptions
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for NSS Server id: %d\n", id)
		res, err := alerts.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	email, _ := d.Get("email").(string)
	if resp == nil && email != "" {
		log.Printf("[INFO] Getting data for NSS Server with email: %s\n", email)
		alertList, err := alerts.GetAll(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, a := range alertList {
			if a.Email == email {
				resp = &a
				break
			}
		}
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("email", resp.Email)
		_ = d.Set("description", resp.Description)
		_ = d.Set("deleted", resp.Deleted)
		_ = d.Set("pt0_severities", resp.Pt0Severities)
		_ = d.Set("secure_severities", resp.SecureSeverities)
		_ = d.Set("manage_severities", resp.ManageSeverities) // fixed typo from "nanage"
		_ = d.Set("comply_severities", resp.ComplySeverities)
		_ = d.Set("system_severities", resp.SystemSeverities)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any NSS Server with email '%s' or id '%d'", email, id))
	}

	return nil
}
