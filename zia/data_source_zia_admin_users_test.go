package zia

/*
import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccDataSourceAdminUsers_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminUsers)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminUsersDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminUsersConfigure(resourceTypeAndName, generatedName, variable.AdminUserLoginName, variable.AdminUserName, variable.AdminUserEmail, variable.AdminPasswordLoginAllowed),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "login_name", resourceTypeAndName, "login_name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "username", resourceTypeAndName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "email", resourceTypeAndName, "email"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "is_password_login_allowed", strconv.FormatBool(variable.AdminPasswordLoginAllowed)),
				),
			},
		},
	})
}
*/
