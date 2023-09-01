package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/adminuserrolemgmt"
)

func TestAccResourceAdminUsersBasic(t *testing.T) {
	var admins adminuserrolemgmt.AdminUsers
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminUsers)
	rEmail := acctest.RandomWithPrefix("tf-acc-test")
	rPassword := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, rEmail, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", fmt.Sprintf(rPassword+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "role.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "admin_scope.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, rEmail, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", rPassword+("Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "role.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "admin_scope.#", "1"),
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

func testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, rEmail, rPassword string) string {
	return fmt.Sprintf(`

data "zia_admin_roles" "super_admin"{
	name = "Super Admin"
}

resource "%s" "%s" {
	login_name                      = "%s@securitygeek.io"
	email                           = "%s@securitygeek.io"
	username                        = "%s"
	password                        = "%sSuper@Secret007"
	comments                        = "Administrator Group"
	is_password_login_allowed       = true
	is_security_report_comm_enabled = true
	is_service_update_comm_enabled  = true
	is_product_update_comm_enabled  = true
	role {
		id = data.zia_admin_roles.super_admin.id
	}
	admin_scope {
		type = "ORGANIZATION"
	}
}

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		resourcetype.AdminUsers,
		generatedName,
		rEmail,
		rEmail,
		variable.AdminUserName,
		rPassword,

		// data source variables
		resourcetype.AdminUsers,
		generatedName,
		resourceTypeAndName,
	)
}
