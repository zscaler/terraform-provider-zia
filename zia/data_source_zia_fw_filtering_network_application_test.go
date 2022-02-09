package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceFWNetworkApplication_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckdataSourceFWNetworkApplicationConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_application.apns", "id"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_application.apns", "locale"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_application.dict", "id"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_application.dict", "locale"),
				),
			},
		},
	})
}

const testAccCheckdataSourceFWNetworkApplicationConfig_basic = `
data "zia_firewall_filtering_network_application" "apns" {
    id = "APNS"
	locale="en-US"
}

data "zia_firewall_filtering_network_application" "dict" {
    id = "DICT"
	locale="en-US"
}
`
