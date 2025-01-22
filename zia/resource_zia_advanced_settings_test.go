package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAdvancedSettings_Basic(t *testing.T) {
	resourceName := "zia_advanced_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceAdvancedSettingsDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with specific values
			{
				Config: testAccResourceAdvancedSettingsConfig(
					false, false, false, false, false, false, false, false, // blocked attributes
					false, 36000), // capture attributes
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_dns_resolution_on_transparent_proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_office365", "false"),
					resource.TestCheckResourceAttr(resourceName, "log_internal_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "enforce_surrogate_ip_for_windows_app", "false"),
					resource.TestCheckResourceAttr(resourceName, "track_http_tunnel_on_http_ports", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_http_tunnel_on_non_http_ports", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_domain_fronting_on_host_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "cascade_url_filtering", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_non_compliant_http_request_on_http_ports", "false"),
					resource.TestCheckResourceAttr(resourceName, "ui_session_timeout", "36000"),
				),
			},
			// Step 2: Update the resource with new values
			{
				Config: testAccResourceAdvancedSettingsConfig(
					true, true, true, true, true, false, false, true, // blocked attributes
					true, 36000), // capture attributes
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_dns_resolution_on_transparent_proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_office365", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_internal_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "enforce_surrogate_ip_for_windows_app", "true"),
					resource.TestCheckResourceAttr(resourceName, "track_http_tunnel_on_http_ports", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_http_tunnel_on_non_http_ports", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_domain_fronting_on_host_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "cascade_url_filtering", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_non_compliant_http_request_on_http_ports", "true"),
					resource.TestCheckResourceAttr(resourceName, "ui_session_timeout", "36000"),
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

func testAccCheckResourceAdvancedSettingsDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

// Helper function to generate test configuration for the resource
func testAccResourceAdvancedSettingsConfig(
	attr1, attr2, attr3, attr4, attr5, attr6, attr7, attr8, attr9 bool, attr10 int, // blocked attributes
) string {
	return fmt.Sprintf(`
resource "zia_advanced_settings" "test" {
  enable_dns_resolution_on_transparent_proxy    	= %t
  enable_office365                					= %t
  log_internal_ip 									= %t
  enforce_surrogate_ip_for_windows_app 				= %t
  track_http_tunnel_on_http_ports               	= %t
  block_http_tunnel_on_non_http_ports              	= %t
  block_domain_fronting_on_host_header              = %t
  cascade_url_filtering                 			= %t
  block_non_compliant_http_request_on_http_ports    = %t
  ui_session_timeout               					= %d
}
`,
		attr1, attr2,
		attr3, attr4,
		attr5, attr6,
		attr7, attr8,
		attr9, attr10)
}
