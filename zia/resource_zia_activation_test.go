package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceActivationStatus(t *testing.T) {
	resourceName := "zia_activation_status.this"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders, // Ensure you have a provider configuration for testing
		CheckDestroy: testAccCheckActivationStatusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceActivationStatusConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckActivationStatusExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckActivationStatusDestroy(s *terraform.State) error {
	// There's nothing to check for destroy as the activation can't be deleted
	return nil
}

func testAccCheckActivationStatusExists(_ string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Implement this function to ensure the resource exists in your real infrastructure
		// You can use the Terraform SDK's helper functions to access and verify resource attributes

		return nil
	}
}

func testAccResourceActivationStatusConfig() string {
	// This mimics a .tf file configuration
	return `
resource "zia_activation_status" "this" {
	status = "ACTIVE"
}`
}
