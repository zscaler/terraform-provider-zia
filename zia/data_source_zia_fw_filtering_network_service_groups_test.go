package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWNetworkServiceGroups_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_network_service_groups.by_id"
	// resourceName2 := "data.zia_firewall_filtering_network_service_groups.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWNetworkServiceGroupsByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWNetworkServiceGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					// resource.TestCheckResourceAttr(resourceName2, "name", rName),
					// resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceFWNetworkServiceGroupsByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	data "zia_firewall_filtering_network_service" "http"{
		name = "HTTP"
	}
	data "zia_firewall_filtering_network_service" "ftp"{
		name = "FTP"
	}
	data "zia_firewall_filtering_network_service" "icmp"{
		name = "ICMP_ANY"
	}
	resource "zia_firewall_filtering_network_service_groups" "testAcc"{
		name = "%s"
		description = "%s"
		services {
			id = [
				data.zia_firewall_filtering_network_service.http.id,
				data.zia_firewall_filtering_network_service.ftp.id,
				data.zia_firewall_filtering_network_service.icmp.id
			]
		}
	}
	data "zia_firewall_filtering_network_service_groups" "by_name" {
		name = zia_firewall_filtering_network_service_groups.testAcc.name
		depends_on = [zia_firewall_filtering_network_service_groups.testAcc]
	}
	data "zia_firewall_filtering_network_service_groups" "by_id" {
		id = zia_firewall_filtering_network_service_groups.testAcc.id
		depends_on = [zia_firewall_filtering_network_service_groups.testAcc]
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
*/
