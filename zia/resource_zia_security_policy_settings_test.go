package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceSecurityPolicySettings_basic(t *testing.T) {
	resourceName := "zia_security_settings.test"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityPolicySettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityPolicySettingsConfig([]string{".example.com"}, []string{".blockme.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityPolicySettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "whitelist_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_urls.0", ".example.com"),
					resource.TestCheckResourceAttr(resourceName, "blacklist_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "blacklist_urls.0", ".blockme.com"),
				),
			},
			{
				Config: testAccResourceSecurityPolicySettingsConfig([]string{".newexample.com"}, []string{".blocknew.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityPolicySettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "whitelist_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "whitelist_urls.0", ".newexample.com"),
					resource.TestCheckResourceAttr(resourceName, "blacklist_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "blacklist_urls.0", ".blocknew.com"),
				),
			},
			// Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSecurityPolicySettingsDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

func testAccCheckSecurityPolicySettingsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Security Policy Settings ID is set")
		}
		return nil
	}
}

func testAccResourceSecurityPolicySettingsConfig(whitelistDomains []string, blacklistDomains []string) string {
	whitelist := ""
	for _, domain := range whitelistDomains {
		whitelist += `"` + domain + `",`
	}
	blacklist := ""
	for _, domain := range blacklistDomains {
		blacklist += `"` + domain + `",`
	}
	whitelist = whitelist[:len(whitelist)-1] // Remove the trailing comma
	blacklist = blacklist[:len(blacklist)-1] // Remove the trailing comma

	config := `resource "zia_security_settings" "test" {
		whitelist_urls = [` + whitelist + `]
		blacklist_urls = [` + blacklist + `]
	}`

	return config
}
