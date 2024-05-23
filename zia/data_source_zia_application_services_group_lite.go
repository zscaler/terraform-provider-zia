package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/appservicegroups"
)

func dataSourceFWApplicationServicesGroupLite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWApplicationServicesGroupLiteRead,
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

func dataSourceFWApplicationServicesGroupLiteRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *appservicegroups.ApplicationServicesGroupLite
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for application services group: %s\n", name)
		res, err := zClient.appservicegroups.GetByName(name)
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
		return fmt.Errorf("couldn't find any device name '%s'", name)
	}

	return nil
}
