package zia

/*
import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlfilteringpolicies"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
)

func TestAccResourceURLFilteringRulesBasic(t *testing.T) {
	var rules urlfilteringpolicies.URLFilteringRule
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_url_filtering_rules.test-url-rule"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckURLFilteringRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLFilteringRuleBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckURLFilteringRuleExists("zia_url_filtering_rules.test-url-rule", &rules),
					resource.TestCheckResourceAttr(resourceName, "name", "tfurl-rule-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "tfurl-rule-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
				),
			},
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckURLFilteringRuleBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zia_url_filtering_rules" "test-url-rule" {
	name = "tfurl-rule-%s"
	description = "tfurl-rule-%s"
	state = "ENABLED"
	action = "ALLOW"
	order = 1
	rank = 7
	url_categories = ["ANY"]
	protocols = ["HTTPS_RULE", "HTTP_RULE"]
	request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}

`, rName, rDesc)
}

func testAccCheckURLFilteringRuleExists(resource string, rule *urlfilteringpolicies.URLFilteringRule) resource.TestCheckFunc {
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

func testAccCheckURLFilteringRuleDestroy(s *terraform.State) error {
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
			return fmt.Errorf("url filtering rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
*/
