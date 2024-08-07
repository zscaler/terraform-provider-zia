package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_notification_templates"
)

func resourceDLPNotificationTemplates() *schema.Resource {
	return &schema.Resource{
		Create: resourceDLPNotificationTemplatesCreate,
		Read:   resourceDLPNotificationTemplatesRead,
		Update: resourceDLPNotificationTemplatesUpdate,
		Delete: resourceDLPNotificationTemplatesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)
				service := zClient.dlp_notification_templates

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("template_id", idInt)
				} else {
					resp, err := dlp_notification_templates.GetByName(service, id)
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Subject line that is displayed within the DLP notification email",
			},
			"attach_content": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "f set to true, the content that is violation is attached to the DLP notification email",
			},
			"plain_text_message": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The template for the plain text UTF-8 message body that must be displayed in the DLP notification email",
			},
			"html_message": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The template for the HTML message body that must be displayed in the DLP notification email",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, TLS will be enabled",
			},
		},
	}
}

func resourceDLPNotificationTemplatesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_notification_templates

	req := expandDLPNotificationTemplates(d)
	log.Printf("[INFO] Creating zia dlp notification templates\n%+v\n", req)

	resp, _, err := dlp_notification_templates.Create(service, &req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia dlp notification templates request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("template_id", resp.ID)

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

	return resourceDLPNotificationTemplatesRead(d, m)
}

func resourceDLPNotificationTemplatesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_notification_templates

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		return fmt.Errorf("no DLP notification template id is set")
	}
	resp, err := dlp_notification_templates.Get(service, id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing dlp notification template %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting dlp notification template :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("template_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("subject", resp.Subject)
	_ = d.Set("attach_content", resp.AttachContent)
	_ = d.Set("plain_text_message", resp.PlainTextMessage)
	_ = d.Set("html_message", resp.HtmlMessage)
	_ = d.Set("tls_enabled", resp.TLSEnabled)

	return nil
}

func resourceDLPNotificationTemplatesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_notification_templates

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] dlp notification template not set: %v\n", id)
	}

	log.Printf("[INFO] Updating dlp notification template ID: %v\n", id)
	req := expandDLPNotificationTemplates(d)
	if _, err := dlp_notification_templates.Get(service, id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := dlp_notification_templates.Update(service, id, &req); err != nil {
		return err
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

	return resourceDLPNotificationTemplatesRead(d, m)
}

func resourceDLPNotificationTemplatesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_notification_templates

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] dlp notification template ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp notification template ID: %v\n", (d.Id()))

	if _, err := dlp_notification_templates.Delete(service, id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] dlp notification template deleted")
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

	return nil
}

func expandDLPNotificationTemplates(d *schema.ResourceData) dlp_notification_templates.DlpNotificationTemplates {
	id, _ := getIntFromResourceData(d, "template_id")
	result := dlp_notification_templates.DlpNotificationTemplates{
		ID:               id,
		Name:             d.Get("name").(string),
		Subject:          d.Get("subject").(string),
		AttachContent:    d.Get("attach_content").(bool),
		PlainTextMessage: d.Get("plain_text_message").(string),
		HtmlMessage:      d.Get("html_message").(string),
		TLSEnabled:       d.Get("tls_enabled").(bool),
	}
	return result
}
