package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceTrafficForwardingGreTunnels_Basic(t *testing.T) {
	rComment := acctest.RandString(5)
	resourceName1 := "data.zia_traffic_forwarding_gre_tunnel_info.get_gre_unnumbered"
	resourceName2 := "data.zia_traffic_forwarding_gre_tunnel_info.get_gre_unnumbered"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTrafficForwardingGreTunnelsUnnumberedBasic(rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingGreTunnels(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "comment", "test-gre-tunnel-unnumbered-"+rComment),
					resource.TestCheckResourceAttr(resourceName1, "source_ip", "121.234.54.85"),
					resource.TestCheckResourceAttr(resourceName1, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName1, "ip_unnumbered", "true"),
					resource.TestCheckResourceAttr(resourceName1, "within_country", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccDataSourceTrafficForwardingGreTunnelsNumberedBasic(rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingGreTunnels(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "comment", "test-gre-tunnel-numbered-"+rComment),
					resource.TestCheckResourceAttr(resourceName2, "source_ip", "121.234.54.86"),
					resource.TestCheckResourceAttr(resourceName2, "country_code", "US"),
					resource.TestCheckResourceAttr(resourceName2, "ip_unnumbered", "false"),
					resource.TestCheckResourceAttr(resourceName2, "within_country", "true"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataSourceTrafficForwardingGreTunnelsUnnumberedBasic(rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test-gre-unnumbered-static-ip"{
	ip_address =  "121.234.54.85"
	routable_ip = true
	geo_override = false
}

data "zia_traffic_forwarding_gre_vip_recommended_list" "test-gre-unnumbered-vip"{
	source_ip = zia_traffic_forwarding_static_ip.test-gre-unnumbered-static-ip.ip_address
	required_count = 2
	depends_on = [ zia_traffic_forwarding_static_ip.test-gre-unnumbered-static-ip ]
}

resource "zia_traffic_forwarding_gre_tunnel" "test-gre-tunnel-unnumbered" {
	source_ip = zia_traffic_forwarding_static_ip.test-gre-unnumbered-static-ip.ip_address
	comment   = "test-gre-tunnel-unnumbered-%s"
	within_country = true
	country_code = "CA"
	ip_unnumbered = true
	primary_dest_vip {
		id = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-unnumbered-vip.list[0].id
		virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-unnumbered-vip.list[0].virtual_ip
	}
	secondary_dest_vip {
		id = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-unnumbered-vip.list[1].id
		virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-unnumbered-vip.list[1].virtual_ip
	}
	depends_on = [ zia_traffic_forwarding_static_ip.test-gre-unnumbered-static-ip, data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-unnumbered-vip ]
}

data "zia_traffic_forwarding_gre_tunnel_info" "get_gre_unnumbered" {
	ip_address = "121.234.54.85"
}
	`, rComment)
}

func testAccDataSourceTrafficForwardingGreTunnelsNumberedBasic(rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test-gre-numbered-static-ip"{
	ip_address =  "121.234.54.86"
	routable_ip = true
	geo_override = false
}

data "zia_traffic_forwarding_gre_vip_recommended_list" "test-gre-numbered-vip"{
	source_ip = zia_traffic_forwarding_static_ip.test-gre-numbered-static-ip.ip_address
	required_count = 2
	depends_on = [ zia_traffic_forwarding_static_ip.test-gre-numbered-static-ip ]
}

resource "zia_traffic_forwarding_gre_tunnel" "test-gre-tunnel-unnumbered" {
	source_ip = zia_traffic_forwarding_static_ip.test-gre-numbered-static-ip.ip_address
	comment   = "test-gre-tunnel-numbered-%s"
	within_country = true
	country_code = "CA"
	ip_unnumbered = false
	primary_dest_vip {
		id = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-numbered-vip.list[0].id
		virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-numbered-vip.list[0].virtual_ip
	}
	secondary_dest_vip {
		id = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-numbered-vip.list[1].id
		virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-numbered-vip.list[1].virtual_ip
	}
	depends_on = [ zia_traffic_forwarding_static_ip.test-gre-static-ip, data.zia_traffic_forwarding_gre_vip_recommended_list.test-gre-numbered-vip ]
}

data "zia_traffic_forwarding_gre_tunnel_info" "get_gre_numbered" {
	ip_address = "121.234.54.86"
}
	`, rComment)
}

func testAccDataSourceTrafficForwardingGreTunnels(source_ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[source_ip]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", source_ip)
		}

		return nil
	}
}
*/
