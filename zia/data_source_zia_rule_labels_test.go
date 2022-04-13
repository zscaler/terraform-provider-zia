package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceRuleLabels_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_rule_labels.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRuleLabelsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceRuleLabels(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-rule-label-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-rule-label-"+rDesc),
				),
			},
		},
	})
}

func testAccCheckRuleLabelsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_rule_labels" "test-rule-label" {
	name = "test-rule-label-%s"
	description = "test-rule-label-%s"
}

data "zia_rule_labels" "by_name" {
	name = zia_rule_labels.test-rule-label.name
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
