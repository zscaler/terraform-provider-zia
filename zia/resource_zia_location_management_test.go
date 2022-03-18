package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceLocationManagementBasic(t *testing.T) {
	var locations locationmanagement.Locations
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.LocationManagement)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, variable.LocationDescription, variable.LocationCountry, variable.LocationTZ, variable.LocationAuthRequired),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationManagementExists(resourceTypeAndName, &locations),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.LocationDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "country", variable.LocationCountry),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tz", variable.LocationTZ),
					resource.TestCheckResourceAttr(resourceTypeAndName, "auth_required", strconv.FormatBool(variable.LocationAuthRequired)),
				),
			},

			// Update test
			{
				Config: testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, variable.LocationDescription, variable.LocationCountry, variable.LocationTZ, variable.LocationAuthRequired),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationManagementExists(resourceTypeAndName, &locations),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.LocationDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "country", variable.LocationCountry),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tz", variable.LocationTZ),
					resource.TestCheckResourceAttr(resourceTypeAndName, "auth_required", strconv.FormatBool(variable.LocationAuthRequired)),
				),
			},
		},
	})
}

func testAccCheckLocationManagementDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.LocationManagement {
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

func testAccCheckLocationManagementExists(resource string, location *locationmanagement.Locations) resource.TestCheckFunc {
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
		receivedLocation, err := apiClient.locationmanagement.GetLocation(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*location = *receivedLocation

		return nil
	}
}

func testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, description, country, tz string, authRequired bool) string {
	return fmt.Sprintf(`
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		LocationManagementResourceHCL(generatedName, description, country, tz, authRequired),

		// data source variables
		resourcetype.LocationManagement,
		generatedName,
		resourceTypeAndName,
	)
}

func LocationManagementResourceHCL(generatedName, description, country, tz string, authRequired bool) string {
	return fmt.Sprintf(`
resource "zia_traffic_forwarding_static_ip" "testAcc_usa_sjc40"{
	ip_address 		= "118.189.211.222"
	routable_ip 	= true
	comment 		= "SJC37 - Static IP"
	geo_override 	= false
}

resource "zia_traffic_forwarding_vpn_credentials" "testAcc_usa_sjc40"{
    type 			= "UFQDN"
    fqdn 			= "usa_sjc40@securitygeek.io"
    comments    	= "Acceptance Test"
    pre_shared_key 	= "newPassword123!"
	depends_on 		= [ zia_traffic_forwarding_static_ip.testAcc_usa_sjc40 ]
}

resource "%s" "%s" {
    name 					= "%s"
    description 			= "%s"
    country 				= "%s"
    tz 						= "%s"
    auth_required 			= "%s"
    idle_time_in_minutes 	= 720
    display_time_unit 		= "HOUR"
    surrogate_ip 			= true
    xff_forward_enabled 	= true
    ofw_enabled 			= true
    ips_control 			= true
    ip_addresses 			= [ zia_traffic_forwarding_static_ip.testAcc_usa_sjc40.ip_address ]
    vpn_credentials {
       id 	= zia_traffic_forwarding_vpn_credentials.testAcc_usa_sjc40.vpn_credental_id
       type = zia_traffic_forwarding_vpn_credentials.testAcc_usa_sjc40.type
    }
    depends_on = [ zia_traffic_forwarding_static_ip.testAcc_usa_sjc40, zia_traffic_forwarding_vpn_credentials.testAcc_usa_sjc40 ]
}
`,
		// resource variables
		resourcetype.LocationManagement,
		generatedName,
		generatedName,
		description,
		country,
		tz,
		strconv.FormatBool(authRequired),
	)
}
