package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservicegroups"
)

func dataSourceFWNetworkServiceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFWNetworkServiceGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_name_l10n_tag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFWNetworkServiceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.networkservicegroups

	var resp *networkservicegroups.NetworkServiceGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting network service group id: %d\n", id)
		res, err := networkservicegroups.GetNetworkServiceGroups(service, id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network service group : %s\n", name)
		res, err := networkservicegroups.GetNetworkServiceGroupsByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)

		if err := d.Set("services", flattenServices(resp.Services)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any network service group with name '%s' or id '%d'", name, id)
	}

	return nil
}

func flattenServices(service []networkservicegroups.Services) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"id":               val.ID,
			"name":             val.Name,
			"description":      val.Description,
			"is_name_l10n_tag": val.IsNameL10nTag,
		}
	}

	return services
}
