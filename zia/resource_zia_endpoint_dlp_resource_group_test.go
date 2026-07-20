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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_group"
)

func TestAccResourceEndpointDLPResourceGroupBasic(t *testing.T) {
	var group endpoint_resource_group.DlpEndpointResourceGroups
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEndpointResourceGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDLPResourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEndpointDLPResourceGroupConfigure(resourceTypeAndName, initialName, variable.DLPEndpointResourceGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPResourceGroupExists(resourceTypeAndName, &group),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointResourceGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "channel", variable.DLPEndpointResourceGroupChannel),
				),
			},

			// Update test
			{
				Config: testAccCheckEndpointDLPResourceGroupConfigure(resourceTypeAndName, updatedName, variable.DLPEndpointResourceGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPResourceGroupExists(resourceTypeAndName, &group),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointResourceGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "channel", variable.DLPEndpointResourceGroupChannel),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEndpointDLPResourceGroupImportStateIDFunc(resourceTypeAndName),
			},
		},
	})
}

// testAccEndpointDLPResourceGroupImportStateIDFunc builds the "<CHANNEL>:<id>"
// import key required by the resource importer, since the read path is
// channel-scoped.
func testAccEndpointDLPResourceGroupImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("didn't find resource: %s", resourceName)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["channel"], rs.Primary.ID), nil
	}
}

func testAccCheckEndpointDLPResourceGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPEndpointResourceGroup {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		channel := endpoint_resource_channel.Channel(rs.Primary.Attributes["channel"])
		list, err := endpoint_resource_group.GetResourceGroupTagsList(context.Background(), service, channel, nil)
		if err != nil {
			continue
		}
		for i := range list {
			if list[i].ID == id {
				return fmt.Errorf("dlp endpoint resource group with id %d exists and wasn't destroyed", id)
			}
		}
	}

	return nil
}

func testAccCheckEndpointDLPResourceGroupExists(resource string, res *endpoint_resource_group.DlpEndpointResourceGroups) resource.TestCheckFunc {
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

		channel := endpoint_resource_channel.Channel(rs.Primary.Attributes["channel"])
		list, err := endpoint_resource_group.GetResourceGroupTagsList(context.Background(), service, channel, nil)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		for i := range list {
			if list[i].ID == id {
				*res = list[i]
				return nil
			}
		}
		return fmt.Errorf("dlp endpoint resource group with id %d not found in channel %q", id, channel)
	}
}

func testAccCheckEndpointDLPResourceGroupConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`
resource "%s" "%s" {
    channel     = "%s"
    name        = "%s"
    description = "%s"
}
`,
		resourcetype.DLPEndpointResourceGroup,
		resourceName,
		variable.DLPEndpointResourceGroupChannel,
		generatedName,
		description,
	)
}
