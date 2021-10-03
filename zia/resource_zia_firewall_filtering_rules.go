package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/filteringrules"
)

func resourceFirewallFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create:   resourceFirewallFilteringRulesCreate,
		Read:     resourceFirewallFilteringRulesRead,
		Update:   resourceFirewallFilteringRulesUpdate,
		Delete:   resourceFirewallFilteringRulesDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rank": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"access_control": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_full_logging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"locations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"location_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"departments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"time_windows": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"src_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"src_ip_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_ip_categories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_countries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_ip_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"nw_service_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"nw_applications": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nw_application_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"app_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"app_service_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
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
			"default_rule": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"predefined": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceFirewallFilteringRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandFirewallFilteringRules(d)
	log.Printf("[INFO] Creating zia firewall filtering rule\n%+v\n", req)

	resp, err := zClient.filteringrules.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia firewall filtering rule request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceFirewallFilteringRulesRead(d, m)
}

func resourceFirewallFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, err := zClient.filteringrules.Get(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting firewall filtering rule:\n%+v\n", resp)

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

	return nil
}

func resourceFirewallFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating firewall filtering rule ID: %v\n", id)
	req := expandFirewallFilteringRules(d)

	if _, err := zClient.filteringrules.Update(id, &req); err != nil {
		return err
	}

	return resourceFirewallFilteringRulesRead(d, m)
}

func resourceFirewallFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting firewall filtering rule ID: %v\n", (d.Id()))

	if _, err := zClient.filteringrules.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] firewall filtering rule deleted")
	return nil
}

func expandFirewallFilteringRules(d *schema.ResourceData) filteringrules.FirewallFilteringRules {
	return filteringrules.FirewallFilteringRules{
		Name:                d.Get("name").(string),
		Order:               d.Get("order").(int),
		Rank:                d.Get("rank").(int),
		Action:              d.Get("action").(string),
		State:               d.Get("state").(string),
		Description:         d.Get("description").(string),
		LastModifiedTime:    d.Get("last_modified_time").(int),
		SrcIps:              d.Get("src_ips").([]string),
		DestAddresses:       d.Get("src_ips").([]string),
		DestIpCategories:    d.Get("dest_ip_categories").([]string),
		DestCountries:       d.Get("dest_countries").([]string),
		NwApplications:      d.Get("nw_applications").([]string),
		DefaultRule:         d.Get("default_rule").(bool),
		Predefined:          d.Get("predefined").(bool),
		Locations:           expandLocations(d),
		LocationsGroups:     expandLocationsGroups(d),
		Departments:         expandDepartments(d),
		Groups:              expandGroups(d),
		Users:               expandUsers(d),
		TimeWindows:         expandTimeWindows(d),
		LastModifiedBy:      expandLastModifiedBy(d),
		SrcIpGroups:         expandSrcIpGroups(d),
		DestIpGroups:        expandDestIpGroups(d),
		NwServices:          expandNwServices(d),
		NwServiceGroups:     expandNwServiceGroups(d),
		NwApplicationGroups: expandNwApplicationGroups(d),
		AppServices:         expandAppServices(d),
		AppServiceGroups:    expandAppServiceGroups(d),
		Labels:              expandLabels(d),
	}
}

