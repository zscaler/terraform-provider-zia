package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/variable"
)

func TestAccDataSourceDLPEngines_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEngines)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPEnginesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDLPEnginesConfigure(resourceTypeAndName, generatedName, generatedName, variable.DLPCustomEngine),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_dlp_engine", strconv.FormatBool(variable.DLPCustomEngine)),
				),
			},
		},
	})
}
