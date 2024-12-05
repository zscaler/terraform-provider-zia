package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceTrafficForwardingStaticIP_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingStaticIP)
	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingStaticIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingStaticIPConfigure(resourceTypeAndName, generatedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comment", resourceTypeAndName, "comment"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_address", resourceTypeAndName, "ip_address"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "routable_ip", strconv.FormatBool(variable.StaticRoutableIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "geo_override", strconv.FormatBool(variable.StaticGeoOverride)),
				),
			},
		},
	})
}
