package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/http_header_control/http_header_action_profile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func resourceHttpHeaderActionProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHttpHeaderActionProfileCreate,
		ReadContext:   resourceHttpHeaderActionProfileRead,
		UpdateContext: resourceHttpHeaderActionProfileUpdate,
		DeleteContext: resourceHttpHeaderActionProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("header_action_profile_id", idInt)
				} else {
					resp, err := http_header_action_profile.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("header_action_profile_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"header_action_profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the HTTP header action profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The HTTP header action profile name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the HTTP header action profile.",
			},
			"slot_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The slot ID assigned to the HTTP header action profile.",
			},
			"profile_ready_for_use": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the HTTP header action profile is ready for use.",
			},
			"http_header_action_profile_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of header key/value pairs applied by the action profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies the header key/value pair.",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The header key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The header value.",
						},
					},
				},
			},
		},
	}
}

func resourceHttpHeaderActionProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandHttpHeaderActionProfile(d)
	log.Printf("[INFO] Creating ZIA HTTP header action profile\n%+v\n", req)

	resp, _, err := http_header_action_profile.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA HTTP header action profile. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("header_action_profile_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceHttpHeaderActionProfileRead(ctx, d, meta)
}

func resourceHttpHeaderActionProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "header_action_profile_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no HTTP header action profile id is set"))
	}

	resp, err := findHttpHeaderActionProfileByID(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		log.Printf("[WARN] Removing HTTP header action profile %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Getting ZIA HTTP header action profile:\n%+v\n", resp)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("header_action_profile_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("slot_id", resp.SlotId)
	_ = d.Set("profile_ready_for_use", resp.ProfileReadyForUse)

	if err := d.Set("http_header_action_profile_keys", flattenHttpHeaderActionProfileKeys(resp.HttpHeaderActionProfileKeys)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceHttpHeaderActionProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "header_action_profile_id")
	if !ok {
		log.Printf("[ERROR] HTTP header action profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating ZIA HTTP header action profile ID: %v\n", id)

	existing, err := findHttpHeaderActionProfileByID(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if existing == nil {
		d.SetId("")
		return nil
	}

	req := expandHttpHeaderActionProfile(d)
	if _, _, err := http_header_action_profile.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceHttpHeaderActionProfileRead(ctx, d, meta)
}

func resourceHttpHeaderActionProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "header_action_profile_id")
	if !ok {
		log.Printf("[ERROR] HTTP header action profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting ZIA HTTP header action profile ID: %v\n", d.Id())
	err := DetachURLFilteringRuleRef(
		ctx,
		zClient,
		id,
		"HTTP header action profile",
		func(r *urlfilteringpolicies.URLFilteringRule) []common.IDNameExtensions {
			return r.HTTPHeaderActionProfiles
		},
		func(r *urlfilteringpolicies.URLFilteringRule, ids []common.IDNameExtensions) {
			r.HTTPHeaderActionProfiles = ids
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := http_header_action_profile.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] ZIA HTTP header action profile deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandHttpHeaderActionProfile(d *schema.ResourceData) http_header_action_profile.HttpHeaderActionProfile {
	id, _ := getIntFromResourceData(d, "header_action_profile_id")
	result := http_header_action_profile.HttpHeaderActionProfile{
		ID:                          id,
		Name:                        d.Get("name").(string),
		Description:                 d.Get("description").(string),
		SlotId:                      d.Get("slot_id").(int),
		ProfileReadyForUse:          d.Get("profile_ready_for_use").(bool),
		HttpHeaderActionProfileKeys: expandHttpHeaderActionProfileKeys(d),
	}
	return result
}

func expandHttpHeaderActionProfileKeys(d *schema.ResourceData) []http_header_action_profile.HttpHeaderActionProfileKeys {
	raw, ok := d.GetOk("http_header_action_profile_keys")
	if !ok {
		return nil
	}
	list := raw.([]interface{})
	keys := make([]http_header_action_profile.HttpHeaderActionProfileKeys, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := http_header_action_profile.HttpHeaderActionProfileKeys{
			Key:   m["key"].(string),
			Value: m["value"].(string),
		}
		if v, ok := m["id"].(int); ok {
			entry.ID = v
		}
		keys = append(keys, entry)
	}
	return keys
}

// findHttpHeaderActionProfileByID resolves a profile by numeric ID from the
// full list. The API does not expose a per-ID lookup endpoint, so the list is
// retrieved once and filtered locally. Returns (nil, nil) when no profile with
// the given ID exists.
func findHttpHeaderActionProfileByID(ctx context.Context, service *zscaler.Service, id int) (*http_header_action_profile.HttpHeaderActionProfile, error) {
	all, err := http_header_action_profile.GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for i := range all {
		if all[i].ID == id {
			return &all[i], nil
		}
	}
	return nil, nil
}
