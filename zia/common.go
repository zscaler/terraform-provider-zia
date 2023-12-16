package zia

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservices"
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
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
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

func setIDsSchemaTypeCustom(maxItems *int, desc string) *schema.Schema {
	ids := &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
	if maxItems != nil && *maxItems > 0 {
		ids.MaxItems = *maxItems
	}
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": ids,
			},
		},
	}
}

func setIdNameSchemaCustom(maxItems int, description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: description,
		MaxItems:    maxItems,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The unique identifier for the resource.",
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the resource.",
				},
			},
		},
	}
}

func setExtIDNameSchemaCustom(maxItems *int, description string) *schema.Schema {
	schema := &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    true,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the application segment.",
				},
				"external_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "External ID of the application segment.",
				},
			},
		},
	}

	if maxItems != nil && *maxItems > 0 {
		schema.MaxItems = *maxItems
	}

	return schema
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

func expandIDNameExtensionsListSingle(d *schema.ResourceData, key string) *common.IDNameExtensions {
	l := expandIDNameExtensionsSet(d, key)
	if len(l) > 0 {
		r := l[0]
		return &r
	}
	return nil
}

func expandIDNameExtensionsSet(d *schema.ResourceData, key string) []common.IDNameExtensions {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.IDNameExtensions
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
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

/*
// Deprecated common helper function

	func expandIDSet(d *schema.ResourceData, key string) []int {
		var ids []int

		if v, ok := d.GetOk(key); ok {
			set := v.(*schema.Set)
			list := set.List()
			for _, item := range list {
				if idMap, ok := item.(map[string]interface{}); ok {
					if id, ok := idMap["id"].(int); ok {
						ids = append(ids, id)
					}
				}
			}
		}
		return ids
	}
*/

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

func flattenIDNameExtensions(list []common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		r := map[string]interface{}{
			"id":   val.ID,
			"name": val.Name,
		}
		if val.Extensions != nil {
			r["extensions"] = val.Extensions
		}
		flattenedList[i] = r
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

func flattenIDExtensionsList(idNameExtension *common.IDNameExtensions) []interface{} {
	flattenedList := make([]interface{}, 0)
	if idNameExtension != nil && (idNameExtension.ID != 0 || idNameExtension.Name != "") {
		flattenedList = append(flattenedList, map[string]interface{}{
			"id":         idNameExtension.ID,
			"name":       idNameExtension.Name,
			"extensions": idNameExtension.Extensions,
		})
	}
	return flattenedList
}

func flattenIDExtensionListIDs(idNameExtensions *common.IDNameExtensions) []interface{} {
	if idNameExtensions == nil || idNameExtensions.ID == 0 && idNameExtensions.Name == "" {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": []int{idNameExtensions.ID},
		},
	}
}

func flattenIDExtensionsListIDs(list []common.IDNameExtensions) []interface{} {
	if list == nil {
		return nil
	}
	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}
	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}

// Flattening function used in the Forwarding Control Policy Resource
func flattenIDNameSet(idName *common.IDName) []interface{} {
	idNameSet := make([]interface{}, 0)
	if idName != nil {
		idNameSet = append(idNameSet, map[string]interface{}{
			"id":   idName.ID,
			"name": idName.Name,
		})
	}
	return idNameSet
}

// expandIDNameSet takes a Terraform set as input and returns a pointer to a common.IDName struct.
func expandIDNameSet(d *schema.ResourceData, key string) *common.IDName {
	idNameList, ok := d.Get(key).(*schema.Set)
	if !ok || idNameList.Len() == 0 {
		return nil
	}

	// Assuming each set can only have one item as per your JSON structure.
	// If it can have multiple, this needs to be adjusted accordingly.
	for _, v := range idNameList.List() {
		item := v.(map[string]interface{})
		return &common.IDName{
			ID:   item["id"].(int),
			Name: item["name"].(string),
		}
	}

	return nil
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
		obj := map[string]interface{}{
			"id":      userDepartment.ID,
			"deleted": userDepartment.Deleted,
		}
		if userDepartment.Name != "" {
			obj["name"] = userDepartment.Name
		}
		if userDepartment.IdpID != 0 {
			obj["idp_id"] = userDepartment.IdpID
		}
		if userDepartment.Comments != "" {
			obj["comments"] = userDepartment.Comments
		}
		department = append(department, obj)
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
			Type:             schema.TypeString,
			ValidateDiagFunc: validateURLFilteringCategories(),
		},
		Optional: true,
	}
}

func getSuperCategories() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Super Category of the URL category. This field is required when creating custom URL categories.",
		ValidateDiagFunc: validateURLSuperCategories(),
	}
}

func getDeviceTrustLevels() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation.",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDeviceTrustLevels(),
		},
		Optional: true,
		ForceNew: true,
	}
}

