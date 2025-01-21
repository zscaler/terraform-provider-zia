package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceATPSecurityExceptions_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceATPSecurityExceptionsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceATPSecurityExceptionsCheck("data.zia_atp_security_exceptions.url_list"),
				),
			},
		},
	})
}

func testAccDataSourceATPSecurityExceptionsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "bypass_urls.#"),
	)
}

var testAccCheckDataSourceATPSecurityExceptionsConfig_basic = `
data "zia_atp_security_exceptions" "url_list" {}
`
