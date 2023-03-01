package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/applicationservicesgroup"
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

	var resp *applicationservicesgroup.ApplicationServicesGroupLite
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for application services group id: %d\n", id)
		res, err := zClient.applicationservicesgroup.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for application services group: %s\n", name)
		res, err := zClient.applicationservicesgroup.GetByName(name)
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
		return fmt.Errorf("couldn't find any device name '%s' or id '%d'", name, id)
	}

	return nil
}
