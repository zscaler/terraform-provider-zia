package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var idmProfileLiteNames = []string{
	"BD_IDM_TEMPLATE01",
}

func TestAccDataSourceDLPIDMProfileLite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDLPIDMProfileLite_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateDLPIDMProfileLiteChecks()...,
				),
			},
		},
	})
}

func generateDLPIDMProfileLiteChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range idmProfileLiteNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zia_dlp_idm_profile_lite.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "template_name"),
		))
	}
	return checks
}

func testAccCheckDataSourceDLPIDMProfileLite_basic() string {
	var configs string
	for _, name := range idmProfileLiteNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zia_dlp_idm_profile_lite" "%s" {
    template_name = "%s"
}
`, resourceName, name)
	}
	return configs
}
