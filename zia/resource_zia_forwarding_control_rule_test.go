package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/forwarding_rules"
)

func TestAccResourceForwardingControlRuleBasic(t *testing.T) {
	var rules forwarding_rules.ForwardingRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ForwardingControlRule)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	// Generate Source IP Group HCL Resource
	sourceIPGroupTypeAndName, _, sourceIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringSourceGroup)
	sourceIPGroupHCL := testAccCheckFWIPSourceGroupsConfigure(sourceIPGroupTypeAndName, sourceIPGroupGeneratedName, variable.FWSRCGroupDescription)

	// Generate Destination IP Group HCL Resource
	dstIPGroupTypeAndName, _, dstIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)
	dstIPGroupHCL := testAccCheckFWIPDestinationGroupsConfigure(dstIPGroupTypeAndName, dstIPGroupGeneratedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingControlRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingControlRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.FowardingControlDescription, variable.FowardingControlType, variable.FWRuleResourceState, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingControlRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FowardingControlDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FowardingControlType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FowardingControlState),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_ip_groups.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckForwardingControlRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.FowardingControlUpdateDescription, variable.FowardingControlType, variable.FowardingControlUpdateState, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingControlRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FowardingControlUpdateDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FowardingControlType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FowardingControlUpdateState),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_ip_groups.0.id.#", "1"),
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

func testAccCheckForwardingControlRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.forwarding_rules

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ForwardingControlRule {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := forwarding_rules.Get(service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("forwarding control rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckForwardingControlRuleExists(resource string, rule *forwarding_rules.ForwardingRules) resource.TestCheckFunc {
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
		service := apiClient.forwarding_rules

		var receivedRule *forwarding_rules.ForwardingRules

		// Integrate retry here
		retryErr := RetryOnError(func() error {
			var innerErr error
			receivedRule, innerErr = forwarding_rules.Get(service, id)
			if innerErr != nil {
				return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, innerErr)
			}
			return nil
		})

		if retryErr != nil {
			return retryErr
		}

		*rule = *receivedRule
		return nil
	}
}

func testAccCheckForwardingControlRuleConfigure(resourceTypeAndName, generatedName, name, description, ruleType, state string, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// source ip group resource
%s

// destination ip group resource
%s

// forwarding control rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		sourceIPGroupHCL,
		dstIPGroupHCL,
		getForwardingControlRuleResourceHCL(generatedName, name, description, ruleType, state, ruleLabelTypeAndName, sourceIPGroupTypeAndName, dstIPGroupTypeAndName),

		// data source variables
		resourcetype.ForwardingControlRule,
		generatedName,
		resourceTypeAndName,
	)
}

func getForwardingControlRuleResourceHCL(generatedName, name, description, ruleType, state string, ruleLabelTypeAndName, sourceIPGroupTypeAndName, dstIPGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
	name = "ZSCALER_PROXY_NW_SERVICES"
}

data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

data "zia_location_groups" "sdwan_usa" {
	name = "SDWAN_USA"
}

data "zia_department_management" "engineering" {
	name = "Engineering"
}

data "zia_department_management" "marketing" {
	name = "Marketing"
}

data "zia_group_management" "engineering" {
	name = "Engineering"
}

data "zia_group_management" "marketing" {
	name = "Marketing"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	type = "%s"
	state = "%s"
	order = 1
	rank = 7
    forward_method = "DIRECT"
	nw_services {
		id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
	}
	location_groups {
		id = [data.zia_location_groups.sdwan_can.id, data.zia_location_groups.sdwan_usa.id]
	}
	groups {
		id = [data.zia_group_management.engineering.id, data.zia_group_management.marketing.id]
	}
	departments {
		id = [data.zia_department_management.engineering.id, data.zia_department_management.marketing.id]
	}
	labels {
		id = ["${%s.id}"]
	}
	src_ip_groups {
		id = ["${%s.id}"]
	}
	dest_ip_groups {
		id = ["${%s.id}"]
	}
	depends_on = [ %s, %s, %s ]
}
		`,
		// resource variables
		resourcetype.ForwardingControlRule,
		generatedName,
		name,
		description,
		ruleType,
		state,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
	)
}
