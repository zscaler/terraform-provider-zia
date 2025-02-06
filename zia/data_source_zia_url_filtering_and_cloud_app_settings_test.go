package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceURLFilteringCloudAppSettings_Basic(t *testing.T) {
	resourceName := "data.zia_url_filtering_and_cloud_app_settings.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceURLFilteringCloludAppSettingsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// Ensure all attributes are set correctly
					resource.TestCheckResourceAttrSet(resourceName, "enable_dynamic_content_cat"),
					resource.TestCheckResourceAttrSet(resourceName, "consider_embedded_sites"),
					resource.TestCheckResourceAttrSet(resourceName, "enforce_safe_search"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_office365"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_msft_o365"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_ucaas_zoom"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_ucaas_logmein"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_ucaas_ring_central"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_ucaas_webex"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_ucaas_talkdesk"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_chatgpt_prompt"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_microsoft_copilot_prompt"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_gemini_prompt"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_poep_prompt"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_meta_prompt"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_per_plexity_prompt"),
					// resource.TestCheckResourceAttrSet(resourceName, "block_skype"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_newly_registered_domains"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_block_override_for_non_auth_user"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_cipa_compliance"),

					// Verify specific values
					resource.TestCheckResourceAttr(resourceName, "enable_dynamic_content_cat", "true"),
					resource.TestCheckResourceAttr(resourceName, "consider_embedded_sites", "false"),
					resource.TestCheckResourceAttr(resourceName, "enforce_safe_search", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_office365", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_msft_o365", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_ucaas_zoom", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_ucaas_logmein", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_ucaas_ring_central", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_ucaas_webex", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_ucaas_talkdesk", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_chatgpt_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_microsoft_copilot_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_gemini_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_poep_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_meta_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_per_plexity_prompt", "false"),
					// resource.TestCheckResourceAttr(resourceName, "block_skype", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_newly_registered_domains", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_block_override_for_non_auth_user", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_cipa_compliance", "false"),
				),
			},
		},
	})
}

var testAccCheckDataSourceURLFilteringCloludAppSettingsConfig_basic = `
data "zia_url_filtering_and_cloud_app_settings" "this" {}
`
