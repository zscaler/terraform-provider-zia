package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_control_settings"
)

func TestAccResourceBrowserControlPolicy_Basic(t *testing.T) {
	resourceName := "zia_browser_control_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceBrowserControlPolicyDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with minimal configuration
			{
				Config: testAccResourceBrowserControlPolicyConfig(
					"DAILY", // plugin_check_frequency
					true,    // bypass_all_browsers
					false,   // allow_all_browsers
					true,    // enable_warnings
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "plugin_check_frequency", "DAILY"),
					// Skip asserting bypass_all_browsers — API may return true as default
					resource.TestCheckResourceAttr(resourceName, "bypass_all_browsers", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_all_browsers", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_warnings", "true"),
				),
			},
			// Step 2: Update the resource with some values
			{
				Config: testAccResourceBrowserControlPolicyConfig(
					"DAILY", // plugin_check_frequency
					true,    // bypass_all_browsers
					false,   // allow_all_browsers
					true,    // enable_warnings
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "plugin_check_frequency", "DAILY"),
					resource.TestCheckResourceAttr(resourceName, "bypass_all_browsers", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_all_browsers", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_warnings", "true"),
				),
			},
			// Step 3: Import the resource and verify the state
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Skip checking specific attributes during import since the API may return defaults
				ImportStateVerifyIgnore: []string{
					"bypass_plugins",
					"bypass_applications",
					"blocked_internet_explorer_versions",
					"blocked_chrome_versions",
					"blocked_firefox_versions",
					"blocked_safari_versions",
					"blocked_opera_versions",
				},
			},
		},
	})
}

func testAccCheckResourceBrowserControlPolicyDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}

// Simplified helper function with only required parameters
func testAccResourceBrowserControlPolicyConfig(
	pluginCheckFrequency string,
	bypassAllBrowsers bool,
	allowAllBrowsers bool,
	enableWarnings bool,
) string {
	return fmt.Sprintf(`
resource "zia_browser_control_policy" "test" {
	plugin_check_frequency = %q
	bypass_all_browsers = %t
	allow_all_browsers = %t
	enable_warnings = %t
}
`,
		pluginCheckFrequency,
		bypassAllBrowsers,
		allowAllBrowsers,
		enableWarnings,
	)
}

// The Browser Control response always carries a smartIsolationProfile object.
// On a tenant with no Cloud Browser Isolation profile every member comes back
// as its zero value, and flattening that into a block used to record a profile
// the user never declared. Every later plan then tried to remove the block,
// which drove an update against the Smart Isolation endpoint and failed with
// 403 NOT_SUBSCRIBED on tenants without the subscription.
func TestFlattenSmartIsolationProfileOmitsEmptyProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile *browser_control_settings.SmartIsolationProfile
		want    int
	}{
		{
			name:    "nil profile",
			profile: nil,
			want:    0,
		},
		{
			name:    "zero-value profile as returned for tenants without a profile",
			profile: &browser_control_settings.SmartIsolationProfile{},
			want:    0,
		},
		{
			name: "zero-value profile with default_profile set",
			profile: &browser_control_settings.SmartIsolationProfile{
				DefaultProfile: true,
			},
			want: 0,
		},
		{
			name: "populated profile",
			profile: &browser_control_settings.SmartIsolationProfile{
				ID:   "d34d1b4d-0000-4000-8000-000000000001",
				Name: "Profile1",
				URL:  "https://redirect.example.com/profile1",
			},
			want: 1,
		},
		{
			name: "profile carrying only an id",
			profile: &browser_control_settings.SmartIsolationProfile{
				ID: "d34d1b4d-0000-4000-8000-000000000001",
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenSmartIsolationProfile(tt.profile)
			if len(got) != tt.want {
				t.Fatalf("expected %d block(s), got %d: %#v", tt.want, len(got), got)
			}
		})
	}
}

