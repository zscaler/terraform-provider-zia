package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/applicationservices"
)

func dataSourceFWApplicationServicesLite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWApplicationServicesLiteRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the application service.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the application service.",
			},
			"name_l10n_tag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceFWApplicationServicesLiteRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.applicationservices

	var resp *applicationservices.ApplicationServicesLite
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for application service name: %s\n", name)
		res, err := applicationservices.GetByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("name_l10n_tag", resp.NameL10nTag)

	} else {
		return fmt.Errorf("couldn't find any application service name '%s'", name)
	}

	return nil
}
