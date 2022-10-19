package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceActivationStatus_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceActivationStatusConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceActivationStatusCheck("data.zia_activation_status.status"),
				),
			},
		},
	})
}

func testAccDataSourceActivationStatusCheck(status string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(status, "status"),
	)
}

var testAccCheckDataSourceActivationStatusConfig_basic = `
data "zia_activation_status" "status" {
}
`
