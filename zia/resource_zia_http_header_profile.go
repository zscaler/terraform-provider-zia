package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/http_header_control/http_header_profile"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func resourceHttpHeaderProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHttpHeaderProfileCreate,
		ReadContext:   resourceHttpHeaderProfileRead,
		UpdateContext: resourceHttpHeaderProfileUpdate,
		DeleteContext: resourceHttpHeaderProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("header_profile_id", idInt)
				} else {
					resp, err := http_header_profile.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("header_profile_id", resp.ID)
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
			"header_profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the HTTP header profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The HTTP header profile name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the HTTP header profile.",
			},
			"slot_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "HTTP header profile slot ID",
			},
			"profile_ready_for_use": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the HTTP header profile is ready for use.",
			},
			"http_header_profile_criteria": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of matching criteria evaluated by the HTTP header profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies the criteria entry.",
						},
						"header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the type of HTTP header",
							ValidateFunc: validation.StringInSlice([]string{
								"USERAGENT",
								"REFERER",
								"ORIGIN",
							}, false),
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The operator applied to the header criteria.",
							ValidateFunc: validation.StringInSlice([]string{
								"UAVERSIONGT",
								"UAVERSIONLT",
								"UAVERSIONEQ",
								"UAVERSIONNEQ",
								"UAVERSIONANY",
							}, false),
						},
						"user_agent": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user agent evaluated by the criteria.",
							ValidateFunc: validation.StringInSlice([]string{
								"OPERA",
								"FIREFOX",
								"MSIE",
								"MSEDGE",
								"CHROME",
								"SAFARI",
								"OTHER",
								"MSCHREDGE",
								"BRAVE",
							}, false),
						},
						"user_agent_bitmap": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user agent bitmap evaluated by the criteria.",
						},
						"user_agent_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user agent version evaluated by the criteria.",
						},
						"category_bitmap": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The URL category bitmap evaluated by the criteria.",
						},
						"cloud_app_bitmap": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The cloud application bitmap evaluated by the criteria.",
						},
					},
				},
			},
		},
	}
}

func resourceHttpHeaderProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandHttpHeaderProfile(d)
	log.Printf("[INFO] Creating ZIA HTTP header profile\n%+v\n", req)

	resp, _, err := http_header_profile.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA HTTP header profile. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("header_profile_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceHttpHeaderProfileRead(ctx, d, meta)
}

func resourceHttpHeaderProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "header_profile_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no HTTP header profile id is set"))
	}

	resp, err := findHttpHeaderProfileByID(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		log.Printf("[WARN] Removing HTTP header profile %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Getting ZIA HTTP header profile:\n%+v\n", resp)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("header_profile_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("slot_id", resp.SlotId)
	_ = d.Set("profile_ready_for_use", resp.ProfileReadyForUse)

	if err := d.Set("http_header_profile_criteria", flattenHttpHeaderProfileCriteria(resp.HttpHeaderProfileCriteria)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceHttpHeaderProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "header_profile_id")
	if !ok {
		log.Printf("[ERROR] HTTP header profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating ZIA HTTP header profile ID: %v\n", id)

	existing, err := findHttpHeaderProfileByID(ctx, service, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if existing == nil {
		d.SetId("")
		return nil
	}

	req := expandHttpHeaderProfile(d)
	if _, _, err := http_header_profile.Update(ctx, service, id, &req); err != nil {
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

	return resourceHttpHeaderProfileRead(ctx, d, meta)
}

func resourceHttpHeaderProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "header_profile_id")
	if !ok {
		log.Printf("[ERROR] HTTP header profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting ZIA HTTP header profile ID: %v\n", d.Id())

	err := DetachURLFilteringRuleRef(
		ctx,
		zClient,
		id,
		"HTTP header profile",
		func(r *urlfilteringpolicies.URLFilteringRule) []common.IDNameExtensions {
			return r.HTTPHeaderProfiles
		},
		func(r *urlfilteringpolicies.URLFilteringRule, ids []common.IDNameExtensions) {
			r.HTTPHeaderProfiles = ids
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err := http_header_profile.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] ZIA HTTP header profile deleted")

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

func expandHttpHeaderProfile(d *schema.ResourceData) http_header_profile.HttpHeaderProfile {
	id, _ := getIntFromResourceData(d, "header_profile_id")
	result := http_header_profile.HttpHeaderProfile{
		ID:                        id,
		Name:                      d.Get("name").(string),
		Description:               d.Get("description").(string),
		SlotId:                    d.Get("slot_id").(int),
		ProfileReadyForUse:        d.Get("profile_ready_for_use").(bool),
		HttpHeaderProfileCriteria: expandHttpHeaderProfileCriteria(d),
	}
	return result
}

func expandHttpHeaderProfileCriteria(d *schema.ResourceData) []http_header_profile.HttpHeaderProfileCriteria {
	raw, ok := d.GetOk("http_header_profile_criteria")
	if !ok {
		return nil
	}
	list := raw.([]interface{})
	criteria := make([]http_header_profile.HttpHeaderProfileCriteria, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		entry := http_header_profile.HttpHeaderProfileCriteria{
			Header:           m["header"].(string),
			Operator:         m["operator"].(string),
			UserAgent:        m["user_agent"].(string),
			UserAgentBitmap:  m["user_agent_bitmap"].(string),
			UserAgentVersion: m["user_agent_version"].(string),
			CategoryBitmap:   nestedSetToStringSlice(m["category_bitmap"]),
			CloudAppBitmap:   nestedSetToStringSlice(m["cloud_app_bitmap"]),
		}
		if v, ok := m["id"].(int); ok {
			entry.Id = v
		}
		criteria = append(criteria, entry)
	}
	return criteria
}

// nestedSetToStringSlice converts a TypeSet value nested inside a block (read
// from the block's map, not from ResourceData) into a []string. SetToStringList
// cannot be used here because it resolves a top-level attribute key, whereas
// category_bitmap / cloud_app_bitmap live inside the http_header_profile_criteria
// block.
func nestedSetToStringSlice(raw interface{}) []string {
	set, ok := raw.(*schema.Set)
	if !ok {
		return nil
	}
	return SetToStringSlice(set)
}

// findHttpHeaderProfileByID resolves a profile by numeric ID from the full
// list. The API does not expose a per-ID lookup endpoint, so the list is
// retrieved once and filtered locally. Returns (nil, nil) when no profile with
// the given ID exists.
func findHttpHeaderProfileByID(ctx context.Context, service *zscaler.Service, id int) (*http_header_profile.HttpHeaderProfile, error) {
	all, err := http_header_profile.GetAll(ctx, service)
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
