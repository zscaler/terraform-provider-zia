package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceURLCategories_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceURLCategoriesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_url_categories.custom05", "id"),
					resource.TestCheckResourceAttrSet(
						"data.zia_url_categories.custom06", "id"),
					resource.TestCheckResourceAttrSet(
						"data.zia_url_categories.custom07", "id"),
				),
			},
		},
	})
}

const testAccCheckDataSourceURLCategoriesConfig_basic = `
data "zia_url_categories" "custom05"{
    id = "CUSTOM_05"
}

data "zia_url_categories" "custom06"{
    id = "CUSTOM_06"
}

data "zia_url_categories" "custom07"{
    id = "CUSTOM_07"
}
`
