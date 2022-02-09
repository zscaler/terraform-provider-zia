package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/rule_labels"
	"github.com/willguibr/terraform-provider-zia/zia/common/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/variable"
)

func TestAccResourceRuleLabelsBasic(t *testing.T) {
	var labels rule_labels.RuleLabels
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuleLabelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRuleLabelsConfigure(resourceTypeAndName, generatedName, variable.RuleLabelDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleLabelsExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.RuleLabelName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.RuleLabelDescription),
				),
			},

			// Update test
			{
				Config: testAccCheckRuleLabelsConfigure(resourceTypeAndName, generatedName, variable.RuleLabelDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleLabelsExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.RuleLabelName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.RuleLabelDescription),
				),
			},
		},
	})
}

func testAccCheckRuleLabelsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.RuleLabels {
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

func testAccCheckRuleLabelsExists(resource string, rule *rule_labels.RuleLabels) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.rule_labels.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckRuleLabelsConfigure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`
// rule label resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		RuleLabelsResourceHCL(generatedName, description),

		// data source variables
		resourcetype.RuleLabels,
		generatedName,
		resourceTypeAndName,
	)
}

func RuleLabelsResourceHCL(generatedName, description string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
    description = "%s"
}
`,
		// resource variables
		resourcetype.RuleLabels,
		generatedName,
		variable.RuleLabelName,
		description,
	)
}
