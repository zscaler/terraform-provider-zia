package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceDCExclusions_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DCExclusions)
	description := variable.DCExclusionsDescription + "-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDCExclusionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDCExclusionsConfigure(resourceTypeAndName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "exclusions.0.id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceTypeAndName, "exclusions.0.dc_id"),
					resource.TestCheckResourceAttrSet(dataSourceTypeAndName, "exclusions.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSourceTypeAndName, "exclusions.0.end_time"),
				),
			},
		},
	})
}
