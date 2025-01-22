package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/devicegroups"
)

func dataSourceDeviceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceGroupsRead,
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

func dataSourceDeviceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *devicegroups.DeviceGroups
	name, ok := d.Get("name").(string)
	if ok && name != "" {
		log.Printf("[INFO] Getting data for device group name: %s\n", name)
		res, err := devicegroups.GetDeviceGroupByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
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
		return diag.FromErr(fmt.Errorf("couldn't find any device group name '%s'", name))
	}

	return nil
}
