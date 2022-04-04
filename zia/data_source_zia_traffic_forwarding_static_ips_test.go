package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceTrafficForwardingStaticIP_Basic(t *testing.T) {
	rComment := acctest.RandString(5)
	resourceName := "data.zia_traffic_forwarding_static_ip.test-static-ip"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingStaticIPBasic(rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingStatic(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "test-static-ip-"+rComment),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "121.234.54.80"),
					resource.TestCheckResourceAttr(resourceName, "routable_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "geo_override", "true"),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingStaticIPBasic(rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test-static-ip"{
	ip_address =  "121.234.54.80"
	routable_ip = true
	geo_override = true
	latitude = -36.848461
	longitude = 174.763336
	comment = "test-static-ip-%s"
}

data "zia_traffic_forwarding_static_ip" "test-static-ip" {
	ip_address = zia_traffic_forwarding_static_ip.test-static-ip.ip_address
}
	`, rComment)
}

func testAccDataSourceTrafficForwardingStatic(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
