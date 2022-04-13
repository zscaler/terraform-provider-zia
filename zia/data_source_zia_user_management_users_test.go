package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
)

func TestAccDataSourceUserManagement_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.Users)
	rEmail := acctest.RandomWithPrefix("tf-acc-test")
	rComments := acctest.RandomWithPrefix("tf-acc-test")
	rPassword := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUserManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserManagementConfigure(resourceTypeAndName, generatedName, rEmail, rPassword, rComments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "email", resourceTypeAndName, fmt.Sprintf(rEmail+"@securitygeek.io")),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comments", resourceTypeAndName, fmt.Sprintf(rComments+"tf-acc-test")),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "password", resourceTypeAndName, fmt.Sprintf(rPassword+"Super@Secret007")),
					// resource.TestCheckResourceAttr(dataSourceTypeAndName, "groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "department.#", "1"),
				),
			},
		},
	})
}
*/
