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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
)

func TestAccResourceSandboxRules_Basic(t *testing.T) {
	var rules sandbox_rules.SandboxRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.SandboxRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSandboxRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSandboxRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.SandboxRuleDescription, variable.SandboxState, variable.SandboxAction, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxRulesExists(resourceTypeAndName, &rules),
					testAccCheckSandboxRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SandboxRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.SandboxState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ba_rule_action", variable.SandboxAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ba_policy_categories.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "file_types.#", "19"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckSandboxRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.SandboxRuleDescription, variable.SandboxState, variable.SandboxAction, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxRulesExists(resourceTypeAndName, &rules),
					testAccCheckSandboxRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.SandboxRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.SandboxState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ba_rule_action", variable.SandboxAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ba_policy_categories.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "file_types.#", "19"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
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

func testAccCheckSandboxRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.SandboxRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := sandbox_rules.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("sandbox rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckSandboxRulesExists(resource string, rule *sandbox_rules.SandboxRules) resource.TestCheckFunc {
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

		receivedRule, err := sandbox_rules.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckSandboxRulesConfigure(resourceTypeAndName, generatedName, name, description, state, action string, ruleLabelTypeAndName, ruleLabelHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s

// sandbox rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		getSandboxRulesResourceHCL(generatedName, name, description, state, action, ruleLabelTypeAndName),

		// data source variables
		resourcetype.SandboxRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getSandboxRulesResourceHCL(generatedName, name, description, state, action string, ruleLabelTypeAndName string) string {
	return fmt.Sprintf(`


data "zia_location_groups" "sdwan_can" {
	name = "SDWAN_CAN"
}

resource "%s" "%s" {
	name = "tf-acc-test-%s"
	description = "%s"
	state = "%s"
	order = 1
	rank  = 7
	ba_rule_action = "%s"
	first_time_enable    = true
	first_time_operation = "ALLOW_SCAN"
    ml_action_enabled    = true
	ba_policy_categories = ["ADWARE_BLOCK", "BOTMAL_BLOCK", "ANONYP2P_BLOCK", "RANSOMWARE_BLOCK"]
    file_types           = [
        "FTCATEGORY_P7Z",
        "FTCATEGORY_MS_WORD",
        "FTCATEGORY_PDF_DOCUMENT",
        "FTCATEGORY_TAR",
        "FTCATEGORY_SCZIP",
        "FTCATEGORY_WINDOWS_EXECUTABLES",
        "FTCATEGORY_HTA",
        "FTCATEGORY_FLASH",
        "FTCATEGORY_RAR",
        "FTCATEGORY_MS_EXCEL",
        "FTCATEGORY_VISUAL_BASIC_SCRIPT",
        "FTCATEGORY_MS_POWERPOINT",
        "FTCATEGORY_WINDOWS_LIBRARY",
        "FTCATEGORY_POWERSHELL",
        "FTCATEGORY_APK",
        "FTCATEGORY_ZIP",
        "FTCATEGORY_BZIP2",
        "FTCATEGORY_JAVA_APPLET",
        "FTCATEGORY_MS_RTF"
    ]
    protocols            = [
        "FOHTTP_RULE",
        "FTP_RULE",
        "HTTPS_RULE",
        "HTTP_RULE"
    ]
	labels {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
		`,
		// resource variables
		resourcetype.SandboxRules,
		generatedName,
		name,
		description,
		state,
		action,
		ruleLabelTypeAndName,
		ruleLabelTypeAndName,
	)
}
