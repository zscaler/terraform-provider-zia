package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceATPSecurityExceptions_basic(t *testing.T) {
	resourceName := "zia_atp_security_exceptions.test"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckATPSecurityExceptionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceATPSecurityExceptionsConfig([]string{".example.com", ".test.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckATPSecurityExceptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bypass_urls.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "bypass_urls.0", ".example.com"),
					resource.TestCheckResourceAttr(resourceName, "bypass_urls.1", ".test.com"),
				),
			},
			{
				Config: testAccResourceATPSecurityExceptionsConfig([]string{".newexample.com"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckATPSecurityExceptionsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bypass_urls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "bypass_urls.0", ".newexample.com"),
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

func testAccCheckATPSecurityExceptionsDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

func testAccCheckATPSecurityExceptionsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ATP Security Exception URL ID is set")
		}
		return nil
	}
}

func testAccResourceATPSecurityExceptionsConfig(domains []string) string {
	config := `resource "zia_atp_security_exceptions" "test" { bypass_urls = [`
	for _, domain := range domains {
		config += `"` + domain + `",`
	}
	config = config[:len(config)-1] // Remove the trailing comma
	config += `] }`

	return config
}
