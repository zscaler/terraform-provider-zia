package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceURLCategories_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceID := "data.zia_url_categories.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceURLCategoriesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceURLCategories(resourceID),
					// resource.TestCheckResourceAttr(resourceName, "configured_name", "test-url-category-"+rName),
					// resource.TestCheckResourceAttr(resourceName, "description", "test-url-category-"+rDesc),
					// resource.TestCheckResourceAttr(resourceName, "custom_category", "true"),
					resource.TestCheckResourceAttr(resourceID, "configured_name", "test-url-category-"+rName),
					resource.TestCheckResourceAttr(resourceID, "description", "test-url-category-"+rDesc),
					resource.TestCheckResourceAttr(resourceID, "custom_category", "true"),
				),
			},
		},
	})
}

func testAccDataSourceURLCategoriesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zia_url_categories" "test" {
	super_category 		= "USER_DEFINED"
	configured_name 	= "test-url-category-%s"
	description 		= "test-url-category-%s"
	custom_category     = "true"
	keywords            = ["microsoft"]
	db_categorized_urls = [".creditkarma.com", ".youku.com"]
	type                = "URL_CATEGORY"
}

data "zia_url_categories" "test" {
	id = zia_url_categories.test.id
}
	`, rName, rDesc)
}

func testAccDataSourceURLCategories(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
