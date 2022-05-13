package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/terraform-provider-zia/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceFWIPSourceGroupsBasic(t *testing.T) {
	var groups ipsourcegroups.IPSourceGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringSourceGroup)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWIPSourceGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWIPSourceGroupsConfigure(resourceTypeAndName, generatedName, variable.FWSRCGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPSourceGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWSRCGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_addresses.#", "3"),
				),
			},

			// Update test
			{
				Config: testAccCheckFWIPSourceGroupsConfigure(resourceTypeAndName, generatedName, variable.FWSRCGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWIPSourceGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FWSRCGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_addresses.#", "3"),
				),
			},
		},
	})
}

func testAccCheckFWIPSourceGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.FWFilteringSourceGroup {
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

func testAccCheckFWIPSourceGroupsExists(resource string, rule *ipsourcegroups.IPSourceGroups) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.ipsourcegroups.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckFWIPSourceGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "tf-acc-test-%s"
	description = "%s"
    ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
  }

  data "%s" "%s" {
	id = "${%s.id}"
  }

`,
		// resource variables
		resourcetype.FWFilteringSourceGroup,
		generatedName,
		generatedName,
		description,

		// data source variables
		resourcetype.FWFilteringSourceGroup,
		generatedName,
		resourceTypeAndName,
	)
}
