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

func resourceATPMaliciousUrls() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceATPMaliciousUrlsRead,
		CreateContext: resourceATPMaliciousUrlsCreate,
		UpdateContext: resourceATPMaliciousUrlsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				urls, err := advancedthreatsettings.GetMaliciousURLs(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("error fetching urls from exception list: %s", err)
				}

				if urls != nil {
					if err := d.Set("malicious_urls", urls.MaliciousUrls); err != nil {
						return nil, fmt.Errorf("error setting urls: %s", err)
					}
				}

				d.SetId("all_urls")

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"malicious_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				MaxItems: 275000,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceATPMaliciousUrlsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	maliciousUrls := expandATPMaliciousURLs(d)
	_, err := advancedthreatsettings.UpdateMaliciousURLs(ctx, service, maliciousUrls)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("all_urls")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceATPMaliciousUrlsRead(ctx, d, meta)
}

func resourceATPMaliciousUrlsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := advancedthreatsettings.GetMaliciousURLs(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("all_urls")
		_ = d.Set("malicious_urls", resp.MaliciousUrls)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read malicious urls"))
	}

	return nil
}

func resourceATPMaliciousUrlsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	maliciousUrls := expandATPMaliciousURLs(d)

	_, err := advancedthreatsettings.UpdateMaliciousURLs(ctx, service, maliciousUrls)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("all_urls")

	// Sleep for 1 seconds before potentially triggering the activation
	time.Sleep(1 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceATPMaliciousUrlsRead(ctx, d, meta)
}

// func expandATPMaliciousURLs(d *schema.ResourceData) []string {
// 	return SetToStringList(d, "malicious_urls")
// }

func expandATPMaliciousURLs(d *schema.ResourceData) advancedthreatsettings.MaliciousURLs {
	return advancedthreatsettings.MaliciousURLs{
		MaliciousUrls: SetToStringList(d, "malicious_urls"),
	}
}
