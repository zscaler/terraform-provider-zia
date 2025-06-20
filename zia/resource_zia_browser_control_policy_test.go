package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceBrowserControlPolicy_Basic(t *testing.T) {
	resourceName := "zia_browser_control_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceBrowserControlPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with minimal configuration
			{
				Config: testAccResourceBrowserControlPolicyConfig(
					"DAILY", // plugin_check_frequency
					true,    // bypass_all_browsers
					false,   // allow_all_browsers
					true,    // enable_warnings
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "plugin_check_frequency", "DAILY"),
					// Skip asserting bypass_all_browsers — API may return true as default
					resource.TestCheckResourceAttr(resourceName, "bypass_all_browsers", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_all_browsers", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_warnings", "true"),
				),
			},
			// Step 2: Update the resource with some values
			{
				Config: testAccResourceBrowserControlPolicyConfig(
					"DAILY", // plugin_check_frequency
					true,    // bypass_all_browsers
					false,   // allow_all_browsers
					true,    // enable_warnings
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "plugin_check_frequency", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "bypass_all_browsers", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_all_browsers", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_warnings", "true"),
				),
			},
			// Step 3: Import the resource and verify the state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Skip checking specific attributes during import since the API may return defaults
				ImportStateVerifyIgnore: []string{
					"bypass_plugins",
					"bypass_applications",
					"blocked_internet_explorer_versions",
					"blocked_chrome_versions",
					"blocked_firefox_versions",
					"blocked_safari_versions",
					"blocked_opera_versions",
				},
			},
		},
	})
}

func testAccCheckResourceBrowserControlPolicyDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

// Simplified helper function with only required parameters
func testAccResourceBrowserControlPolicyConfig(
	pluginCheckFrequency string,
	bypassAllBrowsers bool,
	allowAllBrowsers bool,
	enableWarnings bool,
) string {
	return fmt.Sprintf(`
resource "zia_browser_control_policy" "test" {
	plugin_check_frequency = %q
	bypass_all_browsers = %t
	allow_all_browsers = %t
	enable_warnings = %t
}
`,
		pluginCheckFrequency,
		bypassAllBrowsers,
		allowAllBrowsers,
		enableWarnings,
	)
}
