package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIAGroupManagement_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGroupManagementConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGroupManagementCheck("data.zia_group_management.devops"),
					testAccDataSourceGroupManagementCheck("data.zia_group_management.executives"),
					testAccDataSourceGroupManagementCheck("data.zia_group_management.marketing"),
				),
			},
		},
	})
}

func testAccDataSourceGroupManagementCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceGroupManagementConfig_basic = `
data "zia_group_management" "devops"{
    name = "DevOps"
}

data "zia_group_management" "executives"{
    name = "Executives"
}

data "zia_group_management" "marketing"{
    name = "Marketing"
}
`
