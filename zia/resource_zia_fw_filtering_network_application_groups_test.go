package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkapplications"
)

func TestAccResourceFWNetworkApplicationGroups_basic(t *testing.T) {
	var apGroups networkapplications.NetworkApplicationGroups
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_firewall_filtering_network_application_groups.test-fw-nw-app-group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkApplicationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkApplicationGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkApplicationGroupsExists("zia_firewall_filtering_network_application_groups.test-fw-nw-app-group", &apGroups),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-nw-app-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-nw-app-group-"+rDesc),
					// resource.TestCheckResourceAttr(resourceName, "network_applications", "network_applications"),
				),
			},
		},
	})
}

func testAccCheckFWNetworkApplicationGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_network_application_groups" "test-fw-nw-app-group" {
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
	`, rName, rDesc)
}

func testAccCheckFWNetworkApplicationGroupsExists(resource string, rule *networkapplications.NetworkApplicationGroups) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.networkapplications.GetNetworkApplicationGroups(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkApplicationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_network_application_groups" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.networkapplications.GetNetworkApplicationGroups(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network application group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
