package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func TestAccResourceURLFilteringRules_Basic(t *testing.T) {
	var rules urlfilteringpolicies.URLFilteringRule
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.URLFilteringRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLFilteringRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.URLFilteringRuleDescription, variable.URLFilteringRuleAction, variable.URLFilteringRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.URLFilteringRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.URLFilteringRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.URLFilteringRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "url_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "request_methods.#", "9"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.URLFilteringRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.URLFilteringRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "url_categories.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "request_methods.#", "9"),
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

func testAccCheckURLFilteringRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.URLFilteringRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := urlfilteringpolicies.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("url filtering rule with id %d exists and wasn't destroyed", id)
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
		service := apiClient.Service

		receivedRule, err := urlfilteringpolicies.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckURLFilteringRulesConfigure(resourceTypeAndName, generatedName, name, description, action, state, ruleLabelTypeAndName, ruleLabelHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// url filtering rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		getURLFilteringRuleResourceHCL(generatedName, name, description, action, state, ruleLabelTypeAndName),

		// data source variables
		resourcetype.URLFilteringRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getURLFilteringRuleResourceHCL(generatedName, name, description, action, state, ruleLabelTypeAndName string) string {
	// Generate current time + 5 minutes for validity_start_time
	validityStartTime := time.Now().Add(5 * time.Minute).UTC().Format(time.RFC1123)

	// Generate time 365 days + 5 minutes from now for validity_end_time
	validityEndTime := time.Now().AddDate(1, 0, 0).Add(5 * time.Minute).UTC().Format(time.RFC1123)

	return fmt.Sprintf(`


resource "%s" "%s" {
    name = "tf-acc-test-%s"
    description = "%s"
	action = "%s"
    state = "%s"
    order = 1
	url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
	device_trust_levels = [	"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST" ]
	user_agent_types = [	"OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE", "OTHER" ]
	user_risk_score_levels = ["LOW", "MEDIUM", "HIGH", "CRITICAL"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]

	labels {
		id = ["${%s.id}"]
	}
	enforce_time_validity = true
	validity_time_zone_id = "US/Pacific"
    validity_start_time = "%s"
    validity_end_time = "%s"
}
`,
		// resource variables
		resourcetype.URLFilteringRules,
		generatedName,
		name,
		description,
		action,
		state,
		ruleLabelTypeAndName,
		// validity times
		validityStartTime,
		validityEndTime,
	)
}
