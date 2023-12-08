package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkapplicationgroups"
)

func TestAccResourceFWNetworkApplicationGroupsBasic(t *testing.T) {
	var appGroups networkapplicationgroups.NetworkApplicationGroups
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
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWAppGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "network_applications.#", "11"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWNetworkApplicationGroupsConfigure(resourceTypeAndName, generatedName, variable.FWAppGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWNetworkApplicationGroupsExists(resourceTypeAndName, &appGroups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWAppGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "network_applications.#", "11"),
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

		rule, err := apiClient.networkapplicationgroups.GetNetworkApplicationGroups(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("network application group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckFWNetworkApplicationGroupsExists(resource string, rule *networkapplicationgroups.NetworkApplicationGroups) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.networkapplicationgroups.GetNetworkApplicationGroups(id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWNetworkApplicationGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "tf-acc-test-%s"
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

data "%s" "%s" {
	id = "${%s.id}"
}
`,
		// resource variables
		resourcetype.FWFilteringNetworkAppGroups,
		generatedName,
		generatedName,
		description,

		// data source variables
		resourcetype.FWFilteringNetworkAppGroups,
		generatedName,
		resourceTypeAndName,
	)
}
