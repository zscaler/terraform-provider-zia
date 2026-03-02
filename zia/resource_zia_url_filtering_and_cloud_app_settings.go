package zia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func resourceURLFilteringCloludAppSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceURLFilteringCloludAppSettingsRead,
		CreateContext: resourceURLFilteringCloludAppSettingsCreate,
		UpdateContext: resourceURLFilteringCloludAppSettingsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceURLFilteringCloludAppSettingsRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read url filtering and cloud app settings import: %s", diags[0].Summary)
				}
				d.SetId("app_setting")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"enable_dynamic_content_cat": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value that indicates if dynamic categorization of URLs by analyzing content of uncategorized websites using AI/ML tools is enabled or not.",
			},
			"consider_embedded_sites": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value that indicates if URL filtering rules must be applied to sites that are translated using translation services or not.",
			},
			"enforce_safe_search": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value that indicates whether only safe content must be returned for web, image, and video search.",
			},
			"safe_search_apps": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `A list of applications for which the SafeSearch enforcement applies. You cannot modify this field when the enforceSafeSearch field is disabled.
				See the URL & Cloud App Control Policy Settings for the list of available apps:
				https://help.zscaler.com/zia/url-cloud-app-control-policy-settings#/advancedUrlFilterAndCloudAppSettings-get`,
			},
			"enable_office365": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value that enables or disables Microsoft Office 365 configuration.",
			},
			"enable_msft_o365": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to permit secure local breakout for Office 365 traffic automatically without any manual configuration needed.",
			},
			"enable_ucaas_zoom": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for Zoom traffic, without any manual configuration needed.",
			},
			"enable_ucaas_logmein": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for GoTo traffic, without any manual configuration needed.",
			},
			"enable_ucaas_ring_central": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for RingCentral traffic, without any manual configuration needed.",
			},
			"enable_ucaas_webex": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for Webex traffic, without any manual configuration needed.",
			},
			"enable_ucaas_talkdesk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for Talkdesk traffic, with minimal or no manual configuration needed.",
			},
			"enable_chatgpt_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with ChatGPT by users should be categorized and logged",
			},
			"enable_microsoft_copilot_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Microsoft Copilot by users should be categorized and logged",
			},
			"enable_gemini_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Google Gemini by users should be categorized and logged",
			},
			"enable_poe_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Poe by users should be categorized and logged",
			},
			"enable_meta_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Meta AI by users should be categorized and logged",
			},
			"enable_per_plexity_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Perplexity by users should be categorized and logged",
			},
			"enable_deep_seek_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with DeepSeek by users should be categorized and logged",
			},
			"enable_writer_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Writer by users should be categorized and logged",
			},
			"enable_grok_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Grok by users should be categorized and logged",
			},
			"enable_mistral_ai_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Mistral AI by users should be categorized and logged",
			},
			"enable_claude_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Claude by users should be categorized and logged",
			},
			"enable_grammarly_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Grammarly by users should be categorized and logged",
			},
			"enable_newly_registered_domains": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating whether newly registered and observed domains that are identified within hours of going live are allowed or blocked",
			},
			"enable_block_override_for_non_auth_user": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if authorized users can temporarily override block action on websites by providing their authentication information",
			},
			"enable_cipa_compliance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if the predefined CIPA Compliance Rule is enabled or not. ",
			},
			"zvelo_db_lookup_disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if Zvelo database lookup is disabled.",
			},
			"enable_creative_commons_search_results": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "A Boolean value indicating if Creative Commons search results are enabled.",
			},
		},
	}
}

func resourceURLFilteringCloludAppSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if err := putUrlAndAppSettingsWithFullPayload(ctx, service, d); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("app_setting")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 1 seconds before potentially triggering the activation
		time.Sleep(1 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceURLFilteringCloludAppSettingsRead(ctx, d, meta)
}

func resourceURLFilteringCloludAppSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := urlfilteringpolicies.GetUrlAndAppSettings(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("app_setting")
		_ = d.Set("enable_dynamic_content_cat", resp.EnableDynamicContentCat)
		_ = d.Set("consider_embedded_sites", resp.ConsiderEmbeddedSites)
		_ = d.Set("enforce_safe_search", resp.EnforceSafeSearch)
		_ = d.Set("safe_search_apps", resp.SafeSearchApps)
		_ = d.Set("enable_office365", resp.EnableOffice365)
		_ = d.Set("enable_msft_o365", resp.EnableMsftO365)
		_ = d.Set("enable_ucaas_zoom", resp.EnableUcaasZoom)
		_ = d.Set("enable_ucaas_logmein", resp.EnableUcaasLogMeIn)
		_ = d.Set("enable_ucaas_ring_central", resp.EnableUcaasRingCentral)
		_ = d.Set("enable_ucaas_webex", resp.EnableUcaasWebex)
		_ = d.Set("enable_ucaas_talkdesk", resp.EnableUcaasTalkdesk)
		_ = d.Set("enable_chatgpt_prompt", resp.EnableChatGptPrompt)
		_ = d.Set("enable_microsoft_copilot_prompt", resp.EnableMicrosoftCoPilotPrompt)
		_ = d.Set("enable_gemini_prompt", resp.EnableGeminiPrompt)
		_ = d.Set("enable_poe_prompt", resp.EnablePOEPrompt)
		_ = d.Set("enable_meta_prompt", resp.EnableMetaPrompt)
		_ = d.Set("enable_per_plexity_prompt", resp.EnablePerPlexityPrompt)
		_ = d.Set("enable_deep_seek_prompt", resp.EnableDeepSeekPrompt)
		_ = d.Set("enable_writer_prompt", resp.EnableWriterPrompt)
		_ = d.Set("enable_grok_prompt", resp.EnableGrokPrompt)
		_ = d.Set("enable_mistral_ai_prompt", resp.EnableMistralAIPrompt)
		_ = d.Set("enable_claude_prompt", resp.EnableClaudePrompt)
		_ = d.Set("enable_grammarly_prompt", resp.EnableGrammarlyPrompt)
		_ = d.Set("enable_newly_registered_domains", resp.EnableNewlyRegisteredDomains)
		_ = d.Set("enable_block_override_for_non_auth_user", resp.EnableBlockOverrideForNonAuthUser)
		_ = d.Set("enable_cipa_compliance", resp.EnableCIPACompliance)
		_ = d.Set("zvelo_db_lookup_disabled", resp.ZveloDbLookupDisabled)
		_ = d.Set("enable_creative_commons_search_results", resp.EnableCreativeCommonsSearchResults)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't read url filtering and cloud app settings"))
	}

	return nil
}

func resourceURLFilteringCloludAppSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if err := putUrlAndAppSettingsWithFullPayload(ctx, service, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("app_setting")

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

	return resourceURLFilteringCloludAppSettingsRead(ctx, d, meta)
}

const urlAdvSettingsEndpoint = "/zia/api/v1/advancedUrlFilterAndCloudAppSettings"

