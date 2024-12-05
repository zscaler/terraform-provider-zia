package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipdestinationgroups"
)

func TestAccResourceFWIPDestinationGroupsBasic(t *testing.T) {
	var groups ipdestinationgroups.IPDestinationGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPDestinationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, initialName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDSTGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWDSTGroupTypeDSTNFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "addresses.#", "3"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, updatedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPDestinationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWDSTGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FWDSTGroupTypeDSTNFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "addresses.#", "3"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckFWIPDestinationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FWFilteringDestinationGroup {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := ipdestinationgroups.Get(context.Background(), service, id)

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
		service := apiClient.Service

		receivedRule, err := ipdestinationgroups.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWIPDestinationGroupsConfigure(resourceTypeAndName, generatedName, description, dst_type string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
	type        = "%s"
	addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
  }

  data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// Resource type and name for the destination ip group
		resourcetype.FWFilteringDestinationGroup,
		resourceName,
		generatedName,
		description,
		dst_type,

		// Data source type and name
		resourcetype.FWFilteringDestinationGroup,
		resourceName,

		// Reference to the resource
		resourcetype.FWFilteringDestinationGroup,
		resourceName,
	)
}
