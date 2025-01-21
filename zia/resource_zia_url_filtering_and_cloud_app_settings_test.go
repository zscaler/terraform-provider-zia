package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceURLFilteringCloludAppSettings_Basic(t *testing.T) {
	resourceName := "zia_url_filtering_and_cloud_app_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with specific values
			{
				Config: testAccResourceURLFilteringCloludAppSettingsConfig(
					false, false, false, true, false, false, false, false, // blocked attributes
					false, false, false, false, false, false, false, false, false, false), // capture attributes
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_dynamic_content_cat", "false"),
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
					resource.TestCheckResourceAttr(resourceName, "enable_poep_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_meta_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_per_plexity_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_skype", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_newly_registered_domains", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_cipa_compliance", "false"),
				),
			},
			// Step 2: Update the resource with new values
			{
				Config: testAccResourceURLFilteringCloludAppSettingsConfig(
					true, false, false, true, false, false, false, false, // blocked attributes
					false, false, false, false, false, false, false, false, false, false), // capture attributes
				Check: resource.ComposeTestCheckFunc(
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
					resource.TestCheckResourceAttr(resourceName, "enable_poep_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_meta_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_per_plexity_prompt", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_skype", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_newly_registered_domains", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_cipa_compliance", "false"),
				),
			},
			// Step 3: Import the resource and verify the state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper function to generate test configuration for the resource
func testAccResourceURLFilteringCloludAppSettingsConfig(
	enableDynamicContentCat, considerEmbeddedSites, enforceSafeSearch, enableOffice365, enableMsftO365,
	enableUcaasZoom, enableUcaasLogmein, enableUcaasRingCentral, enableUcaasWebex, enableUcaasTalkdesk,
	enableChatGPTPrompt, enableMicrosoftCopilotPrompt, enablePoepPrompt, enableMetaPrompt, enablePerPlexityPrompt,
	blockSkype, enableNewlyRegisteredDomains, enableCipaCompliance bool,
) string {
	return fmt.Sprintf(`
resource "zia_url_filtering_and_cloud_app_settings" "test" {
    enable_dynamic_content_cat              = %t
    consider_embedded_sites                 = %t
    enforce_safe_search                     = %t
    enable_office365                        = %t
    enable_msft_o365                        = %t
    enable_ucaas_zoom                       = %t
    enable_ucaas_logmein                    = %t
    enable_ucaas_ring_central               = %t
    enable_ucaas_webex                      = %t
    enable_ucaas_talkdesk                   = %t
    enable_chatgpt_prompt                   = %t
    enable_microsoft_copilot_prompt         = %t
    enable_poep_prompt                      = %t
    enable_meta_prompt                      = %t
    enable_per_plexity_prompt               = %t
    block_skype                             = %t
    enable_newly_registered_domains         = %t
    enable_cipa_compliance                  = %t
}
`,
		enableDynamicContentCat, considerEmbeddedSites, enforceSafeSearch, enableOffice365, enableMsftO365,
		enableUcaasZoom, enableUcaasLogmein, enableUcaasRingCentral, enableUcaasWebex, enableUcaasTalkdesk,
		enableChatGPTPrompt, enableMicrosoftCopilotPrompt, enablePoepPrompt, enableMetaPrompt, enablePerPlexityPrompt,
		blockSkype, enableNewlyRegisteredDomains, enableCipaCompliance,
	)
}
