package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccDataSourceFWNetworkServiceGroups_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringNetworkServiceGroups)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWNetworkServiceGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFWNetworkServiceGroupsConfigure(resourceTypeAndName, generatedName, variable.FWNetworkServicesGroupDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "services.#", "2"),
				),
			},
		},
	})
}
