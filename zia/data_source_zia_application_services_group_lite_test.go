package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFWApplicationServicesGroupLite_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceFWApplicationServicesGroupLiteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceFWApplicationServicesGroupLiteCheck("data.zia_firewall_filtering_application_services_group.office365"),
					testAccDataSourceFWApplicationServicesGroupLiteCheck("data.zia_firewall_filtering_application_services_group.zoom"),
					testAccDataSourceFWApplicationServicesGroupLiteCheck("data.zia_firewall_filtering_application_services_group.webex"),
					testAccDataSourceFWApplicationServicesGroupLiteCheck("data.zia_firewall_filtering_application_services_group.ring_central"),
					testAccDataSourceFWApplicationServicesGroupLiteCheck("data.zia_firewall_filtering_application_services_group.logmein"),
				),
			},
		},
	})
}

func testAccDataSourceFWApplicationServicesGroupLiteCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceFWApplicationServicesGroupLiteConfig_basic = `
data "zia_firewall_filtering_application_services_group" "office365"{
    name = "OFFICE365"
}
data "zia_firewall_filtering_application_services_group" "zoom"{
    name = "ZOOM"
}
data "zia_firewall_filtering_application_services_group" "webex"{
    name = "WEBEX"
}
data "zia_firewall_filtering_application_services_group" "ring_central"{
    name = "RINGCENTRAL"
}
data "zia_firewall_filtering_application_services_group" "logmein"{
    name = "LOGMEIN"
}
`
