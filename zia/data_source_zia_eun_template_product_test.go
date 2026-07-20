package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEUNTemplateProduct_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceEUNTemplateProductConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					// Resolved to the default template for the policy type when
					// neither id nor name is provided.
					resource.TestCheckResourceAttrSet("data.zia_eun_template_product.zcc_endpoint_dlp", "id"),
					resource.TestCheckResourceAttrSet("data.zia_eun_template_product.zcc_endpoint_dlp", "name"),
					resource.TestCheckResourceAttrSet("data.zia_eun_template_product.zcc_endpoint_dlp", "type"),
					resource.TestCheckResourceAttr("data.zia_eun_template_product.zcc_endpoint_dlp", "product", "ENDPOINT_DLP"),

					resource.TestCheckResourceAttrSet("data.zia_eun_template_product.browser_url", "id"),
					resource.TestCheckResourceAttrSet("data.zia_eun_template_product.browser_url", "name"),
					resource.TestCheckResourceAttr("data.zia_eun_template_product.browser_url", "product", "URL"),
				),
			},
		},
	})
}

var testAccCheckDataSourceEUNTemplateProductConfig_basic = `
data "zia_eun_template_product" "zcc_endpoint_dlp" {
  template_type = "ZCC"
  product       = "ENDPOINT_DLP"
}

data "zia_eun_template_product" "browser_url" {
  template_type = "BROWSER"
  product       = "URL"
}
`
