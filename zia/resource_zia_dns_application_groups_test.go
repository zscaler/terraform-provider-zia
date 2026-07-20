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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/dns_application_groups"
)

func TestAccResourceDNSApplicationGroupsBasic(t *testing.T) {
	var groups dns_application_groups.DnsApplicationGroup
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DNSApplicationGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSApplicationGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDNSApplicationGroupsConfigure(resourceTypeAndName, initialName, variable.DNSAppGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSApplicationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DNSAppGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dns_applications.#", "8"),
				),
			},

			// Update test
			{
				Config: testAccCheckDNSApplicationGroupsConfigure(resourceTypeAndName, updatedName, variable.DNSAppGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSApplicationGroupsExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DNSAppGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dns_applications.#", "8"),
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

func testAccCheckDNSApplicationGroupsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DNSApplicationGroup {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := dns_application_groups.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("dns application group with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDNSApplicationGroupsExists(resource string, group *dns_application_groups.DnsApplicationGroup) resource.TestCheckFunc {
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

		receivedGroup, err := dns_application_groups.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*group = *receivedGroup

		return nil
	}
}

func testAccCheckDNSApplicationGroupsConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
	name        = "%s"
	description = "%s"
    dns_applications = ["RACKSPACE", "MCAFEE", "SOPHOS", "BRIGHTSPCACE", "CSWG", "INTECH", "SECURESERVER", "FRANCETELECOM"]
  }

data "%s" "%s" {
id = "${%s.%s.id}"
}
`,
		// Resource type and name for the ip source group
		resourcetype.DNSApplicationGroup,
		resourceName,
		generatedName,
		description,

		// Data source type and name
		resourcetype.DNSApplicationGroup,
		resourceName,

		// Reference to the resource
		resourcetype.DNSApplicationGroup,
		resourceName,
	)
}
