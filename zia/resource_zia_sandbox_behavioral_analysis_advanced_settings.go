package zia

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_settings"
)

func resourceSandboxSettings() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSandboxSettingsCreate,
		Read:          resourceSandboxSettingsRead,
		Update:        resourceSandboxSettingsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.sandbox_settings

				// Use the Get method from the SDK to fetch the MD5 file hashes
				hashes, err := sandbox_settings.Get(service)
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

func validateHashes(hashes []string) error {
	for _, hash := range hashes {
		hashType := identifyHashType(hash)
		if hashType != "MD5" {
			return fmt.Errorf("the hash '%s' is a %s type. The sandbox only supports MD5 hashes", hash, hashType)
		}
	}
	return nil
}

func identifyHashType(hash string) string {
	switch len(hash) {
	case 32:
		return "MD5"
	case 40:
		return "SHA1"
	case 64:
		return "SHA256"
	default:
		return "unknown"
	}
}

func resourceSandboxSettingsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.sandbox_settings
	fileHashes := expandAndSortSandboxSettings(d)

	// Validate hashes
	err := validateHashes(fileHashes.FileHashesToBeBlocked)
	if err != nil {
		return err
	}

	_, err = sandbox_settings.Update(service, fileHashes)
	if err != nil {
		return err
	}
	d.SetId("hash_list")
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

	return resourceSandboxSettingsRead(d, m)
}

func resourceSandboxSettingsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.sandbox_settings
	resp, err := sandbox_settings.Get(service)
	if err != nil {
		return err
	}

	if resp != nil {
		d.SetId("hash_list")
		sortedHashes := sortStringSlice(resp.FileHashesToBeBlocked)
		err := d.Set("file_hashes_to_be_blocked", sortedHashes)
		if err != nil {
			return fmt.Errorf("error setting file hashes to be blocked: %s", err)
		}
	} else {
		return fmt.Errorf("couldn't read file hash")
	}
	return nil
}

func resourceSandboxSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.sandbox_settings

	stateHashes := expandAndSortSandboxSettings(d)

	// Validate hashes
	if err := validateHashes(stateHashes.FileHashesToBeBlocked); err != nil {
		return err
	}

	currentSettings, err := sandbox_settings.Get(service)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(stateHashes.FileHashesToBeBlocked, sortStringSlice(currentSettings.FileHashesToBeBlocked)) {
		_, err := sandbox_settings.Update(service, stateHashes)
		if err != nil {
			return err
		}
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

	return resourceSandboxSettingsRead(d, m)
}

func expandAndSortSandboxSettings(d *schema.ResourceData) sandbox_settings.BaAdvancedSettings {
	rawHashes := SetToStringList(d, "file_hashes_to_be_blocked")
	sortedHashes := sortStringSlice(rawHashes)
	return sandbox_settings.BaAdvancedSettings{
		FileHashesToBeBlocked: sortedHashes,
	}
}

func sortStringSlice(slice []string) []string {
	sorted := make([]string, len(slice))
	copy(sorted, slice)
	sort.Strings(sorted)
	return sorted
}
