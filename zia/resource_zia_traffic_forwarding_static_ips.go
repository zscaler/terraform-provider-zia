package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
)

func resourceTrafficForwardingStaticIP() *schema.Resource {
	return &schema.Resource{
		Create:   resourceTrafficForwardingStaticIPCreate,
		Read:     resourceTrafficForwardingStaticIPRead,
		Update:   resourceTrafficForwardingStaticIPUpdate,
		Delete:   resourceTrafficForwardingStaticIPDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// "ip_address": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			"geo_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"routable_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			// "managed_by": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeInt,
			// 				Computed: true,
			// 			},
			// 			"extensions": {
			// 				Type:     schema.TypeMap,
			// 				Computed: true,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			// "last_modified_by": {
			// 	Type:     schema.TypeList,
			// 	Computed: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:     schema.TypeInt,
			// 				Computed: true,
			// 			},
			// 			"extensions": {
			// 				Type:     schema.TypeMap,
			// 				Computed: true,
			// 				Elem: &schema.Schema{
			// 					Type: schema.TypeString,
			// 				},
			// 			},
			// 		},
			// 	},
			// },
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTrafficForwardingStaticIPCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandTrafficForwardingStaticIP(d)
	log.Printf("[INFO] Creating zia static ip\n%+v\n", req)

	resp, _, err := zClient.staticips.CreateStaticIP(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia static ip request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceTrafficForwardingStaticIPRead(d, m)
}

func resourceTrafficForwardingStaticIPRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.staticips.GetStaticIP(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing location management %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting location management:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))

	return nil
}

func resourceTrafficForwardingStaticIPUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating location management ID: %v\n", id)
	req := expandTrafficForwardingStaticIP(d)

	if _, _, err := zClient.staticips.UpdateStaticIP(id, &req); err != nil {
		return err
	}

	return resourceTrafficForwardingStaticIPRead(d, m)
}

func resourceTrafficForwardingStaticIPDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting static ip ID: %v\n", (d.Id()))

	if _, err := zClient.staticips.DeleteStaticIP(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] static ip deleted")
	return nil
}

func expandTrafficForwardingStaticIP(d *schema.ResourceData) staticips.StaticIP {
	return staticips.StaticIP{
		//IpAddress:            d.Get("ip_address").(string),
		IpAddress:   expandStringInSlice(d, "ip_address"),
		GeoOverride: d.Get("geo_override").(bool),
		Latitude:    d.Get("latitude").(int),
		Longitude:   d.Get("longitude").(int),
		RoutableIP:  d.Get("routable_ip").(bool),
		Comment:     d.Get("comment").(string),
		// ManagedBy:      expandManagedBy(d),
		// LastModifiedBy: expandLastModifiedBy(d),
	}
}

/*
func expandManagedBy(d *schema.ResourceData) staticips.ManagedBy {
	managedBy := staticips.ManagedBy{
		ID:         d.Get("id").(int),
		Extensions: d.Get("extensions").(map[string]interface{}),
	}
	return managedBy
}

func expandLastModifiedBy(d *schema.ResourceData) staticips.LastModifiedBy {
	lastModifiedBy := staticips.LastModifiedBy{
		ID:         d.Get("id").(int),
		Extensions: d.Get("extensions").(map[string]interface{}),
	}
	return lastModifiedBy
}
*/
func expandStringInSlice(d *schema.ResourceData, key string) []string {
	applicationSegments := d.Get(key).([]interface{})
	applicationSegmentList := make([]string, len(applicationSegments))
	for i, applicationSegment := range applicationSegments {
		applicationSegmentList[i] = applicationSegment.(string)
	}

	return applicationSegmentList
}
