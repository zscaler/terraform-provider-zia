package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceZIADepartmentManagement_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDeptMgmtConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDeptManagementCheck("data.zia_department_management.engineering"),
					testAccDataSourceDeptManagementCheck("data.zia_department_management.executives"),
					testAccDataSourceDeptManagementCheck("data.zia_department_management.finance"),
					testAccDataSourceDeptManagementCheck("data.zia_department_management.marketing"),
					testAccDataSourceDeptManagementCheck("data.zia_department_management.sales"),
				),
			},
		},
	})
}

func testAccDataSourceDeptManagementCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceDeptMgmtConfig_basic = `
data "zia_department_management" "engineering"{
    name = "Engineering"
}

data "zia_department_management" "executives"{
    name = "Executives"
}

data "zia_department_management" "finance"{
    name = "Finance"
}

data "zia_department_management" "marketing"{
    name = "Marketing"
}

data "zia_department_management" "sales"{
    name = "Sales"
}
`
