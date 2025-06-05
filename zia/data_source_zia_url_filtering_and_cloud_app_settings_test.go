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
					resource.TestCheckResourceAttrSet(resourceName, "block_skype"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_newly_registered_domains"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_block_override_for_non_auth_user"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_cipa_compliance"),
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
					resource.TestCheckResourceAttrSet(resourceName, "block_skype"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_newly_registered_domains"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_block_override_for_non_auth_user"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_cipa_compliance"),
				),
			},
		},
	})
}

var testAccCheckDataSourceURLFilteringCloludAppSettingsConfig_basic = `
data "zia_url_filtering_and_cloud_app_settings" "this" {}
`
