package zia

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
)

func TestAccResourceURLCategoriesBasic(t *testing.T) {
	var categories urlcategories.URLCategory
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.URLCategories)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLCategoriesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLCategoriesConfigure(resourceTypeAndName, initialName, variable.CustomCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesExists(resourceTypeAndName, &categories),
					resource.TestCheckResourceAttr(resourceTypeAndName, "configured_name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_category", strconv.FormatBool(variable.CustomCategory)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", "URL_CATEGORY"),
				),
			},

			// Update test
			{
				Config: testAccCheckURLCategoriesConfigure(resourceTypeAndName, updatedName, variable.CustomCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesExists(resourceTypeAndName, &categories),
					resource.TestCheckResourceAttr(resourceTypeAndName, "configured_name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_category", strconv.FormatBool(variable.CustomCategory)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", "URL_CATEGORY"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckURLCategoriesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.URLCategories {
			continue
		}

		rule, err := urlcategories.Get(context.Background(), service, rs.Primary.ID)

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
		service := apiClient.Service

		receivedRule, err := urlcategories.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckURLCategoriesConfigure(resourceTypeAndName, generatedName string, custom_category bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

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

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the url category
		resourcetype.URLCategories,
		resourceName,
		generatedName,
		generatedName,
		strconv.FormatBool(custom_category),

		// Data source type and name
		resourcetype.URLCategories,
		resourceName,

		// Reference to the resource
		resourcetype.URLCategories,
		resourceName,
	)
}
