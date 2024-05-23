package zia

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDevices_Basic(t *testing.T) {
	var initialDevice map[string]string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDevicesConfig_all,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						// Capture the initial device information
						rs, ok := s.RootModule().Resources["data.zia_devices.all"]
						if !ok {
							return fmt.Errorf("Not found: data.zia_devices.all")
						}

						if rs.Primary.ID == "" {
							log.Printf("[INFO] No resource found for data.zia_devices.all; considering this test successful as the tenant might not have this data.")
							return nil
						}

						initialDevice = map[string]string{
							"id":                rs.Primary.Attributes["id"],
							"name":              rs.Primary.Attributes["name"],
							"device_model":      rs.Primary.Attributes["device_model"],
							"os_type":           rs.Primary.Attributes["os_type"],
							"os_version":        rs.Primary.Attributes["os_version"],
							"owner_name":        rs.Primary.Attributes["owner_name"],
							"device_group_type": rs.Primary.Attributes["device_group_type"],
							"description":       rs.Primary.Attributes["description"],
							"owner_user_id":     rs.Primary.Attributes["owner_user_id"],
						}

						return nil
					},
				),
			},
			{
				Config: testAccCheckDataSourceDevicesConfig(initialDevice["device_model"], "device_model"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDevicesCheck("data.zia_devices.device_model"),
				),
			},
			{
				Config: testAccCheckDataSourceDevicesConfig(initialDevice["os_type"], "os_type"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDevicesCheck("data.zia_devices.os_type"),
				),
			},
			{
				Config: testAccCheckDataSourceDevicesConfig(initialDevice["os_version"], "os_version"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDevicesCheck("data.zia_devices.os_version"),
				),
			},
			{
				Config: testAccCheckDataSourceDevicesConfig(initialDevice["name"], "name"),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDevicesCheck("data.zia_devices.name"),
				),
			},
			// {
			// 	Config: testAccCheckDataSourceDevicesConfig(initialDevice["owner_name"], "owner_name"),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccDataSourceDevicesCheck("data.zia_devices.owner"),
			// 	),
			// },
		},
	})
}

func testAccDataSourceDevicesCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			log.Printf("[INFO] No resource found for %s; considering this test successful as the tenant might not have this data.", name)
			return nil
		}

		// Additional checks can go here
		return nil
	}
}

func testAccCheckDataSourceDevicesConfig(value, attribute string) string {
	return fmt.Sprintf(`
data "zia_devices" "%s" {
	%s = "%s"
}
`, attribute, attribute, value)
}

var testAccCheckDataSourceDevicesConfig_all = `
data "zia_devices" "all" {}
`
