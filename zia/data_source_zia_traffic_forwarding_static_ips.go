package zia

/*
import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
)

func dataSourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingStaticIPRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"geo_override": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_modification_time": {
				Type:     schema.TypeInt,
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
			"last_modified_by": {
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
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTrafficForwardingStaticIPRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *staticips.StaticIP
	/*
		idObj, idSet := d.GetOk("id")
		id, idIsInt := idObj.(int)
		if idSet && idIsInt && id > 0 {
			log.Printf("[INFO] Getting data for gre tunnel id: %d\n", id)
			res, err := zClient.staticips.GetStaticIP(id)
			if err != nil {
				return err
			}
			resp = res
		}

	ipaddress, _ := d.Get("ip_address").(string)
	if resp == nil && ipaddress != "" {
		log.Printf("[INFO] Getting data for location name: %s\n", ipaddress)
		res, err := zClient.staticips.GetStaticByIP(ipaddress)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("ip_address", resp.IpAddress)
		_ = d.Set("geo_override", resp.GeoOverride)
		_ = d.Set("latitude", resp.Latitude)
		_ = d.Set("longitude", resp.Longitude)
		_ = d.Set("routable_ip", resp.RoutableIP)
		_ = d.Set("comment", resp.Comment)
		_ = d.Set("last_modification_time", resp.LastModificationTime)

		if err := d.Set("managed_by", flattenStaticManagedBy(resp.ManagedBy)); err != nil {
			return err
		}

		if err := d.Set("last_modified_by", flattenStaticLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any ip address with id '%s'", ipaddress)
	}

	return nil
}

func flattenStaticManagedBy(managedBy staticips.ManagedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   managedBy.ID,
			"name": managedBy.Name,
		},
	}
}

func flattenStaticLastModifiedBy(managedBy staticips.LastModifiedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   managedBy.ID,
			"name": managedBy.Name,
		},
	}
}
*/
