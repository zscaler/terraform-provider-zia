package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
)

func dataSourceNetworkServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkServicesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_tcp_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"dest_tcp_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"src_udp_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"dest_udp_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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
	}
}

func dataSourceNetworkServicesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *networkservices.NetworkServices
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting network service group id: %d\n", id)
		res, err := zClient.networkservices.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network service group : %s\n", name)
		res, err := zClient.networkservices.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("tag", resp.Tag)
		_ = d.Set("type", resp.Type)
		_ = d.Set("description", resp.Description)
		_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

		if err := d.Set("src_tcp_ports", flattenSrcTCPPorts(resp.SrcTCPPorts)); err != nil {
			return err
		}

		if err := d.Set("dest_tcp_ports", flattenDestTCPPorts(resp.DestTCPPorts)); err != nil {
			return err
		}

		if err := d.Set("src_udp_ports", flattenSrcUDPPorts(resp.SrcUDPPorts)); err != nil {
			return err
		}

		if err := d.Set("dest_udp_ports", flattenDestUDPPorts(resp.DestUDPPorts)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any network service group with name '%s' or id '%d'", name, id)
	}

	return nil
}

func flattenSrcTCPPorts(service []networkservices.SrcTCPPorts) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}

func flattenDestTCPPorts(service []networkservices.DestTCPPorts) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}

func flattenSrcUDPPorts(service []networkservices.SrcUDPPorts) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}

func flattenDestUDPPorts(service []networkservices.DestUDPPorts) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}
