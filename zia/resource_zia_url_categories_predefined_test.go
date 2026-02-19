package zia

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
)

func TestAccResourceURLCategoriesPredefinedBasic(t *testing.T) {
	resourceName := "zia_url_categories_predefined.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLCategoriesPredefinedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccURLCategoriesPredefinedConfig_initial(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesPredefinedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "FINANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttrSet(resourceName, "super_category"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "val"),
				),
			},
			// Update: add and remove list items
			{
				Config: testAccURLCategoriesPredefinedConfig_updated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLCategoriesPredefinedExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "FINANCE"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttrSet(resourceName, "super_category"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "val"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
				},
			},
		},
	})
}

func testAccCheckURLCategoriesPredefinedDestroy(s *terraform.State) error {
	// Predefined categories cannot be deleted — destroy is a no-op that only
	// removes the resource from state. The category will still exist in the API.
	return nil
}

func testAccCheckURLCategoriesPredefinedExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		resp, err := urlcategories.Get(context.Background(), service, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resourceName, err)
		}
		if resp.CustomCategory {
			return fmt.Errorf("expected predefined category but got custom category for ID: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccURLCategoriesPredefinedConfig_initial() string {
	return `
resource "zia_url_categories_predefined" "test" {
	name = "FINANCE"
	urls = [
		".testuniversity-acc.edu",
		".testcollege-acc.org",
	]
	keywords = [
		"tf-acc-finance",
		"tf-acc-university",
	]
	ip_ranges = [
		"203.0.113.0/24",
	]
}
`
}

func testAccURLCategoriesPredefinedConfig_updated() string {
	return `
resource "zia_url_categories_predefined" "test" {
	name = "FINANCE"
	urls = [
		".testcollege-acc.org",
		".testacademy-acc.net",
	]
	keywords = [
		"tf-acc-finance",
	]
	ip_ranges = [
		"203.0.113.0/24",
		"198.51.100.0/24",
	]
}
`
}
