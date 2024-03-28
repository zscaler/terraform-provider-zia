package zia

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/user_authentication_settings"
)

func resourceAuthSettingsUrls() *schema.Resource {
	return &schema.Resource{
		Read:          resourceAuthSettingsUrlsRead,
		Create:        resourceAuthSettingsUrlsCreate,
		Update:        resourceAuthSettingsUrlsUpdate,
		DeleteContext: resourceFuncNoOp,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				urls, err := zClient.user_authentication_settings.Get()
				if err != nil {
					return nil, fmt.Errorf("error fetching urls from exception list: %s", err)
				}

				if urls != nil {
					if err := d.Set("urls", urls.URLs); err != nil {
						return nil, fmt.Errorf("error setting urls: %s", err)
					}
				}

				d.SetId("all_urls")

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				MaxItems: 25000,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAuthSettingsUrlsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	res, err := zClient.user_authentication_settings.Get()
	if err != nil {
		return err
	}
	d.SetId("exempted_urls")
	_ = d.Set("urls", res.URLs)
	return nil
}

func expandAuthSettingsUrls(d *schema.ResourceData) user_authentication_settings.ExemptedUrls {
	return user_authentication_settings.ExemptedUrls{
		URLs: SetToStringList(d, "urls"),
	}
}

func resourceAuthSettingsUrlsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	urls := expandAuthSettingsUrls(d)
	_, err := zClient.user_authentication_settings.Update(urls)
	if err != nil {
		return err
	}
	d.SetId("exempted_urls")
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

	return resourceAuthSettingsUrlsRead(d, m)
}

func resourceAuthSettingsUrlsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	urls := expandAuthSettingsUrls(d)

	_, err := zClient.user_authentication_settings.Update(urls)
	if err != nil {
		return err
	}

	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return resourceAuthSettingsUrlsRead(d, m)
}
