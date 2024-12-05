package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/advancedthreatsettings"
)

func resourceATPSecurityExceptions() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceATPSecurityExceptionsRead,
		CreateContext: resourceATPSecurityExceptionsCreate,
		UpdateContext: resourceATPSecurityExceptionsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				urls, err := advancedthreatsettings.GetSecurityExceptions(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("error fetching urls from bypass list: %s", err)
				}

				if urls != nil {
					if err := d.Set("urls", urls.BypassUrls); err != nil {
						return nil, fmt.Errorf("error setting urls: %s", err)
					}
				}

				d.SetId("all_urls")

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"bypass_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceATPSecurityExceptionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	bypassUrls := expandATPBypassURLs(d)
	_, err := advancedthreatsettings.UpdateSecurityExceptions(ctx, service, bypassUrls)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("bypass_url")

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

	return resourceATPSecurityExceptionsRead(ctx, d, meta)
}

func resourceATPSecurityExceptionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := advancedthreatsettings.GetSecurityExceptions(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("bypass_url")
		_ = d.Set("bypass_urls", resp.BypassUrls)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read bypass urls"))
	}

	return nil
}

func resourceATPSecurityExceptionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	bypassUrls := expandATPBypassURLs(d)

	_, err := advancedthreatsettings.UpdateSecurityExceptions(ctx, service, bypassUrls)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("bypass_url")

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

	return resourceATPSecurityExceptionsRead(ctx, d, meta)
}

func expandATPBypassURLs(d *schema.ResourceData) []string {
	return SetToStringList(d, "bypass_urls")
}
