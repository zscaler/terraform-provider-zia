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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
)

func TestAccResourceFileTypeControlRules_Basic(t *testing.T) {
	var rules filetypecontrol.FileTypeRules
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FileTypeControlRules)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, "tf-acc-test-"+ruleLabelGeneratedName, variable.RuleLabelDescription)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFileTypeControlRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFileTypeControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.FileTypeControlRuleDescription, variable.FileTypeControlRuleAction, variable.FileTypeControlRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileTypeControlRulesExists(resourceTypeAndName, &rules),
					testAccCheckFileTypeControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FileTypeControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "filtering_action", variable.FileTypeControlRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FileTypeControlRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "device_trust_levels.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "file_types.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "protocols.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "labels.0.id.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckFileTypeControlRulesConfigure(resourceTypeAndName, generatedName, generatedName, variable.FileTypeControlRuleDescription, variable.FileTypeControlRuleAction, variable.FileTypeControlRuleState, ruleLabelTypeAndName, ruleLabelHCL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFileTypeControlRulesExists(resourceTypeAndName, &rules),
					testAccCheckFileTypeControlRulesExists(resourceTypeAndName, &rules),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FileTypeControlRuleDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "filtering_action", variable.FileTypeControlRuleAction),
					resource.TestCheckResourceAttr(resourceTypeAndName, "state", variable.FileTypeControlRuleState),
					resource.TestCheckResourceAttr(resourceTypeAndName, "device_trust_levels.#", "4"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "file_types.#", "4"),
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

func testAccCheckFileTypeControlRulesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FileTypeControlRules {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := filetypecontrol.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("file type control rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFileTypeControlRulesExists(resource string, rule *filetypecontrol.FileTypeRules) resource.TestCheckFunc {
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

		receivedRule, err := filetypecontrol.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFileTypeControlRulesConfigure(resourceTypeAndName, generatedName, name, description, action, state string, ruleLabelTypeAndName, ruleLabelHCL string) string {
	return fmt.Sprintf(`
// rule label resource
%s


// file type control rule resource
%s

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		ruleLabelHCL,
		getFileTypeControlRulesResourceHCL(generatedName, name, description, action, state, ruleLabelTypeAndName),

		// data source variables
		resourcetype.FileTypeControlRules,
		generatedName,
		resourceTypeAndName,
	)
}

func getFileTypeControlRulesResourceHCL(generatedName, name, description, action, state string, ruleLabelTypeAndName string) string {
	return fmt.Sprintf(`

data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class   = ["WEB_MAIL"]
}

resource "%s" "%s" {
	name 			    = "tf-acc-test-%s"
	description         = "%s"
	filtering_action    = "%s"
	state 			    = "%s"
	order 			    = 1
	rank 			    = 7
	operation           = "DOWNLOAD"
    active_content      = true
    unscannable         = false
	device_trust_levels = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
	file_types          = ["FTCATEGORY_MS_WORD", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_PDF_DOCUMENT", "FTCATEGORY_MS_EXCEL"]
	protocols           = ["FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
	cloud_applications  = tolist([for app in data.zia_cloud_applications.this.applications : app["app"]])

	labels {
		id = ["${%s.id}"]
	}
	depends_on = [ %s ]
}
		`,
		// resource variables
		resourcetype.FileTypeControlRules,
		generatedName,
		name,
		description,
		action,
		state,
		ruleLabelTypeAndName,
		ruleLabelTypeAndName,
	)
}
