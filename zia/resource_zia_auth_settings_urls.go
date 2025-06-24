package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/user_authentication_settings"
)

func resourceAuthSettingsUrls() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAuthSettingsUrlsRead,
		CreateContext: resourceAuthSettingsUrlsCreate,
		UpdateContext: resourceAuthSettingsUrlsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				urls, err := user_authentication_settings.Get(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("error fetching urls from exception list: %s", err)
				}

				if urls != nil {
					if err := d.Set("urls", urls.URLs); err != nil {
						return nil, fmt.Errorf("error setting urls: %s", err)
					}
				}

				d.SetId("all_urls")

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				MaxItems: 25000,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAuthSettingsUrlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	res, err := user_authentication_settings.Get(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("all_urls")
	_ = d.Set("urls", res.URLs)
	return nil
}

func expandAuthSettingsUrls(d *schema.ResourceData) user_authentication_settings.ExemptedUrls {
	return user_authentication_settings.ExemptedUrls{
		URLs: SetToStringList(d, "urls"),
	}
}

func resourceAuthSettingsUrlsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := meta.(*Client)
	service := zClient.Service

	urls := expandAuthSettingsUrls(d)
	_, err := user_authentication_settings.Update(ctx, service, urls)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("all_urls")
	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceAuthSettingsUrlsRead(ctx, d, meta)
}

func resourceAuthSettingsUrlsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Acquire semaphore before making an API request
	apiSemaphore <- struct{}{}
	defer func() { <-apiSemaphore }() // Release semaphore after the request is done

	zClient := meta.(*Client)
	service := zClient.Service

	urls := expandAuthSettingsUrls(d)

	_, err := user_authentication_settings.Update(ctx, service, urls)
	if err != nil {
		return diag.FromErr(err)
	}

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
		return diag.FromErr(activationErr)
	}
	return resourceAuthSettingsUrlsRead(ctx, d, meta)
}
