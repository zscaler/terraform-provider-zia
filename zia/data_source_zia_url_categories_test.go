package zia

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
)

func TestAccDataSourceZIAURLCategories_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.URLCategories)

	skipAcc := os.Getenv("SKIP_URL_CATEGORIES")
	if skipAcc == "yes" {
		t.Skip("Skipping url categories test as SKIP_URL_CATEGORIES is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLCategoriesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLCategoriesConfigure(resourceTypeAndName, generatedName, variable.CustomCategory),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "configured_name", resourceTypeAndName, "configured_name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "keywords", resourceTypeAndName, "keywords"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
