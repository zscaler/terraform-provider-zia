package zia

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/virtualipaddress"
)

func dataSourceTrafficForwardingGreVipRecommendedList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficForwardingGreVipRecommendedListRead,
		Schema: map[string]*schema.Schema{
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"required_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"within_country_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_private_service_edge": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"include_current_vips": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"subcloud": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"virtual_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_service_edge": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"city": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"latitude": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
						"longitude": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTrafficForwardingGreVipRecommendedListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Gather parameters from schema
	sourceIP := d.Get("source_ip").(string)
	requiredCount := d.Get("required_count").(int)

	// Initialize the list of options
	var options []func(*url.Values)

	// Add the sourceIP as an option
	if sourceIP != "" {
		options = append(options, virtualipaddress.WithSourceIP(sourceIP))
	}

	if v, ok := d.GetOk("routable_ip"); ok {
		options = append(options, virtualipaddress.WithRoutableIP(v.(bool)))
	}
	if v, ok := d.GetOk("within_country_only"); ok {
		options = append(options, virtualipaddress.WithWithinCountryOnly(v.(bool)))
	}
	if v, ok := d.GetOk("include_private_service_edge"); ok {
		options = append(options, virtualipaddress.WithIncludePrivateServiceEdge(v.(bool)))
	}
	if v, ok := d.GetOk("include_current_vips"); ok {
		options = append(options, virtualipaddress.WithIncludeCurrentVips(v.(bool)))
	}
	if v, ok := d.GetOk("latitude"); ok {
		options = append(options, virtualipaddress.WithLatitude(v.(float64)))
	}
	if v, ok := d.GetOk("longitude"); ok {
		options = append(options, virtualipaddress.WithLongitude(v.(float64)))
	}
	if v, ok := d.GetOk("subcloud"); ok {
		options = append(options, virtualipaddress.WithSubcloud(v.(string)))
	}

	// Call the new function with the options
	resp, err := virtualipaddress.GetVIPRecommendedList(ctx, service, options...)
	if err != nil {
		return diag.FromErr(err)
	}

	// Trim the list to the required count, if necessary
	if len(*resp) > requiredCount {
		*resp = (*resp)[:requiredCount]
	}

	d.SetId(sourceIP)
	_ = d.Set("list", flattenVIPList(*resp))

	return nil
}

func flattenVIPList(list []virtualipaddress.GREVirtualIPList) []interface{} {
	result := make([]interface{}, len(list))
	for i, vip := range list {
		result[i] = map[string]interface{}{
			"id":                   vip.ID,
			"virtual_ip":           vip.VirtualIp,
			"private_service_edge": vip.PrivateServiceEdge,
			"datacenter":           vip.DataCenter,
			"city":                 vip.City,
			"region":               vip.Region,
			"latitude":             vip.Latitude,
			"longitude":            vip.Longitude,
		}
	}
	return result
}
