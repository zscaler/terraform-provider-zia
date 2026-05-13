package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ips_control_policies/ips_signature_rules"
)

func TestAccResourceIPSSignatureRulesBasic(t *testing.T) {
	var signature ips_signature_rules.IPSSignatureRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.IPSSignatureRules)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSSignatureRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPSSignatureRulesConfigure(
					resourceTypeAndName,
					initialName,
					variable.IPSSignatureRuleDescription,
					variable.IPSSignatureRuleText,
					variable.IPSSignatureRuleCategoryID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSSignatureRulesExists(resourceTypeAndName, &signature),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.IPSSignatureRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "rule_text", variable.IPSSignatureRuleText),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "category.0.id", strconv.Itoa(variable.IPSSignatureRuleCategoryID)),
				),
			},

			// Update test (name + rule_text)
			{
				Config: testAccCheckIPSSignatureRulesConfigure(
					resourceTypeAndName,
					updatedName,
					variable.IPSSignatureRuleDescription,
					variable.IPSSignatureRuleTextUpdated,
					variable.IPSSignatureRuleCategoryID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPSSignatureRulesExists(resourceTypeAndName, &signature),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "rule_text", variable.IPSSignatureRuleTextUpdated),
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

func testAccCheckIPSSignatureRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.IPSSignatureRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := ips_signature_rules.Get(context.Background(), service, id)
		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}
		if rule != nil {
			return fmt.Errorf("IPS signature rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckIPSSignatureRulesExists(resource string, rule *ips_signature_rules.IPSSignatureRules) resource.TestCheckFunc {
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
		service := apiClient.Service

		receivedRule, err := ips_signature_rules.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckIPSSignatureRulesConfigure(resourceTypeAndName, generatedName, description, ruleText string, categoryID int) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	return fmt.Sprintf(`

resource "%s" "%s" {
  name        = "%s"
  description = "%s"
  rule_text   = %s
  enabled     = true
  category {
    id = %d
  }
}

data "%s" "%s" {
  id = "${%s.%s.id}"
}
`,
		// resource variables
		resourcetype.IPSSignatureRules,
		resourceName,
		generatedName,
		description,
		strconv.Quote(ruleText),
		categoryID,

		// data source variables
		resourcetype.IPSSignatureRules,
		resourceName,
		resourcetype.IPSSignatureRules, resourceName,
	)
}
