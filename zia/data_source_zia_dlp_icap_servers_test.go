package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIADLPICAPServers_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDLPICAPServersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDLPICAPServersCheck("data.zia_dlp_icap_servers.icap_01"),
				),
			},
		},
	})
}

func testAccDataSourceDLPICAPServersCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceDLPICAPServersConfig_basic = `
data "zia_dlp_icap_servers" "icap_01"{
    name = "ZS_BD_ICAP_01"
}
`
