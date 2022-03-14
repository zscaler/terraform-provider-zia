package zia

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlcategories"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceURLCategoriesBasic(t *testing.T) {
	var categories urlcategories.URLCategory
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.URLCategories)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLCategoriesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLCategoriesConfigure(resourceTypeAndName, generatedName, variable.CategoryDescription, variable.CustomCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesExists(resourceTypeAndName, &categories),
					resource.TestCheckResourceAttr(resourceTypeAndName, "configured_name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CategoryDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_category", strconv.FormatBool(variable.CustomCategory)),
				),
			},

			// Update test
			{
				Config: testAccCheckURLCategoriesConfigure(resourceTypeAndName, generatedName, variable.CategoryDescription, variable.CustomCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesExists(resourceTypeAndName, &categories),
					resource.TestCheckResourceAttr(resourceTypeAndName, "configured_name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.CategoryDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_category", strconv.FormatBool(variable.CustomCategory)),
				),
			},
		},
	})
}

func testAccCheckURLCategoriesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.URLCategories {
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

func testAccCheckURLCategoriesConfigure(resourceTypeAndName, generatedName, description string, custom_category bool) string {
	return fmt.Sprintf(`
// rule label resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		URLCategoriesResourceHCL(generatedName, description, custom_category),

		// data source variables
		resourcetype.URLCategories,
		generatedName,
		resourceTypeAndName,
	)
}

func URLCategoriesResourceHCL(generatedName, description string, custom_category bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	super_category 		= "USER_DEFINED"
	configured_name 	= "%s"
	description 		= "%s"
	custom_category     = "%s"
	keywords            = ["microsoft"]
	db_categorized_urls = [".creditkarma.com", ".youku.com"]
	type                = "URL_CATEGORY"
}
`,
		// resource variables
		resourcetype.URLCategories,
		generatedName,
		generatedName,
		description,
		strconv.FormatBool(custom_category),
	)
}
