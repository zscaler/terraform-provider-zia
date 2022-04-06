package zia

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
)

func listIDsSchemaTypeCustom(maxItems int, desc string) *schema.Schema {
	ids := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
	if maxItems > 0 {
		ids.MaxItems = maxItems
	}
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Computed:    true,
		// MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": ids,
			},
		},
	}
}

func listIDsSchemaType(desc string) *schema.Schema {
	return listIDsSchemaTypeCustom(0, desc)
}

func expandIDNameExtensionsMap(m map[string]interface{}, key string) []common.IDNameExtensions {
	setInterface, ok := m[key]
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.IDNameExtensions
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil {
				for _, id := range itemMap["id"].([]interface{}) {
					result = append(result, common.IDNameExtensions{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.IDNameExtensions{}
}

func expandIDNameExtensionsSet(d *schema.ResourceData, key string) []common.IDNameExtensions {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.IDNameExtensions
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil {
				for _, id := range itemMap["id"].([]interface{}) {
					result = append(result, common.IDNameExtensions{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.IDNameExtensions{}
}

func expandUserGroups(d *schema.ResourceData, key string) []common.UserGroups {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.UserGroups
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil {
				for _, id := range itemMap["id"].([]interface{}) {
					result = append(result, common.UserGroups{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []common.UserGroups{}
}

func expandIDNameExtensions(d *schema.ResourceData, key string) *common.IDNameExtensions {
	lastModifiedByObj, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	lastMofiedBy, ok := lastModifiedByObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(lastMofiedBy.List()) > 0 {
		lastModifiedByObj := lastMofiedBy.List()[0]
		lastMofied, ok := lastModifiedByObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &common.IDNameExtensions{
			ID:         lastMofied["id"].(int),
			Name:       lastMofied["name"].(string),
			Extensions: lastMofied["extensions"].(map[string]interface{}),
		}
	}
	return nil
}

func expandUserDepartment(d *schema.ResourceData) *common.UserDepartment {
	departmentObj, ok := d.GetOk("department")
	if !ok {
		return nil
	}
	departments, ok := departmentObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(departments.List()) > 0 {
		departmentObj := departments.List()[0]
		department, ok := departmentObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &common.UserDepartment{
			ID:       department["id"].(int),
			Name:     department["name"].(string),
			IdpID:    department["idp_id"].(int),
			Comments: department["comments"].(string),
			Deleted:  department["deleted"].(bool),
		}
	}
	return nil
}

func flattenIDs(list []common.IDNameExtensions) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]int, len(list))
	for i, item := range list {
		ids[i] = item.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenUserGroupSet(list []common.UserGroups) []interface{} {
	result := make([]interface{}, 1)
	mapIds := make(map[string]interface{})
	ids := make([]int, len(list))
	for i, item := range list {
		ids[i] = item.ID
	}
	mapIds["id"] = ids
	result[0] = mapIds
	return result
}

func flattenIDNameExtensions(list []common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		flattenedList[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}
	return flattenedList
}

func flattenIDExtensions(list []common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		flattenedList[i] = map[string]interface{}{
			"id":         val.ID,
			"name":       val.Name,
			"extensions": val.Extensions,
		}
	}
	return flattenedList
}

func flattenIDExtensionsList(list *common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, 0)
	if list != nil {
		flattenedList = append(flattenedList, map[string]interface{}{
			"id":         list.ID,
			"name":       list.Name,
			"extensions": list.Extensions,
		})
	}
	return flattenedList
}

func flattenLastModifiedBy(lastModifiedBy *common.IDNameExtensions) []interface{} {
	lastModified := make([]interface{}, 0)
	if lastModifiedBy != nil {
		lastModified = append(lastModified, map[string]interface{}{
			"id":         lastModifiedBy.ID,
			"name":       lastModifiedBy.Name,
			"extensions": lastModifiedBy.Extensions,
		})
	}
	return lastModified
}

func flattenCreatedBy(createdBy *common.IDNameExtensions) []interface{} {
	created := make([]interface{}, 0)
	if createdBy != nil {
		created = append(created, map[string]interface{}{
			"id":         createdBy.ID,
			"name":       createdBy.Name,
			"extensions": createdBy.Extensions,
		})
	}
	return created
}

func flattenUserDepartment(userDepartment *common.UserDepartment) []interface{} {
	department := make([]interface{}, 0)
	if userDepartment != nil {
		department = append(department, map[string]interface{}{
			"name":     userDepartment.Name,
			"idp_id":   userDepartment.IdpID,
			"comments": userDepartment.Comments,
			"deleted":  userDepartment.Deleted,
		})
	}
	return department
}

func resourceNetworkPortsSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
				"end": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
				},
			},
		},
	}
}

func dataNetworkPortsSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"start": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"end": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func flattenNetwordPorts(ports []networkservices.NetworkPorts) []interface{} {
	portsObj := make([]interface{}, len(ports))
	for i, val := range ports {
		portsObj[i] = map[string]interface{}{
			"start": val.Start,
			"end":   val.End,
		}
	}
	return portsObj
}

func expandNetworkPorts(d *schema.ResourceData, key string) []networkservices.NetworkPorts {
	var ports []networkservices.NetworkPorts
	if portsInterface, ok := d.GetOk(key); ok {
		portSet, ok := portsInterface.(*schema.Set)
		if !ok {
			log.Printf("[ERROR] conversion failed, destUdpPortsInterface")
			return ports
		}
		ports = make([]networkservices.NetworkPorts, len(portSet.List()))
		for i, val := range portSet.List() {
			portItem := val.(map[string]interface{})
			ports[i] = networkservices.NetworkPorts{
				Start: portItem["start"].(int),
				End:   portItem["end"].(int),
			}
		}
	}
	return ports
}

func getURLCategories() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "List of URL categories for which rule must be applied",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURLFilteringCategories(),
		},
		Optional: true,
		Computed: true,
	}
}

func getURLRequestMethods() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Request method for which the rule must be applied. If not set, rule will be applied to all methods",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURLFilteringRequestMethods(),
		},
		Optional: true,
		// Computed: true,
	}
}

func getURLProtocols() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Supported Protocol criteria",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURLFilteringProtocols(),
		},
		Optional: true,
		// Computed: true,
	}
}

func getLocationManagementCountries() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateLocationManagementCountries(),
		Description:  "Supported Countries",
		Optional:     true,
		Computed:     true,
	}
}

func getLocationManagementTimeZones() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateLocationManagementTimeZones(),
		Description:  "Timezone of the location. If not specified, it defaults to GMT.",
		Optional:     true,
		Computed:     true,
	}
}

func getCloudFirewallDstCountries() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateCloudFirewallDstCountries(),
		},
		Optional: true,
		Computed: true,
	}
}

func getCloudFirewallNwApplications() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateCloudFirewallNwApplications(),
		},
		Optional: true,
		Computed: true,
	}
}

func getCloudFirewallNwServicesTag() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateCloudFirewallNwServicesTag(),
		Optional:     true,
		Computed:     true,
	}
}

func getDLPRuleFileTypes(desc string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "The list of file types to which the DLP policy rule must be applied.",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateDLPRuleFileTypes(),
		},
		Optional: true,
		Computed: true,
	}
}
