package zia

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_control_settings"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/secure_browsing"
)

func resourceBrowserControlPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceBrowserControlPolicyRead,
		CreateContext: resourceBrowserControlPolicyCreate,
		UpdateContext: resourceBrowserControlPolicyUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceBrowserControlPolicyRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to read browser control policy import: %s", diags[0].Summary)
				}
				d.SetId("browser_settings")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"plugin_check_frequency": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies how frequently the service checks browsers and relevant applications to warn users regarding outdated or vulnerable browsers, plugins, and applications. If not set, the warnings are disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"DAILY",
					"WEEKLY",
					"MONTHLY",
					"EVERY_2_HOURS",
					"EVERY_4_HOURS",
					"EVERY_6_HOURS",
					"EVERY_8_HOURS",
					"EVERY_12_HOURS",
				}, false),
			},
			"bypass_plugins": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of plugins that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable plugins are warned",
			},
			"bypass_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of applications that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable applications are warned",
			},
			"blocked_internet_explorer_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Microsoft browser that need to be blocked. If not set, all Microsoft browser versions are allowed.",
			},
			"blocked_chrome_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Google Chrome browser that need to be blocked. If not set, all Google Chrome versions are allowed.",
			},
			"blocked_firefox_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Mozilla Firefox browser that need to be blocked. If not set, all Mozilla Firefox versions are allowed.",
			},
			"blocked_safari_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Apple Safari browser that need to be blocked. If not set, all Apple Safari versions are allowed",
			},
			"blocked_opera_versions": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Opera browser that need to be blocked. If not set, all Opera versions are allowed",
			},
			"bypass_all_browsers": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed:    true,
				Description: "If set to true, all the browsers are bypassed for warnings",
			},
			"allow_all_browsers": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed:    true,
				Description: "A Boolean value that specifies whether or not to allow all the browsers and their respective versions access to the internet",
			},
			"enable_warnings": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed:    true,
				Description: "A Boolean value that specifies if the warnings are enabled",
			},
			"enable_smart_browser_isolation": {
				Type:     schema.TypeBool,
				Optional: true,
				// Computed:    true,
				Description: "A Boolean value that specifies if Smart Browser Isolation is enabled",
			},
			"smart_isolation_profile": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The Cloud Browser Isolation profile applied when Smart Browser Isolation is enabled. Required when `enable_smart_browser_isolation` is `true`. All three of `id`, `name`, and `url` must be populated — typically by referencing the `zia_cloud_browser_isolation_profile` data source.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The universally unique identifier (UUID) for the browser isolation profile.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the browser isolation profile.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The browser isolation profile URL.",
						},
						"default_profile": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this is the tenant's default browser isolation profile. The service sets this field.",
						},
					},
				},
			},
			"smart_isolation_groups": setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of groups for which the rule is applied"),
			"smart_isolation_users":  setIDsSchemaTypeCustom(intPtr(4), "Name-ID pairs of users for which the rule is applied"),
		},
	}
}

func resourceBrowserControlPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if err := validateSmartBrowserIsolation(d); err != nil {
		return diag.FromErr(err)
	}

	if _, _, err := browser_control_settings.UpdateBrowserControlSettings(ctx, service, expandBrowserControlBaseSettings(d)); err != nil {
		return diag.FromErr(err)
	}

	// The base endpoint above does not actually persist the Smart Isolation
	// fields (enable_smart_browser_isolation, smart_isolation_profile,
	// smart_isolation_users, smart_isolation_groups). They must be sent to a
	// separate endpoint, which we only call when the user has actually
	// configured Smart Isolation in their HCL.
	if shouldUpdateSmartIsolation(d) {
		if _, _, err := secure_browsing.UpdateSmartIsolation(ctx, service, expandSmartIsolation(d)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("browser_settings")

	time.Sleep(1 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceBrowserControlPolicyRead(ctx, d, meta)
}

func resourceBrowserControlPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := browser_control_settings.GetBrowserControlSettings(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("browser_settings")
		_ = d.Set("plugin_check_frequency", resp.PluginCheckFrequency)
		_ = d.Set("bypass_plugins", preserveAnyBypassValue(d, "bypass_plugins", resp.BypassPlugins))
		_ = d.Set("bypass_applications", preserveAnyBypassValue(d, "bypass_applications", resp.BypassApplications))
		_ = d.Set("blocked_internet_explorer_versions", preserveAnyBypassValue(d, "blocked_internet_explorer_versions", resp.BlockedInternetExplorerVersions))
		_ = d.Set("blocked_chrome_versions", preserveAnyBypassValue(d, "blocked_chrome_versions", resp.BlockedChromeVersions))
		_ = d.Set("blocked_firefox_versions", preserveAnyBypassValue(d, "blocked_firefox_versions", resp.BlockedFirefoxVersions))
		_ = d.Set("blocked_safari_versions", preserveAnyBypassValue(d, "blocked_safari_versions", resp.BlockedSafariVersions))
		_ = d.Set("blocked_opera_versions", preserveAnyBypassValue(d, "blocked_opera_versions", resp.BlockedOperaVersions))
		_ = d.Set("bypass_all_browsers", resp.BypassAllBrowsers)
		_ = d.Set("allow_all_browsers", resp.AllowAllBrowsers)
		_ = d.Set("enable_warnings", resp.EnableWarnings)
		_ = d.Set("enable_smart_browser_isolation", resp.EnableSmartBrowserIsolation)

		// Smart Isolation users and groups are persisted by a separate
		// endpoint and the base GET does not echo them back. Only overwrite
		// state when the response actually carries values, otherwise preserve
		// whatever the user previously declared to avoid silent drift.
		if len(resp.SmartIsolationGroups) > 0 {
			if err := d.Set("smart_isolation_groups", flattenIDExtensionsListIDs(resp.SmartIsolationGroups)); err != nil {
				return diag.FromErr(err)
			}
		}
		if len(resp.SmartIsolationUsers) > 0 {
			if err := d.Set("smart_isolation_users", flattenIDExtensionsListIDs(resp.SmartIsolationUsers)); err != nil {
				return diag.FromErr(err)
			}
		}

		if err := d.Set("smart_isolation_profile", flattenSmartIsolationProfile(&resp.SmartIsolationProfile)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't read browser control policy"))
	}

	return nil
}

func resourceBrowserControlPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	if err := validateSmartBrowserIsolation(d); err != nil {
		return diag.FromErr(err)
	}

	if _, _, err := browser_control_settings.UpdateBrowserControlSettings(ctx, service, expandBrowserControlBaseSettings(d)); err != nil {
		return diag.FromErr(err)
	}

	if shouldUpdateSmartIsolation(d) {
		if _, _, err := secure_browsing.UpdateSmartIsolation(ctx, service, expandSmartIsolation(d)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("browser_settings")

	time.Sleep(1 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceBrowserControlPolicyRead(ctx, d, meta)
}

// expandBrowserControlBaseSettings builds the payload for the base
// /browserControlSettings PUT endpoint. The Smart Isolation fields are
// intentionally omitted — the base endpoint does not persist them.
func expandBrowserControlBaseSettings(d *schema.ResourceData) browser_control_settings.BrowserControlSettings {
	return browser_control_settings.BrowserControlSettings{
		PluginCheckFrequency:            d.Get("plugin_check_frequency").(string),
		BypassAllBrowsers:               d.Get("bypass_all_browsers").(bool),
		AllowAllBrowsers:                d.Get("allow_all_browsers").(bool),
		EnableWarnings:                  d.Get("enable_warnings").(bool),
		BypassPlugins:                   SetToStringList(d, "bypass_plugins"),
		BypassApplications:              SetToStringList(d, "bypass_applications"),
		BlockedInternetExplorerVersions: SetToStringList(d, "blocked_internet_explorer_versions"),
		BlockedChromeVersions:           SetToStringList(d, "blocked_chrome_versions"),
		BlockedFirefoxVersions:          SetToStringList(d, "blocked_firefox_versions"),
		BlockedSafariVersions:           SetToStringList(d, "blocked_safari_versions"),
		BlockedOperaVersions:            SetToStringList(d, "blocked_opera_versions"),
	}
}

// expandSmartIsolation builds the payload for the Smart Isolation endpoint.
// It carries the same base settings as the browser control payload, plus the
// Smart Isolation-specific fields the base endpoint cannot persist.
func expandSmartIsolation(d *schema.ResourceData) secure_browsing.SmartIsolation {
	return secure_browsing.SmartIsolation{
		PluginCheckFrequency:            d.Get("plugin_check_frequency").(string),
		BypassAllBrowsers:               d.Get("bypass_all_browsers").(bool),
		AllowAllBrowsers:                d.Get("allow_all_browsers").(bool),
		EnableWarnings:                  d.Get("enable_warnings").(bool),
		EnableSmartBrowserIsolation:     d.Get("enable_smart_browser_isolation").(bool),
		BypassPlugins:                   SetToStringList(d, "bypass_plugins"),
		BypassApplications:              SetToStringList(d, "bypass_applications"),
		BlockedInternetExplorerVersions: SetToStringList(d, "blocked_internet_explorer_versions"),
		BlockedChromeVersions:           SetToStringList(d, "blocked_chrome_versions"),
		BlockedFirefoxVersions:          SetToStringList(d, "blocked_firefox_versions"),
		BlockedSafariVersions:           SetToStringList(d, "blocked_safari_versions"),
		BlockedOperaVersions:            SetToStringList(d, "blocked_opera_versions"),
		SmartIsolationProfile:           expandSmartIsolationProfile(d),
		SmartIsolationGroups:            expandIDNameExtensionsSet(d, "smart_isolation_groups"),
		SmartIsolationUsers:             expandIDNameExtensionsSet(d, "smart_isolation_users"),
	}
}

func expandSmartIsolationProfile(d *schema.ResourceData) *secure_browsing.SmartIsolationProfile {
	raw, ok := d.GetOk("smart_isolation_profile")
	if !ok {
		return nil
	}

	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}

	first, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}

	id, _ := first["id"].(string)
	name, _ := first["name"].(string)
	url, _ := first["url"].(string)

	if id == "" && name == "" && url == "" {
		return nil
	}

	return &secure_browsing.SmartIsolationProfile{
		ID:   id,
		Name: name,
		URL:  url,
	}
}

// validateSmartBrowserIsolation enforces the contract that Smart Browser
// Isolation needs both the toggle and a fully-populated profile:
//   - When enable_smart_browser_isolation is true, smart_isolation_profile
//     must be configured.
//   - When smart_isolation_profile is configured, all three of id, name and
//     url must be populated. The zia_cloud_browser_isolation_profile data
//     source returns the values you need.
func validateSmartBrowserIsolation(d *schema.ResourceData) error {
	enabled := d.Get("enable_smart_browser_isolation").(bool)

	rawProfile, _ := d.GetOk("smart_isolation_profile")
	profileList, _ := rawProfile.([]interface{})
	hasProfile := len(profileList) > 0 && profileList[0] != nil

	if enabled && !hasProfile {
		return fmt.Errorf("smart_isolation_profile must be configured when enable_smart_browser_isolation is true")
	}

	if hasProfile {
		first, ok := profileList[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("smart_isolation_profile block is malformed")
		}
		id, _ := first["id"].(string)
		name, _ := first["name"].(string)
		url, _ := first["url"].(string)
		if id == "" || name == "" || url == "" {
			return fmt.Errorf("smart_isolation_profile requires id, name, and url to be set; obtain them from the zia_cloud_browser_isolation_profile data source")
		}
	}

	return nil
}

// shouldUpdateSmartIsolation reports whether the resource needs to send a
// PUT to the Smart Isolation endpoint. It returns true when the user has
// opted into Smart Isolation by setting any of the smart_isolation_* fields
// to a non-zero value, OR when an Update is removing previously configured
// Smart Isolation values (so the API can be brought back in line with the
// declared HCL).
func shouldUpdateSmartIsolation(d *schema.ResourceData) bool {
	if d.Get("enable_smart_browser_isolation").(bool) {
		return true
	}
	if list, ok := d.Get("smart_isolation_profile").([]interface{}); ok && len(list) > 0 {
		return true
	}
	if set, ok := d.Get("smart_isolation_users").(*schema.Set); ok && set.Len() > 0 {
		return true
	}
	if set, ok := d.Get("smart_isolation_groups").(*schema.Set); ok && set.Len() > 0 {
		return true
	}
	if !d.IsNewResource() && d.HasChanges(
		"enable_smart_browser_isolation",
		"smart_isolation_profile",
		"smart_isolation_users",
		"smart_isolation_groups",
	) {
		return true
	}
	return false
}

// preserveAnyBypassValue handles a Browser Control API quirk where the value
// "ANY" in the request is accepted but normalized server-side to either an
// empty list (bypass_plugins / bypass_applications) or ["NONE"]
// (blocked_*_versions). The API treats ANY, NONE and "no entries" as
// equivalent, so without this normalization every plan after apply shows
// drift because the user-declared "ANY" value never round-trips back.
//
// If the API returned the "no entries" sentinel AND the current state value
// (which on a fresh Create reflects the user-declared config, and on a
// subsequent refresh reflects the previously normalized state) is exactly
// ["ANY"], we keep ["ANY"] so subsequent plans converge. Otherwise the API
// value is written through unchanged.
func preserveAnyBypassValue(d *schema.ResourceData, key string, apiList []string) []string {
	if !isBrowserControlNoneSentinel(apiList) {
		return apiList
	}

	set, ok := d.Get(key).(*schema.Set)
	if !ok || set == nil || set.Len() != 1 {
		return apiList
	}

	for _, v := range set.List() {
		if s, ok := v.(string); ok && strings.EqualFold(s, "ANY") {
			return []string{"ANY"}
		}
	}
	return apiList
}

// isBrowserControlNoneSentinel reports whether the API response represents
// "no entries" — either an empty list or a single-element list containing
// only "NONE".
func isBrowserControlNoneSentinel(list []string) bool {
	switch len(list) {
	case 0:
		return true
	case 1:
		return strings.EqualFold(list[0], "NONE")
	default:
		return false
	}
}
