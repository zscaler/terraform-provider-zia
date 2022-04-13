package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDevices_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDevicesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDevicesCheck("data.zia_devices.device_model"),
					testAccDataSourceDevicesCheck("data.zia_devices.os_type"),
					testAccDataSourceDevicesCheck("data.zia_devices.os_version"),
				),
			},
		},
	})
}

func testAccDataSourceDevicesCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceDevicesConfig_basic = `
data "zia_devices" "device_model"{
    device_model = "VMware Virtual Platform"
}
data "zia_devices" "os_type"{
    os_type = "WINDOWS_OS"
}
data "zia_devices" "os_version"{
    os_version = "Microsoft Windows 10 Pro;64 bit"
}
`
