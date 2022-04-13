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
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringLocManagement)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocationManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLocationManagementConfigure(resourceTypeAndName, generatedName, variable.LocName, variable.LocDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "country", resourceTypeAndName, "country"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "tz", resourceTypeAndName, "tz"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "profile", resourceTypeAndName, "profile"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "idle_time_in_minutes", resourceTypeAndName, "idle_time_in_minutes"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "display_time_unit", resourceTypeAndName, "display_time_unit"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "auth_required", strconv.FormatBool(variable.LocAuthRequired)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "surrogate_ip", strconv.FormatBool(variable.LocSurrogateIP)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "xff_forward_enabled", strconv.FormatBool(variable.LocXFF)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ofw_enabled", strconv.FormatBool(variable.LocOFW)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ips_control", strconv.FormatBool(variable.LocIPS)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "vpn_credentials.#", "1"),
				),
			},
		},
	})
}
