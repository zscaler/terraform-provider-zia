package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFWApplicationServicesLite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceFWApplicationServicesLiteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.zoom_meeting"),
				),
			},
		},
	})
}

func testAccDataSourceFWApplicationServicesLiteCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceFWApplicationServicesLiteConfig_basic = `

data "zia_firewall_filtering_application_services" "zoom_meeting"{
    name = "ZOOMMEETING"
}
`
