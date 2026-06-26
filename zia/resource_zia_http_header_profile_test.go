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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/http_header_control/http_header_profile"
)

func TestAccResourceHTTPHeaderProfileBasic(t *testing.T) {
	var profile http_header_profile.HttpHeaderProfile
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.HTTPHeaderProfile)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPHeaderProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckHTTPHeaderProfileConfigure(resourceTypeAndName, initialName, variable.HTTPHeaderProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPHeaderProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.HTTPHeaderProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "http_header_profile_criteria.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
				),
			},

			// Update test
			{
				Config: testAccCheckHTTPHeaderProfileConfigure(resourceTypeAndName, updatedName, variable.HTTPHeaderProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHTTPHeaderProfileExists(resourceTypeAndName, &profile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.HTTPHeaderProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "http_header_profile_criteria.#", "1"),
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

func testAccCheckHTTPHeaderProfileDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.HTTPHeaderProfile {
			continue
		}

		profile, err := http_header_profile.GetByName(context.Background(), service, rs.Primary.Attributes["name"])
		if err == nil && profile != nil {
			return fmt.Errorf("HTTP header profile with id %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckHTTPHeaderProfileExists(resource string, profile *http_header_profile.HttpHeaderProfile) resource.TestCheckFunc {
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

		received, err := http_header_profile.GetByName(context.Background(), service, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*profile = *received

		return nil
	}
}

func testAccCheckHTTPHeaderProfileConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name        = "%s"
    description = "%s"

    http_header_profile_criteria {
        header           = "ORIGIN"
        category_bitmap   = ["GENERAL_AI_ML", "AI_ML_APPS"]
        cloud_app_bitmap  = ["CHATGPT_AI"]
    }
}

data "%s" "%s" {
    name = "${%s.%s.name}"
}
`,
		// resource variables
		resourcetype.HTTPHeaderProfile,
		resourceName,
		generatedName,
		description,

		// data source variables
		resourcetype.HTTPHeaderProfile,
		resourceName,
		// Reference to the resource
		resourcetype.HTTPHeaderProfile, resourceName,
	)
}
