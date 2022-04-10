package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceTrafficForwardingVPNCredentials_Basic(t *testing.T) {
	rEmail := acctest.RandString(5)
	rComment := acctest.RandString(5)
	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")
	resourceName1 := "data.zia_traffic_forwarding_vpn_credentials.test-type-ip"
	resourceName2 := "data.zia_traffic_forwarding_vpn_credentials.test-type-fqdn"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsTypeIPBasic(rIP, rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingVPNCredentials(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "comments", "test-type-ip-"+rComment),
					resource.TestCheckResourceAttr(resourceName1, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName1, "ip_address", rIP),
				),
			},
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsUFQDNBasic(rEmail, rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingVPNCredentials(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "comments", "test-type-fqdn-"+rComment),
					resource.TestCheckResourceAttr(resourceName2, "type", "UFQDN"),
					resource.TestCheckResourceAttr(resourceName2, "fqdn", rEmail+"@securitygeek.io"),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingVPNCredentialsTypeIPBasic(rIP, rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "static_ip"{
	ip_address =  "%s"
	comment = "test-type-ip-%s"
	routable_ip = true
	geo_override = true
	latitude = -36.848461
	longitude = 174.763336
}

resource "zia_traffic_forwarding_vpn_credentials" "test-type-ip"{
	type = "IP"
	ip_address = zia_traffic_forwarding_static_ip.static_ip.ip_address
	comments = "test-type-ip-%s"
	pre_shared_key = "newPassword123!"
	depends_on = [ zia_traffic_forwarding_static_ip.static_ip ]
}

data "zia_traffic_forwarding_vpn_credentials" "test-type-ip" {
	type = "IP"
	ip_address = zia_traffic_forwarding_vpn_credentials.test-type-ip.ip_address
	depends_on = [ zia_traffic_forwarding_vpn_credentials.test-type-ip ]
}
	`, rIP, rComment, rComment)
}

func testAccCheckTrafficForwardingVPNCredentialsUFQDNBasic(rEmail, rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_vpn_credentials" "test-type-fqdn"{
	type = "UFQDN"
	fqdn = "%s@securitygeek.io"
	comments = "test-type-fqdn-%s"
	pre_shared_key = "newPassword123!"
}

data "zia_traffic_forwarding_vpn_credentials" "test-type-fqdn" {
	fqdn = zia_traffic_forwarding_vpn_credentials.test-type-fqdn.fqdn
	depends_on = [ zia_traffic_forwarding_vpn_credentials.test-type-fqdn ]
}
	`, rEmail, rComment)
}

func testAccDataSourceTrafficForwardingVPNCredentials(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
