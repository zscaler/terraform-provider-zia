package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkapplications"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceFWNetworkApplicationGroupsBasic(t *testing.T) {
	var appGroups networkapplications.NetworkApplicationGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringNetworkAppGroups)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkApplicationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkApplicationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWAppGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkApplicationGroupsExists(resourceTypeAndName, &appGroups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.FWAppGroupName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWAppGroupDescription),
				),
			},

			// Update test
			{
				Config: testAccCheckFWNetworkApplicationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWAppGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkApplicationGroupsExists(resourceTypeAndName, &appGroups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.FWAppGroupName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWAppGroupDescription),
				),
			},
		},
	})
}

func testAccCheckFWNetworkApplicationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FWFilteringNetworkAppGroups {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.networkapplications.GetNetworkApplicationGroups(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network application group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFWNetworkApplicationGroupsExists(resource string, rule *networkapplications.NetworkApplicationGroups) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.networkapplications.GetNetworkApplicationGroups(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkApplicationGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`
// network application group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		FWNetworkApplicationGroupsResourceHCL(generatedName, description),

		// data source variables
		resourcetype.FWFilteringNetworkAppGroups,
		generatedName,
		resourceTypeAndName,
	)
}

func FWNetworkApplicationGroupsResourceHCL(generatedName, description string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
    description = "%s"
    network_applications = [
            "YAMMER",
            "OFFICE365",
            "SKYPE_FOR_BUSINESS",
            "OUTLOOK",
            "SHAREPOINT",
            "SHAREPOINT_ADMIN",
            "SHAREPOINT_BLOG",
            "SHAREPOINT_CALENDAR",
            "SHAREPOINT_DOCUMENT",
            "SHAREPOINT_ONLINE",
            "ONEDRIVE"
    ]
}
`,
		// resource variables
		resourcetype.FWFilteringNetworkAppGroups,
		generatedName,
		variable.FWAppGroupName,
		description,
	)
}
