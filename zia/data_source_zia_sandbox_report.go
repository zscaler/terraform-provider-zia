package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_report"
)

func dataSourceSandboxReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSandboxReportRead,
		Schema: map[string]*schema.Schema{
			"md5_hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			"details": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "summary",
				ValidateFunc: validation.StringInSlice([]string{
					"full",
					"summary",
				}, false),
			},
			"summary": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     summaryDetailSchema(),
			},
			"classification": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     classificationSchema(),
			},
			"file_properties": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     filePropertiesSchema(),
			},
			"origin":          originSchema(),
			"system_summary":  systemSummaryDetailSchema(),
			"spyware":         commonSandboxRSSSchema("spyware"),
			"networking":      commonSandboxRSSSchema("networking"),
			"security_bypass": commonSandboxRSSSchema("security_bypass"),
			"exploit":         commonSandboxRSSSchema("exploit"),
			"stealth":         commonSandboxRSSSchema("stealth"),
			"persistence":     commonSandboxRSSSchema("persistence"),
		},
	}
}

func summaryDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func classificationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"score": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"detected_malware": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func filePropertiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"md5": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sha1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"digital_cerificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssdeep": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_ca": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Helper function to create schema for SandboxRSS fields
func commonSandboxRSSSchema(fieldName string) *schema.Schema {
	// No-op use of fieldName to avoid unused parameter warning
	_ = fieldName

	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"risk": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"signature": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"signature_sources": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// Define schemas for new structs
func originSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"risk": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"language": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"country": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func systemSummaryDetailSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"risk": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"signature": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"signature_sources": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

// dataSourceSandboxReportRead reads the sandbox report data source.
func dataSourceSandboxReportRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	md5Hash := d.Get("md5_hash").(string)
	details := d.Get("details").(string)

	resp, err := zClient.sandbox_report.GetReportMD5Hash(md5Hash, details)
	if err != nil {
		return err
	}

	if resp != nil && resp.Details != nil {
		d.SetId(md5Hash)

		if err := d.Set("md5_hash", md5Hash); err != nil {
			return fmt.Errorf("error setting md5_hash: %s", err)
		}

		if resp.Details != nil {
			if err := d.Set("summary", []interface{}{flattenSummaryDetail(resp.Details.Summary)}); err != nil {
				return fmt.Errorf("error setting summary: %s", err)
			}
			if err := d.Set("classification", []interface{}{flattenClassification(resp.Details.Classification)}); err != nil {
				return fmt.Errorf("error setting classification: %s", err)
			}

			if err := d.Set("file_properties", []interface{}{flattenFileProperties(&resp.Details.FileProperties)}); err != nil {
				return fmt.Errorf("error setting file_properties: %s", err)
			}
		}

		if details == "full" {
			// Set additional fields for full details
			if err := d.Set("origin", []interface{}{flattenOrigin(resp.Details.Origin)}); err != nil {
				return fmt.Errorf("error setting origin: %s", err)
			}

			if err := d.Set("system_summary", flattenSystemSummaryDetails(resp.Details.SystemSummary)); err != nil {
				return fmt.Errorf("error setting system_summary: %s", err)
			}

			// Repeat for other SandboxRSS fields like spyware, networking, etc.
			if err := d.Set("spyware", flattenSandboxRSS(resp.Details.Spyware)); err != nil {
				return fmt.Errorf("error setting spyware: %s", err)
			}
			if err := d.Set("networking", flattenSandboxRSS(resp.Details.Networking)); err != nil {
				return fmt.Errorf("error setting networking: %s", err)
			}
			if err := d.Set("security_bypass", flattenSandboxRSS(resp.Details.SecurityBypass)); err != nil {
				return fmt.Errorf("error setting security_bypass: %s", err)
			}
			if err := d.Set("exploit", flattenSandboxRSS(resp.Details.Exploit)); err != nil {
				return fmt.Errorf("error setting exploit: %s", err)
			}
			if err := d.Set("stealth", flattenSandboxRSS(resp.Details.Stealth)); err != nil {
				return fmt.Errorf("error setting stealth: %s", err)
			}
			if err := d.Set("persistence", flattenSandboxRSS(resp.Details.Persistence)); err != nil {
				return fmt.Errorf("error setting persistence: %s", err)
			}
			if err := d.Set("persistence", flattenSandboxRSS(resp.Details.Persistence)); err != nil {
				return fmt.Errorf("error setting persistence: %s", err)
			}
		}
	} else {
		return fmt.Errorf("couldn't find any reports for MD5 hash: %s", md5Hash)
	}

	return nil
}

// Helper functions to flatten details
func flattenSummaryDetail(summary sandbox_report.SummaryDetail) map[string]interface{} {
	return map[string]interface{}{
		"status":     summary.Status,
		"category":   summary.Category,
		"file_type":  summary.FileType,
		"start_time": summary.StartTime,
		"duration":   summary.Duration,
	}
}

func flattenClassification(classification sandbox_report.Classification) map[string]interface{} {
	return map[string]interface{}{
		"type":             classification.Type,
		"category":         classification.Category,
		"score":            classification.Score,
		"detected_malware": classification.DetectedMalware,
	}
}

func flattenFileProperties(fileProperties *sandbox_report.FileProperties) map[string]interface{} {
	return map[string]interface{}{
		"file_type":          fileProperties.FileType,
		"file_size":          fileProperties.FileSize,
		"md5":                fileProperties.MD5,
		"sha1":               fileProperties.SHA1,
		"sha256":             fileProperties.SHA256,
		"issuer":             fileProperties.Issuer,
		"digital_cerificate": fileProperties.DigitalCerificate,
		"ssdeep":             fileProperties.SSDeep,
		"root_ca":            fileProperties.RootCA,
	}
}

func flattenOrigin(origin *sandbox_report.Origin) map[string]interface{} {
	if origin == nil {
		return nil
	}

	return map[string]interface{}{
		"risk":     origin.Risk,
		"language": origin.Language,
		"country":  origin.Country,
	}
}

func flattenSystemSummaryDetails(details []sandbox_report.SystemSummaryDetail) []interface{} {
	if details == nil {
		return nil
	}

	var out []interface{}
	for _, detail := range details {
		m := make(map[string]interface{})
		m["risk"] = detail.Risk
		m["signature"] = detail.Signature
		m["signature_sources"] = detail.SignatureSources

		out = append(out, m)
	}
	return out
}

func flattenSandboxRSS(items []*common.SandboxRSS) []interface{} {
	if items == nil {
		return nil
	}

	var out []interface{}
	for _, item := range items {
		m := make(map[string]interface{})
		m["risk"] = item.Risk
		m["signature"] = item.Signature
		m["signature_sources"] = item.SignatureSources

		out = append(out, m)
	}
	return out
}
