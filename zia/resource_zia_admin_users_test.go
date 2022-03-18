package zia

/*
import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceAdminUsersBasic(t *testing.T) {
	var admins adminuserrolemgmt.AdminUsers
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminUsers)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, variable.AdminUserLoginName, variable.AdminUserName, variable.AdminUserEmail, variable.AdminPasswordLoginAllowed),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", variable.AdminUserLoginName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", variable.AdminUserEmail),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", variable.AdminUserPassword),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_password_login_allowed", strconv.FormatBool(variable.AdminPasswordLoginAllowed)),
				),
			},

			// Update test
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, variable.AdminUserLoginName, variable.AdminUserName, variable.AdminUserEmail, variable.AdminPasswordLoginAllowed),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", variable.AdminUserLoginName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", variable.AdminUserEmail),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", variable.AdminUserPassword),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_password_login_allowed", strconv.FormatBool(variable.AdminPasswordLoginAllowed)),
				),
			},
		},
	})
}

func testAccCheckAdminUsersDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.AdminUsers {
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
		receivedRule, err := apiClient.adminuserrolemgmt.GetAdminUsers(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*admin = *receivedRule

		return nil
	}
}

func testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, loginname, username, email string, password_allowed bool) string {
	return fmt.Sprintf(`
// admin user resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		AdminUsersResourceHCL(generatedName, loginname, username, email, password_allowed),

		// data source variables
		resourcetype.AdminUsers,
		generatedName,
		resourceTypeAndName,
	)
}

func AdminUsersResourceHCL(generatedName, loginname, username, email string, password_allowed bool) string {
	return fmt.Sprintf(`

data "zia_admin_roles" "super_admin" {
	name = "Super Admin"
}

data "zia_location_groups" "corporate_user_traffic_group" {
	name = "Corporate User Traffic Group"
}

resource "%s" "%s" {
	login_name                      = "%s"
	username                        = "%s"
	email                           = "%s"
	password                        = "Password@123!"
	comments                        = "Administrator Group"
	is_password_login_allowed       = "%s"
	is_security_report_comm_enabled = true
	is_service_update_comm_enabled  = true
	is_product_update_comm_enabled  = true
	role {
		id = data.zia_admin_roles.super_admin.id
	}
	admin_scope {
		type = "LOCATION_GROUP"
		scope_entities {
		  id = [data.zia_location_groups.corporate_user_traffic_group.id]
		}
	}
}
`,
		// resource variables
		resourcetype.AdminUsers,
		generatedName,
		variable.AdminUserLoginName,
		variable.AdminUserName,
		variable.AdminUserEmail,
		strconv.FormatBool(password_allowed),
	)
}
*/
