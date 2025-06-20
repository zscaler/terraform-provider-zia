package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_control_settings"
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
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of plugins that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable plugins are warned",
			},
			"bypass_applications": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of applications that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable applications are warned",
			},
			"blocked_internet_explorer_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Microsoft browser that need to be blocked. If not set, all Microsoft browser versions are allowed.",
			},
			"blocked_chrome_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Google Chrome browser that need to be blocked. If not set, all Google Chrome versions are allowed.",
			},
			"blocked_firefox_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Mozilla Firefox browser that need to be blocked. If not set, all Mozilla Firefox versions are allowed.",
			},
			"blocked_safari_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Apple Safari browser that need to be blocked. If not set, all Apple Safari versions are allowed",
			},
			"blocked_opera_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Opera browser that need to be blocked. If not set, all Opera versions are allowed",
			},
			"bypass_all_browsers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, all the browsers are bypassed for warnings",
			},
			"allow_all_browsers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value that specifies whether or not to allow all the browsers and their respective versions access to the internet",
			},
			"enable_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value that specifies if the warnings are enabled",
			},
			"enable_smart_browser_isolation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value that specifies if Smart Browser Isolation is enabled",
			},
			"smart_isolation_profile": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The isolation profile",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The universally unique identifier (UUID) for the browser isolation profile",
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

	policy := browser_control_settings.BrowserControlSettings{
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
	}

	_, _, err := browser_control_settings.UpdateBrowserControlSettings(ctx, service, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("policy")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
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
		_ = d.Set("bypass_plugins", resp.BypassPlugins)
		_ = d.Set("bypass_applications", resp.BypassApplications)
		_ = d.Set("blocked_internet_explorer_versions", resp.BlockedInternetExplorerVersions)
		_ = d.Set("blocked_chrome_versions", resp.BlockedChromeVersions)
		_ = d.Set("blocked_firefox_versions", resp.BlockedFirefoxVersions)
		_ = d.Set("blocked_safari_versions", resp.BlockedSafariVersions)
		_ = d.Set("blocked_opera_versions", resp.BlockedOperaVersions)
		_ = d.Set("bypass_all_browsers", resp.BypassAllBrowsers)
		_ = d.Set("allow_all_browsers", resp.AllowAllBrowsers)
		_ = d.Set("enable_warnings", resp.EnableWarnings)
		_ = d.Set("enable_smart_browser_isolation", resp.EnableSmartBrowserIsolation)

		if err := d.Set("smart_isolation_groups", flattenIDExtensionsListIDs(resp.SmartIsolationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("smart_isolation_users", flattenIDExtensionsListIDs(resp.SmartIsolationUsers)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("smart_isolation_profile", flattenSmartIsolationProfileID(&resp.SmartIsolationProfile)); err != nil {
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

	policy := browser_control_settings.BrowserControlSettings{
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

	_, _, err := browser_control_settings.UpdateBrowserControlSettings(ctx, service, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("browser_settings")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceBrowserControlPolicyRead(ctx, d, meta)
}

func expandSmartIsolationProfile(d *schema.ResourceData) browser_control_settings.SmartIsolationProfile {
	raw, ok := d.GetOk("smart_isolation_profile")
	if !ok {
		return browser_control_settings.SmartIsolationProfile{}
	}

	list := raw.([]interface{})
	if len(list) == 0 {
		return browser_control_settings.SmartIsolationProfile{}
	}

	first := list[0].(map[string]interface{})
	id, _ := first["id"].(string)

	return browser_control_settings.SmartIsolationProfile{
		ID: id,
	}
}

func flattenSmartIsolationProfileID(profile *browser_control_settings.SmartIsolationProfile) []interface{} {
	if profile == nil || profile.ID == "" {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": profile.ID,
		},
	}
}
