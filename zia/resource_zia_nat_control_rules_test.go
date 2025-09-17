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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/nat_control_policies"
)

func TestAccResourceNatControlRules_Basic(t *testing.T) {
	var rules nat_control_policies.NatControlPolicies
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.NatControlRules)

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
		CheckDestroy: testAccCheckNatControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNatControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.NATControlRuleDescription, variable.NATControlRuleState, variable.NATControlRuleLogging, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatControlRulesExists(resourceTypeAndName, &rules),
					testAccCheckNatControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.NATControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.NATControlRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "redirect_port", strconv.Itoa(variable.NATControlRedirectPort)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "redirect_ip", variable.NATControlRedirectIP),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_countries.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_addresses.#", "5"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ips.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckNatControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.NATControlRuleDescription, variable.NATControlRuleState, variable.NATControlRuleLogging, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatControlRulesExists(resourceTypeAndName, &rules),
					testAccCheckNatControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.NATControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "redirect_port", strconv.Itoa(variable.NATControlRedirectPort)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "redirect_ip", variable.NATControlRedirectIP),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_countries.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_addresses.#", "5"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_ips.#", "2"),
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

func testAccCheckNatControlRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.NatControlRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := nat_control_policies.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("nat control rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckNatControlRulesExists(resource string, rule *nat_control_policies.NatControlPolicies) resource.TestCheckFunc {
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

		receivedRule, err := nat_control_policies.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckNatControlRulesConfigure(resourceTypeAndName, generatedName, name, description, state string, enableLogging bool, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// source ip group resource
%s

// destination ip group resource
%s

// firewall ips rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		sourceIPGroupHCL,
		dstIPGroupHCL,
		getNatControlRulesResourceHCL(generatedName, name, description, state, enableLogging, ruleLabelTypeAndName, sourceIPGroupTypeAndName, dstIPGroupTypeAndName),

		// data source variables
		resourcetype.NatControlRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getNatControlRulesResourceHCL(generatedName, name, description, state string, enableLogging bool, ruleLabelTypeAndName, sourceIPGroupTypeAndName, dstIPGroupTypeAndName string) string {
	return fmt.Sprintf(`

data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
	name = "ZSCALER_PROXY_NW_SERVICES"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	state = "%s"
	order = 3
    redirect_port="5000"
    redirect_ip="192.168.100.150"
    src_ips=["192.168.100.0/24", "192.168.200.1"]
    dest_addresses=["3.217.228.0-3.217.231.255", "3.235.112.0-3.235.119.255", "35.80.88.0-35.80.95.255", "server1.acme.com", "*.acme.com"]
    dest_countries=["BR", "CA"]
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
		resourcetype.NatControlRules,
		generatedName,
		name,
		description,
		state,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
		ruleLabelTypeAndName,
		sourceIPGroupTypeAndName,
		dstIPGroupTypeAndName,
	)
}
