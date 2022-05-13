package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zia/gozscaler/devicegroups"
)

func dataSourceDevices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDevicesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the device.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The device name.",
			},
			"device_group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The device group type.",
			},
			"device_model": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The device model.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The operating system (OS).",
			},
			"os_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The operating system version.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The device group's description.",
			},
			"owner_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the device owner (i.e., user).",
			},
			"owner_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The device owner's user name.",
			},
		},
	}
}

func dataSourceDevicesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *devicegroups.Devices
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for device group id: %d\n", id)
		res, err := zClient.devicegroups.GetDevicesByID(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for device group name: %s\n", name)
		res, err := zClient.devicegroups.GetDevicesByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	model, _ := d.Get("device_model").(string)
	if resp == nil && model != "" {
		log.Printf("[INFO] Getting data for device model : %s\n", model)
		res, err := zClient.devicegroups.GetDevicesByModel(model)
		if err != nil {
			return err
		}
		resp = res
	}

	owner, _ := d.Get("owner_name").(string)
	if resp == nil && owner != "" {
		log.Printf("[INFO] Getting data for owner : %s\n", owner)
		res, err := zClient.devicegroups.GetDevicesByOwner(owner)
		if err != nil {
			return err
		}
		resp = res
	}

	osType, _ := d.Get("os_type").(string)
	if resp == nil && osType != "" {
		log.Printf("[INFO] Getting data for OS Type : %s\n", osType)
		res, err := zClient.devicegroups.GetDevicesByOSType(osType)
		if err != nil {
			return err
		}
		resp = res
	}

	osVersion, _ := d.Get("os_version").(string)
	if resp == nil && osVersion != "" {
		log.Printf("[INFO] Getting data for OS Version : %s\n", osVersion)
		res, err := zClient.devicegroups.GetDevicesByOSVersion(osVersion)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("device_group_type", resp.DeviceGroupType)
		_ = d.Set("device_model", resp.DeviceModel)
		_ = d.Set("os_type", resp.OSType)
		_ = d.Set("os_version", resp.OSVersion)
		_ = d.Set("description", resp.Description)
		_ = d.Set("owner_user_id", resp.OwnerUserId)
		_ = d.Set("owner_name", resp.OwnerName)

	} else {
		return fmt.Errorf("couldn't find any device name '%s' or id '%d'", name, id)
	}

	return nil
}
