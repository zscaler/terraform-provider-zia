package zia

/*
import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAdminUsersBasic(t *testing.T) {
	var admins admins.AdminUsers
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminUsers)
	rEmail := acctest.RandomWithPrefix("tf-acc-test")
	rPassword := acctest.RandString(10)
	rPasswordUpdate := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, rEmail, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", (generatedName+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", (rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", (rPassword+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "role.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "admin_scope_entities.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, rEmail, rPasswordUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminUsersExists(resourceTypeAndName, &admins),
					resource.TestCheckResourceAttr(resourceTypeAndName, "login_name", (generatedName+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", (rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username", variable.AdminUserName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", (rPasswordUpdate+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "role.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "admin_scope_entities.#", "1"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccCheckAdminUsersDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.AdminUsers {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		admin, err := admins.GetAdminUsers(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if admin != nil {
			return fmt.Errorf("admin user account with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckAdminUsersExists(resource string, admin *admins.AdminUsers) resource.TestCheckFunc {
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
		service := apiClient.Service

		receivedRule, err := admins.GetAdminUsers(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*admin = *receivedRule

		return nil
	}
}

func testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, rEmail, rPassword string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

data "zia_admin_roles" "super_admin"{
	name = "Super Admin"
}

data "zia_department_management" "engineering" {
	name = "Engineering"
  }

  data "zia_department_management" "sales" {
	name = "Sales"
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
    admin_scope_type = "DEPARTMENT"

    admin_scope_entities {
        id = [ data.zia_department_management.engineering.id, data.zia_department_management.sales.id ]
    }
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the admin user
		resourcetype.AdminUsers,
		resourceName,
		generatedName,
		rEmail,
		variable.AdminUserName,
		rPassword,

		// Data source type and name
		resourcetype.AdminUsers,
		resourceName,

		// Reference to the resource
		resourcetype.AdminUsers,
		resourceName,
	)
}
*/
