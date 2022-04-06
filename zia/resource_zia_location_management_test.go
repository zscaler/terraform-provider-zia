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
)

func TestAccResourceLocationManagement_basic(t *testing.T) {
	var locations locationmanagement.Locations
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_location_management.test_zs_sjc2022_type_ip"
	resourceName2 := "zia_location_management.test_zs_sjc2022_type_ufqdn"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationManagementDestroy,
		Steps: []resource.TestStep{
			{
				// Test Location Management with VPN Credential Type IP
				Config: testAccCheckResourceLocationManagementVPNTypeIP(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationManagementExists("zia_location_management.test_zs_sjc2022_type_ip", &locations),
					resource.TestCheckResourceAttr(resourceName, "name", "test_zs_sjc2022_type_ip-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test_zs_sjc2022_type_ip-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "country", "UNITED_STATES"),
					resource.TestCheckResourceAttr(resourceName, "tz", "UNITED_STATES_AMERICA_LOS_ANGELES"),
					resource.TestCheckResourceAttr(resourceName, "display_time_unit", "HOUR"),
					resource.TestCheckResourceAttr(resourceName, "auth_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "surrogate_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "xff_forward_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ofw_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ips_control", "true"),
				),
			},
			{
				// Test Location Management with VPN Credential Type UFQDN
				Config: testAccCheckResourceLocationManagementVPNTypeUFQDN(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationManagementExists("zia_location_management.test_zs_sjc2022_type_ufqdn", &locations),
					resource.TestCheckResourceAttr(resourceName2, "name", "test_zs_sjc2022_type_ufqdn-"+rName),
					resource.TestCheckResourceAttr(resourceName2, "description", "test_zs_sjc2022_type_ufqdn-"+rDesc),
					resource.TestCheckResourceAttr(resourceName2, "country", "UNITED_STATES"),
					resource.TestCheckResourceAttr(resourceName2, "tz", "UNITED_STATES_AMERICA_LOS_ANGELES"),
					resource.TestCheckResourceAttr(resourceName2, "display_time_unit", "HOUR"),
					resource.TestCheckResourceAttr(resourceName2, "auth_required", "true"),
					resource.TestCheckResourceAttr(resourceName2, "surrogate_ip", "true"),
					resource.TestCheckResourceAttr(resourceName2, "xff_forward_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "ofw_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "ips_control", "true"),
				),
			},
		},
	})
}

func testAccCheckResourceLocationManagementVPNTypeIP(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test_zs_sjc2022_type_ip"{
	comment 		= "Test SJC2022 - Static IP"
	ip_address 		=  "121.234.54.92"
	routable_ip 	= true
	geo_override 	= true
	latitude 		= -36.848461
	longitude 		= 174.763336
}

resource "zia_traffic_forwarding_vpn_credentials" "test_zs_sjc2022_type_ip"{
	comments    	= "Test SJC2022 - VPN Credentials"
	type        	= "IP"
	ip_address  	=  zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ip.ip_address
	pre_shared_key 	= "newPassword123!"
	depends_on 		= [ zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ip ]
}

resource "zia_location_management" "test_zs_sjc2022_type_ip"{
	name 					= "test_zs_sjc2022_type_ip-%s"
	description 			= "test_zs_sjc2022_type_ip-%s"
	country 				= "UNITED_STATES"
	tz 						= "UNITED_STATES_AMERICA_LOS_ANGELES"
	auth_required 			= true
	idle_time_in_minutes 	= 720
	display_time_unit 		= "HOUR"
	surrogate_ip 			= true
	xff_forward_enabled 	= true
	ofw_enabled 			= true
	ips_control 			= true
	ip_addresses 			= [ zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ip.ip_address ]
	depends_on 				= [ zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ip, zia_traffic_forwarding_vpn_credentials.test_zs_sjc2022_type_ip]
	vpn_credentials {
		id = zia_traffic_forwarding_vpn_credentials.test_zs_sjc2022_type_ip.vpn_credental_id
		type = zia_traffic_forwarding_vpn_credentials.test_zs_sjc2022_type_ip.type
		ip_address = zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ip.ip_address
	}
}
	`, rName, rDesc)
}

func testAccCheckResourceLocationManagementVPNTypeUFQDN(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_traffic_forwarding_static_ip" "test_zs_sjc2022_type_ufqdn"{
	comment 		= "Test SJC2022 - Static IP"
	ip_address 		=  "121.234.54.93"
	routable_ip 	= true
	geo_override 	= true
	latitude 		= -36.848461
	longitude 		= 174.763336
}

resource "zia_traffic_forwarding_vpn_credentials" "test_zs_sjc2022_type_ufqdn"{
	comments    	= "Test SJC2022 - VPN Credentials"
	type        	= "UFQDN"
	fqdn  			=  "test_zs_sjc2022_type_ufqdn@securitygeek.io"
	pre_shared_key 	= "newPassword123!"
}

resource "zia_location_management" "test_zs_sjc2022_type_ufqdn"{
	name 					= "test_zs_sjc2022_type_ufqdn-%s"
	description 			= "test_zs_sjc2022_type_ufqdn-%s"
	country 				= "UNITED_STATES"
	tz 						= "UNITED_STATES_AMERICA_LOS_ANGELES"
	auth_required 			= true
	idle_time_in_minutes 	= 720
	display_time_unit 		= "HOUR"
	surrogate_ip 			= true
	xff_forward_enabled 	= true
	ofw_enabled 			= true
	ips_control 			= true
	ip_addresses 			= [ zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ufqdn.ip_address ]
	depends_on 				= [ zia_traffic_forwarding_static_ip.test_zs_sjc2022_type_ufqdn, zia_traffic_forwarding_vpn_credentials.test_zs_sjc2022_type_ufqdn]
	vpn_credentials {
		id = zia_traffic_forwarding_vpn_credentials.test_zs_sjc2022_type_ufqdn.vpn_credental_id
		type = zia_traffic_forwarding_vpn_credentials.test_zs_sjc2022_type_ufqdn.type
	}
}
	`, rName, rDesc)
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

func testAccCheckLocationManagementDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_location_management" {
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
