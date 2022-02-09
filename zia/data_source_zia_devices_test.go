package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDevices_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDevicesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_devices.device_model", "device_model"),
					resource.TestCheckResourceAttrSet(
						"data.zia_devices.os_type", "os_type"),
					resource.TestCheckResourceAttrSet(
						"data.zia_devices.os_version", "os_version"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDevicesConfig_basic = `
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
