package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
)

func dataSourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingStaticIPRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"geo_override": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routableIP": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modification_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"last_modified_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"Comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTrafficForwardingStaticIPRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.staticips.GetStaticIP(d.Id())
	if err != nil {
		return nil
	}

	d.SetId(resp.ID)
	_ = d.Set("ip_address", resp.IpAddress)
	_ = d.Set("geo_override", resp.GeoOverride)
	_ = d.Set("latitude", resp.Latitude)
	_ = d.Set("longitude", resp.Longitude)
	_ = d.Set("routable_ip", resp.RoutableIP)
	_ = d.Set("last_nodification_time", resp.LastModificationTime)
	_ = d.Set("managed_by", flattenManagedBy(resp))
	_ = d.Set("managed_by", flattenLastModifiedBy(resp))
	return nil
}

func flattenManagedBy(managedBy *staticips.StaticIP) []interface{} {
	managed := make([]interface{}, len(managedBy.ManagedBy))
	for i, managedItem := range managedBy.ManagedBy {
		managed[i] = map[string]interface{}{

			"id":         managedItem.ID,
			"name":       managedItem.Name,
			"extensions": managedItem.Extensions,
		}
	}

	return managed
}

func flattenLastModifiedBy(lastModifiedBy *staticips.StaticIP) []interface{} {
	lastModified := make([]interface{}, len(lastModifiedBy.LastModifiedBy))
	for i, lastModifiedByItem := range lastModifiedBy.LastModifiedBy {
		lastModified[i] = map[string]interface{}{

			"id":         lastModifiedByItem.ID,
			"name":       lastModifiedByItem.Name,
			"extensions": lastModifiedByItem.Extensions,
		}
	}

	return lastModified
}
