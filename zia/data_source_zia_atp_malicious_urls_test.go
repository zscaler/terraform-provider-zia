package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceATPMaliciousUrls_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceATPMaliciousUrlsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceATPMaliciousUrlsCheck("data.zia_atp_malicious_urls.url_list"),
				),
			},
		},
	})
}

func testAccDataSourceATPMaliciousUrlsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "malicious_urls.#"),
	)
}

var testAccCheckDataSourceATPMaliciousUrlsConfig_basic = `
data "zia_atp_malicious_urls" "url_list" {}
`
