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
					resource.TestCheckResourceAttr(resourceName, "enable_dns_resolution_on_transparent_proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_office365", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_internal_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "enforce_surrogate_ip_for_windows_app", "true"),
					resource.TestCheckResourceAttr(resourceName, "track_http_tunnel_on_http_ports", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_http_tunnel_on_non_http_ports", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_domain_fronting_on_host_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "cascade_url_filtering", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_non_compliant_http_request_on_http_ports", "true"),
					resource.TestCheckResourceAttr(resourceName, "ui_session_timeout", "36000"),
				),
			},
		},
	})
}

var testAccCheckDataSourceAdvancedSettingsConfig_basic = `
data "zia_advanced_settings" "this" {}
`
