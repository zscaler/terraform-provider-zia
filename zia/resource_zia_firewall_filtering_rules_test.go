package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
)

func TestAccResourceFirewallFilteringRule_Basic(t *testing.T) {
	var rules filteringrules.FirewallFilteringRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FirewallFilteringRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	// Generate Source IP Group HCL Resource
	sourceIPGroupTypeAndName, _, sourceIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringSourceGroup)
	sourceIPGroupHCL := testAccCheckFWIPSourceGroupsConfigure(sourceIPGroupTypeAndName, "tf-acc-test-"+sourceIPGroupGeneratedName, variable.FWSRCGroupDescription)

	// Generate Destination IP Group HCL Resource
	dstIPGroupTypeAndName, _, dstIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)
	dstIPGroupHCL := testAccCheckFWIPDestinationGroupsConfigure(dstIPGroupTypeAndName, "tf-acc-test-"+dstIPGroupGeneratedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFirewallFilteringRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceState, variable.FWRuleEnableLogging, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallFilteringRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "nw_services.#", "1"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
				),
			},
			// Update test
			{
				Config: testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.FWRuleResourceDescription, variable.FWRuleResourceAction, variable.FWRuleResourceStateUpdate, variable.FWRuleEnableLogging, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallFilteringRuleExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWRuleResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "action", variable.FWRuleResourceAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FWRuleResourceStateUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "nw_services.#", "1"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "departments.0.id.#", "2"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "groups.0.id.#", "2"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "time_windows.0.id.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_ip_groups.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "enable_full_logging", strconv.FormatBool(variable.FWRuleEnableLogging)),
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

func testAccCheckFirewallFilteringRuleDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FirewallFilteringRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := filteringrules.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("firewall filtering rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFirewallFilteringRuleExists(resource string, rule *filteringrules.FirewallFilteringRules) resource.TestCheckFunc {
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

		var receivedRule *filteringrules.FirewallFilteringRules

		// Integrate retry here
		retryErr := RetryOnError(func() error {
			var innerErr error
			receivedRule, innerErr = filteringrules.Get(context.Background(), service, id)
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

func testAccCheckFirewallFilteringRuleConfigure(resourceTypeAndName, generatedName, name, description, action, state string, enableLogging bool, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// source ip group resource
%s

// destination ip group resource
%s

// firewall filtering rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		sourceIPGroupHCL,
		dstIPGroupHCL,
		getFirewallFilteringRuleResourceHCL(generatedName, name, description, action, state, enableLogging, ruleLabelTypeAndName, sourceIPGroupTypeAndName, dstIPGroupTypeAndName),

		// data source variables
		resourcetype.FirewallFilteringRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getFirewallFilteringRuleResourceHCL(generatedName, name, description, action, state string, enableLogging bool, ruleLabelTypeAndName, sourceIPGroupTypeAndName, dstIPGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
	name = "ZSCALER_PROXY_NW_SERVICES"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	action = "%s"
	state = "%s"
	order = 1
	enable_full_logging = "%s"
	device_trust_levels = [	"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST" ]
	nw_services {
		id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
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
		resourcetype.FirewallFilteringRules,
		generatedName,
		name,
		description,
		action,
		state,
		strconv.FormatBool(enableLogging),
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
	)
}
