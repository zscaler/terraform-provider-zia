package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zia/gozscaler/dlp_notification_templates"
)

func dataSourceDLPNotificationTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPNotificationTemplatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"attach_content": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"plain_test_message": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"html_message": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDLPNotificationTemplatesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *dlp_notification_templates.DlpNotificationTemplates
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp notifiation template id: %d\n", id)
		res, err := zClient.dlp_notification_templates.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp notifiation template: %s\n", name)
		res, err := zClient.dlp_notification_templates.GetByName(name)
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

	} else {
		return fmt.Errorf("couldn't find any dlp notification template with name '%s' or id '%d'", name, id)
	}

	return nil
}
