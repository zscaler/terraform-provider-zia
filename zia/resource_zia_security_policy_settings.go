package zia

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/security_policy_settings"
)

func resourceSecurityPolicySettings() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSecurityPolicySettingsRead,
		Create: resourceSecurityPolicySettingsCreate,
		Update: resourceSecurityPolicySettingsUpdate,
		Delete: resourceSecurityPolicySettingsDelete,
		Schema: map[string]*schema.Schema{
			"whitelist_urls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    255,
				Description: "Allowlist URLs whose contents will not be scanned. Allows up to 255 URLs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"blacklist_urls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    25000,
				Description: "URLs on the denylist for your organization. Allow up to 25000 URLs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func expandSecurityPolicySettings(d *schema.ResourceData) security_policy_settings.ListUrls {
	return security_policy_settings.ListUrls{
		Black: SetToStringList(d, "blacklist_urls"),
		White: SetToStringList(d, "whitelist_urls"),
	}
}

func resourceSecurityPolicySettingsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	listUrls := expandSecurityPolicySettings(d)
	_, err := zClient.security_policy_settings.UpdateListUrls(listUrls)
	if err != nil {
		return err
	}
	d.SetId("url_list")
	return resourceSecurityPolicySettingsRead(d, m)
}

func resourceSecurityPolicySettingsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	listUrls := expandSecurityPolicySettings(d)

	_, err := zClient.security_policy_settings.UpdateListUrls(listUrls)
	if err != nil {
		return err
	}
	return resourceSecurityPolicySettingsRead(d, m)
}

func resourceSecurityPolicySettingsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.security_policy_settings.GetListUrls()
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("url_id")
		_ = d.Set("whitelist_urls", resp.White)
		_ = d.Set("blacklist_urls", resp.Black)

	} else {
		return fmt.Errorf("couldn't read urls")
	}

	return nil
}

func resourceSecurityPolicySettingsDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
