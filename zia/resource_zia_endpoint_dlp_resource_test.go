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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
)

func TestAccResourceEndpointDLPResourceBasic(t *testing.T) {
	var resource0 endpoint_resource.EndpointResource
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEndpointResource)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDLPResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEndpointDLPResourceConfigure(resourceTypeAndName, initialName, variable.DLPEndpointResourceDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPResourceExists(resourceTypeAndName, &resource0),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "channel", variable.DLPEndpointResourceChannel),
					resource.TestCheckResourceAttr(resourceTypeAndName, "printer.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckEndpointDLPResourceConfigure(resourceTypeAndName, updatedName, variable.DLPEndpointResourceDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPResourceExists(resourceTypeAndName, &resource0),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointResourceDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "channel", variable.DLPEndpointResourceChannel),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccEndpointDLPResourceImportStateIDFunc(resourceTypeAndName),
			},
		},
	})
}

// testAccEndpointDLPResourceImportStateIDFunc builds the "<CHANNEL>:<id>" import
// key required by the resource importer, since the read path is channel-scoped.
func testAccEndpointDLPResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("didn't find resource: %s", resourceName)
		}
		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["channel"], rs.Primary.ID), nil
	}
}

func testAccCheckEndpointDLPResourceDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPEndpointResource {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		channel := endpoint_resource_channel.Channel(rs.Primary.Attributes["channel"])
		resource0, err := endpoint_resource_channel.GetChannel(context.Background(), service, channel, id)
		if err == nil && resource0 != nil {
			return fmt.Errorf("dlp endpoint resource with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckEndpointDLPResourceExists(resource string, res *endpoint_resource.EndpointResource) resource.TestCheckFunc {
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
		received, err := endpoint_resource_channel.GetChannel(context.Background(), service, channel, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*res = *received

		return nil
	}
}

func testAccCheckEndpointDLPResourceConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`
resource "%s" "%s" {
    channel     = "%s"
    name        = "%s"
    description = "%s"

    printer {
        ip_address = "10.10.10.20"
        domain     = "acme.local"
    }
}
`,
		resourcetype.DLPEndpointResource,
		resourceName,
		variable.DLPEndpointResourceChannel,
		generatedName,
		description,
	)
}
