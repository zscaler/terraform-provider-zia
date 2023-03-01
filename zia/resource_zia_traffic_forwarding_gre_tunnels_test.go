package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/gretunnels"
)

func TestAccResourceTrafficForwardingGRETunnelBasic(t *testing.T) {
	var gretunnel gretunnels.GreTunnels
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingGRETunnel)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingGRETunnelDestroy,
		Steps: []resource.TestStep{
			{
				// create gree tunnel
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingGRETunnelExists(resourceTypeAndName, &gretunnel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
			},

			// update
			{
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingGRETunnelExists(resourceTypeAndName, &gretunnel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingGRETunnelDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficForwardingGRETunnel {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.gretunnels.GetGreTunnels(id)

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
		receivedGRETunnel, err := apiClient.gretunnels.GetGreTunnels(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedGRETunnel

		return nil
	}
}

func testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName string, withinCountry, ipUnnumbered bool) string {
	return fmt.Sprintf(`

	// gre tunnel resource

	// static ip resource
	%s

	data "%s" "%s" {
	  id = "${%s.id}"
	}
`,
		// resource variables
		getTrafficForwardingGRETunnel_HCL(generatedName, withinCountry, ipUnnumbered),

		// data source variables
		resourcetype.TrafficForwardingGRETunnel,
		generatedName,
		resourceTypeAndName,
	)
}

func getTrafficForwardingGRETunnel_HCL(generatedName string, withinCountry, ipUnnumbered bool) string {
	return fmt.Sprintf(`

data "zia_traffic_forwarding_gre_vip_recommended_list" "this"{
	source_ip = "104.238.235.100"
	required_count = 2
}

data "zia_gre_internal_ip_range_list" "this"{
    required_count = 10
}

resource "%s" "%s" {
	source_ip = "104.238.235.100"
	comment = "tf-acc-test-%s"
	country_code   = "CA"
    within_country = "%s"
    ip_unnumbered  = "%s"
	primary_dest_vip {
		datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].datacenter
		id = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].id
		virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].virtual_ip
	  }
	  secondary_dest_vip {
		datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].datacenter
		id = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].id
		virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].virtual_ip
	  }
}
`,
		// resource variables
		resourcetype.TrafficForwardingGRETunnel,
		generatedName,
		generatedName,
		strconv.FormatBool(withinCountry),
		strconv.FormatBool(ipUnnumbered),
	)
}
