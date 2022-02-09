package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLocationGroup_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceLocationGroupConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_location_groups.example1", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_location_groups.example2", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_location_groups.example3", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_location_groups.example4", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceLocationGroupConfig_basic = `
data "zia_location_groups" "example1"{
    name = "Corporate User Traffic Group"
}

data "zia_location_groups" "example2"{
    name = "Guest Wifi Group"
}

data "zia_location_groups" "example3"{
    name = "IoT Traffic Group"
}

data "zia_location_groups" "example4"{
    name = "Server Traffic Group"
}
`
