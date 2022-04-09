package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceLocationManagement_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")
	resourceName := "data.zia_location_management.test-sjc2022"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceLocationManagementBasic(rIP, rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceLocationManagement(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-sjc2022-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-sjc2022-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "country", "UNITED_STATES"),
					resource.TestCheckResourceAttr(resourceName, "tz", "UNITED_STATES_AMERICA_LOS_ANGELES"),
					resource.TestCheckResourceAttr(resourceName, "display_time_unit", "HOUR"),
					resource.TestCheckResourceAttr(resourceName, "auth_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "surrogate_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "xff_forward_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ofw_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ips_control", "true"),
				),
			},
		},
	})
}

func testAccCheckDataSourceLocationManagementBasic(rIP, rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test-sjc2022"{
	comment = "Test SJC2022 - Static IP"
	ip_address =  "%s"
	routable_ip = true
	geo_override = true
	latitude = -36.848461
	longitude = 174.763336
}

resource "zia_traffic_forwarding_vpn_credentials" "test-sjc2022"{
	comments    = "Test SJC2022 - VPN Credentials"
	type        = "IP"
	ip_address  =  zia_traffic_forwarding_static_ip.test-sjc2022.ip_address
	pre_shared_key = "newPassword123!"
	depends_on = [ zia_traffic_forwarding_static_ip.test-sjc2022 ]
}

resource "zia_location_management" "test-sjc2022"{
	name = "test-sjc2022-%s"
	description = "test-sjc2022-%s"
	country = "UNITED_STATES"
	tz = "UNITED_STATES_AMERICA_LOS_ANGELES"
	auth_required = true
	idle_time_in_minutes = 720
	display_time_unit = "HOUR"
	surrogate_ip = true
	xff_forward_enabled = true
	ofw_enabled = true
	ips_control = true
	ip_addresses = [ zia_traffic_forwarding_static_ip.test-sjc2022.ip_address ]
	depends_on = [ zia_traffic_forwarding_static_ip.test-sjc2022, zia_traffic_forwarding_vpn_credentials.test-sjc2022 ]
	vpn_credentials {
		id = zia_traffic_forwarding_vpn_credentials.test-sjc2022.vpn_credental_id
		type = zia_traffic_forwarding_vpn_credentials.test-sjc2022.type
		ip_address = zia_traffic_forwarding_static_ip.test-sjc2022.ip_address
	}
}

data "zia_location_management" "test-sjc2022" {
	name = zia_location_management.test-sjc2022.name
}
	`, rIP, rName, rDesc)
}

func testAccDataSourceLocationManagement(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
