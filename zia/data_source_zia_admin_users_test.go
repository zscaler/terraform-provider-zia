package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAdminUsers_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceAdminUsersConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_admin_users.foobar", "login_name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceAdminUsersConfig_basic = `
data "zia_admin_users" "foobar" {
    login_name = "admin@24326813.zscalerthree.net"
}`
