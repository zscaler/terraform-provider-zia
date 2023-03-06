package zia

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipdestinationgroups"
)

func TestAccZIAFWIPDestinationGroupsBasic(t *testing.T) {
	var groups ipdestinationgroups.IPDestinationGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)

	skipAcc := os.Getenv("SKIP_FW_IP_DESTINATION_GROUPS")
	if skipAcc == "yes" {
		t.Skip("Skipping ip destination group test as SKIP_FW_IP_DESTINATION_GROUPS is set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDSTGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWDSTGroupTypeDSTNFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "addresses.#", "3"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDSTGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWDSTGroupTypeDSTNFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "addresses.#", "3"),
				),
			},
		},
	})
}

func testAccCheckFWIPDestinationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FWFilteringDestinationGroup {
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
		receivedRule, err := apiClient.ipdestinationgroups.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, description, dst_type string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "tf-acc-test-%s"
	description = "%s"
	type        = "%s"
	addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
  }

data "%s" "%s" {
	id = "${%s.id}"
  }
`,
		// resource variables
		resourcetype.FWFilteringDestinationGroup,
		generatedName,
		generatedName,
		description,
		dst_type,

		// data source variables
		resourcetype.FWFilteringDestinationGroup,
		generatedName,
		resourceTypeAndName,
	)
}

/*
func TestAccResourceFWIPDestinationGroups_basic(t *testing.T) {
	var destIPGroup ipdestinationgroups.IPDestinationGroups

	rName := acctest.RandString(5)

	skipAcc := os.Getenv("SKIP_FW_IP_DESTINATION_GROUPS")
	if skipAcc == "yes" {
		t.Skip("Skipping FW ip destination group test as SKIP_FW_IP_DESTINATION_GROUPS is set")
	}

	resourceName := "zia_firewall_filtering_destination_groups.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWIPDestinationGroupsBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceName, &destIPGroup),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("destgroup-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("destgroup-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "DSTN_IP"),
					resource.TestCheckResourceAttr(resourceName, "addresses.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFWIPDestinationGroupsBasic(rName string) string {
	return fmt.Sprintf(`
resource "zia_firewall_filtering_destination_groups" "foo" {
	name        = "destgroup-%s"
	description = "destgroup-%s"
	addresses   = ["test1.acme.com"]
	type        = "DSTN_FQDN"
	}
`, rName, rName)
}

func testAccCheckFWIPDestinationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_firewall_filtering_destination_groups" {
			continue
		}
		foundDestGroup := ipdestinationgroups.IPDestinationGroups{
			Name:        rs.Primary.Attributes["name"],
			Description: rs.Primary.Attributes["description"],
			Type:        rs.Primary.Attributes["type"],
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		group, err := apiClient.ipdestinationgroups.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}
		if foundDestGroup.Name != rs.Primary.ID {
			return fmt.Errorf("destination_groups not found")
		}
		if group != nil {
			return fmt.Errorf("ip destination group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFWIPDestinationGroupsExists(n string, destGroup *ipdestinationgroups.IPDestinationGroups) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Destination Group Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no Destination Group ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)

		foundDestGroup := ipdestinationgroups.IPDestinationGroups{
			Name:        rs.Primary.Attributes["name"],
			Description: rs.Primary.Attributes["description"],
			Type:        rs.Primary.Attributes["type"],
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}
		foundDestGroup2, err := apiClient.ipdestinationgroups.Get(id)
		if err != nil {
			return err
		}
		if foundDestGroup2.Name+"~"+foundDestGroup2.Type != rs.Primary.ID {
			return fmt.Errorf("Destination Group not found")
		}

		*destGroup = foundDestGroup
		return nil
	}
}
*/
