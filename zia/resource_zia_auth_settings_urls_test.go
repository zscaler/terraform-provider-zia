package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAuthSettingsUrls_basic(t *testing.T) {
	resourceName := "zia_auth_settings_urls.test"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthSettingsUrlsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAuthSettingsUrlsConfig([]string{".example.com", ".test.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthSettingsUrlsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "urls.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "urls.0", ".example.com"),
					resource.TestCheckResourceAttr(resourceName, "urls.1", ".test.com"),
				),
			},
			{
				Config: testAccResourceAuthSettingsUrlsConfig([]string{".newexample.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAuthSettingsUrlsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "urls.0", ".newexample.com"),
				),
			},
		},
	})
}

func testAccCheckAuthSettingsUrlsDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

func testAccCheckAuthSettingsUrlsExists(n string) resource.TestCheckFunc {
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

func testAccResourceAuthSettingsUrlsConfig(domains []string) string {
	config := `resource "zia_auth_settings_urls" "test" { urls = [`
	for _, domain := range domains {
		config += `"` + domain + `",`
	}
	config = config[:len(config)-1] // Remove the trailing comma
	config += `] }`

	return config
}
