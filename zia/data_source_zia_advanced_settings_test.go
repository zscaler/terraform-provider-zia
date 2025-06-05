package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceAdvancedSettings_Basic(t *testing.T) {
	resourceName := "data.zia_advanced_settings.this"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceAdvancedSettingsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "enable_dns_resolution_on_transparent_proxy"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_office365"),
					resource.TestCheckResourceAttrSet(resourceName, "log_internal_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "enforce_surrogate_ip_for_windows_app"),
					resource.TestCheckResourceAttrSet(resourceName, "track_http_tunnel_on_http_ports"),
					resource.TestCheckResourceAttrSet(resourceName, "block_http_tunnel_on_non_http_ports"),
					resource.TestCheckResourceAttrSet(resourceName, "block_domain_fronting_on_host_header"),
					resource.TestCheckResourceAttrSet(resourceName, "cascade_url_filtering"),
					resource.TestCheckResourceAttrSet(resourceName, "block_non_compliant_http_request_on_http_ports"),
					resource.TestCheckResourceAttr(resourceName, "ui_session_timeout", "36000"),
				),
			},
		},
	})
}

var testAccCheckDataSourceAdvancedSettingsConfig_basic = `
data "zia_advanced_settings" "this" {}
`
