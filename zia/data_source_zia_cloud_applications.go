package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudapplications"
)

func dataSourceCloudApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudApplicationsRead,
		Schema: map[string]*schema.Schema{
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of policy to fetch: `cloud_application_policy` or `cloud_application_ssl_policy`",
			},
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cloud applications",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application enum constant",
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud application name",
						},
						"parent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application category enum constant",
						},
						"parent_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the cloud application category",
						},
					},
				},
			},
			"app_class": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter application by application category",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results to include only a specific application name",
			},
		},
	}
}

func dataSourceCloudApplicationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	policyType := d.Get("policy_type").(string)
	params := map[string]interface{}{}
	if appClass, ok := d.GetOk("app_class"); ok {
		params["appClass"] = appClass.([]interface{})
	}
	appNameFilter, hasAppName := d.GetOk("app_name")

	var rawResp interface{}
	var err error

	switch policyType {
	case "cloud_application_policy":
		rawResp, err = cloudapplications.GetCloudApplicationPolicy(ctx, service, params)
	case "cloud_application_ssl_policy":
		rawResp, err = cloudapplications.GetCloudApplicationSSLPolicy(ctx, service, params)
	default:
		return diag.FromErr(fmt.Errorf("invalid policy_type: %s", policyType))
	}

	if err != nil {
		return diag.FromErr(err)
	}

	// Process the response to support multiple applications
	if resp, ok := rawResp.([]cloudapplications.CloudApplications); ok && len(resp) > 0 {
		d.SetId(policyType)

		// Convert the list of applications into a format Terraform can use
		applications := make([]map[string]interface{}, 0)
		for _, app := range resp {
			if hasAppName && app.AppName != appNameFilter.(string) {
				continue
			}
			application := map[string]interface{}{
				"app":         app.App,
				"app_name":    app.AppName,
				"parent":      app.Parent,
				"parent_name": app.ParentName,
			}
			applications = append(applications, application)
		}

		// Set the applications attribute
		if err := d.Set("applications", applications); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("no data found for provided parameters"))
	}

	return nil
}
