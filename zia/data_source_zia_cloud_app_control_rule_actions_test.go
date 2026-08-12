package zia

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

func TestAccDataSourceCloudAppControlRuleActions_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceCloudAppControlRuleActionsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this1", "id"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this1", "available_actions.0"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this2", "id"),
					resource.TestCheckResourceAttrSet("data.zia_cloud_app_control_rule_actions.this2", "available_actions.0"),
				),
			},
		},
	})
}

var testAccCheckDataSourceCloudAppControlRuleActionsConfig_basic = `
data "zia_cloud_app_control_rule_actions" "this1" {
  type       = "ENTERPRISE_COLLABORATION"
  cloud_apps = ["SLACK"]
}

data "zia_cloud_app_control_rule_actions" "this2" {
  type       = "FILE_SHARE"
  cloud_apps = ["GDRIVE"]
}
`

func TestActionsEndpointUnavailable(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"plain error", errors.New("connection refused"), false},
		{"405 method not allowed", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusMethodNotAllowed}}, true},
		{"404 not found", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}, true},
		{"405 via parsed status only", &errorx.ErrorResponse{Parsed: &errorx.ParsedAPIError{Status: http.StatusMethodNotAllowed}}, true},
		{"wrapped 405", fmt.Errorf("calling actions: %w", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusMethodNotAllowed}}), true},
		{"400 bad request is a real error", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusBadRequest}}, false},
		{"401 unauthorized is a real error", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusUnauthorized}}, false},
		{"403 forbidden is a real error", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusForbidden}}, false},
		{"500 server error is a real error", &errorx.ErrorResponse{Response: &http.Response{StatusCode: http.StatusInternalServerError}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := actionsEndpointUnavailable(tc.err); got != tc.want {
				t.Errorf("actionsEndpointUnavailable(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
