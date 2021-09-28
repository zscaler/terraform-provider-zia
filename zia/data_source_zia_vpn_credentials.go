package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
)

func dataSourceVPNCredentials() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNCredentialsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pre_shared_key": {
				Type:     schema.TypeString,
				Computed: true,
				// Sensitive: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_by": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"location": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVPNCredentialsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *vpncredentials.VPNCredentials
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting data for vpn credential id: %d\n", id)
		res, err := zClient.vpncredentials.GetVPNCredentials(id)
		if err != nil {
			return err
		}
		resp = res
	}
	fqdn, _ := d.Get("fqdn").(string)
	if resp == nil && fqdn != "" {
		log.Printf("[INFO] Getting data for vpn credential fqdn: %s\n", fqdn)
		res, err := zClient.vpncredentials.GetVPNCredentialsByFQDN(fqdn)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("type", resp.Type)
		_ = d.Set("fqdn", resp.FQDN)
		_ = d.Set("pre_shared_key", resp.PreSharedKey)
		_ = d.Set("comments", resp.Comments)
		if err := d.Set("location", flattenLocation(resp.Location)); err != nil {
			return err
		}

		if err := d.Set("managed_by", flattenVPNCredentialManagedBy(resp.ManagedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any vpn credentials with fqdn '%s' or id '%d'", fqdn, id)
	}

	return nil
}

// Want to simplify this. This flattening function will be used in multiple places.
func flattenLocation(location vpncredentials.Location) interface{} {
	return []map[string]interface{}{
		{
			"id":   location.ID,
			"name": location.Name,
		},
	}
}

// Want to simplify this. This flattening function will be used in multiple places.
func flattenVPNCredentialManagedBy(managedBy vpncredentials.ManagedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   managedBy.ID,
			"name": managedBy.Name,
		},
	}
}
