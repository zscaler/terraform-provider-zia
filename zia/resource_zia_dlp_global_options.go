package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_global_options"
)

func resourceDLPGlobalOptions() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDLPGlobalOptionsRead,
		CreateContext: resourceDLPGlobalOptionsCreate,
		UpdateContext: resourceDLPGlobalOptionsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceDLPGlobalOptionsRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read atp malware inspection import: %s", diags[0].Summary)
				}
				d.SetId("dlp_global_options")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of cloud applications to which the DLP policy rule must be applied
				Use the data source zia_cloud_applications to get the list of available cloud applications:
				https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_applications
				`,
			},
			"urls": {
				Type:             schema.TypeList,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: suppressURLCategoriesReorderDiff,
			},
			"http_get_custom_url_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of URL Categories to associate with Inspect HTTP GET Requests",
			},
			"exempt_url_encoded_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether or not URL encoded data from DLP evaluation is exempted",
			},
			"enable_npk_edm_templates": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `"Indicates whether EDM with No Primary Keys is enabled.
				If enabled, you select the fields from the template that must match, fields that are optional, and how many optional fields are required to match.
				Before enabling EDM with No Primary Keys, all existing EDM schema (templates, dictionaries, engines, and policies) must be deleted or unassigned.
				New schema with no primary keys functionality can be created after the feature is enabled."`,
			},
			"enable_npk_edm_templates_for_org": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: `"Read Only. Indicates whether EDM with No Primary Keys is enabled for the organization.
				To enable this field, you must first enable the enableNpkEdmTemplates field."`,
			},
			"enable_inline_dlp_ocr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether optical character recognition (OCR) for Zscaler DLP engines to scan images for text content in data in transit is enabled",
			},
			"enable_casb_ocr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether SaaS Security for Zscaler DLP engines to scan images for text content in data at rest is enabled",
			},
			"enable_email_dlp_ocr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether Outbound Email DLP for Zscaler DLP engines to scan images for text content in outbound emails is sent to external domains",
			},
			"enable_evaluate_all_dlp_rules": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether DLP engines evaluate all rules or stop when a matching rule is found. When enabled, if multiple rules match, the rule engine selects the rule with the most restrictive action.",
			},
			"enable_edm_popular_format": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether EDM for popular formats is enabled. When enabled, the EDM scans SSN, SIN, and CCN dictionaries for a popular format, and the EDM match only succeeds if the entered format matches the popular data format type.",
			},
			"url_categories": setIDsSchemaTypeCustom(nil, "The list of URL categories to which the DLP global options must be applied"),
		},
	}
}

func resourceDLPGlobalOptionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandDLPGlobalSettings(d)
	_, _, err := dlp_global_options.UpdateDLPGlobalOptions(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("dlp_global_options")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDLPGlobalOptionsRead(ctx, d, meta)
}

func resourceDLPGlobalOptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := dlp_global_options.GetDLPGlobalOptions(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("dlp_global_options")
		_ = d.Set("applications", resp.Applications)
		_ = d.Set("urls", reorderURLsToReferenceOrder(d.Get("urls"), resp.Urls))
		_ = d.Set("http_get_custom_url_categories", resp.HttpGetCustomUrlCategories)
		_ = d.Set("exempt_url_encoded_data", resp.ExemptUrlEncodedData)
		_ = d.Set("enable_npk_edm_templates", resp.EnableNpkEdmTemplates)
		_ = d.Set("enable_npk_edm_templates_for_org", resp.EnableNpkEdmTemplatesForOrg)
		_ = d.Set("enable_inline_dlp_ocr", resp.EnableInlineDlpOcr)
		_ = d.Set("enable_casb_ocr", resp.EnableCasbOcr)
		_ = d.Set("enable_email_dlp_ocr", resp.EnableEmailDlpOcr)
		_ = d.Set("enable_evaluate_all_dlp_rules", resp.EnableEvaluateAllDlpRules)
		_ = d.Set("enable_edm_popular_format", resp.EnableEdmPopularFormat)

		if err := d.Set("url_categories", flattenIDExtensionsListIDs(resp.URLCategories)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read advanced threat protection setting options"))
	}

	return nil
}

func resourceDLPGlobalOptionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandDLPGlobalSettings(d)

	_, _, err := dlp_global_options.UpdateDLPGlobalOptions(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("dlp_global_options")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceDLPGlobalOptionsRead(ctx, d, meta)
}

func expandDLPGlobalSettings(d *schema.ResourceData) dlp_global_options.WebDlpGlobal {
	result := dlp_global_options.WebDlpGlobal{
		ExemptUrlEncodedData:        d.Get("exempt_url_encoded_data").(bool),
		EnableNpkEdmTemplates:       d.Get("enable_npk_edm_templates").(bool),
		EnableNpkEdmTemplatesForOrg: d.Get("enable_npk_edm_templates_for_org").(bool),
		EnableInlineDlpOcr:          d.Get("enable_inline_dlp_ocr").(bool),
		EnableCasbOcr:               d.Get("enable_casb_ocr").(bool),
		EnableEmailDlpOcr:           d.Get("enable_email_dlp_ocr").(bool),
		EnableEvaluateAllDlpRules:   d.Get("enable_evaluate_all_dlp_rules").(bool),
		EnableEdmPopularFormat:      d.Get("enable_edm_popular_format").(bool),
		Applications:                SetToStringList(d, "applications"),
		HttpGetCustomUrlCategories:  SetToStringList(d, "http_get_custom_url_categories"),
		Urls:                        ListToStringList(d, "urls"),
		URLCategories:               expandIDNameExtensionsSet(d, "url_categories"),
	}
	return result
}
