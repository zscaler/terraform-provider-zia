package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ips_control_policies/ips_signature_rules"
)

func dataSourceIPSSignatureRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPSSignatureRulesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "System-generated identifier for the custom IPS signature rule.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Custom IPS signature rule name. Matched case-insensitively against the API result set when used as a lookup.",
			},
			"rule_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule text in Suricata/Snort syntax that defines the custom IPS signature.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about the custom signature rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the custom signature rule is enabled and ready to be used in IPS Control rules via the assigned threat category.",
			},
			"category": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Threat category assigned to the custom signature rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier of the threat category.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the threat category (e.g. ADVANCED_SECURITY).",
						},
						"is_name_l10n_tag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the name is a localization tag rather than a literal label.",
						},
					},
				},
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the custom signature rule is marked as deleted.",
			},
			"promote_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp (in seconds) when the rule was promoted. 0 if not yet promoted.",
			},
			"rule_text_mod_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp (in seconds) when the rule text was last modified.",
			},
			"dynamic_validation_submitted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the rule was submitted for dynamic validation.",
			},
			"dynamic_validation_rejected": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether dynamic validation rejected the rule.",
			},
			"dynamic_validation_succeeded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether dynamic validation succeeded for the rule.",
			},
			"disabled_from_zscm": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the rule was disabled from Zscaler Cloud Management.",
			},
			"dynamic_val_reject_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Reject code returned by dynamic validation.",
			},
		},
	}
}

func dataSourceIPSSignatureRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *ips_signature_rules.IPSSignatureRules

	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	if idProvided {
		log.Printf("[INFO] Getting data for IPS signature rule id: %d\n", id)
		res, err := ips_signature_rules.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting IPS signature rule by id %d: %w", id, err))
		}
		resp = res
	}

	if resp == nil && nameStr != "" {
		log.Printf("[INFO] Getting data for IPS signature rule name: %s\n", nameStr)
		res, err := ips_signature_rules.GetByName(ctx, service, nameStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting IPS signature rule by name %q: %w", nameStr, err))
		}
		resp = res
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("rule_text", resp.RuleText)
	_ = d.Set("description", resp.Description)
	_ = d.Set("enabled", resp.Enabled)
	_ = d.Set("deleted", resp.Deleted)
	_ = d.Set("promote_time", resp.PromoteTime)
	_ = d.Set("rule_text_mod_time", resp.RuleTextModTime)
	_ = d.Set("dynamic_validation_submitted", resp.DynamicValidationSubmitted)
	_ = d.Set("dynamic_validation_rejected", resp.DynamicValidationRejected)
	_ = d.Set("dynamic_validation_succeeded", resp.DynamicValidationSucceeded)
	_ = d.Set("disabled_from_zscm", resp.DisabledFromZSCM)
	_ = d.Set("dynamic_val_reject_code", resp.DynamicValRejectCode)

	if err := d.Set("category", flattenIPSSignatureCategory(resp.Category)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting category: %s", err))
	}

	log.Printf("[DEBUG] IPS signature rule found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}
