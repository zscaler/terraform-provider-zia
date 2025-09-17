package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFTPControlPolicy_Basic(t *testing.T) {
	resourceName := "data.zia_ftp_control_policy.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "zia_ftp_control_policy" "this" {}
`,
				Check: resource.ComposeTestCheckFunc(
					// Boolean flags
					resource.TestCheckResourceAttr(resourceName, "ftp_over_http_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ftp_enabled", "true"),

					// Set element counts
					resource.TestCheckResourceAttr(resourceName, "url_categories.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "urls.#", "2"),

					// Set element values (unordered)
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "PROFESSIONAL_SERVICES"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "AI_ML_APPS"),
					resource.TestCheckTypeSetElemAttr(resourceName, "url_categories.*", "GENERAL_AI_ML"),

					resource.TestCheckTypeSetElemAttr(resourceName, "urls.*", "test10.acme.com"),
				),
			},
		},
	})
}
