package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLocationLite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLocationLiteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLocationLiteCheck("data.zia_location_lite.road_warrior"),
				),
			},
		},
	})
}

func testAccDataSourceLocationLiteCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceLocationLiteConfig_basic = `
data "zia_location_lite" "road_warrior"{
    name = "Road Warrior"
}
`
