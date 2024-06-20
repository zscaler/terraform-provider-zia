package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_notification_templates"
)

func dataSourceDLPNotificationTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPNotificationTemplatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier for a DLP notification template",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The DLP notification template name",
			},
			"subject": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Subject line that is displayed within the DLP notification email",
			},
			"attach_content": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "f set to true, the content that is violation is attached to the DLP notification email",
			},
			"plain_text_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template for the plain text UTF-8 message body that must be displayed in the DLP notification email",
			},
			"html_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template for the HTML message body that must be displayed in the DLP notification email",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, TLS will be enabled",
			},
		},
	}
}

func dataSourceDLPNotificationTemplatesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.dlp_notification_templates

	var resp *dlp_notification_templates.DlpNotificationTemplates
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp notifiation template id: %d\n", id)
		res, err := dlp_notification_templates.Get(service, id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp notifiation template: %s\n", name)
		res, err := dlp_notification_templates.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("subject", resp.Subject)
		_ = d.Set("attach_content", resp.AttachContent)
		_ = d.Set("plain_text_message", resp.PlainTextMessage)
		_ = d.Set("html_message", resp.HtmlMessage)
		_ = d.Set("tls_enabled", resp.TLSEnabled)

	} else {
		return fmt.Errorf("couldn't find any dlp notification template with name '%s' or id '%d'", name, id)
	}

	return nil
}
