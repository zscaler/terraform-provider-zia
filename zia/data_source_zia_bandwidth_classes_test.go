package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
)

func TestAccDataSourceBandwdithClasses_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.BandwdithClasses)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBandwdithClassesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBandwdithClassesConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "web_applications.#", "3"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "urls.#", "3"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "url_categories.#", "2"),
				),
			},
		},
	})
}
