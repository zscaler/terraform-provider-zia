package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccdataSourceUserManagement_Basic(t *testing.T) {
	rName := acctest.RandString(10)
	rComments := acctest.RandString(15)
	resourceName := "data.zia_user_management.test-user-account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccdataSourceUserManagementBasic(rName, rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccdataSourceUserManagement(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "testAcc TF User"),
					resource.TestCheckResourceAttr(resourceName, "email", rName+"@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName, "comments", "test-user-account-"+rComments),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccdataSourceUserManagementBasic(rName, rComments string) string {
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
	name 		= "testAcc TF User"
	email 		= "%s@securitygeek.io"
	password 	= "yty4kuq_dew!eux3AGD-124"
	comments	= "test-user-account-%s"
	groups {
		id = [ data.zia_group_management.normal_internet.id,
			data.zia_group_management.devops.id ]
	}
	department {
		id = data.zia_department_management.engineering.id
	}
}

data "zia_user_management" "test-user-account" {
	name = zia_user_management.test-user-account.name
}
	`, rName, rComments)
}

func testAccdataSourceUserManagement(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
*/
