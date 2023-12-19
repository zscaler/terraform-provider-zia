package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var profileNames = []string{
	"BD_SA_Profile1_ZIA", "BD_SA_Profile2_ZIA", "BD SA Profile ZIA", "BD  SA Profile ZIA", "BD   SA   Profile  ZIA",
}

func TestAccDataSourceCBIProfile_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCBIProfile_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateCBIProfileChecks()...,
				),
			},
		},
	})
}

func generateCBIProfileChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range profileNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zia_cloud_browser_isolation_profile.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "name"),
		))
	}
	return checks
}

func testAccCheckDataSourceCBIProfile_basic() string {
	var configs string
	for _, name := range profileNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zia_cloud_browser_isolation_profile" "%s" {
    name = "%s"
}
`, resourceName, name)
	}
	return configs
}
