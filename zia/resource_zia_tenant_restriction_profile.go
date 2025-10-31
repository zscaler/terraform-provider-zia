package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/tenancy_restriction"
)

func resourceTenantRestrictionProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantRestrictionProfileCreate,
		ReadContext:   resourceTenantRestrictionProfileRead,
		UpdateContext: resourceTenantRestrictionProfileUpdate,
		DeleteContext: resourceTenantRestrictionProfileDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("profile_id", idInt)
				} else {
					resp, err := tenancy_restriction.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("profile_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System-generated tenant profile ID",
			},
			"profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "System-generated tenant profile ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tenant restriction profile name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the profile",
			},
			"app_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Restricted tenant profile application type.
				See the Tenancy Restriction Profile API for the list of available application types:
				https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-get`,
			},
			"item_type_primary": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Tenant profile primary item type.
				See the Tenancy Restriction Profile API for the list of available items:
				https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-get`,
			},
			"item_type_secondary": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Tenant profile secondary item type.
				See the Tenancy Restriction Profile API for the list of available items:
				https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-get`,
			},
			"restrict_personal_o365_domains": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to restrict personal domains for Office 365",
			},
			"allow_google_consumers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to allow Google consumers",
			},
			"ms_login_services_tr_v2": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to decide between v1 and v2 for tenant restriction on MSLOGINSERVICES",
			},
			"allow_google_visitors": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to allow Google visitors",
			},
			"allow_gcp_cloud_storage_read": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Flag to allow or disallow cloud storage resources for GCP",
			},
			"item_data_primary": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Tenant profile primary item data",
			},
			"item_data_secondary": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of certifications to be included or excluded for the profile.",
			},
			"item_value": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `Tenant profile item value for YouTube category.
				See the Tenancy Restriction Profile API for the list of available item values:
				https://help.zscaler.com/zia/cloud-app-control-policy#/tenancyRestrictionProfile-get`,
			},
		},
	}
}

func resourceTenantRestrictionProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandTenantRestrictionProfile(d)
	log.Printf("[INFO] Creating ZIA tenant restriction profile\n%+v\n", req)

	resp, _, err := tenancy_restriction.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA tenant restriction profile request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("profile_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceTenantRestrictionProfileRead(ctx, d, meta)
}

func resourceTenantRestrictionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "profile_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia tenant restriction profile id is set"))
	}
	resp, err := tenancy_restriction.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing tenant restriction profile %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting tenant restriction profile:\n%+v\n", resp)

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

	return nil
}

func resourceTenantRestrictionProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "profile_id")
	if !ok {
		log.Printf("[ERROR] tenant restriction profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia tenant restriction profile ID: %v\n", id)
	req := expandTenantRestrictionProfile(d)
	if _, err := tenancy_restriction.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := tenancy_restriction.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceTenantRestrictionProfileRead(ctx, d, meta)
}

func resourceTenantRestrictionProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "profile_id")
	if !ok {
		log.Printf("[ERROR] tenant restriction profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia tenant restriction profile ID: %v\n", (d.Id()))

	if _, err := tenancy_restriction.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia tenant restriction profile deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandTenantRestrictionProfile(d *schema.ResourceData) tenancy_restriction.TenancyRestrictionProfile {
	id, _ := getIntFromResourceData(d, "profile_id")
	result := tenancy_restriction.TenancyRestrictionProfile{
		ID:                          id,
		Name:                        d.Get("name").(string),
		Description:                 d.Get("description").(string),
		AppType:                     d.Get("app_type").(string),
		ItemTypePrimary:             d.Get("item_type_primary").(string),
		ItemTypeSecondary:           d.Get("item_type_secondary").(string),
		RestrictPersonalO365Domains: d.Get("restrict_personal_o365_domains").(bool),
		AllowGoogleConsumers:        d.Get("allow_google_consumers").(bool),
		MsLoginServicesTrV2:         d.Get("ms_login_services_tr_v2").(bool),
		AllowGoogleVisitors:         d.Get("allow_google_visitors").(bool),
		AllowGcpCloudStorageRead:    d.Get("allow_gcp_cloud_storage_read").(bool),
		ItemDataPrimary:             SetToStringList(d, "item_data_primary"),
		ItemDataSecondary:           SetToStringList(d, "item_data_secondary"),
		ItemValue:                   SetToStringList(d, "item_value"),
	}
	return result
}
