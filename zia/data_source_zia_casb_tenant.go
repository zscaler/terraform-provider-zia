package zia

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/saas_security_api"
)

func dataSourceCasbTenant() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCasbTenantRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "SaaS Security API email label ID",
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "SaaS Security API email label name",
			},
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Return only active tenants",
			},
			"include_deleted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include deleted tenants in the results",
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter tenants by application type",
			},
			"app": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter tenants by application",
			},
			"scan_config_tenants_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Return only tenants with scan config",
			},
			"include_bucket_ready_s3_tenants": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include S3 tenants ready for bucket creation",
			},
			"filter_by_feature": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter tenants by supported features",
			},
			"modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"last_tenant_validation_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"saas_application": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"enterprise_tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"tenant_webhook_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"tenant_deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"re_auth": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Color to apply to the email label",
			},
			"features_supported": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"zscaler_app_tenant_id": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.",
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
					},
				},
			},
		},
	}
}

func dataSourceCasbTenantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var matched *saas_security_api.CasbTenants
	queryParams := make(map[string]interface{})

	// Required filtering
	if id, ok := getIntFromResourceData(d, "tenant_id"); ok {
		queryParams["tenantId"] = strconv.Itoa(id)
	}
	if name, ok := d.GetOk("tenant_name"); ok {
		queryParams["tenantName"] = name.(string)
	}

	// Optional filtering parameters
	if v, ok := d.GetOk("active_only"); ok {
		queryParams["activeOnly"] = v.(bool)
	}
	if v, ok := d.GetOk("include_deleted"); ok {
		queryParams["includeDeleted"] = v.(bool)
	}
	if v, ok := d.GetOk("app_type"); ok {
		queryParams["appType"] = v.(string)
	}
	if v, ok := d.GetOk("app"); ok {
		queryParams["app"] = strings.ToUpper(v.(string))
	}

	if v, ok := d.GetOk("scan_config_tenants_only"); ok {
		queryParams["scanConfigTenantsOnly"] = v.(bool)
	}
	if v, ok := d.GetOk("include_bucket_ready_s3_tenants"); ok {
		queryParams["includeBucketReadyS3Tenants"] = v.(bool)
	}
	if v, ok := d.GetOk("filter_by_feature"); ok {
		rawList := v.([]interface{})
		var list []string
		for _, item := range rawList {
			list = append(list, item.(string))
		}
		queryParams["filterByFeature"] = list
	}

	// Fetch filtered tenants
	tenants, err := saas_security_api.GetCasbTenantLite(ctx, service, queryParams)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve casb tenant: %w", err))
	}

	// Resolve matching tenant
	id, idOk := getIntFromResourceData(d, "tenant_id")
	name, _ := d.Get("tenant_name").(string)

	for _, tenant := range tenants {
		if idOk && tenant.TenantID == id {
			matched = &tenant
			break
		}
		if name != "" && tenant.TenantName == name {
			matched = &tenant
			break
		}
	}

	if matched == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any casb tenant with name '%v' or id '%d'", name, id))
	}

	d.SetId(fmt.Sprintf("%d", matched.TenantID))
	_ = d.Set("tenant_id", matched.TenantID)
	_ = d.Set("tenant_name", matched.TenantName)
	_ = d.Set("last_tenant_validation_time", matched.LastTenantValidationTime)
	_ = d.Set("saas_application", matched.SaaSApplication)
	_ = d.Set("enterprise_tenant_id", matched.EnterpriseTenantID)
	_ = d.Set("tenant_webhook_enabled", matched.TenantWebhookEnabled)
	_ = d.Set("tenant_deleted", matched.TenantDeleted)
	_ = d.Set("re_auth", matched.ReAuth)
	_ = d.Set("features_supported", matched.FeaturesSupported)
	_ = d.Set("status", matched.Status)
	_ = d.Set("modified_time", matched.ModifiedTime)

	if err := d.Set("zscaler_app_tenant_id", flattenIDNameSet(matched.ZscalerAppTenantID)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
