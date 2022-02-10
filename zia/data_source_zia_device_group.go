package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/devicegroups"
)

func dataSourceDeviceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDeviceGroupsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifer for the device group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The device group name.",
			},
			"group_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device group type.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device group's description.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operating system (OS).",
			},
			"predefined": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is a predefined device group. If this value is set to true, the group is predefined.",
			},
			"device_names": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The names of devices that belong to the device group. The device names are comma-separated.",
			},
			"device_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of devices within the group.",
			},
		},
	}
}

func dataSourceDeviceGroupsRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *devicegroups.DeviceGroups
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for device group id: %d\n", id)
		res, err := zClient.devicegroups.GetDeviceGroups(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for device group name: %s\n", name)
		res, err := zClient.devicegroups.GetDeviceGroupByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("group_type", resp.GroupType)
		_ = d.Set("description", resp.Description)
		_ = d.Set("os_type", resp.OSType)
		_ = d.Set("predefined", resp.Predefined)
		_ = d.Set("device_names", resp.DeviceNames)
		_ = d.Set("device_count", resp.DeviceCount)

	} else {
		return fmt.Errorf("couldn't find any device group name '%s' or id '%d'", name, id)
	}

	return nil
}
