package zia

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/ftp_control_policy"
)

func resourceFTPControlPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceFTPControlPolicyRead,
		CreateContext: resourceFTPControlPolicyCreate,
		UpdateContext: resourceFTPControlPolicyUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				diags := resourceFTPControlPolicyRead(ctx, d, meta)
				if diags.HasError() {
					return nil, fmt.Errorf("import read error: %v", diags)
				}
				d.SetId("advanced_settings")
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"ftp_over_http_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to enable FTP over HTTP",
			},
			"ftp_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether to enable native FTP. When enabled, users can connect to native FTP sites and download files.",
			},
			"urls": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Domains or URLs included for the FTP Control settings",
			},
			"url_categories": getURLCategories(),
		},
	}
}

func resourceFTPControlPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFTPControlSettings(d)
	_, _, err := ftp_control_policy.UpdateFTPControlPolicy(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("advanced_settings")

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

	return resourceFTPControlPolicyRead(ctx, d, meta)
}

func resourceFTPControlPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := ftp_control_policy.GetFTPControlPolicy(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp != nil {

		d.SetId("ftp_settings")
		_ = d.Set("ftp_over_http_enabled", resp.FtpOverHttpEnabled)
		_ = d.Set("ftp_enabled", resp.FtpEnabled)
		_ = d.Set("url_categories", resp.UrlCategories)
		_ = d.Set("urls", resp.Urls)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read ftp control settings"))
	}

	return nil
}

func resourceFTPControlPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFTPControlSettings(d)

	_, _, err := ftp_control_policy.UpdateFTPControlPolicy(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ftp_settings")

	time.Sleep(1 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceFTPControlPolicyRead(ctx, d, meta)
}

func expandFTPControlSettings(d *schema.ResourceData) ftp_control_policy.FTPControlPolicy {

	result := ftp_control_policy.FTPControlPolicy{
		FtpOverHttpEnabled: d.Get("ftp_over_http_enabled").(bool),
		FtpEnabled:         d.Get("ftp_enabled").(bool),
		UrlCategories:      SetToStringList(d, "url_categories"),
		Urls:               SetToStringList(d, "urls"),
	}
	return result
}