func expandLocations(d *schema.ResourceData) []filteringrules.Locations {
	var locations []filteringrules.Locations
	if locationInterface, ok := d.GetOk("locations"); ok {
		location := locationInterface.([]interface{})
		locations = make([]filteringrules.Locations, len(location))
		for i, location := range location {
			locationItem := location.(map[string]interface{})
			locations[i] = filteringrules.Locations{
				ID:         locationItem["id"].(int),
				Extensions: locationItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return locations
}

func expandLocationsGroups(d *schema.ResourceData) []filteringrules.LocationsGroups {
	var locationGroups []filteringrules.LocationsGroups
	if locationGroupInterface, ok := d.GetOk("location_groups"); ok {
		location := locationGroupInterface.([]interface{})
		locationGroups = make([]filteringrules.LocationsGroups, len(locationGroups))
		for i, locationGroup := range location {
			locationItem := locationGroup.(map[string]interface{})
			locationGroups[i] = filteringrules.LocationsGroups{
				ID:         locationItem["id"].(int),
				Extensions: locationItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return locationGroups
}

func expandDepartments(d *schema.ResourceData) []filteringrules.Departments {
	var departments []filteringrules.Departments
	if departmentsInterface, ok := d.GetOk("departments"); ok {
		department := departmentsInterface.([]interface{})
		departments = make([]filteringrules.Departments, len(departments))
		for i, dept := range department {
			departmentItem := dept.(map[string]interface{})
			departments[i] = filteringrules.Departments{
				ID:         departmentItem["id"].(int),
				Extensions: departmentItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return departments
}

func expandGroups(d *schema.ResourceData) []filteringrules.Groups {
	var groups []filteringrules.Groups
	if groupsInterface, ok := d.GetOk("groups"); ok {
		group := groupsInterface.([]interface{})
		groups = make([]filteringrules.Groups, len(groups))
		for i, fwgroup := range group {
			groupItem := fwgroup.(map[string]interface{})
			groups[i] = filteringrules.Groups{
				ID:         groupItem["id"].(int),
				Extensions: groupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return groups
}

func expandUsers(d *schema.ResourceData) []filteringrules.Users {
	var users []filteringrules.Users
	if groupsInterface, ok := d.GetOk("groups"); ok {
		user := groupsInterface.([]interface{})
		users = make([]filteringrules.Users, len(user))
		for i, fwuser := range user {
			userItem := fwuser.(map[string]interface{})
			users[i] = filteringrules.Users{
				ID:         userItem["id"].(int),
				Extensions: userItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return users
}

func expandTimeWindows(d *schema.ResourceData) []filteringrules.TimeWindows {
	var timeWindows []filteringrules.TimeWindows
	if timeWindowInterface, ok := d.GetOk("time_windows"); ok {
		time := timeWindowInterface.([]interface{})
		timeWindows = make([]filteringrules.TimeWindows, len(time))
		for i, timeWindow := range time {
			timeWindowItem := timeWindow.(map[string]interface{})
			timeWindows[i] = filteringrules.TimeWindows{
				ID:         timeWindowItem["id"].(int),
				Extensions: timeWindowItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return timeWindows
}

func expandLastModifiedBy(d *schema.ResourceData) []filteringrules.LastModifiedBy {
	var lastModifiedBy []filteringrules.LastModifiedBy
	if lastModifiedByInterface, ok := d.GetOk("time_windows"); ok {
		modifiedBy := lastModifiedByInterface.([]interface{})
		lastModifiedBy = make([]filteringrules.LastModifiedBy, len(modifiedBy))
		for i, lastModified := range modifiedBy {
			lastModifiedItem := lastModified.(map[string]interface{})
			lastModifiedBy[i] = filteringrules.LastModifiedBy{
				ID:         lastModifiedItem["id"].(int),
				Extensions: lastModifiedItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return lastModifiedBy
}

func expandSrcIpGroups(d *schema.ResourceData) []filteringrules.SrcIpGroups {
	var srcIPGroups []filteringrules.SrcIpGroups
	if srcIPGroupInterface, ok := d.GetOk("src_ip_groups"); ok {
		srcIPGroup := srcIPGroupInterface.([]interface{})
		srcIPGroups = make([]filteringrules.SrcIpGroups, len(srcIPGroup))
		for i, srcIP := range srcIPGroup {
			srcIPGroupItem := srcIP.(map[string]interface{})
			srcIPGroups[i] = filteringrules.SrcIpGroups{
				ID:         srcIPGroupItem["id"].(int),
				Extensions: srcIPGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return srcIPGroups
}

func expandDestIpGroups(d *schema.ResourceData) []filteringrules.DestIpGroups {
	var destIPGroups []filteringrules.DestIpGroups
	if destIPGroupInterface, ok := d.GetOk("src_ip_groups"); ok {
		destIPGroup := destIPGroupInterface.([]interface{})
		destIPGroups = make([]filteringrules.DestIpGroups, len(destIPGroup))
		for i, destIP := range destIPGroup {
			destIPGroupItem := destIP.(map[string]interface{})
			destIPGroups[i] = filteringrules.DestIpGroups{
				ID:         destIPGroupItem["id"].(int),
				Extensions: destIPGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return destIPGroups
}

func expandNwServices(d *schema.ResourceData) []filteringrules.NwServices {
	var nwServices []filteringrules.NwServices
	if nwServiceInterface, ok := d.GetOk("nw_services"); ok {
		nwService := nwServiceInterface.([]interface{})
		nwServices = make([]filteringrules.NwServices, len(nwService))
		for i, services := range nwService {
			nwServicesItem := services.(map[string]interface{})
			nwServices[i] = filteringrules.NwServices{
				ID:         nwServicesItem["id"].(int),
				Extensions: nwServicesItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return nwServices
}

func expandNwServiceGroups(d *schema.ResourceData) []filteringrules.NwServiceGroups {
	var nwServiceGroups []filteringrules.NwServiceGroups
	if nwServiceGroupsInterface, ok := d.GetOk("nw_service_groups"); ok {
		nwServiceGroup := nwServiceGroupsInterface.([]interface{})
		nwServiceGroups = make([]filteringrules.NwServiceGroups, len(nwServiceGroup))
		for i, serviceGroups := range nwServiceGroup {
			serviceGroupItem := serviceGroups.(map[string]interface{})
			nwServiceGroups[i] = filteringrules.NwServiceGroups{
				ID:         serviceGroupItem["id"].(int),
				Extensions: serviceGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return nwServiceGroups
}

func expandNwApplicationGroups(d *schema.ResourceData) []filteringrules.NwApplicationGroups {
	var nwApplicationGroups []filteringrules.NwApplicationGroups
	if nwApplicationGroupInterface, ok := d.GetOk("nw_application_groups"); ok {
		nwApplicationGroup := nwApplicationGroupInterface.([]interface{})
		nwApplicationGroups = make([]filteringrules.NwApplicationGroups, len(nwApplicationGroup))
		for i, appGroups := range nwApplicationGroup {
			appGroupItem := appGroups.(map[string]interface{})
			nwApplicationGroups[i] = filteringrules.NwApplicationGroups{
				ID:         appGroupItem["id"].(int),
				Extensions: appGroupItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return nwApplicationGroups
}

func expandAppServices(d *schema.ResourceData) []filteringrules.AppServices {
	var appServices []filteringrules.AppServices
	if appServicesInterface, ok := d.GetOk("app_services"); ok {
		appService := appServicesInterface.([]interface{})
		appServices = make([]filteringrules.AppServices, len(appService))
		for i, services := range appService {
			servicesItem := services.(map[string]interface{})
			appServices[i] = filteringrules.AppServices{
				ID:         servicesItem["id"].(int),
				Extensions: servicesItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return appServices
}

func expandAppServiceGroups(d *schema.ResourceData) []filteringrules.AppServiceGroups {
	var appServices []filteringrules.AppServiceGroups
	if appServicesInterface, ok := d.GetOk("app_service_groups"); ok {
		appService := appServicesInterface.([]interface{})
		appServices = make([]filteringrules.AppServiceGroups, len(appService))
		for i, services := range appService {
			servicesItem := services.(map[string]interface{})
			appServices[i] = filteringrules.AppServiceGroups{
				ID:         servicesItem["id"].(int),
				Extensions: servicesItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return appServices
}

func expandLabels(d *schema.ResourceData) []filteringrules.Labels {
	var labels []filteringrules.Labels
	if labelsInterface, ok := d.GetOk("labels"); ok {
		label := labelsInterface.([]interface{})
		labels = make([]filteringrules.Labels, len(label))
		for i, fwlabels := range label {
			labelItem := fwlabels.(map[string]interface{})
			labels[i] = filteringrules.Labels{
				ID:         labelItem["id"].(int),
				Extensions: labelItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return labels
}
