package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAdminUsers_Basic(t *testing.T) {
	rComments := acctest.RandString(5)
	resourceName1 := "data.zia_admin_users.test-admin-loginname"
	resourceName2 := "data.zia_admin_users.test-admin-username"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAdminUsersBasic(rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAdminUsersLoginName(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "login_name", "test-admin-user@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName1, "username", "testAcc Tf Admin"),
					resource.TestCheckResourceAttr(resourceName1, "email", "test-admin-user@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName1, "comments", "test-admin-user-"+rComments),
					resource.TestCheckResourceAttr(resourceName1, "is_password_login_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName1, "is_security_report_comm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName1, "is_service_update_comm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName1, "is_exec_mobile_app_enabled", "false"),
				),
			},
			{
				Config: testAccCheckDataSourceAdminUsersBasic(rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceAdminUsersUsername(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "login_name", "test-admin-user@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName2, "username", "testAcc Tf Admin"),
					resource.TestCheckResourceAttr(resourceName2, "email", "test-admin-user@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName2, "comments", "test-admin-user-"+rComments),
					resource.TestCheckResourceAttr(resourceName2, "is_password_login_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName2, "is_security_report_comm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "is_service_update_comm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "is_exec_mobile_app_enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckDataSourceAdminUsersBasic(rComments string) string {
	return fmt.Sprintf(`

data "zia_admin_roles" "super_admin" {
	name = "Super Admin"
}

resource "zia_admin_users" "test-admin-user" {
	login_name                      = "test-admin-user@securitygeek.io"
	username                        = "testAcc Tf Admin"
	email                           = "test-admin-user@securitygeek.io"
	comments                        = "test-admin-user-%s"
	password                        = "yty4kuq_dew!eux3AGD"
	is_password_login_allowed       = true
	is_security_report_comm_enabled = true
	is_service_update_comm_enabled  = true
	is_exec_mobile_app_enabled      = false
	role {
		id = data.zia_admin_roles.super_admin.id
	}
	admin_scope {
		type = "ORGANIZATION"
	}
}

data "zia_admin_users" "test-admin-loginname" {
	login_name = zia_admin_users.test-admin-user.login_name
	depends_on = [ zia_admin_users.test-admin-user ]
}

data "zia_admin_users" "test-admin-username" {
	username = zia_admin_users.test-admin-user.username
	depends_on = [ zia_admin_users.test-admin-user ]
}
	`, rComments)
}

func testAccDataSourceAdminUsersLoginName(login_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[login_name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", login_name)
		}

		return nil
	}
}

func testAccDataSourceAdminUsersUsername(username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[username]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", username)
		}

		return nil
	}
}
*/
