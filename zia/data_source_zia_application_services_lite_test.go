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
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.skype_business"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.one_drive"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.exchange_online"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.m365_common"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.zoom_meeting"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.webex_meeting"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.webex_teams"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.webex_calling"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.ring_central_meeting"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.go_to_meeting"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.logmein_meeting"),
					testAccDataSourceFWApplicationServicesLiteCheck("data.zia_firewall_filtering_application_services.logmein_rescue"),
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
data "zia_firewall_filtering_application_services" "skype_business"{
    name = "SKYPEFORBUSINESS"
}
data "zia_firewall_filtering_application_services" "one_drive"{
    name = "FILE_SHAREPT_ONEDRIVE"
}
data "zia_firewall_filtering_application_services" "exchange_online"{
    name = "EXCHANGEONLINE"
}
data "zia_firewall_filtering_application_services" "m365_common"{
    name = "M365COMMON"
}
data "zia_firewall_filtering_application_services" "zoom_meeting"{
    name = "ZOOMMEETING"
}

data "zia_firewall_filtering_application_services" "webex_meeting"{
    name = "WEBEXMEETING"
}

data "zia_firewall_filtering_application_services" "webex_teams"{
    name = "WEBEXTEAMS"
}

data "zia_firewall_filtering_application_services" "webex_calling"{
    name = "WEBEXCALLING"
}

data "zia_firewall_filtering_application_services" "ring_central_meeting"{
    name = "RINGCENTRALMEETING"
}

data "zia_firewall_filtering_application_services" "go_to_meeting"{
    name = "GOTOMEETING"
}

data "zia_firewall_filtering_application_services" "goto_meeting_inroom"{
    name = "GOTOMEETING_INROOM"
}

data "zia_firewall_filtering_application_services" "logmein_meeting"{
    name = "LOGMEINMEETING"
}

data "zia_firewall_filtering_application_services" "logmein_rescue"{
    name = "LOGMEINRESCUE"
}
`
