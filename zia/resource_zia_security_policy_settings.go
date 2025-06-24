package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_policy_settings"
)

func resourceSecurityPolicySettings() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSecurityPolicySettingsRead,
		CreateContext: resourceSecurityPolicySettingsCreate,
		UpdateContext: resourceSecurityPolicySettingsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				// Use the GetListUrls method to fetch both whitelist and blacklist URLs.
				resp, err := security_policy_settings.GetListUrls(ctx, service)
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
				Description: "URLs on the denylist for your organization. Allow up to 275000 URLs.",
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

func resourceSecurityPolicySettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := meta.(*Client)
	service := zClient.Service
	listUrls := expandSecurityPolicySettings(d)
	_, err := security_policy_settings.UpdateListUrls(ctx, service, listUrls)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("all_urls")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSecurityPolicySettingsRead(ctx, d, meta)
}

func resourceSecurityPolicySettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := meta.(*Client)
	service := zClient.Service
	listUrls := expandSecurityPolicySettings(d)

	_, err := security_policy_settings.UpdateListUrls(ctx, service, listUrls)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSecurityPolicySettingsRead(ctx, d, meta)
}

func resourceSecurityPolicySettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := security_policy_settings.GetListUrls(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("all_urls")
		_ = d.Set("whitelist_urls", resp.White)
		_ = d.Set("blacklist_urls", resp.Black)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read urls"))
	}

	return nil
}
