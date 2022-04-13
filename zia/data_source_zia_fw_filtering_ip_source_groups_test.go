package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWIPSourceGroups_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_ip_source_groups.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWIPSourceGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWIPSourceGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-src-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-src-group-"+rDesc),
				),
			},
		},
	})
}

func testAccDataSourceFWIPSourceGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_ip_source_groups" "test"{
	name = "test-fw-src-group-%s"
	description = "test-fw-src-group-%s"
	ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}

data "zia_firewall_filtering_ip_source_groups" "test" {
	name = zia_firewall_filtering_ip_source_groups.test.name
}
	`, rName, rDesc)
}

func testAccDataSourceFWIPSourceGroups(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
