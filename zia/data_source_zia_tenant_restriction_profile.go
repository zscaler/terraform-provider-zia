package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/tenancy_restriction"
)

func dataSourceTenantRestrictionProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTenantRestrictionProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "The unique identifier for the tenant restriction profile",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The tenant restriction profile name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Additional information about the profile",
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Restricted tenant profile application type",
			},
			"item_type_primary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant profile primary item type",
			},
			"item_type_secondary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant profile secondary item type",
			},
			"restrict_personal_o365_domains": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to restrict personal domains for Office 365",
			},
			"allow_google_consumers": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to allow Google consumers",
			},
			"ms_login_services_tr_v2": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to decide between v1 and v2 for tenant restriction on MSLOGINSERVICES",
			},
			"allow_google_visitors": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to allow Google visitors",
			},
			"allow_gcp_cloud_storage_read": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag to allow or disallow cloud storage resources for GCP",
			},
			"item_data_primary": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tenant profile primary item data",
			},
			"item_data_secondary": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of certifications to be included or excluded for the profile.",
			},
			"item_value": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tenant profile item value for YouTube category",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The time the tenant was last modified",
			},
			"last_modified_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The user who last modified the tenant",
			},
		},
	}
}

func dataSourceTenantRestrictionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *tenancy_restriction.TenancyRestrictionProfile
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for tenant restriction profile id: %d\n", id)
		res, err := tenancy_restriction.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for tenant restriction profile name: %s\n", name)
		res, err := tenancy_restriction.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)
		_ = d.Set("app_type", resp.AppType)
		_ = d.Set("item_type_primary", resp.ItemTypePrimary)
		_ = d.Set("item_type_secondary", resp.ItemTypeSecondary)
		_ = d.Set("restrict_personal_o365_domains", resp.RestrictPersonalO365Domains)
		_ = d.Set("allow_google_consumers", resp.AllowGoogleConsumers)
		_ = d.Set("ms_login_services_tr_v2", resp.MsLoginServicesTrV2)
		_ = d.Set("allow_google_visitors", resp.AllowGoogleVisitors)
		_ = d.Set("allow_gcp_cloud_storage_read", resp.AllowGcpCloudStorageRead)
		_ = d.Set("item_data_primary", resp.ItemDataPrimary)
		_ = d.Set("item_data_secondary", resp.ItemDataSecondary)
		_ = d.Set("item_value", resp.ItemValue)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("last_modified_user_id", resp.LastModifiedUserID)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any tenant restriction profile name '%s' or id '%d'", name, id))
	}

	return nil
}
