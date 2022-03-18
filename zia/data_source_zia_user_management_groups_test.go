package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroupManagement_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceGroupManagementConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_group_management.devops", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_group_management.executives", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_group_management.sales", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_group_management.marketing", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceGroupManagementConfig_basic = `
data "zia_group_management" "devops"{
    name = "DevOps"
}
data "zia_group_management" "executives"{
    name = "Executives"
}
data "zia_group_management" "sales"{
    name = "Sales"
}
data "zia_group_management" "marketing"{
    name = "Marketing"
}
`
