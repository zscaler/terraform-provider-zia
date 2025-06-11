package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceForwardingControlProxies_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ForwardingControlProxies)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingControlProxiesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingControlProxiesConfigure(resourceTypeAndName, generatedName, variable.ProxyDescription, variable.ProxyType, variable.ProxyAddress, variable.ProxyPort, variable.ProxyInsertXauHeader, variable.ProxyBase64EncodeXauHeader),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "address", resourceTypeAndName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "port", resourceTypeAndName, "port"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "insert_xau_header", strconv.FormatBool(variable.ProxyInsertXauHeader)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "base64_encode_xau_header", strconv.FormatBool(variable.ProxyBase64EncodeXauHeader)),
				),
			},
		},
	})
}
