package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIAAuthSettingsUrls_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAuthSettingsUrlsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAuthSettingsUrlsCheck("data.zia_auth_settings_urls.all_urls"),
				),
			},
		},
	})
}

func testAccDataSourceAuthSettingsUrlsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc()
}

var testAccCheckDataSourceAuthSettingsUrlsConfig_basic = `
data "zia_auth_settings_urls" "all_urls" {}
`
