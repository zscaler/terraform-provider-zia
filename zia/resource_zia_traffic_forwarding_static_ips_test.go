package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
)

func TestAccResourceTrafficForwardingStaticIPBasic(t *testing.T) {
	var static staticips.StaticIP
	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")
	rComment := acctest.RandString(5)
	resourceName := "zia_traffic_forwarding_static_ip.test-static-ip"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingStaticIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceTrafficForwardingStaticIPBasic(rIP, rComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists("zia_traffic_forwarding_static_ip.test-static-ip", &static),
					resource.TestCheckResourceAttr(resourceName, "comment", "test-static-ip-"+rComment),
					resource.TestCheckResourceAttr(resourceName, "ip_address", rIP),
					resource.TestCheckResourceAttr(resourceName, "routable_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "geo_override", "true"),
					resource.TestCheckResourceAttr(resourceName, "latitude", "-36.848461"),
					resource.TestCheckResourceAttr(resourceName, "longitude", "174.763336"),
				),
			},
		},
	})
}

func testAccCheckResourceTrafficForwardingStaticIPBasic(rIP, rComment string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test-static-ip"{
	ip_address =  "%s"
	routable_ip = true
	geo_override = true
	latitude = -36.848461
	longitude = 174.763336
	comment = "test-static-ip-%s"
}
	`, rIP, rComment)
}

func testAccCheckTrafficForwardingStaticIPExists(resource string, static *staticips.StaticIP) resource.TestCheckFunc {
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
		receivedStatic, err := apiClient.staticips.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*static = *receivedStatic

		return nil
	}
}

func testAccCheckTrafficForwardingStaticIPDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_traffic_forwarding_static_ip" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.staticips.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("static ip address with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
