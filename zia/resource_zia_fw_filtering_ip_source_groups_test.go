package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipsourcegroups"
)

func TestAccResourceFWIPSourceGroups_basic(t *testing.T) {
	var groups ipsourcegroups.IPSourceGroups
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_firewall_filtering_ip_source_groups.test-fw-src-group"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPSourceGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWIPSourceGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPSourceGroupsExists("zia_firewall_filtering_ip_source_groups.test-fw-src-group", &groups),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-src-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-src-group-"+rDesc),
					// resource.TestCheckResourceAttr(resourceName, "ip_addresses", "ip_addresses"),
				),
			},
		},
	})
}

func testAccCheckFWIPSourceGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_ip_source_groups" "test-fw-src-group"{
	name = "test-fw-src-group-%s"
	description = "test-fw-src-group-%s"
	ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}
	`, rName, rDesc)
}

func testAccCheckFWIPSourceGroupsExists(resource string, group *ipsourcegroups.IPSourceGroups) resource.TestCheckFunc {
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
		receivedGroup, err := apiClient.ipsourcegroups.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*group = *receivedGroup

		return nil
	}
}

func testAccCheckFWIPSourceGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_ip_source_groups" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.ipsourcegroups.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("ip source group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}
