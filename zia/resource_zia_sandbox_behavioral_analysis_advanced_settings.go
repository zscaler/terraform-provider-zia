package zia

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
)

func resourceSandboxSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSandboxSettingsCreate,
		ReadContext:   resourceSandboxSettingsRead,
		UpdateContext: resourceSandboxSettingsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				// Use the Get method from the SDK to fetch the MD5 file hashes
				hashes, err := sandbox_settings.Get(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("error fetching MD5 file hashes: %s", err)
				}

				// Assuming BaAdvancedSettings struct contains a slice of MD5 hashes under the attribute name "FileHashesToBeBlocked"
				if hashes != nil {
					if err := d.Set("file_hashes_to_be_blocked", hashes.FileHashesToBeBlocked); err != nil {
						return nil, fmt.Errorf("error setting file_hashes_to_be_blocked: %s", err)
					}
				}

				// Set an ID for the imported resource. Since this resource seems to encompass a global setting rather than individual items,
				// a static ID like "sandbox_settings" could be used. Adjust if your use case requires.
				d.SetId("sandbox_settings")

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"file_hashes_to_be_blocked": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A custom list of unique MD5 file hashes that must be blocked by Sandbox. A maximum of 10000 MD5 file hashes can be blocked",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSandboxSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	fileHashes := expandAndSortSandboxSettings(d)

	// Validate hashes
	err := validateHashes(fileHashes.FileHashesToBeBlocked)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = sandbox_settings.Update(ctx, service, fileHashes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("sandbox_settings")
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

	return resourceSandboxSettingsRead(ctx, d, meta)
}

func resourceSandboxSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := sandbox_settings.Get(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	// ✅ Always set static ID regardless of API response
	d.SetId("sandbox_settings")

	hashes := []string{}
	if resp != nil && resp.FileHashesToBeBlocked != nil {
		hashes = sortStringSlice(resp.FileHashesToBeBlocked)
	}

	if err := d.Set("file_hashes_to_be_blocked", hashes); err != nil {
		return diag.FromErr(fmt.Errorf("error setting file_hashes_to_be_blocked: %s", err))
	}

	return nil
}

func resourceSandboxSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	stateHashes := expandAndSortSandboxSettings(d)

	// Validate hashes
	if err := validateHashes(stateHashes.FileHashesToBeBlocked); err != nil {
		return diag.FromErr(err)
	}

	currentSettings, err := sandbox_settings.Get(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(stateHashes.FileHashesToBeBlocked, sortStringSlice(currentSettings.FileHashesToBeBlocked)) {
		_, err := sandbox_settings.Update(ctx, service, stateHashes)
		if err != nil {
			return diag.FromErr(err)
		}
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

	return resourceSandboxSettingsRead(ctx, d, meta)
}

func expandAndSortSandboxSettings(d *schema.ResourceData) sandbox_settings.BaAdvancedSettings {
	rawHashes := SetToStringList(d, "file_hashes_to_be_blocked")
	sortedHashes := sortStringSlice(rawHashes)
	return sandbox_settings.BaAdvancedSettings{
		FileHashesToBeBlocked: sortedHashes,
	}
}
