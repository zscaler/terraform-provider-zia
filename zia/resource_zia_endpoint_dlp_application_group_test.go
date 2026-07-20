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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_application_groups"
)

func TestAccResourceEndpointDLPApplicationGroupBasic(t *testing.T) {
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEndpointApplicationGroup)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDLPApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEndpointDLPApplicationGroupConfigure(resourceTypeAndName, initialName, variable.DLPEndpointApplicationGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPApplicationGroupExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointApplicationGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "channel", "APPLICATION_FILE_ACCESS"),
				),
			},

			// Update test
			{
				Config: testAccCheckEndpointDLPApplicationGroupConfigure(resourceTypeAndName, updatedName, variable.DLPEndpointApplicationGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPApplicationGroupExists(resourceTypeAndName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointApplicationGroupDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "channel", "APPLICATION_FILE_ACCESS"),
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

func testAccCheckEndpointDLPApplicationGroupDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPEndpointApplicationGroup {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		list, err := endpoint_application_groups.GetAll(context.Background(), service)
		if err != nil {
			continue
		}
		for i := range list {
			if list[i].GroupID == id {
				return fmt.Errorf("endpoint application group with id %d exists and wasn't destroyed", id)
			}
		}
	}

	return nil
}

func testAccCheckEndpointDLPApplicationGroupExists(resource string) resource.TestCheckFunc {
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

		list, err := endpoint_application_groups.GetAll(context.Background(), service)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		for i := range list {
			if list[i].GroupID == id {
				return nil
			}
		}
		return fmt.Errorf("endpoint application group with id %d not found", id)
	}
}

func testAccCheckEndpointDLPApplicationGroupConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`
resource "%s" "%s" {
    name        = "%s"
    description = "%s"
}
`,
		resourcetype.DLPEndpointApplicationGroup,
		resourceName,
		generatedName,
		description,
	)
}
