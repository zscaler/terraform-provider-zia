package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
)

func TestAccResourceFWNetworkServiceGroupsBasic(t *testing.T) {
	var groups networkservices.NetworkServiceGroups
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_firewall_filtering_network_service_groups.test-fw-nw-svc-group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServiceGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceFWNetworkServiceGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServiceGroupsExists("zia_firewall_filtering_network_service_groups.test-fw-nw-svc-group", &groups),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-nw-svc-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-nw-svc-group-"+rDesc),
					// resource.TestCheckResourceAttr(resourceName, "services", "services"),
				),
			},
		},
	})
}

func testAccCheckResourceFWNetworkServiceGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

data "zia_firewall_filtering_network_service" "icmp_any" {
	name = "ICMP_ANY"
}

data "zia_firewall_filtering_network_service" "dns" {
	name = "DNS"
}

resource "zia_firewall_filtering_network_service_groups" "test-fw-nw-svc-group"{
	name = "test-fw-nw-svc-group-%s"
	description = "test-fw-nw-svc-group-%s"
	services {
		id = [
			data.zia_firewall_filtering_network_service.icmp_any.id,
			data.zia_firewall_filtering_network_service.dns.id,
		]
	}
}
	`, rName, rDesc)
}

func testAccCheckFWNetworkServiceGroupsExists(resource string, rule *networkservices.NetworkServiceGroups) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.networkservices.GetNetworkServiceGroups(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkServiceGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_network_service_groups" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.networkservices.GetNetworkServiceGroups(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network services group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
