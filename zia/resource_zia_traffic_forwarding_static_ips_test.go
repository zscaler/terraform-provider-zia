package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceTrafficForwardingStaticIPBasic(t *testing.T) {
	var static staticips.StaticIP
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringStaticIP)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingStaticIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, variable.StaticIPAddress, variable.StaticRoutableIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists(resourceTypeAndName, &static),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", variable.StaticIPComment),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", variable.StaticIPAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
				),
			},

			// Update test
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, variable.StaticIPAddress, variable.StaticRoutableIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists(resourceTypeAndName, &static),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", variable.StaticIPComment),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", variable.StaticIPAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingStaticIPDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficFilteringStaticIP {
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

func testAccCheckTrafficForwardingStaticIPExists(resource string, rule *staticips.StaticIP) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.staticips.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, address string, routableIP bool) string {
	return fmt.Sprintf(`
// network application group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		TrafficForwardingStaticIPResourceHCL(generatedName, address, routableIP),

		// data source variables
		resourcetype.TrafficFilteringStaticIP,
		generatedName,
		resourceTypeAndName,
	)
}

func TrafficForwardingStaticIPResourceHCL(generatedName, address string, routableIP bool) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	comment = "%s"
    ip_address =  "%s"
    routable_ip = "%s"
    geo_override = true
    latitude = -36.848461
    longitude = 174.763336
}
`,
		// resource variables
		resourcetype.TrafficFilteringStaticIP,
		generatedName,
		variable.StaticIPComment,
		address,
		strconv.FormatBool(routableIP),
	)
}
