package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
)

func resourceFWNetworkServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkServicesCreate,
		Read:   resourceNetworkServicesRead,
		Update: resourceNetworkServicesUpdate,
		Delete: resourceNetworkServicesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("network_service_id", id)
				} else {
					resp, err := zClient.networkservices.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						d.Set("network_service_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_service_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag":            getCloudFirewallNwServicesTag(),
			"src_tcp_ports":  resourceNetworkPortsSchema("src tcp ports"),
			"dest_tcp_ports": resourceNetworkPortsSchema("dest tcp ports"),
			"src_udp_ports":  resourceNetworkPortsSchema("src udp ports"),
			"dest_udp_ports": resourceNetworkPortsSchema("dest udp ports"),
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"STANDARD",
					"PREDEFINED",
					"CUSTOM",
				}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"is_name_l10n_tag": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNetworkServicesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandNetworkServices(d)
	log.Printf("[INFO] Creating network services\n%+v\n", req)

	resp, err := zClient.networkservices.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia network services request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("network_service_id", resp.ID)
	return resourceNetworkServicesRead(d, m)
}

func resourceNetworkServicesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		return fmt.Errorf("no network services id is set")
	}
	resp, err := zClient.networkservices.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia network services %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting network services :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("network_service_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("tag", resp.Tag)
	_ = d.Set("description", resp.Description)
	_ = d.Set("type", resp.Type)
	_ = d.Set("is_name_l10n_tag", resp.IsNameL10nTag)

	if err := d.Set("src_tcp_ports", flattenNetwordPorts(resp.SrcTCPPorts)); err != nil {
		return err
	}
	if err := d.Set("dest_tcp_ports", flattenNetwordPorts(resp.DestTCPPorts)); err != nil {
		return err
	}

	if err := d.Set("src_udp_ports", flattenNetwordPorts(resp.SrcUDPPorts)); err != nil {
		return err
	}

	if err := d.Set("dest_udp_ports", flattenNetwordPorts(resp.DestUDPPorts)); err != nil {
		return err
	}

	return nil
}

func resourceNetworkServicesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		log.Printf("[ERROR] network service ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network service ID: %v\n", id)
	req := expandNetworkServices(d)

	if _, _, err := zClient.networkservices.Update(id, &req); err != nil {
		return err
	}

	return resourceNetworkServicesRead(d, m)
}

func resourceNetworkServicesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		log.Printf("[ERROR] network service id ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting network service ID: %v\n", (d.Id()))

	if _, err := zClient.networkservices.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] network service deleted")
	return nil
}

func expandNetworkServices(d *schema.ResourceData) networkservices.NetworkServices {
	id, _ := getIntFromResourceData(d, "network_service_id")
	result := networkservices.NetworkServices{
		ID:            id,
		Name:          d.Get("name").(string),
		Tag:           d.Get("tag").(string),
		Description:   d.Get("description").(string),
		Type:          d.Get("type").(string),
		IsNameL10nTag: d.Get("is_name_l10n_tag").(bool),
	}
	srcTcpPorts := expandNetworkPorts(d, "src_tcp_ports")
	if srcTcpPorts != nil {
		result.SrcTCPPorts = srcTcpPorts
	}

	destTcpPorts := expandNetworkPorts(d, "dest_tcp_ports")
	if destTcpPorts != nil {
		result.DestTCPPorts = destTcpPorts
	}

	SrcUdpPorts := expandNetworkPorts(d, "src_udp_ports")
	if SrcUdpPorts != nil {
		result.SrcUDPPorts = SrcUdpPorts
	}

	DestUdpPorts := expandNetworkPorts(d, "dest_udp_ports")
	if DestUdpPorts != nil {
		result.DestUDPPorts = DestUdpPorts
	}

	return result
}
