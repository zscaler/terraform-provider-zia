package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDeviceGroups_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDeviceGroupsConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.ios", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.android", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.windows", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.mac", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.linux", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.no_client_connector", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_device_groups.cloud_browser_isolation", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDeviceGroupsConfig_basic = `
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
data "zia_device_groups" "cloud_browser_isolation"{
    name = "Cloud Browser Isolation"
}
`
