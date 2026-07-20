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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_custom_apps"
)

func TestAccResourceEndpointDLPCustomAppsBasic(t *testing.T) {
	var app endpoint_custom_apps.EndpointApplications
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEndpointCustomApps)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointDLPCustomAppsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEndpointDLPCustomAppsConfigure(resourceTypeAndName, initialName, variable.DLPEndpointCustomAppDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPCustomAppsExists(resourceTypeAndName, &app),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointCustomAppDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "application.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "application.0.os_type", variable.DLPEndpointCustomAppOsType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "application.0.file_name", variable.DLPEndpointCustomAppFileName),
				),
			},
			// Update test
			{
				Config: testAccCheckEndpointDLPCustomAppsConfigure(resourceTypeAndName, updatedName, variable.DLPEndpointCustomAppDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointDLPCustomAppsExists(resourceTypeAndName, &app),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPEndpointCustomAppDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "application.0.os_type", variable.DLPEndpointCustomAppOsType),
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

func testAccCheckEndpointDLPCustomAppsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPEndpointCustomApps {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		app, err := findCustomAppByID(context.Background(), service, id)
		if err == nil && app != nil {
			return fmt.Errorf("endpoint dlp custom app with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckEndpointDLPCustomAppsExists(resource string, res *endpoint_custom_apps.EndpointApplications) resource.TestCheckFunc {
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

		received, err := findCustomAppByID(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		if received == nil {
			return fmt.Errorf("endpoint dlp custom app with id %d not found", id)
		}
		*res = *received

		return nil
	}
}

func testAccCheckEndpointDLPCustomAppsConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`
resource "%s" "%s" {
    name        = "%s"
    description = "%s"
    channel     = "APPLICATION_FILE_ACCESS"

    application {
        os_type            = "%s"
        file_name          = "%s"
        original_file_name = "%s"
    }
}
`,
		resourcetype.DLPEndpointCustomApps,
		resourceName,
		generatedName,
		description,
		variable.DLPEndpointCustomAppOsType,
		variable.DLPEndpointCustomAppFileName,
		variable.DLPEndpointCustomAppFileName,
	)
}
