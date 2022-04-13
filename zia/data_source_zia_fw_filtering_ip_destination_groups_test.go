package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWIPDestinationGroups_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_destination_groups.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWIPDestinationGroupsBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWIPDestinationGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-fw-dst-group-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-fw-dst-group-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "type", "DSTN_FQDN"),
				),
			},
		},
	})
}

func testAccDataSourceFWIPDestinationGroupsBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_firewall_filtering_destination_groups" "test" {
	name        = "test-fw-dst-group-%s"
	description = "test-fw-dst-group-%s"
	type        = "DSTN_FQDN"
	addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
  }

data "zia_firewall_filtering_destination_groups" "test" {
	name = zia_firewall_filtering_destination_groups.test.name
}
	`, rName, rDesc)
}

func testAccDataSourceFWIPDestinationGroups(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
