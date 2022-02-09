package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlfilteringpolicies"
	"github.com/willguibr/terraform-provider-zia/zia/common/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/variable"
)

func TestAccResourceURLFilteringRulesBasic(t *testing.T) {
	var rule urlfilteringpolicies.URLFilteringRule
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.URLFilteringRules)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLFilteringRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, variable.URLFilteringRuleDescription, variable.URLFilteringRuleState, variable.URLFilteringRuleAction),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRulesExists(resourceTypeAndName, &rule),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.URLFilteringRuleResourceName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.URLFilteringRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.URLFilteringRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.URLFilteringRuleAction),
				),
			},

			// Update test
			{
				Config: testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, variable.URLFilteringRuleDescription, variable.URLFilteringRuleState, variable.URLFilteringRuleAction),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRulesExists(resourceTypeAndName, &rule),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.URLFilteringRuleResourceName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.URLFilteringRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.URLFilteringRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.URLFilteringRuleAction),
				),
			},
		},
	})
}

func testAccCheckURLFilteringRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.URLFilteringRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.urlfilteringpolicies.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckURLFilteringRulesExists(resource string, rule *urlfilteringpolicies.URLFilteringRule) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.urlfilteringpolicies.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, description, state, action string) string {
	return fmt.Sprintf(`
// rule resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		RuleResourceHCL(generatedName, description, state, action),

		// data source variables
		resourcetype.URLFilteringRules,
		generatedName,
		resourceTypeAndName,
	)
}

func RuleResourceHCL(generatedName, description, state, action string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  name                 		= "%s"
  description          		= "%s"
  state 			   		= "%s"
  action 					= "%s"
  order 					= 2
  url_categories 			= ["ANY"]
  protocols 				= ["ANY_RULE"]
  request_methods 			= [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}
`,
		// resource variables
		resourcetype.URLFilteringRules,
		generatedName,
		variable.URLFilteringRuleResourceName,
		description,
		state,
		action,
	)
}
