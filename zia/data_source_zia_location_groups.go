package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/location/locationgroups"
)

func dataSourceLocationGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocationGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dynamic_location_group_criteria": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_string": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"match_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"countries": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"city": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_string": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"match_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"managed_by": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"extensions": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"enforce_authentication": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enforce_aup": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enforce_firewall_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_xff_forwarding": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_caution": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_bandwidth_control": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"profiles": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_mod_user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"last_mod_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"predefined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceLocationGroupRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	service := zClient.locationgroups

	var resp *locationgroups.LocationGroup
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting time window id: %d\n", id)
		res, err := locationgroups.GetLocationGroup(service, id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for location name: %s\n", name)
		res, err := locationgroups.GetLocationGroupByName(service, name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("deleted", resp.Deleted)
		_ = d.Set("group_type", resp.GroupType)
		_ = d.Set("comments", resp.Comments)
		_ = d.Set("last_mod_time", resp.LastModTime)
		_ = d.Set("predefined", resp.Predefined)

		if err := d.Set("dynamic_location_group_criteria", flattenDynamicLocationGroupCriteria(resp.DynamicLocationGroupCriteria)); err != nil {
			return err
		}

		if err := d.Set("locations", flattenGroupsLocations(resp)); err != nil {
			return err
		}

		if err := d.Set("last_mod_user", flattenLastModUser(resp.LastModUser)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any location group with name '%s'", name)
	}

	return nil
}

func flattenDynamicLocationGroupCriteria(dynamicLocationGroup *locationgroups.DynamicLocationGroupCriteria) interface{} {
	if dynamicLocationGroup == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"countries":                dynamicLocationGroup.Countries,
			"enforce_authentication":   dynamicLocationGroup.EnforceAuthentication,
			"enforce_aup":              dynamicLocationGroup.EnforceAup,
			"enforce_firewall_control": dynamicLocationGroup.EnforceFirewallControl,
			"enable_xff_forwarding":    dynamicLocationGroup.EnableXffForwarding,
			"enable_caution":           dynamicLocationGroup.EnableCaution,
			"enable_bandwidth_control": dynamicLocationGroup.EnableBandwidthControl,
			"profiles":                 dynamicLocationGroup.Profiles,
			"name":                     flattenDynamicGroupName(dynamicLocationGroup.Name),
			"city":                     flattenDynamicGroupCity(dynamicLocationGroup.City),
			"managed_by":               flattenLocationGroupManagedBy(dynamicLocationGroup.ManagedBy),
		},
	}
}

func flattenDynamicGroupName(name *locationgroups.Name) interface{} {
	if name == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"match_string": name.MatchString,
			"match_type":   name.MatchType,
		},
	}
}

func flattenDynamicGroupCity(city *locationgroups.City) interface{} {
	if city == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"match_string": city.MatchString,
			"match_type":   city.MatchType,
		},
	}
}

func flattenLocationGroupManagedBy(managedBy []locationgroups.ManagedBy) []interface{} {
	managed := make([]interface{}, len(managedBy))
	for i, val := range managedBy {
		managed[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return managed
}

func flattenGroupsLocations(locationGroup *locationgroups.LocationGroup) []interface{} {
	locations := make([]interface{}, len(locationGroup.Locations))
	for i, val := range locationGroup.Locations {
		locations[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return locations
}

func flattenLastModUser(lastModUser *locationgroups.LastModUser) interface{} {
	if lastModUser == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"id":         lastModUser.ID,
			"name":       lastModUser.Name,
			"extensions": lastModUser.Extensions,
		},
	}
}
