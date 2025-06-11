package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceVZENCluster_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ServiceEdgeCluster)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVZENClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVZENClusterConfigure(
					resourceTypeAndName,
					generatedName,
					variable.VzenStatus,
					variable.VzenType,
					variable.VzenIPAddress,
					variable.VzenSubnetMask,
					variable.VzenDefaultGateway,
					variable.VzenIpSecEnabled,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "status", resourceTypeAndName, "status"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_address", resourceTypeAndName, "ip_address"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "subnet_mask", resourceTypeAndName, "subnet_mask"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "default_gateway", resourceTypeAndName, "default_gateway"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_sec_enabled", resourceTypeAndName, "ip_sec_enabled"),
				),
			},
		},
	})
}
