package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/cloudbrowserisolation"
)

func dataSourceCBIProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCBIProfileRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The universally unique identifier (UUID) for the browser isolation profile",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the browser isolation profile",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The browser isolation profile URL",
			},
			"default_profile": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "(Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field",
			},
		},
	}
}

func dataSourceCBIProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.user_authentication_settings

	var resp *cloudbrowserisolation.IsolationProfile
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for cloud browser isolation profile name %s\n", name)
		res, err := cloudbrowserisolation.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}
	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("url", resp.URL)
		_ = d.Set("default_profile", resp.DefaultProfile)

	} else {
		return fmt.Errorf("couldn't find any cloud browser isolation profile with name '%s'", name)
	}

	return nil
}
