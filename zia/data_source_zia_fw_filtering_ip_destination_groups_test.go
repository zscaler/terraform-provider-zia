package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceFWIPDestinationGroups_ByIdAndName(t *testing.T) {
	rName := acctest.RandString(15)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_firewall_filtering_destination_groups.by_id"
	resourceName2 := "data.zia_firewall_filtering_destination_groups.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFWIPDestinationGroupsByID(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWIPDestinationGroups(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName2, "name", rName),
					resource.TestCheckResourceAttr(resourceName2, "description", rDesc),
				),
				PreventPostDestroyRefresh: true,
			},
		},
	})
}

func testAccDataSourceFWIPDestinationGroupsByID(rName, rDesc string) string {
	return fmt.Sprintf(`
	resource "zia_firewall_filtering_destination_groups" "example" {
		name        = "%s"
		description = "%s"
		type        = "DSTN_FQDN"
		addresses = [ "test1.acme.com" ]
	  }
	data "zia_firewall_filtering_destination_groups" "by_name" {
		name = zia_firewall_filtering_destination_groups.example.name
	}
	data "zia_firewall_filtering_destination_groups" "by_id" {
		id = zia_firewall_filtering_destination_groups.example.id
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
