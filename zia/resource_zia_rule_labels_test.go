package zia

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/rule_labels"
)

func TestAccResourceRuleLabelsBasic(t *testing.T) {
	var labels rule_labels.RuleLabels
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRuleLabelsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRuleLabelsConfigure(resourceTypeAndName, initialName, variable.RuleLabelDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleLabelsExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.RuleLabelDescription),
				),
			},

			// Update test
			{
				Config: testAccCheckRuleLabelsConfigure(resourceTypeAndName, updatedName, variable.RuleLabelDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleLabelsExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.RuleLabelDescription),
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
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name = "%s"
    description = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.RuleLabels,
		resourceName,
		generatedName,
		description,

		// data source variables
		resourcetype.RuleLabels,
		resourceName,
		// Reference to the resource
		resourcetype.RuleLabels, resourceName,
	)
}
