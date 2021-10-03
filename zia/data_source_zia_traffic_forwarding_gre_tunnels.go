package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/gretunnels"
)

func dataSourceTrafficForwardingGreTunnels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTrafficForwardingGreTunnelsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_dest_vip": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"virtual_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_service_edge": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
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
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"secondary_dest_vip": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"virtual_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_service_edge": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
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
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"internal_ip_range": {
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
			"last_modification_time": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"within_country": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_unnumbered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceTrafficForwardingGreTunnelsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *gretunnels.GreTunnels
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for gre tunnel id: %d\n", id)
		res, err := zClient.gretunnels.GetGreTunnels(id)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("source_ip", resp.SourceIP)
		_ = d.Set("internal_ip_range", resp.InternalIpRange)
		_ = d.Set("last_modification_time", resp.LastModificationTime)
		_ = d.Set("within_country", resp.WithinCountry)
		_ = d.Set("comment", resp.Comment)
		_ = d.Set("ip_unnumbered", resp.IPUnnumbered)
		if err := d.Set("primary_dest_vip", flattenGrePrimaryDestVip(resp.PrimaryDestVip)); err != nil {
			return err
		}

		if err := d.Set("secondary_dest_vip", flattenGreSecondaryDestVip(resp.SecondaryDestVip)); err != nil {
			return err
		}

		if err := d.Set("managed_by", flattenGreManagedBy(resp.ManagedBy)); err != nil {
			return err
		}

		if err := d.Set("last_modified_by", flattenGreLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any gre tunnel with id '%d'", id)
	}

	return nil
}

func flattenGrePrimaryDestVip(primaryDestVip gretunnels.PrimaryDestVip) interface{} {
	return []map[string]interface{}{
		{
			"id":                   primaryDestVip.ID,
			"virtual_ip":           primaryDestVip.VirtualIP,
			"private_service_edge": primaryDestVip.PrivateServiceEdge,
			"datacenter":           primaryDestVip.Datacenter,
			"latitude":             primaryDestVip.Latitude,
			"longitude":            primaryDestVip.Longitude,
			"city":                 primaryDestVip.City,
			"country_code":         primaryDestVip.CountryCode,
			"region":               primaryDestVip.Region,
		},
	}
}

func flattenGreSecondaryDestVip(secondaryDestVip gretunnels.SecondaryDestVip) interface{} {
	return []map[string]interface{}{
		{
			"id":                   secondaryDestVip.ID,
			"virtual_ip":           secondaryDestVip.VirtualIP,
			"private_service_edge": secondaryDestVip.PrivateServiceEdge,
			"datacenter":           secondaryDestVip.Datacenter,
			"latitude":             secondaryDestVip.Latitude,
			"longitude":            secondaryDestVip.Longitude,
			"city":                 secondaryDestVip.City,
			"country_code":         secondaryDestVip.CountryCode,
			"region":               secondaryDestVip.Region,
		},
	}
}

func flattenGreManagedBy(managedBy gretunnels.ManagedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   managedBy.ID,
			"name": managedBy.Name,
		},
	}
}

func flattenGreLastModifiedBy(managedBy gretunnels.LastModifiedBy) interface{} {
	return []map[string]interface{}{
		{
			"id":   managedBy.ID,
			"name": managedBy.Name,
		},
	}
}
