package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_control_settings"
)

func dataSourceBrowserControlPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrowserControlPolicyRead,
		Schema: map[string]*schema.Schema{
			"plugin_check_frequency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies how frequently the service checks browsers and relevant applications to warn users regarding outdated or vulnerable browsers, plugins, and applications. If not set, the warnings are disabled",
			},
			"bypass_plugins": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of plugins that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable plugins are warned.",
			},
			"bypass_applications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of applications that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable applications are warned.",
			},
			"blocked_internet_explorer_versions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Microsoft browser that need to be blocked. If not set, all Microsoft browser versions are allowed.",
			},
			"blocked_chrome_versions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Google Chrome browser that need to be blocked. If not set, all Google Chrome versions are allowed.",
			},
			"blocked_firefox_versions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Mozilla Firefox browser that need to be blocked. If not set, all Mozilla Firefox versions are allowed.",
			},
			"blocked_safari_versions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Apple Safari browser that need to be blocked. If not set, all Apple Safari versions are allowed",
			},
			"blocked_opera_versions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Versions of Opera browser that need to be blocked. If not set, all Opera versions are allowed",
			},
			"bypass_all_browsers": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, all the browsers are bypassed for warnings.",
			},
			"allow_all_browsers": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that specifies whether or not to allow all the browsers and their respective versions access to the internet",
			},
			"enable_warnings": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that specifies if the warnings are enabled",
			},
			"enable_smart_browser_isolation": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value that specifies if Smart Browser Isolation is enabled",
			},
			"smart_isolation_profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The isolation profile ID",
			},
			"smart_isolation_profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The universally unique identifier (UUID) for the browser isolation profile",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the browser isolation profile",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The browser isolation profile URL",
						},
						"default_profile": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this is a default browser isolation profile. Zscaler sets this field.",
						},
					},
				},
			},
		},
	}
}

func dataSourceBrowserControlPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		if err := d.Set("smart_isolation_groups", flattenIDNameExtensions(resp.SmartIsolationGroups)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("smart_isolation_users", flattenIDNameExtensions(resp.SmartIsolationUsers)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("smart_isolation_profile", flattenSmartIsolationProfile(&resp.SmartIsolationProfile)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read browser control setting options"))
	}

	return nil
}

func flattenSmartIsolationProfile(cbiProfile *browser_control_settings.SmartIsolationProfile) []interface{} {
	if cbiProfile == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"id":              cbiProfile.ID,
		"name":            cbiProfile.Name,
		"url":             cbiProfile.URL,
		"default_profile": cbiProfile.DefaultProfile,
	}}
}
