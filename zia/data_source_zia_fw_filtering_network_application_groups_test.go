package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWNetworkApplicationGroups_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_network_application_groups.by_id"
	resourceName2 := "data.zia_firewall_filtering_network_application_groups.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWNetworkApplicationGroupsByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWNetworkApplicationGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceFWNetworkApplicationGroupsByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zia_firewall_filtering_network_application_groups" "o365" {
		name = "%s"
		description = "%s"
		network_applications = [
			"YAMMER",
			"OFFICE365",
			"SKYPE_FOR_BUSINESS",
			"OUTLOOK",
			"SHAREPOINT",
			"SHAREPOINT_ADMIN",
			"SHAREPOINT_BLOG",
			"SHAREPOINT_CALENDAR",
			"SHAREPOINT_DOCUMENT",
			"SHAREPOINT_ONLINE",
			"ONEDRIVE",
		]
	}
	data "zia_firewall_filtering_network_application_groups" "by_name" {
		name = zia_firewall_filtering_network_application_groups.o365.name
	}
	data "zia_firewall_filtering_network_application_groups" "by_id" {
		id = zia_firewall_filtering_network_application_groups.o365.id
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
