package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
)

func TestAccResourceAdminUsers_basic(t *testing.T) {
	var admins adminuserrolemgmt.AdminUsers
	rComments := acctest.RandString(5)
	rEmail := acctest.RandString(5)
	rPassword := acctest.RandString(20)
	resourceName := "zia_admin_users.test-admin-account"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceAdminUsersBasic(rEmail, rComments, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists("zia_admin_users.test-admin-account", &admins),
					resource.TestCheckResourceAttr(resourceName, "login_name", "test-admin-"+rEmail+"@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName, "username", "testAcc Tf Admin"),
					resource.TestCheckResourceAttr(resourceName, "email", "test-admin-"+rEmail+"@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName, "comments", "test-admin-account-"+rComments),
					resource.TestCheckResourceAttr(resourceName, "password", "yty4kuq_dew!eux3AGD-"+rPassword),
					resource.TestCheckResourceAttr(resourceName, "is_password_login_allowed", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_security_report_comm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_service_update_comm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_exec_mobile_app_enabled", "false"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckResourceAdminUsersBasic(rEmail, rComments, rPassword string) string {
	return fmt.Sprintf(`

data "zia_admin_roles" "super_admin" {
	name = "Super Admin"
}

resource "zia_admin_users" "test-admin-account" {
	login_name                      = "test-admin-%s@securitygeek.io"
	username                        = "testAcc Tf Admin"
	email                           = "test-admin-%s@securitygeek.io"
	comments                        = "test-admin-account-%s"
	password                        = "yty4kuq_dew!eux3AGD-%s"
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
	`, rEmail, rEmail, rComments, rPassword)
}

func testAccCheckAdminUsersExists(resource string, admin *adminuserrolemgmt.AdminUsers) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedAccount, err := apiClient.adminuserrolemgmt.GetAdminUsers(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*admin = *receivedAccount

		return nil
	}
}

func testAccCheckAdminUsersDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_admin_users" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		admin, err := apiClient.adminuserrolemgmt.GetAdminUsers(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if admin != nil {
			return fmt.Errorf("admin user account with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
