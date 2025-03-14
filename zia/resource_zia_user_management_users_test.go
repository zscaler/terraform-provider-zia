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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/users"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceUserManagementBasic(t *testing.T) {
	var users users.Users
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.Users)
	rEmail := acctest.RandomWithPrefix("tf-acc-test")
	rComments := acctest.RandomWithPrefix("tf-acc-test")

	rPassword := acctest.RandString(10)
	rPasswordUpdate := acctest.RandString(10)
	name := "tf-acc-test " + generatedName
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserManagementConfigure(resourceTypeAndName, name, rEmail, rPassword, rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserManagementExists(resourceTypeAndName, &users),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", name),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", (rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", (rPassword+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", rComments),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "department.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckUserManagementConfigure(resourceTypeAndName, name, rEmail, rPasswordUpdate, rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserManagementExists(resourceTypeAndName, &users),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", name),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", (rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "password", (rPasswordUpdate+"Super@Secret007")),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", rComments),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "department.#", "1"),
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

func testAccCheckUserManagementDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.Users {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		users, err := users.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if users != nil {
			return fmt.Errorf("user account with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckUserManagementExists(resource string, users *users.Users) resource.TestCheckFunc {
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

		receivedUser, err := users.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*users = *receivedUser

		return nil
	}
}

func testAccCheckUserManagementConfigure(resourceTypeAndName, name, rEmail, rPassword, rComments string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

data "zia_group_management" "marketing" {
	name = "Marketing"
}

data "zia_group_management" "sales" {
	name = "Sales"
}

data "zia_department_management" "finance" {
	name = "Finance"
}

resource "%s" "%s" {
	name 		= "%s"
	email 		= "%s@securitygeek.io"
	password 	= "%sSuper@Secret007"
	comments	= "%s"
	groups {
		id = [data.zia_group_management.marketing.id,
		      data.zia_group_management.sales.id ]
	}
	department {
		id = data.zia_department_management.finance.id
	}
	depends_on = [ data.zia_group_management.marketing, data.zia_group_management.sales, data.zia_department_management.finance ]
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.Users,
		resourceName,
		name,
		rEmail,
		rPassword,
		rComments,

		// data source variables
		resourcetype.Users,
		resourceName,

		resourcetype.Users,
		resourceName,
	)
}
*/
