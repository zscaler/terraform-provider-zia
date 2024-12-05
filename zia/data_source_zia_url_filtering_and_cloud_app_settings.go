package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func dataSourceURLFilteringCloludAppSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceURLFilteringCloludAppSettingsRead,
		Schema: map[string]*schema.Schema{
			"enable_dynamic_content_cat": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates if dynamic categorization of URLs by analyzing content of uncategorized websites using AI/ML tools is enabled or not.",
			},
			"consider_embedded_sites": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates if URL filtering rules must be applied to sites that are translated using translation services or not.",
			},
			"enforce_safe_search": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that indicates whether only safe content must be returned for web, image, and video search.",
			},
			"enable_office365": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that enables or disables Microsoft Office 365 configuration.",
			},
			"enable_msft_o365": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to permit secure local breakout for Office 365 traffic automatically without any manual configuration needed.",
			},
			"enable_ucaas_zoom": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for Zoom traffic, without any manual configuration needed.",
			},
			"enable_ucaas_logmein": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for GoTo traffic, without any manual configuration needed.",
			},
			"enable_ucaas_ring_central": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for RingCentral traffic, without any manual configuration needed.",
			},
			"enable_ucaas_webex": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for Webex traffic, without any manual configuration needed.",
			},
			"enable_ucaas_talkdesk": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the Zscaler service is allowed to automatically permit secure local breakout for Talkdesk traffic, with minimal or no manual configuration needed.",
			},
			"enable_chatgpt_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with ChatGPT by users should be categorized and logged",
			},
			"enable_microsoft_copilot_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Microsoft Copilot by users should be categorized and logged",
			},
			"enable_gemini_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Google Gemini by users should be categorized and logged",
			},
			"enable_poep_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Poe by users should be categorized and logged",
			},
			"enable_meta_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Meta AI by users should be categorized and logged",
			},
			"enable_per_plexity_prompt": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the use of generative AI prompts with Perplexity by users should be categorized and logged",
			},
			"block_skype": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether access to Skype is blocked or not.",
			},
			"enable_newly_registered_domains": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating whether newly registered and observed domains that are identified within hours of going live are allowed or blocked",
			},
			"enable_block_override_for_non_auth_user": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if authorized users can temporarily override block action on websites by providing their authentication information",
			},
			"enable_cipa_compliance": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value indicating if the predefined CIPA Compliance Rule is enabled or not. ",
			},
		},
	}
}

func dataSourceURLFilteringCloludAppSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		_ = d.Set("enable_poep_prompt", resp.EnablePOEPrompt)
		_ = d.Set("enable_meta_prompt", resp.EnableMetaPrompt)
		_ = d.Set("enable_per_plexity_prompt", resp.EnablePerPlexityPrompt)
		_ = d.Set("block_skype", resp.BlockSkype)
		_ = d.Set("enable_newly_registered_domains", resp.EnableNewlyRegisteredDomains)
		_ = d.Set("enable_block_override_for_non_auth_user", resp.EnableBlockOverrideForNonAuthUser)
		_ = d.Set("enable_cipa_compliance", resp.EnableCIPACompliance)
	} else {
		return diag.FromErr(fmt.Errorf("couldn't read url filtering and cloud app settings"))
	}

	return nil
}
