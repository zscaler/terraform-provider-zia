package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
)

func resourceFWIPDestinationGroups() *schema.Resource {
	return &schema.Resource{
		Create:   resourceFWIPDestinationGroupsCreate,
		Read:     resourceFWIPDestinationGroupsRead,
		Update:   resourceFWIPDestinationGroupsUpdate,
		Delete:   resourceFWIPDestinationGroupsDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifer for the destination IP group",
			},
			"ip_destination_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifer for the destination IP group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination IP group name",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs)",
				ValidateFunc: validation.StringInSlice([]string{
					"DSTN_IP",
					"DSTN_FQDN",
				}, false),
			},
			"addresses": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Destination IP addresses within the group",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional information about the destination IP group",
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"ip_categories": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Destination IP address URL categories. You can identify destinations based on the URL category of the domain.",
			},
			"countries": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "Destination IP address counties. You can identify destinations based on the location of a server.",
			},
		},
	}
}

func resourceFWIPDestinationGroupsCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandIPDestinationGroups(d)
	log.Printf("[INFO] Creating zia ip destination groups\n%+v\n", req)

	resp, err := zClient.ipdestinationgroups.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia ip destination groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("ip_destination_id", resp.ID)
	return resourceFWIPDestinationGroupsRead(d, m)
}

func resourceFWIPDestinationGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "ip_destination_id")
	if !ok {
		return fmt.Errorf("no ip destination groups id is set")
	}
	resp, err := zClient.ipdestinationgroups.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia ip destination groups %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting zia ip destination groups:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("ip_destination_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("type", resp.Type)
	_ = d.Set("addresses", resp.Addresses)
	_ = d.Set("description", resp.Description)
	_ = d.Set("ip_categories", resp.IPCategories)
	_ = d.Set("countries", resp.Countries)

	return nil
}

func resourceFWIPDestinationGroupsUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "ip_destination_id")
	if !ok {
		log.Printf("[ERROR] ip destination groups  ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia ip destination groups ID: %v\n", id)
	req := expandIPDestinationGroups(d)

	if _, _, err := zClient.ipdestinationgroups.Update(id, &req); err != nil {
		return err
	}

	return resourceFWIPDestinationGroupsRead(d, m)
}

func resourceFWIPDestinationGroupsDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "ip_destination_id")
	if !ok {
		log.Printf("[ERROR] ip destination groups ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia ip destination groups ID: %v\n", (d.Id()))

	if _, err := zClient.ipdestinationgroups.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] zia ip destination groups deleted")
	return nil
}

func expandIPDestinationGroups(d *schema.ResourceData) ipdestinationgroups.IPDestinationGroups {
	return ipdestinationgroups.IPDestinationGroups{
		Name:         d.Get("name").(string),
		Type:         d.Get("type").(string),
		Description:  d.Get("description").(string),
		Addresses:    SetToStringList(d, "addresses"),
		IPCategories: SetToStringList(d, "ip_categories"),
		Countries:    SetToStringList(d, "countries"),
	}
}
