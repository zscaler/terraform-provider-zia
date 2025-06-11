package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceFTPControlPolicy_Basic(t *testing.T) {
	resourceName := "zia_ftp_control_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceFTPControlPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with specific values (disabled state)
			{
				Config: testAccResourceFTPControlPolicyConfig(false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ftp_over_http_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ftp_enabled", "false"),
				),
			},
			// Step 2: Update the resource with full configuration (enabled + categories + urls)
			{
				Config: testAccResourceFTPControlPolicyConfig(true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "ftp_over_http_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ftp_enabled", "true"),

					resource.TestCheckResourceAttr(resourceName, "url_categories.#", "7"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "HOBBIES_AND_LEISURE"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "HEALTH"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "HISTORY"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "INSURANCE"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "IMAGE_HOST"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "INTERNET_SERVICES"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "GOVERNMENT"),

					resource.TestCheckResourceAttr(resourceName, "urls.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "urls.*", "test1.acme.com"),
					resource.TestCheckTypeSetElemAttr(resourceName, "urls.*", "test10.acme.com"),
				),
			},
			// Step 3: Import the resource and verify state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckResourceFTPControlPolicyDestroy(s *terraform.State) error {
	// Add destroy validation here if needed
	return nil
}

func testAccResourceFTPControlPolicyConfig(ftpEnabled, ftpOverHttpEnabled bool) string {
	base := fmt.Sprintf(`
resource "zia_ftp_control_policy" "test" {
  ftp_enabled            = %t
  ftp_over_http_enabled  = %t
`, ftpEnabled, ftpOverHttpEnabled)

	// only include categories + urls if FTP is enabled
	if ftpEnabled || ftpOverHttpEnabled {
		base += `
  url_categories         = ["HOBBIES_AND_LEISURE","HEALTH","HISTORY","INSURANCE","IMAGE_HOST","INTERNET_SERVICES","GOVERNMENT"]
  urls                   = ["test1.acme.com", "test10.acme.com"]
`
	}

	base += "}\n"
	return base
}
