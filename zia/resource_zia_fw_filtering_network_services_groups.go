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
	log.Printf("[INFO] Creating network service groups\n%+v\n", req)

	resp, err := zClient.networkservices.CreateNetworkServiceGroups(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia network service groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("network_service_group_id", resp.ID)
	return resourceFWNetworkServiceGroupsRead(d, m)
}

func resourceFWNetworkServiceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_group_id")
	if !ok {
		return fmt.Errorf("no network service groups id is set")
	}
	resp, err := zClient.networkservices.GetNetworkServiceGroups(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing zia network service groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting network service groups :\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("network_service_group_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)

	if err := d.Set("services", flattenServices(resp.Services)); err != nil {
		return err
	}

	return nil
}

func resourceFWNetworkServiceGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_group_id")
	if !ok {
		log.Printf("[ERROR] network service groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating network service groups ID: %v\n", id)
	req := expandNetworkServiceGroups(d)

	if _, _, err := zClient.networkservices.UpdateNetworkServiceGroups(id, &req); err != nil {
		return err
	}

	return resourceFWNetworkServiceGroupsRead(d, m)
}

func resourceFWNetworkServiceGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "network_service_group_id")
	if !ok {
		log.Printf("[ERROR] network service groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting network service groups ID: %v\n", (d.Id()))

	if _, err := zClient.networkservices.DeleteNetworkServiceGroups(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] network service groups deleted")
	return nil
}

func expandNetworkServiceGroups(d *schema.ResourceData) networkservices.NetworkServiceGroups {
	id, _ := getIntFromResourceData(d, "network_service_group_id")
	result := networkservices.NetworkServiceGroups{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		// Need to expand Services as part of Network Service Groups
	}

	return result
}
