package zia

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_settings"
)

func resourceSandboxSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceSandboxSettingsCreate,
		Read:   resourceSandboxSettingsRead,
		Update: resourceSandboxSettingsUpdate,
		Delete: resourceSandboxSettingsDelete,
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

func resourceSandboxSettingsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	fileHashes := expandAndSortSandboxSettings(d)

	_, err := zClient.sandbox_settings.Update(fileHashes)
	if err != nil {
		return err
	}
	d.SetId("hash_list")
	return resourceSandboxSettingsRead(d, m)
}

func resourceSandboxSettingsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	resp, err := zClient.sandbox_settings.Get()
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

	stateHashes := expandAndSortSandboxSettings(d)
	currentSettings, err := zClient.sandbox_settings.Get()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(stateHashes.FileHashesToBeBlocked, sortStringSlice(currentSettings.FileHashesToBeBlocked)) {
		_, err := zClient.sandbox_settings.Update(stateHashes)
		if err != nil {
			return err
		}
	}

	return resourceSandboxSettingsRead(d, m)
}

func resourceSandboxSettingsDelete(d *schema.ResourceData, m interface{}) error {
	return nil // Since there is no DELETE method for this API.
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
