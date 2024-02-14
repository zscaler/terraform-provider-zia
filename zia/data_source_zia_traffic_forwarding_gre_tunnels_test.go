package zia

/*
import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
)

func TestAccDataSourceTrafficForwardingGreTunnels_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingGRETunnel)

	rIP, _ := acctest.RandIpAddress("104.238.235.0/24")
	staticIPTypeAndName, _, staticIPGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingStaticIP)
	staticIPResourceHCL := testAccCheckTrafficForwardingStaticIPConfigure(staticIPTypeAndName, staticIPGeneratedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingGRETunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingGRETunnelConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, variable.GRETunnelWithinCountry, variable.GRETunnelIPUnnumbered),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "source_ip", resourceTypeAndName, "source_ip"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comment", resourceTypeAndName, "comment"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "within_country", strconv.FormatBool(variable.GRETunnelWithinCountry)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_unnumbered", strconv.FormatBool(variable.GRETunnelIPUnnumbered)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
*/
