package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
)

func resourceBandwdithClasses() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwdithClassesCreate,
		ReadContext:   resourceBandwdithClassesRead,
		UpdateContext: resourceBandwdithClassesUpdate,
		DeleteContext: resourceBandwdithClassesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("class_id", idInt)
				} else {
					resp, err := bandwidth_classes.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("class_id", resp.ID)
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
			"class_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bandwidth classname.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bandwidth classname.",
			},
			"is_name_l10n_tag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The bandwidth classname.",
			},
			"file_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The bandwidth classname.",
			},
			"urls": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"web_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBandwdithClassesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandBandwidthClasses(d)

	log.Printf("[INFO] Creating ZIA bandwidth classs\n%+v\n", req)

	resp, _, err := bandwidth_classes.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA bandwidth classs request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("class_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceBandwdithClassesRead(ctx, d, meta)
}

func resourceBandwdithClassesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "class_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no bandwidth classs id is set"))
	}
	resp, err := bandwidth_classes.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia bandwidth classs %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia bandwidth classs:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("class_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("type", resp.Type)
	_ = d.Set("file_size", resp.FileSize)
	_ = d.Set("urls", resp.Urls)
	_ = d.Set("url_categories", resp.UrlCategories)
	_ = d.Set("web_applications", resp.WebApplications)
	_ = d.Set("applications", resp.Applications)
	_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

	return nil
}

func resourceBandwdithClassesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "class_id")
	if !ok {
		log.Printf("[ERROR] bandwidth class ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia bandwidth class ID: %v\n", id)

	req := expandBandwidthClasses(d)
	if _, err := bandwidth_classes.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := bandwidth_classes.Update(ctx, service, id, &req); err != nil {
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

	return resourceBandwdithClassesRead(ctx, d, meta)
}

func resourceBandwdithClassesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "class_id")
	if !ok {
		log.Printf("[ERROR] bandwidth class ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia bandwidth class ID: %v\n", (d.Id()))

	if _, err := bandwidth_classes.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia bandwidth class deleted")

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

func expandBandwidthClasses(d *schema.ResourceData) bandwidth_classes.BandwidthClasses {
	id, _ := getIntFromResourceData(d, "class_id")
	result := bandwidth_classes.BandwidthClasses{
		ID:              id,
		Name:            d.Get("name").(string),
		Type:            d.Get("type").(string),
		FileSize:        d.Get("file_size").(string),
		IsNameL10nTag:   d.Get("is_name_l10n_tag").(bool),
		Urls:            SetToStringList(d, "urls"),
		UrlCategories:   SetToStringList(d, "url_categories"),
		Applications:    SetToStringList(d, "applications"),
		WebApplications: SetToStringList(d, "web_applications"),
	}
	return result
}
