package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWNetworkApplicationGroups_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_network_application_groups.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWNetworkApplicationGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWNetworkApplicationGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-nw-app-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-nw-app-group-"+rDesc),
					// resource.TestCheckResourceAttr(resourceName, "network_applications", "network_applications"),
				),
			},
		},
	})
}

func testAccDataSourceFWNetworkApplicationGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_network_application_groups" "test" {
	name        = "test-fw-nw-app-group-%s"
	description = "test-fw-nw-app-group-%s"
	network_applications  = [ "YAMMER",
                            "OFFICE365",
                            "SKYPE_FOR_BUSINESS",
                            "OUTLOOK",
                            "SHAREPOINT",
                            "SHAREPOINT_ADMIN",
                            "SHAREPOINT_BLOG",
                            "SHAREPOINT_CALENDAR",
                            "SHAREPOINT_DOCUMENT",
                            "SHAREPOINT_ONLINE",
                            "ONEDRIVE"
            ]
}
data "zia_firewall_filtering_network_application_groups" "test" {
	name = zia_firewall_filtering_network_application_groups.test.name
}
	`, rName, rDesc)
}

func testAccDataSourceFWNetworkApplicationGroups(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
