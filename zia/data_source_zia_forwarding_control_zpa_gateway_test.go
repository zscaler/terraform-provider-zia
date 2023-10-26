package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
)

func TestAccDataForwardingControlZPAGateway_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ForwardingControlZPAGateway)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingControlZPAGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingControlZPAGatewayConfigure(resourceTypeAndName, generatedName, variable.FowardingControlZPAGWDescription, variable.FowardingControlType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "zpa_server_group.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "zpa_application_segments.#", "1"),
				),
			},
		},
	})
}
