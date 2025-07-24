package zia

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceBandwdithClassesWebConferencing_Basic(t *testing.T) {
	resourceName := "zia_bandwidth_classes_web_conferencing.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceBandwdithClassesWebConferencingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBandwdithClassesWebConferencingConfig(
					[]string{"WEBEX", "GOTOMEETING", "LIVEMEETING"},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "BANDWIDTH_CAT_WEBCONF"),
					resource.TestCheckResourceAttr(resourceName, "type", "BANDWIDTH_CAT_WEBCONF"),
					resource.TestCheckResourceAttr(resourceName, "applications.#", "3"),
				),
			},
			{
				Config: testAccResourceBandwdithClassesWebConferencingConfig(
					[]string{"INTERCALL", "CONNECT"},
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "applications.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckResourceBandwdithClassesWebConferencingDestroy(s *terraform.State) error {
	// No-op since resource updates a static object.
	return nil
}

func testAccResourceBandwdithClassesWebConferencingConfig(apps []string) string {
	var appsHCL string
	for _, app := range apps {
		appsHCL += fmt.Sprintf(`"%s",`, app)
	}
	return fmt.Sprintf(`
resource "zia_bandwidth_classes_web_conferencing" "test" {
  name         = "BANDWIDTH_CAT_WEBCONF"
  type         = "BANDWIDTH_CAT_WEBCONF"
  applications = [%s]
}
`, strings.TrimRight(appsHCL, ","))
}
