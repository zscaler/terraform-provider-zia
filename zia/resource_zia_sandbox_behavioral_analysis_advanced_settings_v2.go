package zia

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_settings"
)

func resourceSandboxSettingsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSandboxSettingsV2Create,
		ReadContext:   resourceSandboxSettingsV2Read,
		UpdateContext: resourceSandboxSettingsV2Update,
		DeleteContext: resourceFuncNoOp,
		CustomizeDiff: sandboxSettingsV2CustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				payload, err := sandbox_settings.Getv2(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("error fetching sandbox behavioral analysis advanced settings: %s", err)
				}

				if payload != nil && len(payload.Md5HashValueList) > 0 {
					if err := d.Set("md5_hash_value_list", flattenMd5HashValueList(payload.Md5HashValueList)); err != nil {
						return nil, fmt.Errorf("error setting md5_hash_value_list: %s", err)
					}
				}

				d.SetId("sandbox_settings")
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"md5_hash_value_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A custom list of MD5 hash values with metadata for sandbox blocking. Each entry includes a URL, optional comment, and type (e.g. MALWARE).",
				Set:         md5HashValueListHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL or hash identifier for the MD5 entry.",
						},
						"url_comment": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional comment describing the URL or hash entry.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the entry (CUSTOM_FILEHASH_ALLOW,  CUSTOM_FILEHASH_DENY, MALWARE).",
							ValidateFunc: validation.StringInSlice([]string{
								"CUSTOM_FILEHASH_ALLOW",
								"CUSTOM_FILEHASH_DENY",
								"MALWARE",
							}, false),
						},
					},
				},
			},
		},
	}
}

// sandboxSettingsV2CustomizeDiff suppresses drift when config has empty md5_hash_value_list blocks
// (url=="" && type=="") but state has no blocks—both represent "no hashes".
func sandboxSettingsV2CustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	oldRaw, newRaw := d.GetChange("md5_hash_value_list")
	oldEffective := rawToEffectiveMd5HashValueList(oldRaw)
	newEffective := rawToEffectiveMd5HashValueList(newRaw)
	if md5HashValueListsEqual(oldEffective, newEffective) {
		if err := d.Clear("md5_hash_value_list"); err != nil {
			log.Printf("[WARN] failed to clear md5_hash_value_list diff: %s", err)
		}
	}
	return nil
}

// rawToEffectiveMd5HashValueList converts schema raw value to effective Md5HashValue slice
// (filters empty blocks the same way expandMd5HashValueList does).
func rawToEffectiveMd5HashValueList(raw interface{}) []sandbox_settings.Md5HashValue {
	if raw == nil {
		return nil
	}
	var list []interface{}
	switch v := raw.(type) {
	case *schema.Set:
		list = v.List()
	case []interface{}:
		list = v
	default:
		return nil
	}
	return expandMd5HashValueList(list)
}

func resourceSandboxSettingsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service
	payload := expandMd5HashValueListPayload(d)

	_, err := sandbox_settings.Updatev2(ctx, service, payload)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("sandbox_settings")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSandboxSettingsV2Read(ctx, d, meta)
}

func resourceSandboxSettingsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := sandbox_settings.Getv2(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("sandbox_settings")

	if resp != nil && len(resp.Md5HashValueList) > 0 {
		if err := d.Set("md5_hash_value_list", flattenMd5HashValueList(resp.Md5HashValueList)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting md5_hash_value_list: %s", err))
		}
		return nil
	}

	// API returned empty list. Only overwrite state if it has real (non-empty) entries.
	// This preserves empty blocks from config in state, avoiding a false diff between
	// "0 set elements" (cleared state) and "1 element with all-empty fields" (config's empty block).
	stateEffective := rawToEffectiveMd5HashValueList(d.Get("md5_hash_value_list"))
	if len(stateEffective) > 0 {
		if err := d.Set("md5_hash_value_list", []interface{}{}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting md5_hash_value_list: %s", err))
		}
	}

	return nil
}

func resourceSandboxSettingsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	statePayload := expandMd5HashValueListPayload(d)

	currentSettings, err := sandbox_settings.Getv2(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	if !md5HashValueListsEqual(statePayload.Md5HashValueList, currentSettings.Md5HashValueList) {
		_, err := sandbox_settings.Updatev2(ctx, service, statePayload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSandboxSettingsV2Read(ctx, d, meta)
}

func expandMd5HashValueListPayload(d *schema.ResourceData) sandbox_settings.Md5HashValueListPayload {
	rawSet := d.Get("md5_hash_value_list")
	if rawSet == nil {
		return sandbox_settings.Md5HashValueListPayload{
			Md5HashValueList: []sandbox_settings.Md5HashValue{},
		}
	}
	set, ok := rawSet.(*schema.Set)
	if !ok {
		// Handle []interface{} from config or migration
		if list, ok := rawSet.([]interface{}); ok && len(list) > 0 {
			return sandbox_settings.Md5HashValueListPayload{
				Md5HashValueList: expandMd5HashValueList(list),
			}
		}
		return sandbox_settings.Md5HashValueListPayload{
			Md5HashValueList: []sandbox_settings.Md5HashValue{},
		}
	}
	if set.Len() == 0 {
		return sandbox_settings.Md5HashValueListPayload{
			Md5HashValueList: []sandbox_settings.Md5HashValue{},
		}
	}
	return sandbox_settings.Md5HashValueListPayload{
		Md5HashValueList: expandMd5HashValueList(set.List()),
	}
}

// md5HashValueListHash computes a hash for a single md5_hash_value_list element (used by TypeSet).
func md5HashValueListHash(v interface{}) int {
	if v == nil {
		return schema.HashString("")
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return schema.HashString(fmt.Sprintf("%v", v))
	}
	url := ""
	if u, ok := m["url"].(string); ok {
		url = u
	}
	urlComment := ""
	if uc, ok := m["url_comment"].(string); ok {
		urlComment = uc
	}
	typ := ""
	if t, ok := m["type"].(string); ok {
		typ = t
	}
	return schema.HashString(fmt.Sprintf("%s|%s|%s", url, urlComment, typ))
}

// md5HashValueListsEqual compares two Md5HashValue slices regardless of order.
func md5HashValueListsEqual(a, b []sandbox_settings.Md5HashValue) bool {
	if len(a) != len(b) {
		return false
	}
	// Sort both by url+type for deterministic comparison
	sortedA := make([]sandbox_settings.Md5HashValue, len(a))
	copy(sortedA, a)
	sort.Slice(sortedA, func(i, j int) bool {
		keyI := sortedA[i].URL + "|" + sortedA[i].Type
		keyJ := sortedA[j].URL + "|" + sortedA[j].Type
		return keyI < keyJ
	})
	sortedB := make([]sandbox_settings.Md5HashValue, len(b))
	copy(sortedB, b)
	sort.Slice(sortedB, func(i, j int) bool {
		keyI := sortedB[i].URL + "|" + sortedB[i].Type
		keyJ := sortedB[j].URL + "|" + sortedB[j].Type
		return keyI < keyJ
	})
	for i := range sortedA {
		if sortedA[i].URL != sortedB[i].URL ||
			sortedA[i].URLComment != sortedB[i].URLComment ||
			sortedA[i].Type != sortedB[i].Type {
			return false
		}
	}
	return true
}

func expandMd5HashValueList(v []interface{}) []sandbox_settings.Md5HashValue {
	out := make([]sandbox_settings.Md5HashValue, 0, len(v))
	for _, raw := range v {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		url := getString(item["url"])
		typ := getString(item["type"])
		// Skip empty blocks - e.g. md5_hash_value_list {} with no url/type
		if url == "" && typ == "" {
			continue
		}
		entry := sandbox_settings.Md5HashValue{
			URL:        url,
			URLComment: getString(item["url_comment"]),
			Type:       typ,
		}
		out = append(out, entry)
	}
	// Return empty slice (not nil) so API receives {"md5HashValueList":[]} not {}
	return out
}

func flattenMd5HashValueList(list []sandbox_settings.Md5HashValue) []interface{} {
	if len(list) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(list))
	for _, v := range list {
		m := map[string]interface{}{
			"url":         v.URL,
			"url_comment": v.URLComment,
			"type":        v.Type,
		}
		out = append(out, m)
	}
	return out
}
