package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
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
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "email", resourceTypeAndName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comments", resourceTypeAndName, "comments"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "department.#", "1"),
				),
			},
		},
	})
}
