package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceLocationManagementBasic(t *testing.T) {
	var locations locationmanagement.Locations
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringLocManagement)

	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")
	staticIPTypeAndName, _, staticIPGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringStaticIP)
	staticIPResourceHCL := testAccCheckTrafficForwardingStaticIPConfigure(staticIPTypeAndName, staticIPGeneratedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride)

	rSharedKey := acctest.RandString(20)
	vpnCredentialTypeAndName, _, vpnCredentialGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringVPNCredentials)
	vpnCredentialResourceHCL := testAccCheckTrafficForwardingVPNCredentialsIPConfigure(vpnCredentialTypeAndName, vpnCredentialGeneratedName, vpnCredentialGeneratedName, variable.VPNCredentialTypeIP, rSharedKey)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, vpnCredentialResourceHCL, vpnCredentialTypeAndName, variable.LocAuthRequired, variable.LocSurrogateIP, variable.LocXFF, variable.LocOFW, variable.LocIPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationManagementExists(resourceTypeAndName, &locations),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "country", "UNITED_STATES"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tz", "UNITED_STATES_AMERICA_LOS_ANGELES"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "profile", "CORPORATE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "display_time_unit", "HOUR"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "auth_required", strconv.FormatBool(variable.LocAuthRequired)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "surrogate_ip", strconv.FormatBool(variable.LocSurrogateIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "xff_forward_enabled", strconv.FormatBool(variable.LocXFF)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ofw_enabled", strconv.FormatBool(variable.LocOFW)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ips_control", strconv.FormatBool(variable.LocIPS)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_addresses.#", "1"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "vpn_credentials.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, vpnCredentialResourceHCL, vpnCredentialTypeAndName, variable.LocAuthRequired, variable.LocSurrogateIP, variable.LocXFF, variable.LocOFW, variable.LocIPS),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationManagementExists(resourceTypeAndName, &locations),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "country", "UNITED_STATES"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tz", "UNITED_STATES_AMERICA_LOS_ANGELES"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "profile", "CORPORATE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "display_time_unit", "HOUR"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "auth_required", strconv.FormatBool(variable.LocAuthRequired)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "surrogate_ip", strconv.FormatBool(variable.LocSurrogateIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "xff_forward_enabled", strconv.FormatBool(variable.LocXFF)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ofw_enabled", strconv.FormatBool(variable.LocOFW)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ips_control", strconv.FormatBool(variable.LocIPS)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_addresses.#", "1"),
					// resource.TestCheckResourceAttr(resourceTypeAndName, "vpn_credentials.#", "1"),
				),
			},
		},
	})
}

func testAccCheckLocationManagementDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficFilteringLocManagement {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.locationmanagement.GetLocation(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("location with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckLocationManagementExists(resource string, rule *locationmanagement.Locations) resource.TestCheckFunc {
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
		receivedLoc, err := apiClient.locationmanagement.GetLocation(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedLoc

		return nil
	}
}

func testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, vpnCredentialResourceHCL, vpnCredentialTypeAndName string, authRequired, surrogateIP, xffEnabled, ofwEnabled, ipsEnabled bool) string {
	return fmt.Sprintf(`

// static ip resource
%s

// location management resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		staticIPResourceHCL,
		// vpnCredentialResourceHCL,
		getLocationManagementHCL(generatedName, staticIPTypeAndName, vpnCredentialTypeAndName, authRequired, surrogateIP, xffEnabled, ofwEnabled, ipsEnabled),

		// data source variables
		resourcetype.TrafficFilteringLocManagement,
		generatedName,
		resourceTypeAndName,
	)
}

func getLocationManagementHCL(generatedName, staticIPTypeAndName, vpnCredentialTypeAndName string, authRequired, surrogateIP, xffEnabled, ofwEnabled, ipsEnabled bool) string {
	return fmt.Sprintf(`


resource "%s" "%s" {
	name 					= "tf-acc-test-%s"
	description 			= "tf-acc-test-%s"
	country 				= "UNITED_STATES"
	tz 						= "UNITED_STATES_AMERICA_LOS_ANGELES"
	auth_required 			= "%s"
	surrogate_ip 			= "%s"
	xff_forward_enabled 	= "%s"
	ofw_enabled 			= "%s"
	ips_control 			= "%s"
	idle_time_in_minutes 	= 720
	display_time_unit 		= "HOUR"
	profile					= "CORPORATE"
	ip_addresses			= [ "${%s.ip_address}"]
	depends_on = [ %s ]
}
`,

		// resource variables
		resourcetype.TrafficFilteringLocManagement,
		generatedName,
		generatedName,
		generatedName,
		strconv.FormatBool(authRequired),
		strconv.FormatBool(surrogateIP),
		strconv.FormatBool(xffEnabled),
		strconv.FormatBool(ofwEnabled),
		strconv.FormatBool(ipsEnabled),
		staticIPTypeAndName,
		// vpnCredentialTypeAndName,
		// vpnCredentialTypeAndName,
		staticIPTypeAndName,
	)
}
