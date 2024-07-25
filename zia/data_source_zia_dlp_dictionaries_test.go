package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/variable"
)

func TestAccDataSourceDLPDictionaries_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPDictionaries)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPDictionariesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDLPDictionariesConfigure(resourceTypeAndName, generatedName, variable.DLPDictionaryDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "phrases.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "patterns.#", "2"),
				),
			},
		},
	})
}
