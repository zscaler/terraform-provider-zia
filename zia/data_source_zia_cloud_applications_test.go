package zia

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCloudApplications_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCloudApplicationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zia_cloud_applications.policy", "applications.#"),
					resource.TestCheckResourceAttr("data.zia_cloud_applications.policy", "policy_type", "cloud_application_policy"),
					testAccDataSourceCloudApplicationsCheck("data.zia_cloud_applications.policy", "WEB_MAIL"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_applications.ssl_policy", "applications.#"),
					resource.TestCheckResourceAttr("data.zia_cloud_applications.ssl_policy", "policy_type", "cloud_application_ssl_policy"),
					testAccDataSourceCloudApplicationsCheck("data.zia_cloud_applications.ssl_policy", "SOCIAL_NETWORKING"),
				),
			},
		},
	})
}

func testAccDataSourceCloudApplicationsCheck(resourceName, expectedAppCategory string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appsCount, ok := rs.Primary.Attributes["applications.#"]
		if !ok || appsCount == "0" {
			return fmt.Errorf("%s: No applications found", resourceName)
		}

		// Convert appsCount string to integer
		count, err := strconv.Atoi(appsCount)
		if err != nil {
			return err
		}

		// Check for a specific application category in the list of applications
		found := false
		for i := 0; i < count; i++ {
			parentAttr := fmt.Sprintf("applications.%d.parent", i)
			parent, ok := rs.Primary.Attributes[parentAttr]
			if ok && parent == expectedAppCategory {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("%s: Expected application category '%s' not found", resourceName, expectedAppCategory)
		}

		return nil
	}
}

func testAccCheckDataSourceCloudApplicationsConfig() string {
	return `
data "zia_cloud_applications" "policy" {
    policy_type = "cloud_application_policy"
    app_class   = ["WEB_MAIL"]
}

data "zia_cloud_applications" "ssl_policy" {
    policy_type = "cloud_application_ssl_policy"
    app_class   = ["SOCIAL_NETWORKING"]
}
`
}
