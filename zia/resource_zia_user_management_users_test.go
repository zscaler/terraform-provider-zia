package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func TestAccResourceUserManagement_basic(t *testing.T) {
	var users usermanagement.Users
	rComments := acctest.RandString(5)
	rPassword := acctest.RandString(20)
	resourceName := "zia_user_management.test-user-account"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceUserManagementBasic(rComments, rPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserManagementExists("zia_user_management.test-user-account", &users),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc TF User"),
					resource.TestCheckResourceAttr(resourceName, "email", "test-user-account@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName, "comments", "test-user-account-"+rComments),
					resource.TestCheckResourceAttr(resourceName, "password", "yty4kuq_dew!eux3AGD-"+rPassword),
				),
			},
		},
	})
}

func testAccCheckResourceUserManagementBasic(rComments, rPassword string) string {
	return fmt.Sprintf(`

data "zia_group_management" "normal_internet" {
	name = "Normal_Internet"
}

data "zia_group_management" "devops" {
	name = "DevOps"
}

data "zia_department_management" "engineering" {
	name = "Engineering"
}

resource "zia_user_management" "test-user-account" {
	name = "testAcc TF User"
	email = "test-user-account@securitygeek.io"
	password = "yty4kuq_dew!eux3AGD-%s"
	comments = "test-user-account-%s"
	groups {
	 id = [ data.zia_group_management.normal_internet.id,
			data.zia_group_management.devops.id ]
	 }
	department {
	 id = data.zia_department_management.engineering.id
	 }
}
	`, rPassword, rComments)
}

func testAccCheckUserManagementExists(resource string, users *usermanagement.Users) resource.TestCheckFunc {
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
		receivedAccount, err := apiClient.usermanagement.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*users = *receivedAccount

		return nil
	}
}

func testAccCheckAdminUsersDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_user_management" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		admin, err := apiClient.usermanagement.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if admin != nil {
			return fmt.Errorf("user account with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
