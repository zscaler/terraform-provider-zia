package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/staticips"
)

func TestAccResourceTrafficForwardingStaticIPBasic(t *testing.T) {
	var static staticips.StaticIP
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingStaticIP)
	rIP, _ := acctest.RandIpAddress("104.238.235.0/24")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingStaticIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists(resourceTypeAndName, &static),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", rIP),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "geo_override", strconv.FormatBool(variable.StaticGeoOverride)),
				),
			},

			// Update static ip address information
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists(resourceTypeAndName, &static),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", rIP),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "geo_override", strconv.FormatBool(variable.StaticGeoOverride)),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingStaticIPDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficForwardingStaticIP {
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

func testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, ipAddress string, routableIP, geoOverride bool) string {
	return fmt.Sprintf(`

// location management resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		getTrafficForwardingStaticIPConfigure(generatedName, ipAddress, routableIP, geoOverride),

		// data source variables
		resourcetype.TrafficForwardingStaticIP,
		generatedName,
		resourceTypeAndName,
	)
}

func getTrafficForwardingStaticIPConfigure(generatedName, ipAddress string, routableIP, geoOverride bool) string {
	return fmt.Sprintf(`

resource "%s" "%s" {
	comment = "tf-acc-test-%s"
    ip_address =  "%s"
    routable_ip = "%s"
    geo_override = "%s"
}
`,
		// resource variables
		resourcetype.TrafficForwardingStaticIP,
		generatedName,
		generatedName,
		ipAddress,
		strconv.FormatBool(routableIP),
		strconv.FormatBool(geoOverride),
	)
}
