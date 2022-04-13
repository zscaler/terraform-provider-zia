package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/rule_labels"
)

func TestAccResourceRuleLabelsBasic(t *testing.T) {
	var labels rule_labels.RuleLabels
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_rule_labels.test-rule-label"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuleLabelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRuleLabelsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleLabelsExists(resourceName, &labels),
					resource.TestCheckResourceAttr(resourceName, "name", "test-rule-label-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-rule-label-"+rDesc),
				),
			},
		},
	})
}

func testAccResourceRuleLabelsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_rule_labels" "test-rule-label" {
	name = "test-rule-label-%s"
	description = "test-rule-label-%s"
}
	`, rName, rDesc)
}

func testAccCheckRuleLabelsExists(resource string, label *rule_labels.RuleLabels) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedLabel, err := apiClient.rule_labels.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*label = *receivedLabel

		return nil
	}
}

func testAccCheckRuleLabelsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_rule_labels" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.rule_labels.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
