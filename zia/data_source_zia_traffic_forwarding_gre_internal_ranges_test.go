package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTrafficGreInternalIPRangeList_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceTrafficGreInternalIPRangeList_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficGreInternalIPRangeListCheck("data.zia_gre_internal_ip_range_list.internal_gre"),
				),
			},
		},
	})
}

func testAccDataSourceTrafficGreInternalIPRangeListCheck(required_count string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(required_count, "id"),
		resource.TestCheckResourceAttrSet(required_count, "required_count"),
	)
}

var testAccCheckDataSourceTrafficGreInternalIPRangeList_basic = `
data "zia_gre_internal_ip_range_list" "internal_gre"{
    required_count = 10
}
`
