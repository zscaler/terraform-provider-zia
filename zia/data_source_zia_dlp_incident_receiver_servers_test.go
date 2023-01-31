package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDLPIncidentReceiverServers_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDLPIncidentReceiverServersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDLPIncidentReceiverServersCheck("data.zia_dlp_incident_receiver_servers.receiver_01"),
					testAccDataSourceDLPIncidentReceiverServersCheck("data.zia_dlp_incident_receiver_servers.receiver_02"),
				),
			},
		},
	})
}

func testAccDataSourceDLPIncidentReceiverServersCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceDLPIncidentReceiverServersConfig_basic = `
data "zia_dlp_incident_receiver_servers" "receiver_01"{
    name = "ZS_BD_INC_RECEIVER_01"
}

data "zia_dlp_incident_receiver_servers" "receiver_02"{
    name = "ZS_BD_INC_RECEIVER_02"
}
`
