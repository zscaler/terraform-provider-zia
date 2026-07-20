package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_global_options"
)

func dataSourceDLPGlobalOptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPGlobalOptionsRead,
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of cloud applications exempted from DLP evaluation",
			},
			"urls": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of URLs exempted from DLP evaluation",
			},
			"http_get_custom_url_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of URL Categories to associate with Inspect HTTP GET Requests",
			},
			"exempt_url_encoded_data": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether or not URL encoded data from DLP evaluation is exempted",
			},
			"enable_npk_edm_templates": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `"Indicates whether EDM with No Primary Keys is enabled.
				If enabled, you select the fields from the template that must match, fields that are optional, and how many optional fields are required to match.
				Before enabling EDM with No Primary Keys, all existing EDM schema (templates, dictionaries, engines, and policies) must be deleted or unassigned.
				New schema with no primary keys functionality can be created after the feature is enabled."`,
			},
			"enable_npk_edm_templates_for_org": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `"Read Only. Indicates whether EDM with No Primary Keys is enabled for the organization.
				To enable this field, you must first enable the enableNpkEdmTemplates field."`,
			},
			"enable_inline_dlp_ocr": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether optical character recognition (OCR) for Zscaler DLP engines to scan images for text content in data in transit is enabled",
			},
			"enable_casb_ocr": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether SaaS Security for Zscaler DLP engines to scan images for text content in data at rest is enabled",
			},
			"enable_email_dlp_ocr": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether Outbound Email DLP for Zscaler DLP engines to scan images for text content in outbound emails is sent to external domains",
			},
			"enable_evaluate_all_dlp_rules": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether DLP engines evaluate all rules or stop when a matching rule is found. When enabled, if multiple rules match, the rule engine selects the rule with the most restrictive action.",
			},
			"enable_edm_popular_format": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether EDM for popular formats is enabled. When enabled, the EDM scans SSN, SIN, and CCN dictionaries for a popular format, and the EDM match only succeeds if the entered format matches the popular data format type.",
			},
			"url_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of custom URL categories exempted from DLP evaluation",
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
							Description: "Identifier that uniquely identifies an entity",
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
		},
	}
}

func dataSourceDLPGlobalOptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := dlp_global_options.GetDLPGlobalOptions(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("dlp_global_options")
		_ = d.Set("applications", resp.Applications)
		_ = d.Set("urls", resp.Urls)
		_ = d.Set("http_get_custom_url_categories", resp.HttpGetCustomUrlCategories)
		_ = d.Set("exempt_url_encoded_data", resp.ExemptUrlEncodedData)
		_ = d.Set("enable_npk_edm_templates", resp.EnableNpkEdmTemplates)
		_ = d.Set("enable_npk_edm_templates_for_org", resp.EnableNpkEdmTemplatesForOrg)
		_ = d.Set("enable_inline_dlp_ocr", resp.EnableInlineDlpOcr)
		_ = d.Set("enable_casb_ocr", resp.EnableCasbOcr)
		_ = d.Set("enable_email_dlp_ocr", resp.EnableEmailDlpOcr)
		_ = d.Set("enable_evaluate_all_dlp_rules", resp.EnableEvaluateAllDlpRules)
		_ = d.Set("enable_edm_popular_format", resp.EnableEdmPopularFormat)
		if err := d.Set("url_categories", flattenIDExtensions(resp.URLCategories)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't read DLP global options"))
	}

	return nil
}
