package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIAFWTimeWindow_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceFWTimeWindowConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWTimeWindowCheck("data.zia_firewall_filtering_time_window.work_hours"),
					testAccDataSourceFWTimeWindowCheck("data.zia_firewall_filtering_time_window.weekends"),
					testAccDataSourceFWTimeWindowCheck("data.zia_firewall_filtering_time_window.off_hours"),
				),
			},
		},
	})
}

func testAccDataSourceFWTimeWindowCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceFWTimeWindowConfig_basic = `
data "zia_firewall_filtering_time_window" "work_hours"{
    name = "Work hours"
}
data "zia_firewall_filtering_time_window" "weekends"{
    name = "Weekends"
}
data "zia_firewall_filtering_time_window" "off_hours"{
    name = "Off hours"
}
`
