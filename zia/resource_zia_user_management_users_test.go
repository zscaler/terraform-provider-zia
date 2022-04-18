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
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
)

func TestAccResourceUserManagementBasic(t *testing.T) {
	var users usermanagement.Users
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.Users)
	rEmail := acctest.RandomWithPrefix("tf-acc-test")
	rComments := acctest.RandomWithPrefix("tf-acc-test")

	rPassword := acctest.RandString(10)
	name := "testAcc TF User " + generatedName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserManagementConfigure(resourceTypeAndName, generatedName, name, rEmail, rPassword, rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserManagementExists(resourceTypeAndName, &users),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", name),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", fmt.Sprintf(rPassword+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", rComments),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "department.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckUserManagementConfigure(resourceTypeAndName, generatedName, name, rEmail, rPassword, rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserManagementExists(resourceTypeAndName, &users),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", name),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", fmt.Sprintf(rPassword+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", rComments),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "groups.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "department.#", "1"),
				),
			},
		},
	})
}

func testAccCheckUserManagementDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.Users {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		users, err := apiClient.usermanagement.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if users != nil {
			return fmt.Errorf("user account with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
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
		receivedUser, err := apiClient.usermanagement.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*users = *receivedUser

		return nil
	}
}

func testAccCheckUserManagementConfigure(resourceTypeAndName, generatedName, name, rEmail, rPassword, rComments string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name 		= "%s"
	email 		= "%s@securitygeek.io"
	password 	= "%sSuper@Secret007"
	comments	= "%s"
	groups {
		id = 26348357
	}
	groups {
		id = 24392492
	}
	department {
		id = 25684245
	}
}

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		resourcetype.Users,
		generatedName,
		name,
		rEmail,
		rPassword,
		rComments,

		// data source variables
		resourcetype.Users,
		rEmail,
		resourceTypeAndName,
	)
}
