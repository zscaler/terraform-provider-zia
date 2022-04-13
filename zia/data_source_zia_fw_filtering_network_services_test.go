package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWNetworkServices_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_network_service.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWNetworkServicesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWNetworkServices(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-nw-svc-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-nw-svc-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "type", "CUSTOM"),
				),
			},
		},
	})
}

func testAccDataSourceFWNetworkServicesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_network_service" "test" {
	name        = "test-fw-nw-svc-%s"
	description = "test-fw-nw-svc-%s"
	type = "CUSTOM"
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

data "zia_firewall_filtering_network_service" "test" {
	name = zia_firewall_filtering_network_service.test.name
}
	`, rName, rDesc)
}

func testAccDataSourceFWNetworkServices(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
