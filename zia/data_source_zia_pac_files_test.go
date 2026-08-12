package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePacFiles_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourcePacFilesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// A blank query returns the whole list, content included.
					resource.TestCheckResourceAttrSet("data.zia_pac_files.all", "pac_files.0.id"),
					resource.TestCheckResourceAttrSet("data.zia_pac_files.all", "pac_files.0.name"),
					resource.TestCheckResourceAttrSet("data.zia_pac_files.all", "pac_files.0.pac_content"),

					// The filtered list omits the PAC content but keeps the
					// rest of the fields.
					resource.TestCheckResourceAttrSet("data.zia_pac_files.no_content", "pac_files.0.id"),
					resource.TestCheckResourceAttrSet("data.zia_pac_files.no_content", "pac_files.0.name"),
					resource.TestCheckResourceAttr("data.zia_pac_files.no_content", "pac_files.0.pac_content", ""),

					// Looked up by the name resolved from the full list, which
					// must return exactly one matching file.
					resource.TestCheckResourceAttr("data.zia_pac_files.by_name", "pac_files.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.zia_pac_files.by_name", "pac_files.0.id",
						"data.zia_pac_files.all", "pac_files.0.id",
					),
				),
			},
		},
	})
}

var testAccCheckDataSourcePacFilesConfig_basic = `
data "zia_pac_files" "all" {}

data "zia_pac_files" "no_content" {
  filter = "pac_content"
}

data "zia_pac_files" "by_name" {
  name = data.zia_pac_files.all.pac_files.0.name
}
`
