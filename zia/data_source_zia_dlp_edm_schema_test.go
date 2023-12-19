package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var edmProjectNames = []string{
	"BD_EDM_TEMPLATE01", "BD_EDM_TEMPLATE02",
}

func TestAccDataSourceDLPEDMSchema_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDLPEDMSchema_basic(),
				Check: resource.ComposeTestCheckFunc(
					generateDLPEDMSchemaChecks()...,
				),
			},
		},
	})
}

func generateDLPEDMSchemaChecks() []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for _, name := range edmProjectNames {
		resourceName := createValidResourceName(name)
		checkName := fmt.Sprintf("data.zia_dlp_edm_schema.%s", resourceName)
		checks = append(checks, resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(checkName, "id"),
			resource.TestCheckResourceAttrSet(checkName, "project_name"),
		))
	}
	return checks
}

func testAccCheckDataSourceDLPEDMSchema_basic() string {
	var configs string
	for _, name := range edmProjectNames {
		resourceName := createValidResourceName(name)
		configs += fmt.Sprintf(`
data "zia_dlp_edm_schema" "%s" {
    project_name = "%s"
}
`, resourceName, name)
	}
	return configs
}
