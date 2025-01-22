package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnels"
)

func TestAccResourceTrafficForwardingGRETunnelBasic(t *testing.T) {
	var gretunnel gretunnels.GreTunnels
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingGRETunnel)

	randomIP, _ := acctest.RandIpAddress("104.238.235.0/24")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingGRETunnelDestroy,
		Steps: []resource.TestStep{
			{
				// create gree tunnel
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, variable.GRETunnelComment, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered, randomIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingGRETunnelExists(resourceTypeAndName, &gretunnel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", variable.GRETunnelComment),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
			},

			// update
			{
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, variable.GRETunnelComment, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered, randomIP),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingGRETunnelExists(resourceTypeAndName, &gretunnel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", variable.GRETunnelComment),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"country_code",
					"within_country",
				},
			},
		},
	})
}

func testAccCheckTrafficForwardingGRETunnelDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficForwardingGRETunnel {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := gretunnels.GetGreTunnels(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("gre tunnel with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckTrafficForwardingGRETunnelExists(resource string, rule *gretunnels.GreTunnels) resource.TestCheckFunc {
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

		receivedGRETunnel, err := gretunnels.GetGreTunnels(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedGRETunnel

		return nil
	}
}

func testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, comment string, withinCountry, ipUnnumbered bool, randomIP string) string {
	return fmt.Sprintf(`

	// gre tunnel resource
	%s

	data "%s" "%s" {
	  id = "${%s.id}"
	}
`,
		// resource variables
		getTrafficForwardingGRETunnel_HCL(generatedName, comment, withinCountry, ipUnnumbered, randomIP),

		// data source variables
		resourcetype.TrafficForwardingGRETunnel,
		generatedName,
		resourceTypeAndName,
	)
}

func getTrafficForwardingGRETunnel_HCL(generatedName, comment string, withinCountry, ipUnnumbered bool, randomIP string) string {
	return fmt.Sprintf(`

data "zia_gre_internal_ip_range_list" "this"{
    required_count = 1
}
resource "zia_traffic_forwarding_static_ip" "this"{
    ip_address =  "%s"
	comment        = "GRE Tunnel Created with Terraform"
    routable_ip = true
    geo_override = true
    latitude = 43.7063
    longitude = -79.4202
}

resource "%s" "%s" {
	comment        = "%s"
	internal_ip_range = data.zia_gre_internal_ip_range_list.this.list[0].start_ip_address
	source_ip = zia_traffic_forwarding_static_ip.this.ip_address
	country_code   = "CA"
    within_country = "%s"
    ip_unnumbered  = "%s"
	lifecycle {
		ignore_changes = [
		internal_ip_range,
		]
	}
	depends_on = [zia_traffic_forwarding_static_ip.this]
}
`,
		// resource variables
		randomIP,
		resourcetype.TrafficForwardingGRETunnel,
		generatedName,
		comment,
		strconv.FormatBool(withinCountry),
		strconv.FormatBool(ipUnnumbered),
	)
}
