package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
)

func TestAccResourceFWIPDestinationGroups_basic(t *testing.T) {
	var groups ipdestinationgroups.IPDestinationGroups
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_firewall_filtering_destination_groups.test-fw-dst-group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFWIPDestinationGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists("zia_firewall_filtering_destination_groups.test-fw-dst-group", &groups),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-dst-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-dst-group-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "type", "DSTN_FQDN"),
				),
			},
		},
	})
}

func testAccResourceFWIPDestinationGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_destination_groups" "test-fw-dst-group" {
	name        = "test-fw-dst-group-%s"
	description = "test-fw-dst-group-%s"
	type        = "DSTN_FQDN"
	addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
  }
	`, rName, rDesc)
}

func testAccCheckFWIPDestinationGroupsExists(resource string, rule *ipdestinationgroups.IPDestinationGroups) resource.TestCheckFunc {
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
		receivedGroup, err := apiClient.ipdestinationgroups.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedGroup

		return nil
	}
}

func testAccCheckFWIPDestinationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_destination_groups" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.ipdestinationgroups.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("ip destination group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
