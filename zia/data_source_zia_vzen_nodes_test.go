package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceVZENNode_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ServiceEdgeNode)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVZENNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVZENNodeConfigure(resourceTypeAndName, "tf-acc-test-"+generatedName, variable.VzenStatus, variable.VzenNodeType, variable.VzenNodeIPAddress, variable.VzenNodeSubnetMask, variable.VzenNodeDefaultGateway, variable.VzenNodeInProduction, variable.VzenNodeDeploymentMode, variable.VZenSKUType, variable.VzenOnDemandSupportTunnel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "status", resourceTypeAndName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_address", resourceTypeAndName, "ip_address"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "subnet_mask", resourceTypeAndName, "subnet_mask"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "default_gateway", resourceTypeAndName, "default_gateway"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "deployment_mode", resourceTypeAndName, "deployment_mode"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "vzen_sku_type", resourceTypeAndName, "vzen_sku_type"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "in_production", strconv.FormatBool(variable.VzenNodeInProduction)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "on_demand_support_tunnel_enabled", strconv.FormatBool(variable.VzenOnDemandSupportTunnel)),
				),
			},
		},
	})
}
