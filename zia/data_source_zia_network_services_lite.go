package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
)

func dataSourceNetworkServicesLite() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkServicesLiteRead,
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
		},
	}
}

func dataSourceNetworkServicesLiteRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *networkservices.NetworkServiceGroups
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting network service group id: %d\n", id)
		res, err := zClient.networkservices.GetNetworkServiceGroupsLite(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting network service group : %s\n", name)
		res, err := zClient.networkservices.GetNetworkServiceGroupsLiteByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("description", resp.Description)

		if err := d.Set("src_tcp_ports", flattenSrcTCPPorts(resp)); err != nil {
			return err
		}

		if err := d.Set("dest_tcp_ports", flattenDestTCPPorts(resp)); err != nil {
			return err
		}

		if err := d.Set("src_udp_ports", flattenSrcUDPPorts(resp)); err != nil {
			return err
		}

		if err := d.Set("dest_udp_ports", flattenDestUDPPorts(resp)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any network service group with name '%s' or id '%d'", name, id)
	}

	return nil
}

func flattenSrcTCPPorts(service []networkservices.NetworkServices) []interface{} {
	services := make([]interface{}, len(service))
	for i, val := range service {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}

func flattenDestTCPPorts(service *networkservices.NetworkServices) []interface{} {
	services := make([]interface{}, len(service.DestTCPPorts))
	for i, val := range service.DestTCPPorts {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}

func flattenSrcUDPPorts(service *networkservices.NetworkServices) []interface{} {
	services := make([]interface{}, len(service.SrcUDPPorts))
	for i, val := range service.SrcUDPPorts {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}

func flattenDestUDPPorts(service *networkservices.NetworkServices) []interface{} {
	services := make([]interface{}, len(service.DestUDPPorts))
	for i, val := range service.DestUDPPorts {
		services[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}

	return services
}
