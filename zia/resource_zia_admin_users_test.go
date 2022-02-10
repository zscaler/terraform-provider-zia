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
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, variable.AdminUserLoginName, variable.AdminUserName, variable.AdminUserEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", variable.AdminUserLoginName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", variable.AdminUserEmail),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", variable.AdminUserPassword),
				),
			},

			// Update test
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, variable.AdminUserLoginName, variable.AdminUserName, variable.AdminUserEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", variable.AdminUserLoginName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", variable.AdminUserEmail),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", variable.AdminUserPassword),
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

func testAccCheckAdminUsersConfigure(resourceTypeAndName, loginname, adminusername, email, password string) string {
	return fmt.Sprintf(`
// admin user resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		AdminUsersResourceHCL(loginname, adminusername, email, password),

		// data source variables
		resourcetype.AdminUsers,
		loginname,
		resourceTypeAndName,
	)
}

func AdminUsersResourceHCL(loginname, adminusername, email, password string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	login_name                      = "%s"
	username                       = "%s"
	email                           = "%s"
	password                        = "%s"
	comments                        = "Administrator Group"
	is_password_login_allowed       = true
	is_security_report_comm_enabled = true
	is_service_update_comm_enabled  = true
	is_product_update_comm_enabled  = true
`,
		// resource variables
		resourcetype.AdminUsers,
		variable.AdminUserLoginName,
		loginname,
		adminusername,
		email,
		password,
	)
}
*/
