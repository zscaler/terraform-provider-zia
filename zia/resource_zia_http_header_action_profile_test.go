package zia

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/http_header_control/http_header_action_profile"
)

func TestAccResourceHTTPHeaderActionProfileBasic(t *testing.T) {
	var profile http_header_action_profile.HttpHeaderActionProfile
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.HTTPHeaderActionProfile)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPHeaderActionProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHTTPHeaderActionProfileConfigure(resourceTypeAndName, initialName, variable.HTTPHeaderActionProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPHeaderActionProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.HTTPHeaderActionProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "http_header_action_profile_keys.#", "2"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
				),
			},

			// Update test
			{
				Config: testAccCheckHTTPHeaderActionProfileConfigure(resourceTypeAndName, updatedName, variable.HTTPHeaderActionProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPHeaderActionProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.HTTPHeaderActionProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "http_header_action_profile_keys.#", "2"),
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

func testAccCheckHTTPHeaderActionProfileDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.HTTPHeaderActionProfile {
			continue
		}

		profile, err := http_header_action_profile.GetByName(context.Background(), service, rs.Primary.Attributes["name"])
		if err == nil && profile != nil {
			return fmt.Errorf("HTTP header action profile with id %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckHTTPHeaderActionProfileExists(resource string, profile *http_header_action_profile.HttpHeaderActionProfile) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		received, err := http_header_action_profile.GetByName(context.Background(), service, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*profile = *received

		return nil
	}
}

func testAccCheckHTTPHeaderActionProfileConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name        = "%s"
    description = "%s"

    http_header_action_profile_keys {
        key   = "X-Forwarded-For"
        value = "10.0.0.1"
    }

    http_header_action_profile_keys {
        key   = "X-Custom-Header"
        value = "example"
    }
}

data "%s" "%s" {
    name = "${%s.%s.name}"
}
`,
		// resource variables
		resourcetype.HTTPHeaderActionProfile,
		resourceName,
		generatedName,
		description,

		// data source variables
		resourcetype.HTTPHeaderActionProfile,
		resourceName,
		// Reference to the resource
		resourcetype.HTTPHeaderActionProfile, resourceName,
	)
}