func TestSmartIsolationProfileConfigured(t *testing.T) {
	tests := []struct {
		name string
		raw  interface{}
		want bool
	}{
		{
			name: "nil",
			raw:  nil,
			want: false,
		},
		{
			name: "no block",
			raw:  []interface{}{},
			want: false,
		},
		{
			name: "block of empty strings, as written by earlier releases",
			raw: []interface{}{map[string]interface{}{
				"id": "", "name": "", "url": "", "default_profile": false,
			}},
			want: false,
		},
		{
			name: "block with only default_profile set",
			raw: []interface{}{map[string]interface{}{
				"id": "", "name": "", "url": "", "default_profile": true,
			}},
			want: false,
		},
		{
			name: "populated block",
			raw: []interface{}{map[string]interface{}{
				"id": "d34d1b4d", "name": "Profile1", "url": "https://example.com",
			}},
			want: true,
		},
		{
			name: "block carrying only a name",
			raw: []interface{}{map[string]interface{}{
				"id": "", "name": "Profile1", "url": "",
			}},
			want: true,
		},
		{
			name: "nil element",
			raw:  []interface{}{nil},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := smartIsolationProfileConfigured(tt.raw); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestShouldRecordSmartIsolationProfile(t *testing.T) {
	populatedBlock := []interface{}{map[string]interface{}{
		"id": "d34d1b4d", "name": "Profile1", "url": "https://example.com",
	}}
	emptyBlock := []interface{}{map[string]interface{}{
		"id": "", "name": "", "url": "", "default_profile": false,
	}}

	tests := []struct {
		name         string
		apiProfile   []interface{}
		stateProfile interface{}
		want         bool
	}{
		{
			name:         "service reported a profile",
			apiProfile:   populatedBlock,
			stateProfile: []interface{}{},
			want:         true,
		},
		{
			name:         "service reported a profile that replaces the one in state",
			apiProfile:   populatedBlock,
			stateProfile: populatedBlock,
			want:         true,
		},
		{
			name:         "no profile either side",
			apiProfile:   nil,
			stateProfile: []interface{}{},
			want:         true,
		},
		{
			name:         "clears an all-empty block left by an earlier release",
			apiProfile:   nil,
			stateProfile: emptyBlock,
			want:         true,
		},
		{
			name:         "keeps a real profile the service did not echo back",
			apiProfile:   nil,
			stateProfile: populatedBlock,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRecordSmartIsolationProfile(tt.apiProfile, tt.stateProfile); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

// A tenant that leaves Smart Isolation off must not trigger a call to the
// Smart Isolation endpoint, which answers 403 NOT_SUBSCRIBED without a Cloud
// Browser Isolation subscription.
func TestShouldUpdateSmartIsolation(t *testing.T) {
	tests := []struct {
		name string
		raw  map[string]interface{}
		want bool
	}{
		{
			name: "isolation disabled and no profile declared",
			raw: map[string]interface{}{
				"enable_smart_browser_isolation": false,
			},
			want: false,
		},
		{
			name: "isolation disabled with an all-empty profile block",
			raw: map[string]interface{}{
				"enable_smart_browser_isolation": false,
				"smart_isolation_profile": []interface{}{map[string]interface{}{
					"id": "", "name": "", "url": "",
				}},
			},
			want: false,
		},
		{
			name: "isolation enabled",
			raw: map[string]interface{}{
				"enable_smart_browser_isolation": true,
			},
			want: true,
		},
		{
			name: "profile declared while the toggle is off",
			raw: map[string]interface{}{
				"enable_smart_browser_isolation": false,
				"smart_isolation_profile": []interface{}{map[string]interface{}{
					"id": "d34d1b4d", "name": "Profile1", "url": "https://example.com",
				}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceBrowserControlPolicy().Schema, tt.raw)
			if got := shouldUpdateSmartIsolation(d); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

// An all-empty profile block must not fail validation when the toggle is off,
// otherwise state written by earlier releases could not be cleaned up.
func TestValidateSmartBrowserIsolationAllowsEmptyProfileWhenDisabled(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceBrowserControlPolicy().Schema, map[string]interface{}{
		"enable_smart_browser_isolation": false,
	})
	if err := validateSmartBrowserIsolation(d); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateSmartBrowserIsolationRequiresProfileWhenEnabled(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceBrowserControlPolicy().Schema, map[string]interface{}{
		"enable_smart_browser_isolation": true,
	})
	if err := validateSmartBrowserIsolation(d); err == nil {
		t.Fatal("expected an error when the toggle is on with no profile")
	}
}
