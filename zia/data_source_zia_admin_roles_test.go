package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAdminRoles_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAdminRolesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAdminRolesCheck("data.zia_admin_roles.super_admin"),
					testAccDataSourceAdminRolesCheck("data.zia_admin_roles.velocloud"),
					testAccDataSourceAdminRolesCheck("data.zia_admin_roles.silverpeak"),
					testAccDataSourceAdminRolesCheck("data.zia_admin_roles.engineering"),
				),
			},
		},
	})
}

func testAccDataSourceAdminRolesCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceAdminRolesConfig_basic = `
data "zia_admin_roles" "super_admin" {
	name = "Super Admin"
}

data "zia_admin_roles" "velocloud" {
	name = "SDWAN-VeloCloud"
}

data "zia_admin_roles" "silverpeak" {
	name = "SDWAN-SilverPeak"
}

data "zia_admin_roles" "engineering" {
	name = "Engineering_Role"
}
`
