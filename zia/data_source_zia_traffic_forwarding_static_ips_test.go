package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceTrafficForwardingStaticIP_ByIdAndIP(t *testing.T) {
	rComment := acctest.RandString(15)
	resourceName := "data.zia_traffic_forwarding_static_ip.by_id"
	resourceName2 := "data.zia_traffic_forwarding_static_ip.by_ip"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTrafficForwardingStaticIPById(rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingStaticIP(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", rComment),
					resource.TestCheckResourceAttr(resourceName, "routable_ip", "true"),
				),
				PreventPostDestroyRefresh: true,
			},
			{

				Config: testAccDataSourceTrafficForwardingStaticIPByIP(rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingStaticIP(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "comment", rComment),
					resource.TestCheckResourceAttr(resourceName2, "routable_ip", "true"),
				),
			},
		},
	})
}

func testAccDataSourceTrafficForwardingStaticIPById(rComment string) string {
	return fmt.Sprintf(`
	resource "zia_traffic_forwarding_static_ip" "testAcc"{
		ip_address =  "118.189.211.221"
		routable_ip = true
		comment = "%s"
		geo_override = true
		latitude = -36.848461
		longitude = 174.763336
	}
data "zia_traffic_forwarding_static_ip" "by_id" {
	id = zia_traffic_forwarding_static_ip.testAcc.id
}
	`, rComment)
}

func testAccDataSourceTrafficForwardingStaticIPByIP(rComment string) string {
	return fmt.Sprintf(`
	resource "zia_traffic_forwarding_static_ip" "testAcc"{
		ip_address =  "118.189.211.221"
		routable_ip = true
		comment = "%s"
		geo_override = true
		latitude = -36.848461
		longitude = 174.763336
	}
data "zia_traffic_forwarding_static_ip" "by_ip" {
	ip_address = zia_traffic_forwarding_static_ip.testAcc.ip_address
}
	`, rComment)
}

func testAccDataSourceTrafficForwardingStaticIP(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
