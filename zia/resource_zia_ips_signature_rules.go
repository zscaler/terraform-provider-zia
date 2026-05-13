package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ips_control_policies/ips_signature_rules"
)

func resourceIPSSignatureRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPSSignatureRulesCreate,
		ReadContext:   resourceIPSSignatureRulesRead,
		UpdateContext: resourceIPSSignatureRulesUpdate,
		DeleteContext: resourceIPSSignatureRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("signature_id", idInt)
				} else {
					resp, err := ips_signature_rules.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("signature_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "System-generated identifier for the custom IPS signature rule.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "Custom IPS signature rule name.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
				Description:      "Additional information about the custom signature rule.",
			},
			"rule_text": {
				Type:             schema.TypeString,
				Required:         true,
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
				Description:      "The rule text in Suricata/Snort syntax that defines the custom IPS signature. The provider calls the IPS signature rule text validation endpoint before every create and update. Whitespace, indentation, and trailing newlines (e.g. introduced by HCL heredocs) are normalized before comparison so the same logical rule produced by `\"...\"`, `<<EOT`, or `<<-EOT` does not cause spurious diffs.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the custom signature rule is enabled and ready to be used in IPS Control rules via the assigned threat category.",
			},
			"category": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Threat category assigned to the custom signature rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Unique identifier of the threat category.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the threat category (e.g. ADVANCED_SECURITY).",
						},
						"is_name_l10n_tag": {
							Type:        schema.TypeBool,
							Optional:    true,
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

func resourceIPSSignatureRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	req := expandIPSSignatureRules(d)

	if err := validateIPSSignatureRuleText(ctx, service, req.RuleText); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating ZIA IPS signature rule: %s\n", req.Name)
	resp, _, err := ips_signature_rules.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA IPS signature rule. ID: %d\n", resp.ID)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("signature_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceIPSSignatureRulesRead(ctx, d, meta)
}

func resourceIPSSignatureRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "signature_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no IPS signature rule id is set"))
	}

	resp, err := ips_signature_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia IPS signature rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Got zia IPS signature rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("signature_id", resp.ID)
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

	return nil
}

func resourceIPSSignatureRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "signature_id")
	if !ok {
		log.Printf("[ERROR] IPS signature rule ID not set: %v\n", id)
	}

	req := expandIPSSignatureRules(d)
	req.ID = id

	if err := validateIPSSignatureRuleText(ctx, service, req.RuleText); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating zia IPS signature rule ID: %d\n", id)
	if _, err := ips_signature_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := ips_signature_rules.Update(ctx, service, id, &req); err != nil {
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

	return resourceIPSSignatureRulesRead(ctx, d, meta)
}

func resourceIPSSignatureRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "signature_id")
	if !ok {
		log.Printf("[ERROR] IPS signature rule ID not set: %v\n", id)
	}

	log.Printf("[INFO] Deleting zia IPS signature rule ID: %v\n", d.Id())
	if _, err := ips_signature_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia IPS signature rule deleted")

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

func expandIPSSignatureRules(d *schema.ResourceData) ips_signature_rules.IPSSignatureRules {
	id, _ := getIntFromResourceData(d, "signature_id")
	result := ips_signature_rules.IPSSignatureRules{
		ID:          id,
		Name:        d.Get("name").(string),
		RuleText:    d.Get("rule_text").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		Category:    expandIPSSignatureCategory(d, "category"),
	}
	return result
}

func expandIPSSignatureCategory(d *schema.ResourceData, key string) *ips_signature_rules.IPSSignatureCategory {
	raw, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	cat := &ips_signature_rules.IPSSignatureCategory{}
	if v, ok := m["id"].(int); ok {
		cat.ID = v
	}
	if v, ok := m["name"].(string); ok {
		cat.Name = v
	}
	if v, ok := m["is_name_l10n_tag"].(bool); ok {
		cat.IsNameL10nTag = v
	}
	return cat
}

func flattenIPSSignatureCategory(c *ips_signature_rules.IPSSignatureCategory) []interface{} {
	if c == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":               c.ID,
			"name":             c.Name,
			"is_name_l10n_tag": c.IsNameL10nTag,
		},
	}
}

// validateIPSSignatureRuleText calls the API's IPS signature rule text
// validation endpoint before every Create / Update. The SDK function
// returns:
//
//   - err != nil  → the API responded HTTP 4xx with an INVALID_INPUT_ARGUMENT
//     envelope; we surface the diagnostic verbatim so the user sees the
//     server-reported syntax error.
//   - err == nil and validation.Status != 0  → API reported a validation
//     failure inline (e.g. duplicate signature); surface a Go error built
//     from the diagnostic fields so the caller can refuse the create /
//     update before any state is committed.
//   - err == nil and validation.Status == 0  → the rule is well-formed.
func validateIPSSignatureRuleText(ctx context.Context, service *zscaler.Service, ruleText string) error {
	v, err := ips_signature_rules.ValidateIPSSignatureRuleText(ctx, service, ruleText)
	if err != nil {
		return fmt.Errorf("IPS signature rule text validation failed: %w", err)
	}
	if v == nil {
		return nil
	}
	if v.Status == 0 && v.ErrMsg == "" {
		return nil
	}
	return fmt.Errorf(
		"IPS signature rule text validation failed: status=%d errPosition=%d errParameter=%q errMsg=%q errSuggestion=%q",
		v.Status, v.ErrPosition, v.ErrParameter, v.ErrMsg, v.ErrSuggestion,
	)
}
