package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
)

func dataSourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingStaticIPRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"geo_override": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"city": {
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
			"comment": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Additional information about this static IP address",
			},
		},
	}
}

func dataSourceTrafficForwardingStaticIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *staticips.StaticIP
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for static ip id: %d\n", id)
		res, err := staticips.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	ipAddress, _ := d.Get("ip_address").(string)
	if resp == nil && ipAddress != "" {
		log.Printf("[INFO] Getting data for static ip : %s\n", ipAddress)
		res, err := staticips.GetByIPAddress(ctx, service, ipAddress)
		if err != nil {
			return diag.FromErr(err)
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
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenStaticLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("city", flattenCity(resp.City)); err != nil {
			return diag.FromErr(err)
		}

	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any static ip address with id '%d'", id))
	}

	return nil
}

func flattenStaticManagedBy(managedBy *staticips.ManagedBy) interface{} {
	if managedBy == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":         managedBy.ID,
			"name":       managedBy.Name,
			"extensions": managedBy.Extensions,
		},
	}
}

func flattenStaticLastModifiedBy(lastModifiedBy *staticips.LastModifiedBy) interface{} {
	if lastModifiedBy == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":         lastModifiedBy.ID,
			"name":       lastModifiedBy.Name,
			"extensions": lastModifiedBy.Extensions,
		},
	}
}

func flattenCity(city *staticips.City) interface{} {
	if city == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":   city.ID,
			"name": city.Name,
		},
	}
}
