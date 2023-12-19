package zia

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservicegroups"
)

func TestAccResourceFWNetworkServiceGroupsBasic(t *testing.T) {
	var services networkservicegroups.NetworkServiceGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringNetworkServiceGroups)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServiceGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, initialName, variable.FWNetworkServicesGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServiceGroupsExists(resourceTypeAndName, &services),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWNetworkServicesGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "services.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, updatedName, variable.FWNetworkServicesGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServiceGroupsExists(resourceTypeAndName, &services),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWNetworkServicesGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "services.#", "1"),
				),
			},
		},
	})
}

func testAccCheckFWNetworkServiceGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FWFilteringNetworkAppGroups {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.networkservicegroups.GetNetworkServiceGroups(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network services group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFWNetworkServiceGroupsExists(resource string, rule *networkservicegroups.NetworkServiceGroups) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.networkservicegroups.GetNetworkServiceGroups(id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

data "zia_firewall_filtering_network_service" "example1" {
	name = "ICMP_ANY"
  }

data "zia_firewall_filtering_network_service" "example2" {
	name = "TCP_ANY"
  }

resource "%s" "%s" {
    name = "tf-acc-test-%s"
    description = "%s"
    services {
        id = [
            data.zia_firewall_filtering_network_service.example1.id,
            data.zia_firewall_filtering_network_service.example2.id,
        ]
    }
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the network services group
		resourcetype.FWFilteringNetworkServiceGroups,
		resourceName,
		generatedName,
		description,

		// Data source type and name
		resourcetype.FWFilteringNetworkServiceGroups,
		resourceName,

		// Reference to the resource
		resourcetype.FWFilteringNetworkServiceGroups,
		resourceName,
	)
}
