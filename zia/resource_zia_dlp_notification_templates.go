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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_notification_templates"
)

func resourceDLPNotificationTemplates() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDLPNotificationTemplatesCreate,
		ReadContext:   resourceDLPNotificationTemplatesRead,
		UpdateContext: resourceDLPNotificationTemplatesUpdate,
		DeleteContext: resourceDLPNotificationTemplatesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("template_id", idInt)
				} else {
					resp, err := dlp_notification_templates.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("template_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for a DLP notification template",
			},
			"template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for a DLP notification template",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
				Description:  "The DLP notification template name",
			},
			"subject": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The Subject line that is displayed within the DLP notification email",
				ValidateDiagFunc: stringIsMultiLine,        // Validates that it's a valid multi-line string
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"attach_content": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "f set to true, the content that is violation is attached to the DLP notification email",
			},
			"plain_text_message": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The template for the plain text UTF-8 message body that must be displayed in the DLP notification email",
				ValidateDiagFunc: stringIsMultiLine,        // Validates that it's a valid multi-line string
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"html_message": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The template for the HTML message body that must be displayed in the DLP notification email",
				ValidateDiagFunc: stringIsMultiLine,        // Validates that it's a valid multi-line string
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, TLS will be enabled",
			},
		},
	}
}

func resourceDLPNotificationTemplatesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandDLPNotificationTemplates(d)
	log.Printf("[INFO] Creating zia dlp notification templates\n%+v\n", req)

	resp, _, err := dlp_notification_templates.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia dlp notification templates request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("template_id", resp.ID)

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

	return resourceDLPNotificationTemplatesRead(ctx, d, meta)
}

func resourceDLPNotificationTemplatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no DLP notification template id is set"))
	}
	resp, err := dlp_notification_templates.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing dlp notification template %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting dlp notification template :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("template_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("attach_content", resp.AttachContent)
	_ = d.Set("subject", normalizeMultiLineString(resp.Subject))
	_ = d.Set("plain_text_message", normalizeMultiLineString(resp.PlainTextMessage))
	_ = d.Set("html_message", normalizeMultiLineString(resp.HtmlMessage))
	_ = d.Set("tls_enabled", resp.TLSEnabled)

	return nil
}

func resourceDLPNotificationTemplatesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] dlp notification template not set: %v\n", id)
	}

	log.Printf("[INFO] Updating dlp notification template ID: %v\n", id)
	req := expandDLPNotificationTemplates(d)
	if _, err := dlp_notification_templates.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := dlp_notification_templates.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
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

	return resourceDLPNotificationTemplatesRead(ctx, d, meta)
}

func resourceDLPNotificationTemplatesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] dlp notification template ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp notification template ID: %v\n", (d.Id()))

	if _, err := dlp_notification_templates.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] dlp notification template deleted")
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

	return nil
}

func expandDLPNotificationTemplates(d *schema.ResourceData) dlp_notification_templates.DlpNotificationTemplates {
	id, _ := getIntFromResourceData(d, "template_id")
	result := dlp_notification_templates.DlpNotificationTemplates{
		ID:               id,
		Name:             d.Get("name").(string),
		AttachContent:    d.Get("attach_content").(bool),
		Subject:          unescapeTerraformVariables(d.Get("subject").(string)),
		PlainTextMessage: unescapeTerraformVariables(d.Get("plain_text_message").(string)),
		HtmlMessage:      unescapeTerraformVariables(d.Get("html_message").(string)),
		TLSEnabled:       d.Get("tls_enabled").(bool),
	}
	return result
}
