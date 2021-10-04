package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
)

func resourceTrafficForwardingVPNCredentials() *schema.Resource {
	return &schema.Resource{
		Create:   resourceTrafficForwardingVPNCredentialsCreate,
		Read:     resourceTrafficForwardingVPNCredentialsRead,
		Update:   resourceTrafficForwardingVPNCredentialsUpdate,
		Delete:   resourceTrafficForwardingVPNCredentialsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"vpn_credental_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pre_shared_key": {
				Type:     schema.TypeString,
				Optional: true,
				// Sensitive: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			/*
				"managed_by": {
					Type:     schema.TypeSet,
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
							"extensions": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"location": {
					Type:     schema.TypeSet,
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
							"extensions": {
								Type:     schema.TypeMap,
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			*/
		},
	}
}

func resourceTrafficForwardingVPNCredentialsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandVPNCredentials(d)
	log.Printf("[INFO] Creating zia vpn credentials\n%+v\n", req)

	resp, _, err := zClient.vpncredentials.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia vpn credentials request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("vpn_credental_id", resp.ID)

	return resourceTrafficForwardingVPNCredentialsRead(d, m)
}

func resourceTrafficForwardingVPNCredentialsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "vpn_credental_id")
	if !ok {
		return fmt.Errorf("no Traffic Forwarding zia vpn credentials id is set")
	}
	resp, err := zClient.vpncredentials.Get(id)

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing vpn credentials %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting vpn credentials:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("vpn_credental_id", resp.ID)
	_ = d.Set("type", resp.Type)
	_ = d.Set("fqdn", resp.FQDN)
	_ = d.Set("pre_shared_key", resp.PreSharedKey)
	_ = d.Set("comments", resp.Comments)

	return nil
}

func resourceTrafficForwardingVPNCredentialsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "vpn_credental_id")
	if !ok {
		log.Printf("[ERROR] vpn credentials ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating vpn credentials ID: %v\n", id)
	req := expandVPNCredentials(d)

	if _, _, err := zClient.vpncredentials.Update(id, &req); err != nil {
		return err
	}

	return resourceTrafficForwardingVPNCredentialsRead(d, m)
}

func resourceTrafficForwardingVPNCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "vpn_credental_id")
	if !ok {
		log.Printf("[ERROR] vpn credentials ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting vpn credentials ID: %v\n", (d.Id()))

	if err := zClient.vpncredentials.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] vpn credentials deleted")
	return nil
}

func expandVPNCredentials(d *schema.ResourceData) vpncredentials.VPNCredentials {
	id, _ := getIntFromResourceData(d, "vpn_credental_id")
	result := vpncredentials.VPNCredentials{
		ID:           id,
		Type:         d.Get("type").(string),
		FQDN:         d.Get("fqdn").(string),
		PreSharedKey: d.Get("pre_shared_key").(string),
		Comments:     d.Get("comments").(string),
	}

	/*
		location := expandLocation(d)
		if location != nil {
			result.Location = location
		}
		managedBy := expandManagedBy(d)
		if managedBy != nil {
			result.ManagedBy = managedBy
		}
	*/
	return result
}

/*
func expandLocation(d *schema.ResourceData) *vpncredentials.Location {
	locationObj, ok := d.GetOk("location")
	if !ok {
		return nil
	}
	locations, ok := locationObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(locations.List()) > 0 {
		locationObj := locations.List()[0]
		location, ok := locationObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &vpncredentials.Location{
			ID:         location["id"].(int),
			Name:       location["name"].(string),
			Extensions: location["extensions"].(map[string]interface{}),
		}
	}
	return nil
}

func expandManagedBy(d *schema.ResourceData) *vpncredentials.ManagedBy {
	managedByObj, ok := d.GetOk("managed_by")
	if !ok {
		return nil
	}
	managedBy, ok := managedByObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(managedBy.List()) > 0 {
		managedByObj := managedBy.List()[0]
		managedBy, ok := managedByObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &vpncredentials.ManagedBy{
			ID:         managedBy["id"].(int),
			Name:       managedBy["name"].(string),
			Extensions: managedBy["extensions"].(map[string]interface{}),
		}
	}
	return nil
}
*/
