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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/alerts"
)

func resourceSubscriptionAlerts() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionAlertsCreate,
		ReadContext:   resourceSubscriptionAlertsRead,
		UpdateContext: resourceSubscriptionAlertsUpdate,
		DeleteContext: resourceSubscriptionAlertsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				rawID := d.Id()

				// Attempt direct lookup via ID
				if idInt, err := strconv.Atoi(rawID); err == nil {
					if sub, err := alerts.Get(ctx, service, idInt); err == nil && sub != nil {
						d.Set("alert_id", sub.ID)
						d.SetId(strconv.Itoa(sub.ID))
						return []*schema.ResourceData{d}, nil
					}
				}

				// Fallback to searching all entries
				alertsList, err := alerts.GetAll(ctx, service)
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve alert subscriptions: %w", err)
				}

				for _, a := range alertsList {
					if strconv.Itoa(a.ID) == rawID || a.Email == rawID {
						d.Set("alert_id", a.ID)
						d.Set("email", a.Email)
						d.SetId(strconv.Itoa(a.ID))
						return []*schema.ResourceData{d}, nil
					}
				}

				return nil, fmt.Errorf("alert subscription not found with ID or email '%s'", rawID)
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "System-generated identifier for the alert subscription",
			},
			"alert_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "System-generated identifier for the alert subscription",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address of the alert recipient",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the status of the alert subscription",
			},
			"pt0_severities":    getAlertSubscriptionSeverity(),
			"secure_severities": getAlertSubscriptionSeverity(),
			"manage_severities": getAlertSubscriptionSeverity(),
			"comply_severities": getAlertSubscriptionSeverity(),
			"system_severities": getAlertSubscriptionSeverity(),
		},
	}
}

func resourceSubscriptionAlertsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandSubscriptionAlert(d)
	log.Printf("[INFO] Creating ZIA subscription alerts\n%+v\n", req)

	resp, _, err := alerts.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA subscription alerts request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("alert_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSubscriptionAlertsRead(ctx, d, meta)
}

func resourceSubscriptionAlertsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "alert_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no alert subscriptions id is set"))
	}
	resp, err := alerts.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia alert subscriptions %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia alert subscriptions:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("email", resp.Email)
	_ = d.Set("description", resp.Description)
	_ = d.Set("pt0_severities", resp.Pt0Severities)
	_ = d.Set("secure_severities", resp.SecureSeverities)
	_ = d.Set("manage_severities", resp.ManageSeverities)
	_ = d.Set("comply_severities", resp.ComplySeverities)
	_ = d.Set("system_severities", resp.SystemSeverities)

	return nil
}

func resourceSubscriptionAlertsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "alert_id")
	if !ok {
		log.Printf("[ERROR] subscription alert ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia subscription alert ID: %v\n", id)
	req := expandSubscriptionAlert(d)
	if _, err := alerts.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := alerts.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceSubscriptionAlertsRead(ctx, d, meta)
}

func resourceSubscriptionAlertsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "alert_id")
	if !ok {
		log.Printf("[ERROR] subscription alert ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia subscription alert ID: %v\n", (d.Id()))

	if _, err := alerts.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia subscription alert deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandSubscriptionAlert(d *schema.ResourceData) alerts.AlertSubscriptions {
	id, _ := getIntFromResourceData(d, "alert_id")
	result := alerts.AlertSubscriptions{
		ID:               id,
		Email:            d.Get("email").(string),
		Description:      d.Get("description").(string),
		Pt0Severities:    SetToStringList(d, "pt0_severities"),
		SecureSeverities: SetToStringList(d, "secure_severities"),
		ManageSeverities: SetToStringList(d, "manage_severities"),
		ComplySeverities: SetToStringList(d, "comply_severities"),
		SystemSeverities: SetToStringList(d, "system_severities"),
	}
	return result
}
