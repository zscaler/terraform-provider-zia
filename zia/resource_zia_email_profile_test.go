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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/email_profiles"
)

func TestAccResourceEmailProfile(t *testing.T) {
	var emailProfile email_profiles.EmailProfiles
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.EmailProfile)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEmailProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckEmailProfileConfigure(resourceTypeAndName, initialName, variable.EmailProfileDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailProfileExists(resourceTypeAndName, &emailProfile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.EmailProfileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "emails.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckEmailProfileConfigure(resourceTypeAndName, updatedName, variable.EmailProfileDescriptionUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmailProfileExists(resourceTypeAndName, &emailProfile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.EmailProfileDescriptionUpdate),
					resource.TestCheckResourceAttr(resourceTypeAndName, "emails.#", "2"),
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

func testAccCheckEmailProfileDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.EmailProfile {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := email_profiles.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("email profile with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckEmailProfileExists(resource string, profile *email_profiles.EmailProfiles) resource.TestCheckFunc {
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

		receivedProfile, err := email_profiles.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*profile = *receivedProfile

		return nil
	}
}

func testAccCheckEmailProfileConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	return fmt.Sprintf(`

resource "%s" "%s" {
    name        = "%s"
    description = "%s"
    emails      = ["user1@example.com", "user2@example.com"]
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.EmailProfile,
		resourceName,
		generatedName,
		description,

		// data source variables
		resourcetype.EmailProfile,
		resourceName,
		// Reference to the resource
		resourcetype.EmailProfile, resourceName,
	)
}
