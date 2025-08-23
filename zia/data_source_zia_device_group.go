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
				Description: "The device group name. If not provided, all device groups will be returned.",
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
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all device groups when no name is specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifer for the device group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceDeviceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	name, hasName := d.Get("name").(string)

	if hasName && name != "" {
		// Get specific device group by name
		log.Printf("[INFO] Getting data for device group name: %s\n", name)
		resp, err := devicegroups.GetDeviceGroupByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("group_type", resp.GroupType)
		_ = d.Set("description", resp.Description)
		_ = d.Set("os_type", resp.OSType)
		_ = d.Set("predefined", resp.Predefined)
		_ = d.Set("device_names", resp.DeviceNames)
		_ = d.Set("device_count", resp.DeviceCount)

	} else {
		// Get all device groups
		log.Printf("[INFO] Getting all device groups\n")
		allDeviceGroups, err := devicegroups.GetAllDevicesGroups(ctx, service)
		if err != nil {
			return diag.FromErr(err)
		}

		if len(allDeviceGroups) == 0 {
			return diag.FromErr(fmt.Errorf("no device groups found"))
		}

		// Set the first device group's ID as the data source ID
		d.SetId(fmt.Sprintf("%d", allDeviceGroups[0].ID))

		// Populate the list with all device groups
		deviceGroupList := make([]map[string]interface{}, len(allDeviceGroups))
		for i, deviceGroup := range allDeviceGroups {
			deviceGroupList[i] = map[string]interface{}{
				"id":           deviceGroup.ID,
				"name":         deviceGroup.Name,
				"group_type":   deviceGroup.GroupType,
				"description":  deviceGroup.Description,
				"os_type":      deviceGroup.OSType,
				"predefined":   deviceGroup.Predefined,
				"device_names": deviceGroup.DeviceNames,
				"device_count": deviceGroup.DeviceCount,
			}
		}

		_ = d.Set("list", deviceGroupList)

		// Also set the first device group as the main attributes for backward compatibility
		resp := allDeviceGroups[0]
		_ = d.Set("name", resp.Name)
		_ = d.Set("group_type", resp.GroupType)
		_ = d.Set("description", resp.Description)
		_ = d.Set("os_type", resp.OSType)
		_ = d.Set("predefined", resp.Predefined)
		_ = d.Set("device_names", resp.DeviceNames)
		_ = d.Set("device_count", resp.DeviceCount)
	}

	return nil
}
