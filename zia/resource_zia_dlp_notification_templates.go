package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_notification_templates"
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

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("template_id", id)
				} else {
					resp, err := zClient.dlp_notification_templates.GetByName(id)
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_content": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"plain_text_message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"html_message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tls_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDLPNotificationTemplatesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandDLPNotificationTemplates(d)
	log.Printf("[INFO] Creating zia dlp notification templates\n%+v\n", req)

	resp, _, err := zClient.dlp_notification_templates.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia dlp notification templates request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("template_id", resp.ID)

	return resourceDLPNotificationTemplatesRead(d, m)
}

func resourceDLPNotificationTemplatesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		return fmt.Errorf("no DLP notification template id is set")
	}
	resp, err := zClient.dlp_notification_templates.Get(id)

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

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] dlp notification template not set: %v\n", id)
	}

	log.Printf("[INFO] Updating dlp notification template ID: %v\n", id)
	req := expandDLPNotificationTemplates(d)

	if _, _, err := zClient.dlp_notification_templates.Update(id, &req); err != nil {
		return err
	}

	return resourceDLPNotificationTemplatesRead(d, m)
}

func resourceDLPNotificationTemplatesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "template_id")
	if !ok {
		log.Printf("[ERROR] dlp notification template ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting dlp notification template ID: %v\n", (d.Id()))

	if _, err := zClient.dlp_notification_templates.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] dlp notification template deleted")
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
