package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
)

func TestAccResourceTrafficForwardingVPNCredentials_basic(t *testing.T) {
	var credentials vpncredentials.VPNCredentials
	rName := acctest.RandString(5)
	rComment := acctest.RandString(5)
	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")
	resourceName := "zia_traffic_forwarding_vpn_credentials.test-type-ip"
	resourceName2 := "zia_traffic_forwarding_vpn_credentials.test-type-fqdn"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingVPNCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsTypeIP(rIP, rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingVPNCredentialsExists("zia_traffic_forwarding_vpn_credentials.test-type-ip", &credentials),
					resource.TestCheckResourceAttr(resourceName, "comments", "test-type-ip-"+rComment),
					resource.TestCheckResourceAttr(resourceName, "type", "IP"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", rIP),
					resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "newPassword123!"),
				),
			},

			// Test VPN Credential Type UFQDN
			{
				Config: testAccCheckTrafficForwardingStaticIPTypeUFQDN(rName, rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingVPNCredentialsExists("zia_traffic_forwarding_vpn_credentials.test-type-fqdn", &credentials),
					resource.TestCheckResourceAttr(resourceName2, "comments", "test-type-fqdn-"+rComment),
					resource.TestCheckResourceAttr(resourceName2, "type", "UFQDN"),
					resource.TestCheckResourceAttr(resourceName2, "fqdn", rName+"@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceName2, "pre_shared_key", "newPassword123!"),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingVPNCredentialsTypeIP(rIP, rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "static_ip"{
	ip_address =  "%s"
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
	depends_on = [zia_traffic_forwarding_static_ip.static_ip]
}
	`, rIP, rComment)
}

func testAccCheckTrafficForwardingStaticIPTypeUFQDN(rName, rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_vpn_credentials" "test-type-fqdn"{
	type = "UFQDN"
	fqdn = "%s@securitygeek.io"
	comments = "test-type-fqdn-%s"
	pre_shared_key = "newPassword123!"
}
	`, rName, rComment)
}

func testAccCheckTrafficForwardingVPNCredentialsExists(resource string, credentials *vpncredentials.VPNCredentials) resource.TestCheckFunc {
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
		receivedVpnCredentials, err := apiClient.vpncredentials.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*credentials = *receivedVpnCredentials

		return nil
	}
}

func testAccCheckTrafficForwardingVPNCredentialsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_traffic_forwarding_vpn_credentials" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.vpncredentials.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("vpn credentials with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
