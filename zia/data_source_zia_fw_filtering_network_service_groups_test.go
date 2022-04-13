package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWNetworkServiceGroups_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_network_service_groups.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceFWNetworkServiceGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWNetworkServiceGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-nw-svc-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-nw-svc-group-"+rDesc),
					// resource.TestCheckResourceAttr(resourceName, "services", "services"),
				),
			},
		},
	})
}

func testAccCheckDataSourceFWNetworkServiceGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

data "zia_firewall_filtering_network_service" "icmp_any" {
	name = "ICMP_ANY"
}

data "zia_firewall_filtering_network_service" "dns" {
	name = "DNS"
}

resource "zia_firewall_filtering_network_service_groups" "test"{
	name = "test-fw-nw-svc-group-%s"
	description = "test-fw-nw-svc-group-%s"
	services {
		id = [
			data.zia_firewall_filtering_network_service.icmp_any.id,
			data.zia_firewall_filtering_network_service.dns.id,
		]
	}
}

data "zia_firewall_filtering_network_service_groups" "test" {
	name = zia_firewall_filtering_network_service_groups.test.name
}
	`, rName, rDesc)
}

func testAccDataSourceFWNetworkServiceGroups(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