func getURLRequestMethods() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Request method for which the rule must be applied. If not set, rule will be applied to all methods",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateURLFilteringRequestMethods(),
		},
		Optional: true,
	}
}

func getURLProtocols() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Supported Protocol criteria",
		Required:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateURLFilteringProtocols(), // Use ValidateDiagFunc here
		},
	}
}

func getLocationManagementCountries() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		ValidateDiagFunc: validateLocationManagementCountries(),
		Description:      "Supported Countries",
		Optional:         true,
		Computed:         true,
	}
}

func getLocationManagementTimeZones() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		ValidateDiagFunc: validateLocationManagementTimeZones(),
		Description:      "Timezone of the location. If not specified, it defaults to GMT.",
		Optional:         true,
		Computed:         true,
	}
}

func getDestinationCountries() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateDestinationCountries,
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
			Type:             schema.TypeString,
			ValidateDiagFunc: validateCloudFirewallNwApplications(),
		},
		Optional: true,
		Computed: true,
	}
}

func getCloudFirewallNwServicesTag() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		ValidateDiagFunc: validateCloudFirewallNwServicesTag(),
		Optional:         true,
		Computed:         true,
	}
}

func sortOrders(ruleOrderMap map[int]orderWithState) RuleIDOrderPairList {
	pl := make(RuleIDOrderPairList, len(ruleOrderMap))
	i := 0
	for k, v := range ruleOrderMap {
		pl[i] = RuleIDOrderPair{k, v.order}
		i++
	}
	sort.Sort(pl)
	return pl
}

type orderWithState struct {
	order int
	done  bool
}

type listrules struct {
	orders  map[string]map[int]orderWithState
	orderer map[string]int
	sync.Mutex
}

var rules = listrules{
	orders: make(map[string]map[int]orderWithState),
}

type RuleIDOrderPair struct {
	ID    int
	Order int
}

type RuleIDOrderPairList []RuleIDOrderPair

func (p RuleIDOrderPairList) Len() int { return len(p) }
func (p RuleIDOrderPairList) Less(i, j int) bool {
	if p[i].Order == p[j].Order {
		return p[i].ID < p[j].ID
	}
	return p[i].Order < p[j].Order
}
func (p RuleIDOrderPairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func reorderAll(resourceType string, getCount func() (int, error), updateOrder func(id, order int) error) {
	ticker := time.NewTicker(time.Second * 15) // create a ticker that ticks every half minute
	defer ticker.Stop()                        // stop the ticker when the loop ends
	numResources := []int{0, 0, 0}
	for {
		select {
		case <-ticker.C:
			rules.Lock()
			size := len(rules.orders[resourceType])
			done := true
			// first check if all rules creation is done
			for _, v := range rules.orders[resourceType] {
				if !v.done {
					done = false
				}
			}
			numResources[0], numResources[1], numResources[2] = numResources[1], numResources[2], size
			if done && numResources[0] == numResources[1] && numResources[1] == numResources[2] {
				// No changes after a while (4 runs), so reorder, and return
				count, _ := getCount()
				// sort by order (ascending)
				sorted := sortOrders(rules.orders[resourceType])
				log.Printf("[INFO] sorting filtering rule after tick; sorted:%v", sorted)
				for _, v := range sorted {
					if v.Order <= count {
						if err := updateOrder(v.ID, v.Order); err != nil {
							log.Printf("[ERROR] couldn't reorder the rule after tick, the order may not have taken place: %v\n", err)
						}
					}
				}
				rules.Unlock()
				return
			}
			rules.Unlock()
		default:
			time.Sleep(time.Second * 10)
		}
	}
}

func markOrderRuleAsDone(id int, resourceType string) {
	rules.Lock()
	r := rules.orders[resourceType][id]
	r.done = true
	rules.orders[resourceType][id] = r
	rules.Unlock()
}

func reorder(order, id int, resourceType string, getCount func() (int, error), updateOrder func(id, order int) error) {
	rules.Lock()
	shouldCallReorder := false
	if len(rules.orders) == 0 {
		rules.orders = map[string]map[int]orderWithState{}
		rules.orderer = map[string]int{}
	}
	if _, ok := rules.orderer[resourceType]; ok {
		shouldCallReorder = false
	} else {
		rules.orderer[resourceType] = id
		shouldCallReorder = true
	}
	if len(rules.orders[resourceType]) == 0 {
		rules.orders[resourceType] = map[int]orderWithState{}
	}
	rules.orders[resourceType][id] = orderWithState{order, shouldCallReorder}
	rules.Unlock()
	if shouldCallReorder {
		log.Printf("[INFO] starting to reorder the rules, delegating to rule:%d, order:%d", id, order)
		// one resource will wait until all resources are done and reorder then return
		reorderAll(resourceType, getCount, updateOrder)
	}
}
