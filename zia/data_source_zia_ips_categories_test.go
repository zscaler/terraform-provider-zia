package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceIPSCategories_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceIPSCategoriesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// Looked up by name.
					resource.TestCheckResourceAttr("data.zia_ips_categories.by_name", "name", "ADSPYWARE"),
					resource.TestCheckResourceAttrSet("data.zia_ips_categories.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.zia_ips_categories.by_name", "back_end_name"),
					resource.TestCheckResourceAttr("data.zia_ips_categories.by_name", "predefined", "true"),

					// Looked up by the id resolved above, which must return the
					// same category.
					resource.TestCheckResourceAttr("data.zia_ips_categories.by_id", "name", "ADSPYWARE"),
					resource.TestCheckResourceAttrPair(
						"data.zia_ips_categories.by_id", "id",
						"data.zia_ips_categories.by_name", "id",
					),

					// A blank query returns the whole list.
					resource.TestCheckResourceAttrSet("data.zia_ips_categories.all", "categories.0.id"),
					resource.TestCheckResourceAttrSet("data.zia_ips_categories.all", "categories.0.name"),
				),
			},
		},
	})
}

var testAccCheckDataSourceIPSCategoriesConfig_basic = `
data "zia_ips_categories" "by_name" {
  name = "ADSPYWARE"
}

data "zia_ips_categories" "by_id" {
  id = data.zia_ips_categories.by_name.id
}

data "zia_ips_categories" "all" {}
`
