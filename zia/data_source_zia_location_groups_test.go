package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIALocationGroup_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLocationGroupConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLocationGroupCheck("data.zia_location_groups.corp_user"),
					testAccDataSourceLocationGroupCheck("data.zia_location_groups.guest_wifi"),
					testAccDataSourceLocationGroupCheck("data.zia_location_groups.iot_traffic"),
					testAccDataSourceLocationGroupCheck("data.zia_location_groups.server_traffic"),
				),
			},
		},
	})
}

func testAccDataSourceLocationGroupCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceLocationGroupConfig_basic = `
data "zia_location_groups" "corp_user"{
    name = "Corporate User Traffic Group"
}
data "zia_location_groups" "guest_wifi"{
    name = "Guest Wifi Group"
}
data "zia_location_groups" "iot_traffic"{
    name = "IoT Traffic Group"
}
data "zia_location_groups" "server_traffic"{
    name = "Server Traffic Group"
}
`
