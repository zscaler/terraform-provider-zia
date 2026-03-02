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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/email_profiles"
)

func resourceEmailProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailProfileCreate,
		ReadContext:   resourceEmailProfileRead,
		UpdateContext: resourceEmailProfileUpdate,
		DeleteContext: resourceEmailProfileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("email_profile_id", idInt)
				} else {
					resp, err := email_profiles.GetEmailProfileByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("email_profile_id", resp.ID)
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
			"email_profile_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "The name of the email recipient profile.",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString,
				DiffSuppressFunc: noChangeInMultiLineText,
				Description:      "The description of the email recipient profile.",
			},
			"emails": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The list of recipient email addresses.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceEmailProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}
	service := zClient.Service

	req := expandEmailProfile(d)
	log.Printf("[INFO] Creating ZIA email profile\n%+v\n", req)

	resp, _, err := email_profiles.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA email profile request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("email_profile_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceEmailProfileRead(ctx, d, meta)
}

func resourceEmailProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "email_profile_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no email_profile_id is set"))
	}
	resp, err := email_profiles.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia_email_profile %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("email_profile_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("emails", resp.Emails)

	return nil
}

func resourceEmailProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "email_profile_id")
	if !ok {
		log.Printf("[ERROR] email profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia email profile ID: %v\n", id)
	req := expandEmailProfile(d)

	if _, err := email_profiles.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, _, err := email_profiles.Update(ctx, service, id, &req); err != nil {
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

	return resourceEmailProfileRead(ctx, d, meta)
}

func resourceEmailProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "email_profile_id")
	if !ok {
		log.Printf("[ERROR] email profile ID not set: %v\n", id)
	}

	if _, err := email_profiles.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	if _, err := email_profiles.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia_email_profile %d deleted", id)

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

func expandEmailProfile(d *schema.ResourceData) email_profiles.EmailProfiles {
	id, _ := getIntFromResourceData(d, "email_profile_id")
	result := email_profiles.EmailProfiles{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: getString(d.Get("description")),
		Emails:      SetToStringList(d, "emails"),
	}
	return result
}
