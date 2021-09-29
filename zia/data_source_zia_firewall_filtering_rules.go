package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/fwfilteringrules"
)

func dataSourceFirewallFilteringRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallFilteringRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rank": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_control": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_full_logging": {
				Type:     schema.TypeBool,
				Computed: true,
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
			"location_groups": {
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
			"departments": {
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
			"groups": {
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
			"users": {
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
			"time_windows": {
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
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_modified_by": {
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
			"src_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"src_ip_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_ip_categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_countries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_ip_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_services": {
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
			"nw_service_groups": {
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
			"nw_applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_application_groups": {
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
			"app_services": {
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
			"app_service_groups": {
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
			"labels": {
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
			"default_rule": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"predefined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceFirewallFilteringRuleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *fwfilteringrules.FirewallFilteringRules
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting data for rule id: %d\n", id)
		res, err := zClient.fwfilteringrules.GetFirewallFilteringRules(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for rule : %s\n", name)
		res, err := zClient.fwfilteringrules.GetFirewallFilteringRulesByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.Name)
		_ = d.Set("order", resp.Order)
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("access_control", resp.AccessControl)
		_ = d.Set("enable_full_logging", resp.EnableFullLogging)
		_ = d.Set("action", resp.Action)
		_ = d.Set("state", resp.State)
		_ = d.Set("description", resp.Description)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("src_ips", resp.SrcIps)
		_ = d.Set("dest_addresses", resp.DestAddresses)
		_ = d.Set("dest_ip_categories", resp.DestIpCategories)
		_ = d.Set("dest_countries", resp.DestCountries)
		_ = d.Set("nw_applications", resp.NwApplications)
		_ = d.Set("default_rule", resp.DefaultRule)
		_ = d.Set("predefined", resp.Predefined)

		if err := d.Set("locations", flattenLocations(resp.Locations)); err != nil {
			return err
		}

		if err := d.Set("location_groups", flattenLocationGroups(resp.LocationsGroups)); err != nil {
			return err
		}

		if err := d.Set("departments", flattenDepartments(resp.Departments)); err != nil {
			return err
		}

		if err := d.Set("groups", flattenGroups(resp.Groups)); err != nil {
			return err
		}

		if err := d.Set("users", flattenUsers(resp.Users)); err != nil {
			return err
		}

		if err := d.Set("time_windows", flattenTimeWindows(resp.TimeWindows)); err != nil {
			return err
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
			return err
		}

		if err := d.Set("src_ip_groups", flattenSrcIPGroups(resp.SrcIpGroups)); err != nil {
			return err
		}

		if err := d.Set("dest_ip_groups", flattenDestIPGroups(resp.DestIpGroups)); err != nil {
			return err
		}

		if err := d.Set("nw_services", flattenNWServices(resp.NwServices)); err != nil {
			return err
		}

		if err := d.Set("nw_service_groups", flattenNWServiceGroups(resp.NwServiceGroups)); err != nil {
			return err
		}

		if err := d.Set("nw_application_groups", flattenNWApplicationGroups(resp.NwApplicationGroups)); err != nil {
			return err
		}

		if err := d.Set("app_services", flattenAppServices(resp.AppServices)); err != nil {
			return err
		}

		if err := d.Set("app_services", flattenAppServiceGroups(resp.AppServiceGroups)); err != nil {
			return err
		}

		if err := d.Set("labels", flattenLabels(resp.Labels)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any user with name '%s' or id '%d'", name, id)
	}

	return nil
}

func flattenLocations(locations []fwfilteringrules.Locations) []interface{} {
	location := make([]interface{}, len(locations))
	for i, val := range locations {
		location[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return location
}

func flattenLocationGroups(locationGroups []fwfilteringrules.LocationsGroups) []interface{} {
	locationGroup := make([]interface{}, len(locationGroups))
	for i, val := range locationGroups {
		locationGroup[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return locationGroup
}

func flattenDepartments(departments []fwfilteringrules.Departments) []interface{} {
	department := make([]interface{}, len(departments))
	for i, val := range departments {
		department[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return department
}

func flattenGroups(groups []fwfilteringrules.Groups) []interface{} {
	group := make([]interface{}, len(groups))
	for i, val := range groups {
		group[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return group
}

func flattenUsers(users []fwfilteringrules.Users) []interface{} {
	user := make([]interface{}, len(users))
	for i, val := range users {
		user[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return user
}

func flattenTimeWindows(timeWindows []fwfilteringrules.TimeWindows) []interface{} {
	timeWindow := make([]interface{}, len(timeWindows))
	for i, val := range timeWindows {
		timeWindow[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return timeWindow
}

func flattenLastModifiedBy(lastModifiedBy []fwfilteringrules.LastModifiedBy) []interface{} {
	lastModified := make([]interface{}, len(lastModifiedBy))
	for i, val := range lastModifiedBy {
		lastModified[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return lastModified
}

func flattenSrcIPGroups(srcIPGroups []fwfilteringrules.SrcIpGroups) []interface{} {
	srcIpGroup := make([]interface{}, len(srcIPGroups))
	for i, val := range srcIPGroups {
		srcIpGroup[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return srcIpGroup
}

func flattenDestIPGroups(destIPGroups []fwfilteringrules.DestIpGroups) []interface{} {
	destIpGroup := make([]interface{}, len(destIPGroups))
	for i, val := range destIPGroups {
		destIpGroup[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return destIpGroup
}

func flattenNWServices(nwServices []fwfilteringrules.NwServices) []interface{} {
	nwService := make([]interface{}, len(nwServices))
	for i, val := range nwServices {
		nwService[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return nwService
}

func flattenNWServiceGroups(nwServiceGroups []fwfilteringrules.NwServiceGroups) []interface{} {
	NwServiceGroup := make([]interface{}, len(nwServiceGroups))
	for i, val := range nwServiceGroups {
		NwServiceGroup[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return NwServiceGroup
}

func flattenNWApplicationGroups(nwApplicationGroups []fwfilteringrules.NwApplicationGroups) []interface{} {
	nwApplicationGroup := make([]interface{}, len(nwApplicationGroups))
	for i, val := range nwApplicationGroups {
		nwApplicationGroup[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return nwApplicationGroup
}

func flattenAppServices(appServices []fwfilteringrules.AppServices) []interface{} {
	appService := make([]interface{}, len(appServices))
	for i, val := range appServices {
		appService[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return appService
}

func flattenAppServiceGroups(appServiceGroups []fwfilteringrules.AppServiceGroups) []interface{} {
	appServiceGroup := make([]interface{}, len(appServiceGroups))
	for i, val := range appServiceGroups {
		appServiceGroup[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return appServiceGroup
}

func flattenLabels(labels []fwfilteringrules.Labels) []interface{} {
	label := make([]interface{}, len(labels))
	for i, val := range labels {
		label[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}

	return label
}
