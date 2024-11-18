package zia

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/security_policy_settings"
)

func resourceSecurityPolicySettings() *schema.Resource {
	return &schema.Resource{
		Read:   resourceSecurityPolicySettingsRead,
		Create: resourceSecurityPolicySettingsCreate,
		Update: resourceSecurityPolicySettingsUpdate,
		Delete: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.security_policy_settings

				// Use the GetListUrls method to fetch both whitelist and blacklist URLs.
				resp, err := security_policy_settings.GetListUrls(service)
				if err != nil {
					return []*schema.ResourceData{}, err
				}

				// Set the whitelist and blacklist URLs in the Terraform state.
				if err := d.Set("whitelist_urls", resp.White); err != nil {
					return []*schema.ResourceData{}, fmt.Errorf("error setting whitelist_urls: %s", err)
				}
				if err := d.Set("blacklist_urls", resp.Black); err != nil {
					return []*schema.ResourceData{}, fmt.Errorf("error setting blacklist_urls: %s", err)
				}

				// Set a generic ID since we're not differentiating based on import type.
				d.SetId("all_urls")

				return []*schema.ResourceData{d}, nil
			},
		},

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
				MaxItems:    275000,
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

	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := m.(*Client)
	service := zClient.security_policy_settings
	listUrls := expandSecurityPolicySettings(d)
	_, err := security_policy_settings.UpdateListUrls(service, listUrls)
	if err != nil {
		return err
	}
	d.SetId("url_list")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSecurityPolicySettingsRead(d, m)
}

func resourceSecurityPolicySettingsUpdate(d *schema.ResourceData, m interface{}) error {

	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := m.(*Client)
	service := zClient.security_policy_settings
	listUrls := expandSecurityPolicySettings(d)

	_, err := security_policy_settings.UpdateListUrls(service, listUrls)
	if err != nil {
		return err
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return activationErr
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSecurityPolicySettingsRead(d, m)
}

func resourceSecurityPolicySettingsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.security_policy_settings

	resp, err := security_policy_settings.GetListUrls(service)
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
