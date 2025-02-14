package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
)

func TestAccResourceTrafficForwardingStaticIPBasic(t *testing.T) {
	var static staticips.StaticIP
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingStaticIP)
	rIP, _ := acctest.RandIpAddress("104.238.235.0/24")

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingStaticIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, initialName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists(resourceTypeAndName, &static),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", rIP),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "geo_override", strconv.FormatBool(variable.StaticGeoOverride)),
				),
			},

			// Update static ip address information
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, updatedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingStaticIPExists(resourceTypeAndName, &static),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", rIP),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "geo_override", strconv.FormatBool(variable.StaticGeoOverride)),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTrafficForwardingStaticIPDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficForwardingStaticIP {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := staticips.Get(context.Background(), service, id)

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
		service := apiClient.Service

		receivedRule, err := staticips.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, ipAddress string, routableIP, geoOverride bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`

resource "%s" "%s" {
	comment = "%s"
    ip_address =  "%s"
    routable_ip = "%s"
    geo_override = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.TrafficForwardingStaticIP,
		resourceName,
		generatedName,
		ipAddress,
		strconv.FormatBool(routableIP),
		strconv.FormatBool(geoOverride),

		// Data source type and name
		resourcetype.TrafficForwardingStaticIP, resourceName,

		// Reference to the resource
		resourcetype.TrafficForwardingStaticIP, resourceName,
	)
}
