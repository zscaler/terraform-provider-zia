package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFWNetworkServices_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceFWNetworkServicesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_service.http", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_service.ftp", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_firewall_filtering_network_service.icmp", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceFWNetworkServicesConfig_basic = `
data "zia_firewall_filtering_network_service" "http"{
    name = "HTTP"
}
data "zia_firewall_filtering_network_service" "ftp"{
    name = "FTP"
}
data "zia_firewall_filtering_network_service" "icmp"{
    name = "ICMP_ANY"
}
`
