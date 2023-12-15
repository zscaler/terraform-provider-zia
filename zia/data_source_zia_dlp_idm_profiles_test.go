package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var idmProfileNames = []string{
	"BD_IDM_TEMPLATE01",
}

func TestAccDataSourceDLPIDMProfiles_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDLPIDMProfiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateDLPIDMProfilesChecks()...,
				),
			},
		},
	})
}

func generateDLPIDMProfilesChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range idmProfileNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zia_dlp_idm_profiles.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "profile_name"),
		))
	}
	return checks
}

func testAccCheckDataSourceDLPIDMProfiles_basic() string {
	var configs string
	for _, name := range idmProfileNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zia_dlp_idm_profiles" "%s" {
    profile_name = "%s"
}
`, resourceName, name)
	}
	return configs
}
