package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceATPMaliciousUrls_basic(t *testing.T) {
	resourceName := "zia_atp_malicious_urls.test"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckATPMaliciousUrlsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceATPMaliciousUrlsConfig([]string{".example.com", ".test.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckATPMaliciousUrlsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "malicious_urls.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "malicious_urls.0", ".example.com"),
					resource.TestCheckResourceAttr(resourceName, "malicious_urls.1", ".test.com"),
				),
			},
			{
				Config: testAccResourceATPMaliciousUrlsConfig([]string{".newexample.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckATPMaliciousUrlsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "malicious_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "malicious_urls.0", ".newexample.com"),
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

func testAccCheckATPMaliciousUrlsDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

func testAccCheckATPMaliciousUrlsExists(n string) resource.TestCheckFunc {
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

func testAccResourceATPMaliciousUrlsConfig(domains []string) string {
	config := `resource "zia_atp_malicious_urls" "test" { malicious_urls = [`
	for _, domain := range domains {
		config += `"` + domain + `",`
	}
	config = config[:len(config)-1] // Remove the trailing comma
	config += `] }`

	return config
}
