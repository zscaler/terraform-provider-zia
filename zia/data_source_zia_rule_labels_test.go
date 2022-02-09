package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceRuleLabels_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_rule_labels.by_id"
	resourceName2 := "data.zia_rule_labels.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRuleLabelsByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRuleLabels(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceRuleLabelsByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zia_rule_labels" "testAcc" {
		name = "%s"
		description = "%s"
	}
	data "zia_rule_labels" "by_name" {
		name = zia_rule_labels.testAcc.name
	}
	data "zia_rule_labels" "by_id" {
		id = zia_rule_labels.testAcc.id
	}
	`, rName, rDesc)
}

func testAccDataSourceRuleLabels(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
