package zia

/*
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
)

func TestAccDataSourceAdminUsers_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminUsers)
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
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "login_name", resourceTypeAndName, rEmail+"@securitygeek.io"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "email", resourceTypeAndName, rEmail+"@securitygeek.io"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "username", resourceTypeAndName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "password", resourceTypeAndName, rPassword+("Super@Secret007")),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "role.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "admin_scope.#", "1"),
				),
			},
		},
	})
}
*/
