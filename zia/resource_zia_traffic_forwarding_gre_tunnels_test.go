package zia

/*
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
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/gretunnels"
)

func TestAccResourceTrafficForwardingGRETunnelBasic(t *testing.T) {
	var gretunnel gretunnels.GreTunnels
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingGRETunnel)

	rIP, _ := acctest.RandIpAddress("104.238.235.0/24")
	staticIPTypeAndName, _, staticIPGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingStaticIP)
	staticIPResourceHCL := testAccCheckTrafficForwardingStaticIPConfigure(staticIPTypeAndName, staticIPGeneratedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingGRETunnelDestroy,
		Steps: []resource.TestStep{
			{
				// create gree tunnel
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingGRETunnelExists(resourceTypeAndName, &gretunnel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
				ExpectNonEmptyPlan: true,
			},

			// update
			{
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingGRETunnelExists(resourceTypeAndName, &gretunnel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comment", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
				ExpectNonEmptyPlan: true,
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

func testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName string, withinCountry, ipUnnumbered bool) string {
	return fmt.Sprintf(`

	// static ip resource
	%s

	// gre tunnel resource
	%s

	data "%s" "%s" {
	  id = "${%s.id}"
	}
`,
		// resource variables
		staticIPResourceHCL,
		getTrafficForwardingGRETunnel_HCL(generatedName, staticIPTypeAndName, withinCountry, ipUnnumbered),

		// data source variables
		resourcetype.TrafficForwardingGRETunnel,
		generatedName,
		resourceTypeAndName,
	)
}

func getTrafficForwardingGRETunnel_HCL(generatedName, staticIPTypeAndName string, withinCountry, ipUnnumbered bool) string {

	return fmt.Sprintf(`

resource "%s" "%s" {
	source_ip = "${%s.ip_address}"
	comment = 	"%s"
	country_code   = "CA"
    within_country = "%s"
    ip_unnumbered  = "%s"
}
`,
		// resource variables
		resourcetype.TrafficForwardingGRETunnel,
		generatedName,
		staticIPTypeAndName,
		generatedName,
		strconv.FormatBool(withinCountry),
		strconv.FormatBool(ipUnnumbered),
	)
}
*/
