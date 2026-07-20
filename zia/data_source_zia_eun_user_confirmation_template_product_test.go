package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEUNUserConfirmationTemplateProduct_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceEUNUserConfirmationTemplateProductConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// Resolved to the default template for the policy type when
					// neither id nor name is provided.
					resource.TestCheckResourceAttrSet("data.zia_eun_user_confirmation_template_product.endpoint_dlp", "id"),
					resource.TestCheckResourceAttrSet("data.zia_eun_user_confirmation_template_product.endpoint_dlp", "name"),
					resource.TestCheckResourceAttr("data.zia_eun_user_confirmation_template_product.endpoint_dlp", "product", "ENDPOINT_DLP"),

					resource.TestCheckResourceAttrSet("data.zia_eun_user_confirmation_template_product.inline", "id"),
					resource.TestCheckResourceAttr("data.zia_eun_user_confirmation_template_product.inline", "product", "INLINE"),
				),
			},
		},
	})
}

var testAccCheckDataSourceEUNUserConfirmationTemplateProductConfig_basic = `
data "zia_eun_user_confirmation_template_product" "endpoint_dlp" {
  product = "ENDPOINT_DLP"
}

data "zia_eun_user_confirmation_template_product" "inline" {
  product = "INLINE"
}
`