// putUrlAndAppSettingsWithFullPayload GETs current settings, overlays Terraform values,
// and PUTs the complete payload. The API requires the full payload; the SDK struct's
// omitempty causes false values to be omitted, leading to "Request body is invalid".
func putUrlAndAppSettingsWithFullPayload(ctx context.Context, service *zscaler.Service, d *schema.ResourceData) error {
	current, err := urlfilteringpolicies.GetUrlAndAppSettings(ctx, service)
	if err != nil {
		return err
	}

	// Convert current state to map to preserve all fields (including API-only ones)
	data, err := json.Marshal(current)
	if err != nil {
		return err
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	// Overlay Terraform-managed fields (use API JSON keys)
	payload["enableDynamicContentCat"] = d.Get("enable_dynamic_content_cat").(bool)
	payload["considerEmbeddedSites"] = d.Get("consider_embedded_sites").(bool)
	payload["enforceSafeSearch"] = d.Get("enforce_safe_search").(bool)
	payload["safeSearchApps"] = SetToStringList(d, "safe_search_apps")
	payload["enableOffice365"] = d.Get("enable_office365").(bool)
	payload["enableMsftO365"] = d.Get("enable_msft_o365").(bool)
	payload["enableUcaasZoom"] = d.Get("enable_ucaas_zoom").(bool)
	payload["enableUcaasLogMeIn"] = d.Get("enable_ucaas_logmein").(bool)
	payload["enableUcaasRingCentral"] = d.Get("enable_ucaas_ring_central").(bool)
	payload["enableUcaasWebex"] = d.Get("enable_ucaas_webex").(bool)
	payload["enableUcaasTalkdesk"] = d.Get("enable_ucaas_talkdesk").(bool)
	payload["enableChatGptPrompt"] = d.Get("enable_chatgpt_prompt").(bool)
	payload["enableMicrosoftCoPilotPrompt"] = d.Get("enable_microsoft_copilot_prompt").(bool)
	payload["enableGeminiPrompt"] = d.Get("enable_gemini_prompt").(bool)
	payload["enablePOEPrompt"] = d.Get("enable_poe_prompt").(bool)
	payload["enableMetaPrompt"] = d.Get("enable_meta_prompt").(bool)
	payload["enablePerPlexityPrompt"] = d.Get("enable_per_plexity_prompt").(bool)
	payload["enableDeepSeekPrompt"] = d.Get("enable_deep_seek_prompt").(bool)
	payload["enableWriterPrompt"] = d.Get("enable_writer_prompt").(bool)
	payload["enableGrokPrompt"] = d.Get("enable_grok_prompt").(bool)
	payload["enableMistralAIPrompt"] = d.Get("enable_mistral_ai_prompt").(bool)
	payload["enableClaudePrompt"] = d.Get("enable_claude_prompt").(bool)
	payload["enableGrammarlyPrompt"] = d.Get("enable_grammarly_prompt").(bool)
	payload["enableNewlyRegisteredDomains"] = d.Get("enable_newly_registered_domains").(bool)
	payload["enableBlockOverrideForNonAuthUser"] = d.Get("enable_block_override_for_non_auth_user").(bool)
	payload["enableCIPACompliance"] = d.Get("enable_cipa_compliance").(bool)
	payload["zveloDbLookupDisabled"] = d.Get("zvelo_db_lookup_disabled").(bool)
	payload["enableCreativeCommonsSearchResults"] = d.Get("enable_creative_commons_search_results").(bool)

	// Validate CIPA compliance (same as SDK)
	if payload["enableCIPACompliance"] == true {
		if payload["enableNewlyRegisteredDomains"] == true ||
			payload["considerEmbeddedSites"] == true ||
			payload["enforceSafeSearch"] == true ||
			payload["enableDynamicContentCat"] == true {
			return fmt.Errorf("enableCIPACompliance cannot be enabled with the following options: EnableNewlyRegisteredDomains, ConsiderEmbeddedSites, EnforceSafeSearch, EnableDynamicContentCat")
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Use low-level PUT to send full JSON (map includes all keys; SDK struct would omit false)
	_, _, _, err = service.Client.ExecuteRequest(ctx, "PUT", urlAdvSettingsEndpoint, bytes.NewReader(body), nil, "application/json")
	return err
}

func expandURLFilteringCloudAppSettings(d *schema.ResourceData) urlfilteringpolicies.URLAdvancedPolicySettings {
	result := urlfilteringpolicies.URLAdvancedPolicySettings{
		EnableDynamicContentCat:            d.Get("enable_dynamic_content_cat").(bool),
		ConsiderEmbeddedSites:              d.Get("consider_embedded_sites").(bool),
		EnforceSafeSearch:                  d.Get("enforce_safe_search").(bool),
		SafeSearchApps:                     SetToStringList(d, "safe_search_apps"),
		EnableOffice365:                    d.Get("enable_office365").(bool),
		EnableMsftO365:                     d.Get("enable_msft_o365").(bool),
		EnableUcaasZoom:                    d.Get("enable_ucaas_zoom").(bool),
		EnableUcaasLogMeIn:                 d.Get("enable_ucaas_logmein").(bool),
		EnableUcaasRingCentral:             d.Get("enable_ucaas_ring_central").(bool),
		EnableUcaasWebex:                   d.Get("enable_ucaas_webex").(bool),
		EnableUcaasTalkdesk:                d.Get("enable_ucaas_talkdesk").(bool),
		EnableChatGptPrompt:                d.Get("enable_chatgpt_prompt").(bool),
		EnableMicrosoftCoPilotPrompt:       d.Get("enable_microsoft_copilot_prompt").(bool),
		EnableGeminiPrompt:                 d.Get("enable_gemini_prompt").(bool),
		EnablePOEPrompt:                    d.Get("enable_poe_prompt").(bool),
		EnableMetaPrompt:                   d.Get("enable_meta_prompt").(bool),
		EnablePerPlexityPrompt:             d.Get("enable_per_plexity_prompt").(bool),
		EnableDeepSeekPrompt:               d.Get("enable_deep_seek_prompt").(bool),
		EnableWriterPrompt:                 d.Get("enable_writer_prompt").(bool),
		EnableGrokPrompt:                   d.Get("enable_grok_prompt").(bool),
		EnableMistralAIPrompt:              d.Get("enable_mistral_ai_prompt").(bool),
		EnableClaudePrompt:                 d.Get("enable_claude_prompt").(bool),
		EnableGrammarlyPrompt:              d.Get("enable_grammarly_prompt").(bool),
		EnableNewlyRegisteredDomains:       d.Get("enable_newly_registered_domains").(bool),
		EnableBlockOverrideForNonAuthUser:  d.Get("enable_block_override_for_non_auth_user").(bool),
		EnableCIPACompliance:               d.Get("enable_cipa_compliance").(bool),
		ZveloDbLookupDisabled:              d.Get("zvelo_db_lookup_disabled").(bool),
		EnableCreativeCommonsSearchResults: d.Get("enable_creative_commons_search_results").(bool),
	}
	return result
}
