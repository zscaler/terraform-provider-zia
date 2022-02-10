package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceFWIPDestinationGroupsBasic(t *testing.T) {
	var groups ipdestinationgroups.IPDestinationGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.FWDSTGroupName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDSTGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWDSTGroupTypeDSTNFQDN),
				),
			},

			// Update test
			{
				Config: testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.FWDSTGroupName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDSTGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWDSTGroupTypeDSTNFQDN),
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
// ip destination group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		FWIPDestinationGroupsResourceHCL(generatedName, description, dst_type),

		// data source variables
		resourcetype.FWFilteringDestinationGroup,
		generatedName,
		resourceTypeAndName,
	)
}

func FWIPDestinationGroupsResourceHCL(generatedName, description, dst_type string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
	type        = "%s"
	addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
  }
`,
		// resource variables
		resourcetype.FWFilteringDestinationGroup,
		generatedName,
		variable.FWDSTGroupName,
		description,
		dst_type,
	)
}
