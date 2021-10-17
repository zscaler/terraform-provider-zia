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

func resourceFWNetworkServiceGroups() *schema.Resource {
	return &schema.Resource{
		Create:   resourceFWNetworkServiceGroupsCreate,
		Read:     resourceFWNetworkServiceGroupsRead,
		Update:   resourceFWNetworkServiceGroupsUpdate,
		Delete:   resourceFWNetworkServiceGroupsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_service_group_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"src_tcp_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"end": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"dest_tcp_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"end": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"src_udp_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"end": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"dest_udp_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"end": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
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
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_name_l10n_tag": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFWNetworkServiceGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandNetworkServiceGroups(d)
	log.Printf("[INFO] Creating network services\n%+v\n", req)

	resp, err := zClient.networkservices.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia network services request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("network_service_id", resp.ID)
	return resourceFWNetworkServiceGroupsRead(d, m)
}

func resourceFWNetworkServiceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		return fmt.Errorf("no network services id is set")
	}
	resp, err := zClient.networkservices.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
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

	return nil
}

func resourceFWNetworkServiceGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_id")
	if !ok {
		log.Printf("[ERROR] network service ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network service ID: %v\n", id)
	req := expandNetworkServiceGroups(d)

	if _, _, err := zClient.networkservices.Update(id, &req); err != nil {
		return err
	}

	return resourceFWNetworkServiceGroupsRead(d, m)
}

func resourceFWNetworkServiceGroupsDelete(d *schema.ResourceData, m interface{}) error {
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

func expandNetworkServiceGroups(d *schema.ResourceData) networkservices.NetworkServices {
	id, _ := getIntFromResourceData(d, "network_service_id")
	result := networkservices.NetworkServices{
		ID:            id,
		Name:          d.Get("name").(string),
		Tag:           d.Get("tag").(string),
		Description:   d.Get("description").(string),
		Type:          d.Get("type").(string),
		IsNameL10nTag: d.Get("is_name_l10n_tag").(bool),
	}
	srcTcpPorts := expandSrcTcpPorts(d)
	if srcTcpPorts != nil {
		result.SrcTCPPorts = srcTcpPorts
	}

	destTcpPorts := expandDestTCPPorts(d)
	if destTcpPorts != nil {
		result.DestTCPPorts = destTcpPorts
	}

	SrcUdpPorts := expandSrcUdpPorts(d)
	if SrcUdpPorts != nil {
		result.SrcUDPPorts = SrcUdpPorts
	}

	DestUdpPorts := expandDestUDPPorts(d)
	if DestUdpPorts != nil {
		result.DestUDPPorts = DestUdpPorts
	}

	return result
}
