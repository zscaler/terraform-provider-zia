package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFWTimeWindow_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceFWTimeWindowConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_time_window.work_hours", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_time_window.weekends", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_time_window.off_hours", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceFWTimeWindowConfig_basic = `
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
