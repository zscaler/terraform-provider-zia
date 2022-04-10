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

func TestAccResourceFWNetworkServicesBasic(t *testing.T) {
	var services networkservices.NetworkServices
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_firewall_filtering_network_service.test-fw-nw-svc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServicesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkServicesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServicesExists("zia_firewall_filtering_network_service.test-fw-nw-svc", &services),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-nw-svc-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-nw-svc-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "src_tcp_ports.0.start", "5000"),
					resource.TestCheckResourceAttr(resourceName, "src_tcp_ports.1.start", "5001"),
					resource.TestCheckResourceAttr(resourceName, "dest_tcp_ports.0.start", "5000"),
					resource.TestCheckResourceAttr(resourceName, "dest_tcp_ports.1.start", "5001"),
				),
			},
		},
	})
}

func testAccCheckFWNetworkServicesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_network_service" "test-fw-nw-svc" {
	name        = "test-fw-nw-svc-%s"
	description = "test-fw-nw-svc-%s"
	type = "CUSTOM"
	src_tcp_ports {
		start = 5000
	}
	src_tcp_ports {
		start = 5001
	}
    dest_tcp_ports {
        start = 5000
    }
    dest_tcp_ports {
        start = 5001
    }
}
	`, rName, rDesc)
}

func testAccCheckFWNetworkServicesExists(resource string, rule *networkservices.NetworkServices) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.networkservices.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkServicesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_network_service" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.networkservices.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network services with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
