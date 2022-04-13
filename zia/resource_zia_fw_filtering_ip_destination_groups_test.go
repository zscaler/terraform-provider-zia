package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
)

func TestAccResourceFWIPDestinationGroups_basic(t *testing.T) {
	var groups ipdestinationgroups.IPDestinationGroups
	rName1 := acctest.RandString(5)
	rDesc1 := acctest.RandString(20)
	rName2 := acctest.RandString(5)
	rDesc2 := acctest.RandString(20)
	rName3 := acctest.RandString(5)
	rDesc3 := acctest.RandString(20)
	rName4 := acctest.RandString(5)
	rDesc4 := acctest.RandString(20)
	resourceName1 := "zia_firewall_filtering_destination_groups.test-fw-dst-fqdn-group"
	resourceName2 := "zia_firewall_filtering_destination_groups.test-fw-dst-ip-group"
	resourceName3 := "zia_firewall_filtering_destination_groups.test-fw-dst-domain-group"
	resourceName4 := "zia_firewall_filtering_destination_groups.test-fw-dst-other-group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFWIPDestinationGroupsDstFQDNBasic(rName1, rDesc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists("zia_firewall_filtering_destination_groups.test-fw-dst-fqdn-group", &groups),
					resource.TestCheckResourceAttr(resourceName1, "name", "test-fw-dst-fqdn-group-"+rName1),
					resource.TestCheckResourceAttr(resourceName1, "description", "test-fw-dst-fqdn-group-"+rDesc1),
					resource.TestCheckResourceAttr(resourceName1, "type", "DSTN_FQDN"),
				),
			},
			// Test IP Destination IP Group
			{
				Config: testAccResourceFWIPDestinationGroupsDstIPBasic(rName2, rDesc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists("zia_firewall_filtering_destination_groups.test-fw-dst-ip-group", &groups),
					resource.TestCheckResourceAttr(resourceName2, "name", "test-fw-dst-ip-group-"+rName2),
					resource.TestCheckResourceAttr(resourceName2, "description", "test-fw-dst-ip-group-"+rDesc2),
					resource.TestCheckResourceAttr(resourceName2, "type", "DSTN_IP"),
				),
			},
			// Test IP Destination Domain Group
			{
				Config: testAccResourceFWIPDestinationGroupsDstDomainBasic(rName3, rDesc3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists("zia_firewall_filtering_destination_groups.test-fw-dst-domain-group", &groups),
					resource.TestCheckResourceAttr(resourceName3, "name", "test-fw-dst-domain-group-"+rName3),
					resource.TestCheckResourceAttr(resourceName3, "description", "test-fw-dst-domain-group-"+rDesc3),
					resource.TestCheckResourceAttr(resourceName3, "type", "DSTN_DOMAIN"),
				),
			},
			// Test IP Destination Other Group
			{
				Config: testAccResourceFWIPDestinationGroupsDstOtherBasic(rName4, rDesc4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists("zia_firewall_filtering_destination_groups.test-fw-dst-other-group", &groups),
					resource.TestCheckResourceAttr(resourceName4, "name", "test-fw-dst-other-group-"+rName4),
					resource.TestCheckResourceAttr(resourceName4, "description", "test-fw-dst-other-group-"+rDesc4),
					resource.TestCheckResourceAttr(resourceName4, "type", "DSTN_OTHER"),
				),
			},
		},
	})
}

func testAccResourceFWIPDestinationGroupsDstFQDNBasic(rName1, rDesc1 string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_destination_groups" "test-fw-dst-fqdn-group" {
	name        = "test-fw-dst-fqdn-group-%s"
	description = "test-fw-dst-fqdn-group-%s"
	type        = "DSTN_FQDN"
	addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
  }
	`, rName1, rDesc1)
}

func testAccResourceFWIPDestinationGroupsDstIPBasic(rName2, rDesc2 string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_destination_groups" "test-fw-dst-ip-group" {
	name        = "test-fw-dst-ip-group-%s"
	description = "test-fw-dst-ip-group-%s"
	type        = "DSTN_IP"
	addresses = [ "192.168.100.1", "192.168.100.2", "192.168.100.3" ]
  }
	`, rName2, rDesc2)
}

func testAccResourceFWIPDestinationGroupsDstDomainBasic(rName3, rDesc3 string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_destination_groups" "test-fw-dst-domain-group" {
	name        = "test-fw-dst-domain-group-%s"
	description = "test-fw-dst-domain-group-%s"
	type        = "DSTN_DOMAIN"
	addresses 	= [ ".terraformtest1.com", ".terraformtest2.com", ".terraformtest3.com" ]
  }
	`, rName3, rDesc3)
}

func testAccResourceFWIPDestinationGroupsDstOtherBasic(rName4, rDesc4 string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_destination_groups" "test-fw-dst-other-group" {
	name        	= "test-fw-dst-other-group-%s"
	description 	= "test-fw-dst-other-group-%s"
	type        	= "DSTN_OTHER"
	ip_categories 	= [ "CUSTOM_01", "CUSTOM_03" ]
	countries 		= [ "COUNTRY_US", "COUNTRY_CA" ]

  }
	`, rName4, rDesc4)
}

func testAccCheckFWIPDestinationGroupsExists(resource string, rule *ipdestinationgroups.IPDestinationGroups) resource.TestCheckFunc {
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
		receivedGroup, err := apiClient.ipdestinationgroups.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedGroup

		return nil
	}
}

func testAccCheckFWIPDestinationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_destination_groups" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.ipdestinationgroups.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("ip destination group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
