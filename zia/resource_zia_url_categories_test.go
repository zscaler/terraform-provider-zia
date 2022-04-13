package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlcategories"
)

func TestAccResourceURLCategoriesBasic(t *testing.T) {
	var categories urlcategories.URLCategory
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_url_categories.test-url-category"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLCategoriesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLCategoriesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesExists("zia_url_categories.test-url-category", &categories),
					resource.TestCheckResourceAttr(resourceName, "configured_name", "test-url-category-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-url-category-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "custom_category", "true"),
				),
			},
		},
	})
}

func testAccCheckURLCategoriesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zia_url_categories" "test-url-category" {
	super_category 		= "USER_DEFINED"
	configured_name 	= "test-url-category-%s"
	description 		= "test-url-category-%s"
	custom_category     = "true"
	keywords            = ["microsoft"]
	db_categorized_urls = [".creditkarma.com", ".youku.com"]
	type                = "URL_CATEGORY"
}
`, rName, rDesc)
}

func testAccCheckURLCategoriesExists(resource string, rule *urlcategories.URLCategory) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedRule, err := apiClient.urlcategories.Get(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckURLCategoriesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_url_categories" {
			continue
		}

		rule, err := apiClient.urlcategories.Get(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("id %s already exists", rs.Primary.ID)
		}

		if rule != nil {
			return fmt.Errorf("url category with id %s exists and wasn't destroyed", rs.Primary.ID)
		}
	}

	return nil
}
