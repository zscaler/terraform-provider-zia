package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTrafficFWGreInternalIPRangeList_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceTrafficFWGreInternalIPRangeList_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_gre_internal_ip_range_list.internal_gre", "required_count"),
				),
			},
		},
	})
}

const testAccCheckDataSourceTrafficFWGreInternalIPRangeList_basic = `
data "zia_gre_internal_ip_range_list" "internal_gre"{
    required_count = 10
}
`
