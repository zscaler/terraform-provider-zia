package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceURLFilteringRules_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_url_filtering_rules.by_name"
	// resourceID := "data.zia_url_filtering_rules.by_id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckURLFilteringRulesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceURLFilteringRules(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-url-rule-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-url-rule-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "action", "ALLOW"),
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
					// resource.TestCheckResourceAttr(resourceID, "name", "test-url-rule-"+rName),
					// resource.TestCheckResourceAttr(resourceID, "description", "test-url-rule-"+rDesc),
					// resource.TestCheckResourceAttr(resourceID, "action", "ALLOW"),
					// resource.TestCheckResourceAttr(resourceID, "state", "ENABLED"),
				),
				// PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccCheckURLFilteringRulesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "zia_url_filtering_rules" "test-url-rule" {
	name = "test-url-rule-%s"
	description = "test-url-rule-%s"
	state = "ENABLED"
	action = "ALLOW"
	order = 1
	rank = 7
	url_categories = ["ANY"]
	protocols = ["HTTPS_RULE", "HTTP_RULE"]
	request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}

data "zia_url_filtering_rules" "by_name" {
	name = zia_url_filtering_rules.test-url-rule.name
}
	`, rName, rDesc)
}

func testAccDataSourceURLFilteringRules(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
*/
