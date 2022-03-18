package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccDataSourceLocationManagement_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.LocationManagement)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, variable.LocationDescription, variable.LocationCountry, variable.LocationTZ, variable.LocationAuthRequired),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "country", resourceTypeAndName, "country"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "tz", resourceTypeAndName, "tz"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "auth_required", strconv.FormatBool(variable.LocationAuthRequired)),
				),
			},
		},
	})
}
