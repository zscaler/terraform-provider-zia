package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIADeviceGroups_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDeviceGroupsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDeviceGroupsCheck("data.zia_device_groups.ios"),
					testAccDataSourceDeviceGroupsCheck("data.zia_device_groups.android"),
					testAccDataSourceDeviceGroupsCheck("data.zia_device_groups.windows"),
					testAccDataSourceDeviceGroupsCheck("data.zia_device_groups.mac"),
					testAccDataSourceDeviceGroupsCheck("data.zia_device_groups.linux"),
					testAccDataSourceDeviceGroupsCheck("data.zia_device_groups.no_client_connector"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceGroupsCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceDeviceGroupsConfig_basic = `
data "zia_device_groups" "ios"{
    name = "IOS"
}
data "zia_device_groups" "android"{
    name = "Android"
}
data "zia_device_groups" "windows"{
    name = "Windows"
}
data "zia_device_groups" "mac"{
    name = "Mac"
}
data "zia_device_groups" "linux"{
    name = "Linux"
}
data "zia_device_groups" "no_client_connector"{
    name = "No Client Connector"
}
`
