package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
)

func TestAccResourceFWNetworkServicesBasic(t *testing.T) {
	var services networkservices.NetworkServices
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringNetworkServices)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServicesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkServicesConfigure(resourceTypeAndName, initialName, variable.FWNetworkServicesDescription, variable.FWNetworkServicesType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServicesExists(resourceTypeAndName, &services),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWNetworkServicesDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWNetworkServicesType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_tcp_ports.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_tcp_ports.#", "3"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWNetworkServicesConfigure(resourceTypeAndName, updatedName, variable.FWNetworkServicesDescription, variable.FWNetworkServicesType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkServicesExists(resourceTypeAndName, &services),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWNetworkServicesDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWNetworkServicesType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "src_tcp_ports.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dest_tcp_ports.#", "3"),
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

func testAccCheckFWNetworkServicesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FWFilteringNetworkServices {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := networkservices.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network services group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
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
		service := apiClient.Service

		receivedRule, err := networkservices.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkServicesConfigure(resourceTypeAndName, generatedName, description, svc_type string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
	type 		= "%s"
	src_tcp_ports {
	  start = 5000
	}
	src_tcp_ports {
	  start = 5001
	}
	src_tcp_ports {
	  start = 5002
	  end = 5005
	}
	dest_tcp_ports {
	  start = 5000
	}
	dest_tcp_ports {
	  start = 5001
	}
	dest_tcp_ports {
	  start = 5003
	  end = 5005
	}
  }

  data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the network services
		resourcetype.FWFilteringNetworkServices,
		resourceName,
		generatedName,
		description,
		svc_type,

		// Data source type and name
		resourcetype.FWFilteringNetworkServices,
		resourceName,

		// Reference to the resource
		resourcetype.FWFilteringNetworkServices,
		resourceName,
	)
}
