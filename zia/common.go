package zia

import (
	"context"
	"fmt"
	"log"
	"maps"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudappcontrol"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
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
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
	if maxItems != nil && *maxItems > 0 {
		ids.MaxItems = *maxItems
	}
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Computed:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": ids,
			},
		},
	}
}

// setCustomKeyIDsSchema mirrors setIDsSchemaTypeCustom (single block containing a
// set of IDs) but lets the caller name the inner list attribute when the API key
// is not the generic "id" (e.g. "zapp_id", "group_id").
func setCustomKeyIDsSchema(innerKey, desc string) *schema.Schema {
	ids := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				innerKey: ids,
			},
		},
	}
}

// Used for Computed Attributes
func setIDsSchemaTypeCustomSpecial(maxItems *int, desc string) *schema.Schema {
	ids := &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
	if maxItems != nil && *maxItems > 0 {
		ids.MaxItems = *maxItems
	}
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Computed:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": ids,
			},
		},
	}
}

func setSingleIDSchemaTypeCustom(desc string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// Computed:    true,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Optional: true,
				},
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
					Optional:    true,
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

func setIDExternalIDCustom(maxItems *int, desc string) *schema.Schema {
	idList := &schema.Schema{
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeInt},
		Required: true,
	}
	if maxItems != nil && *maxItems > 0 {
		idList.MaxItems = *maxItems
	}
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		MaxItems:    1,
		Description: desc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": idList,
			},
		},
	}
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

// func expandIDNameExtensionsSetSingle(d *schema.ResourceData, key string) *common.IDCustom {
// 	if v, ok := d.GetOk(key); ok {
// 		setList := v.(*schema.Set).List()
// 		if len(setList) > 0 {
// 			if idMap, ok := setList[0].(map[string]interface{}); ok {
// 				return &common.IDCustom{
// 					ID: idMap["id"].(int),
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

func expandIDNameExtensionsSetSingle(d *schema.ResourceData, key string) *common.IDCustom {
	if v, ok := d.GetOk(key); ok {
		log.Printf("[DEBUG] expandIDNameExtensionsSetSingle key=%s raw=%#v", key, v)

		setList := v.(*schema.Set).List()
		if len(setList) > 0 {
			if idMap, ok := setList[0].(map[string]interface{}); ok {
				log.Printf("[DEBUG] expandIDNameExtensionsSetSingle extracted id: %v", idMap["id"])
				return &common.IDCustom{
					ID: idMap["id"].(int),
				}
			}
		}
	}
	log.Printf("[DEBUG] expandIDNameExtensionsSetSingle key=%s returned nil", key)
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

// TEMPORARY FUNCTION UNTIL NEXT GO SDK RELEASE
func expandCloudApplicationInstanceSet(d *schema.ResourceData, key string) []cloudappcontrol.CloudAppInstances {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []cloudappcontrol.CloudAppInstances
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
					result = append(result, cloudappcontrol.CloudAppInstances{
						ID: id.(int),
					})
				}
			}
		}
		return result
	}
	return []cloudappcontrol.CloudAppInstances{}
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

func flattenZPAAppSegmentsSimple(list []common.ZPAAppSegments) []interface{} {
	var flattenedList []interface{}
	for _, segment := range list {
		m := make(map[string]interface{})
		m["name"] = segment.Name
		m["external_id"] = segment.ExternalID

		flattenedList = append(flattenedList, m)
	}
	return flattenedList
}

func expandZPAAppSegmentSet(d *schema.ResourceData, key string) []common.ZPAAppSegments {
	setInterface, exists := d.GetOk(key)
	if !exists {
		return nil
	}

	inputSet := setInterface.(*schema.Set).List()
	var result []common.ZPAAppSegments
	for _, item := range inputSet {
		itemMap := item.(map[string]interface{})
		segment := common.ZPAAppSegments{
			Name:       itemMap["name"].(string),
			ExternalID: itemMap["external_id"].(string),
		}

		result = append(result, segment)
	}
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

func flattenCustomIDSet(customID *common.IDCustom) []interface{} {
	if customID == nil || customID.ID == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": customID.ID,
		},
	}
}

func flattenCustomIDNameSet(customID *common.IDCustom) []interface{} {
	if customID == nil || customID.ID == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id":   customID.ID,
			"name": customID.Name,
		},
	}
}

// expandSingleIDNameExtensions reads a single-element "id"-only set (as produced
// by setSingleIDSchemaTypeCustom) and returns the corresponding
// *common.IDNameExtensions, or nil when unset.
func expandSingleIDNameExtensions(d *schema.ResourceData, key string) *common.IDNameExtensions {
	if v, ok := d.GetOk(key); ok {
		setList := v.(*schema.Set).List()
		if len(setList) > 0 {
			if idMap, ok := setList[0].(map[string]interface{}); ok {
				if id, ok := idMap["id"].(int); ok && id != 0 {
					return &common.IDNameExtensions{ID: id}
				}
			}
		}
	}
	return nil
}

// flattenSingleIDNameExtensions renders a *common.IDNameExtensions into the
// single-element "id"-only set shape expected by setSingleIDSchemaTypeCustom.
func flattenSingleIDNameExtensions(v *common.IDNameExtensions) []interface{} {
	if v == nil || v.ID == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"id": v.ID,
		},
	}
}

// expandReceiver reads the shared "receiver" block used by the DLP web and
// endpoint DLP rule resources and returns the corresponding SDK Receiver.
func expandReceiver(d *schema.ResourceData, key string) *dlp_web_rules.Receiver {
	receiverSet, ok := d.Get(key).(*schema.Set)
	if !ok || receiverSet.Len() == 0 {
		return nil
	}

	receiverList := receiverSet.List()
	if len(receiverList) == 0 {
		return nil
	}

	item := receiverList[0].(map[string]interface{})

	idStr := item["id"].(string)
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		idInt = 0
	}

	receiver := &dlp_web_rules.Receiver{
		ID:   idInt,
		Name: item["name"].(string),
		Type: item["type"].(string),
	}

	if tenantList, ok := item["tenant"].([]interface{}); ok && len(tenantList) > 0 {
		if tenantMap, ok := tenantList[0].(map[string]interface{}); ok {
			tenantIDStr := tenantMap["id"].(string)
			tenantIDInt, err := strconv.Atoi(tenantIDStr)
			if err != nil {
				tenantIDInt = 0
			}
			receiver.Tenant = &common.IDNameExtensions{
				ID:   tenantIDInt,
				Name: tenantMap["name"].(string),
			}
		}
	}

	return receiver
}

// flattenReceiverResource renders the SDK Receiver into the resource-side
// "receiver" block shape (string IDs) shared by the DLP web and endpoint DLP
// rule resources.
func flattenReceiverResource(receiver *dlp_web_rules.Receiver) []interface{} {
	if receiver == nil {
		return nil
	}

	if receiver.ID == 0 && receiver.Name == "" && receiver.Type == "" && receiver.Tenant == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":   strconv.Itoa(receiver.ID),
		"name": receiver.Name,
		"type": receiver.Type,
	}

	if receiver.Tenant != nil {
		tenant := map[string]interface{}{
			"id":   strconv.Itoa(receiver.Tenant.ID),
			"name": receiver.Tenant.Name,
		}
		result["tenant"] = []interface{}{tenant}
	}

	return []interface{}{result}
}

// flattenReceiver renders the SDK Receiver into the data-source-side "receiver"
// block shape (integer IDs) shared by the DLP web and endpoint DLP rule data
// sources.
func flattenReceiver(receiver *dlp_web_rules.Receiver) []map[string]interface{} {
	if receiver == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":   receiver.ID,
		"name": receiver.Name,
		"type": receiver.Type,
	}

	if receiver.Tenant != nil {
		result["tenant"] = flattenIDExtensionsList(receiver.Tenant)
	}

	return []map[string]interface{}{result}
}

func flattenIDExtensionsListIDs(list []common.IDNameExtensions) []interface{} {
	if len(list) == 0 {
		// Return an empty slice instead of nil
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		// Again return []interface{}{} instead of nil
		return []interface{}{}
	}

	// The rest remains the same
	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}

// TEMPORARY FUNCTION UNTIL NEXT GO SDK RELEASE
func flattenIDCloudAppInstance(list []cloudappcontrol.CloudAppInstances) []interface{} {
	if len(list) == 0 {
		// Return an empty slice instead of nil
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		// Again return []interface{}{} instead of nil
		return []interface{}{}
	}

	// The rest remains the same
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

func flattenIDNameExternalSet(idName *common.IDNameExternalID) []interface{} {
	idNameSet := make([]interface{}, 0)
	if idName != nil {
		idNameSet = append(idNameSet, map[string]interface{}{
			"id":   idName.ID,
			"name": idName.Name,
		})
	}
	return idNameSet
}

func flattenCommonNSS(list []common.CommonNSS) []interface{} {
	flattenedList := make([]interface{}, len(list))
	for i, val := range list {
		flattenedList[i] = map[string]interface{}{
			"id":          val.ID,
			"pid":         val.PID,
			"name":        val.Name,
			"description": val.Description,
			"deleted":     val.Deleted,
			"getl_id":     val.GetlID,
		}
	}
	return flattenedList
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

// Common Flattening function to support Workload Groups across other resources
func flattenWorkloadGroups(workloadGroups []common.IDName) []interface{} {
	if workloadGroups == nil {
		return nil
	}

	wgList := make([]interface{}, len(workloadGroups))
	for i, wg := range workloadGroups {
		wgMap := make(map[string]interface{})
		wgMap["id"] = wg.ID
		wgMap["name"] = wg.Name
		wgList[i] = wgMap
	}

	return wgList
}

// Common expand function to support Workload Groups across other resources
func expandWorkloadGroupsIDName(d *schema.ResourceData, key string) []common.IDName {
	// Retrieve the set from the resource data
	if v, ok := d.GetOk(key); ok {
		workloadGroupsSet := v.(*schema.Set)
		// Initialize the slice to hold the expanded workload groups
		workloadGroups := make([]common.IDName, 0, workloadGroupsSet.Len())

		// Iterate over the set and construct the slice of common.IDName
		for _, wgMapInterface := range workloadGroupsSet.List() {
			wgMap := wgMapInterface.(map[string]interface{})
			wg := common.IDName{
				ID:   wgMap["id"].(int),
				Name: wgMap["name"].(string),
			}
			workloadGroups = append(workloadGroups, wg)
		}

		return workloadGroups
	}

	// Return an empty slice if the key is not set
	return []common.IDName{}
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

func flattenLastModifiedByExternalID(lastModifiedBy *common.IDNameExternalID) []interface{} {
	lastModified := make([]interface{}, 0)
	if lastModifiedBy != nil {
		lastModified = append(lastModified, map[string]interface{}{
			"id":          lastModifiedBy.ID,
			"name":        lastModifiedBy.Name,
			"extensions":  lastModifiedBy.Extensions,
			"external_id": lastModifiedBy.ExternalID,
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

func flattenUserGroups(list []common.UserGroups) []interface{} {
	if len(list) == 0 {
		// Return an empty slice instead of nil
		return []interface{}{}
	}

	ids := []int{}
	for _, item := range list {
		if item.ID == 0 && item.Name == "" {
			continue
		}
		ids = append(ids, item.ID)
	}

	if len(ids) == 0 {
		// Again return []interface{}{} instead of nil
		return []interface{}{}
	}

	// The rest remains the same
	return []interface{}{
		map[string]interface{}{
			"id": ids,
		},
	}
}

func expandUserGroups(d *schema.ResourceData, key string) []common.UserGroups {
	setInterface, ok := d.GetOk(key)
	if ok {
		set := setInterface.(*schema.Set)
		var result []common.UserGroups
		for _, item := range set.List() {
			itemMap, _ := item.(map[string]interface{})
			if itemMap != nil && itemMap["id"] != nil {
				set := itemMap["id"].(*schema.Set)
				for _, id := range set.List() {
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

func getFileTypeProtocols() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Protocol for the given rule. This field is not applicable to the Lite API.",
		Required:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateFileTypeControlProtocols(),
		},
	}
}

func getSandboxRuleProtocols() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Protocol for the given rule. This field is not applicable to the Lite API.",
		Required:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateSandboxRuleProtocols(),
		},
	}
}

func getDNSRuleProtocols() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Protocol for the given rule. This field is not applicable to the Lite API.",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDNSRuleProtocols(),
		},
	}
}

func getUserRiskScoreLevels() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateUserRiskScoreLevels(),
		},
		Optional: true,
	}
}

func getUserAgentTypes() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Supported User Agent Types",
		Optional:    true,
		// MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateUserAgentTypes(),
		},
	}
}

func getAppControlType() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "Supported App Control Types",
		Optional:    true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateAppControlType(),
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

func getISOCountryCodes() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.",
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateISOCountryCodes,
		},
		Optional: true,
		Computed: true,
	}
}

func getCloudApplications() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateCloudApplications(),
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

func getFileTypes() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "File type categories for which the policy is applied. If not set, the rule is applied across all file types.",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateFileTypes(),
		},
		Optional: true,
	}
}

func getSandboxFileTypes() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "File type categories for which the policy is applied. If not set, the rule is applied across all file types.",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateSandboxRuleFileTypes(),
		},
		Required: true,
	}
}

func getBaPolicyCategories() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "The threat categories to which the rule applies",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateSandboxPolicyCategories(),
		},
		Optional: true,
	}
}

func getDnsRuleRequestTypes() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "DNS request types to which the rule applies",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDnsRuleRequestTypes(),
		},
		Optional: true,
	}
}

func getSSLInspectionPlatforms() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Supported Protocol criteria",
		Optional:    true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateSSLInspectionPlatforms(), // Use ValidateDiagFunc here
		},
	}
}

func getAdminRolePermissions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Request method for which the rule must be applied. If not set, rule will be applied to all methods",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateAdminRolePermissions(),
		},
	}
}

func getAlertSubscriptionSeverity() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Lists the severity levels of the Patient 0, Secure Alert class, System Alerts class",
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateAlertSubscriptionSeverity(),
		},
		Optional: true,
	}
}

func getCasbRuleCollaborationScope() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Collaboration scope for the rule",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateCasbRuleCollaborationScope(),
		},
	}
}

func getCasbRuleComponents() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "List of components for which the rule is applied. Zscaler service inspects these components for sensitive data.",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateCasbRuleComponents(),
		},
	}
}

func getRiskProfileCertifications() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "List of certifications to be included or excluded for the profile",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateRiskProfileCertifications(),
		},
	}
}

func getRiskProfileIndex() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "The risk index number of the cloud applications. It represents the risk score assigned to each cloud application based on the risk attribute values.",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeInt,
			ValidateDiagFunc: validateRiskProfileIndex(),
		},
	}
}

func getRiskProfileEncryptionInTransit() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "Filters applications based on their support for encrypting data in transit",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validateRiskProfileEncryptionInTransit(),
		},
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
	order OrderRule
	done  bool
}

type listrules struct {
	orders      map[string]map[int]orderWithState
	orderer     map[string]int
	reorderDone map[string]chan struct{}
	revisions   map[string]uint64
	sync.Mutex
}

var rules = listrules{
	orders:      make(map[string]map[int]orderWithState),
	reorderDone: make(map[string]chan struct{}),
}

type RuleIDOrderPair struct {
	ID    int
	Order OrderRule
}

type RuleIDOrderPairList []RuleIDOrderPair

func (p RuleIDOrderPairList) Len() int { return len(p) }
func (p RuleIDOrderPairList) Less(i, j int) bool {
	if p[i].Order == p[j].Order {
		return p[i].ID < p[j].ID
	}
	return p[i].Order.Rank < p[j].Order.Rank || p[i].Order.Rank == p[j].Order.Rank && p[i].Order.Order < p[j].Order.Order
}
func (p RuleIDOrderPairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// reorderTickInterval controls how often the reorder ticker fires.
// Default is 30 seconds. Tests can override this to speed up reorder cycles.
var reorderTickInterval = 30 * time.Second

// maxNoProgressTicks bounds how long reorderAll will keep retrying when
// the loop has stopped making forward progress (no new rules registering
// AND no successful PUTs AND skipped set not shrinking). Each tick is
// reorderTickInterval (30s default), so 6 ticks ≈ 3 min. This protects
// against configurations that can never converge (e.g. declared orders
// with permanent gaps, or some Create calls failing so the maximum
// declared order can never be reached). It is NOT used while creates
// are still in flight — see the progress detection in reorderAll.
const maxNoProgressTicks = 6

// maxStuckOnSkippedTicks is the FAST early-exit when we've reconciled
// every in-range rule and the only remaining work is rules whose
// declared order exceeds the API's current orderable count. This is
// the deadlock case caused by terraform parallelism: 10 rules POST,
// 10 register and block on waitForReorder, but some of those declared
// orders > apiOrderable, so reorder has nothing to do for them. We
// must return so waitForReorder unblocks and the next batch of POSTs
// can extend apiOrderable. Without this, we spent ~3 minutes per
// batch on the slower maxNoProgressTicks safety net for what is
// actually a normal mid-apply state. 2 ticks ≈ 60s lets the cache
// catch up to in-flight PUTs before we declare "nothing more to do".
const maxStuckOnSkippedTicks = 2

// maxNonImprovingPutTicks is the oscillation breaker (SUP-4231). When the
// declared orders are unsatisfiable — duplicate order values in the tenant,
// or unmanaged rules interleaved with managed ones — every PUT displaces a
// neighbouring rule, so each pass keeps issuing PUTs forever without the
// at-target count ever surpassing its best. Because putsIssued > 0 counts
// as progress for maxNoProgressTicks, no other exit condition can fire and
// the cycle livelocks until terraform's own context deadline kills it
// (observed as a 2.5h hang). After this many consecutive passes that issue
// PUTs without improving on the best at-target count, the loop logs the
// still-misplaced rules and returns so the apply can complete; the
// unplaced rules surface as drift on the next plan.
const maxNonImprovingPutTicks = 6

// reorderAll runs as a background goroutine and reconciles the desired
// order of every rule registered for resourceType against the API's
// current state. It is intentionally diff-based:
//
//   - Each pass calls getCurrent ONCE to fetch {ruleID -> {Order, Rank}}
//     for every rule the API knows about.
//   - For each registered rule we compare the API's current Order/Rank
//     to the desired Order/Rank and SKIP the rule entirely when they
//     already match — no per-rule GET, no PUT, no API traffic at all.
//   - Only rules whose position actually drifted from desired call into
//     updateOrder. This is the property the user has been asking for
//     since SUP-3988 was first opened: don't PUT a rule whose order
//     hasn't changed.
//
// If a registered Order exceeds the number of API rules with Order >= 1,
// release this batch without order PUTs so more creates can extend the list.
// The next registration starts a new cycle with the same declared targets.
// This avoids repeated intermediate moves and rejected out-of-range PUTs.
//
// Progress detection: a pass made progress if size grew (new rule
// registered), OR a PUT was issued successfully, OR the deferred set
// shrank. Without progress, noProgressTicks accumulates and the loop
// eventually returns so waitForReorder can unblock — preventing the
// goroutine from leaking when declared orders are unreachable.
func reorderAll(
	resourceType string,
	getCurrent func() (map[int]OrderRule, error),
	updateOrder func(id int, order OrderRule) error,
	beforeReorder func(),
) uint64 {
	var revision uint64
	ticker := time.NewTicker(reorderTickInterval)
	defer ticker.Stop()
	lastSize := -1
	cleanPasses := 0
	noProgressTicks := 0
	stuckOnSkippedTicks := 0
	nonImprovingPutTicks := 0
	bestAtTarget := -1
	prevSkipped := -1
	prevAtTarget := -1
	readFailures := 0
	cyclePuts := 0
	for {
		select {
		case <-ticker.C:
			rules.Lock()
			size := len(rules.orders[resourceType])
			allDone := true
			for _, v := range rules.orders[resourceType] {
				if !v.done {
					allDone = false
					break
				}
			}
			if !(allDone && size > 0) {
				rules.Unlock()
				continue
			}
			sorted := sortOrders(rules.orders[resourceType])
			revision = rules.revisions[resourceType]
			rules.Unlock()

			// One GetAll per pass. This is the cache-friendly call
			// path: with the SDK cache-invalidation fix, any prior
			// PUT in this pass invalidates the parent collection so
			// subsequent passes always see the latest state.
			current, err := getCurrent()
			if err != nil {
				log.Printf("[ERROR] reorderAll: getCurrent failed for %s: %v", resourceType, err)
				// A failed read skips every exit condition below, so a read
				// that keeps failing would loop forever. Give up after
				// maxNoProgressTicks consecutive failures; any successful read
				// resets the count, so transient failures behave as before.
				// The waiting Create/Update then reads the rule back and
				// reports the API error to Terraform.
				readFailures++
				if readFailures >= maxNoProgressTicks {
					log.Printf("[WARN] reorderAll: %s — could not read the current rule order for %d consecutive passes, giving up. Last error: %v", resourceType, readFailures, err)
					return revision
				}
				continue
			}
			readFailures = 0
			count := 0
			for _, c := range current {
				if c.Order >= 1 {
					count++
				}
			}

			if beforeReorder != nil {
				beforeReorder()
			}

			// Terraform creates rules in batches. If this batch includes a
			// position that does not exist yet, let the next batch extend the
			// list before moving any rules. Moving the reachable subset now
			// makes later insertions displace it and repeats that work for
			// every batch. Keep every declared target unchanged; never invent
			// temporary positions or reserve slots for unmanaged rules.
			deferred := 0
			for _, v := range sorted {
				if v.Order.Order > count {
					deferred++
				}
			}
			if deferred > 0 {
				log.Printf("[INFO] reorderAll: deferring placement for %s: %d declared orders exceed the current %d orderable rules; releasing this batch without order PUTs so later creates can extend the list", resourceType, deferred, count)
				return revision
			}

			puterrs := 0
			skipped := 0
			putsIssued := 0
			alreadyAtTarget := 0
			var putDetails []string

			// Rules to place this pass, in declared order, with their target.
			var candidates []placement
			for _, v := range sorted {
				desired := v.Order

				// Rule POSTed but not yet visible to GetAll
				// (cache lag / API eventual consistency).
				// Skip and try next tick.
				if _, exists := current[v.ID]; !exists {
					skipped++
					continue
				}

				// Out-of-range desired position. Defer — usually
				// means more rules are still being POSTed and
				// count will grow on the next tick.
				if desired.Order < 1 || desired.Order > count {
					log.Printf("[INFO] reorderAll: deferring rule %d for %s (desired order %d, orderable range 1..%d)", v.ID, resourceType, desired.Order, count)
					skipped++
					continue
				}
				candidates = append(candidates, placement{id: v.ID, target: desired})
			}

			move := func(c placement) {
				cur := current[c.id]
				log.Printf("[INFO] reorderAll: %s — moving rule %d from order %d (rank %d) to order %d (rank %d)", resourceType, c.id, cur.Order, cur.Rank, c.target.Order, c.target.Rank)
				if err := updateOrder(c.id, c.target); err != nil {
					log.Printf("[ERROR] reorderAll: PUT for rule %d failed in %s: %v", c.id, resourceType, err)
					puterrs++
					return
				}
				putDetails = append(putDetails, fmt.Sprintf("rule %d (API order %d, declared order %d)", c.id, cur.Order, c.target.Order))
				putsIssued++
				// The move shifts the rules between the two positions; keep
				// this pass's view in step so the following rules are compared
				// against where they are now, not where they were when the list
				// was read. The next pass reads the list again.
				applyMove(current, c.id, c.target)
			}

			// Phase 1: leave the longest run of rules already in the right
			// relative order where they are, and move only the others. Moving
			// those into place also brings the run to its positions, so each
			// misplaced rule is moved once (a new rule, or two swapped rules)
			// instead of every rule it displaces.
			keep := inOrderRun(candidates, current)
			moved := map[int]bool{}
			if preferDirectPlacement(candidates, current, keep) {
				log.Printf("[INFO] reorderAll: %s — direct placement reaches the same declared targets with fewer order updates", resourceType)
			} else {
				for _, c := range candidates {
					if keep[c.id] || atTarget(current[c.id], c.target) {
						continue
					}
					move(c)
					moved[c.id] = true
				}
			}
			// Phase 2: move anything still not at its target (for example
			// when rules not registered in this cycle shifted the run).
			// Diff check: rules already at the right Order AND Rank get no
			// GET and no PUT, so a stable pass costs ONE GetAll and ZERO PUTs.
			for _, c := range candidates {
				if atTarget(current[c.id], c.target) {
					if !moved[c.id] {
						alreadyAtTarget++
					}
					continue
				}
				move(c)
			}

			log.Printf("[INFO] reorderAll: %s pass — registered=%d apiOrderable=%d alreadyAtTarget=%d putsIssued=%d skipped=%d errors=%d", resourceType, size, count, alreadyAtTarget, putsIssued, skipped, puterrs)

			// A cycle that has not moved any rule has no PUT in
			// flight, so one pass with every rule in place is
			// enough (for example order-only updates after a rule
			// above them was deleted).
			cyclePuts += putsIssued
			if putsIssued == 0 && skipped == 0 && puterrs == 0 && cyclePuts == 0 {
				log.Printf("[INFO] reorder complete for %s: %d rules, all already in place", resourceType, size)
				return revision
			}

			// Convergence: a clean pass = nothing to do AND no
			// errors AND size didn't grow. We require two
			// consecutive clean passes so any in-flight PUT from
			// the previous pass has time to be reflected by GetAll.
			cleanPass := putsIssued == 0 && skipped == 0 && puterrs == 0 && size == lastSize
			if cleanPass {
				cleanPasses++
				if cleanPasses >= 2 {
					log.Printf("[INFO] reorder complete for %s: %d rules", resourceType, size)
					return revision
				}
			} else {
				cleanPasses = 0
			}

			// Fast deadlock-breaker. The terraform-parallelism
			// deadlock looks like this:
			//
			//   - N rules registered and marked done.
			//   - All in-range rules are at their target (we did
			//     the PUTs).
			//   - K rules have declared orders > apiOrderable,
			//     so they're permanently skipped this cycle.
			//
			// With Create blocked on waitForReorder, the next
			// batch of POSTs can't extend apiOrderable until we
			// return. Returning here lets terraform schedule the
			// next batch; the new rule's registration will start
			// a fresh reorder cycle that picks up the previously
			// skipped ones (they stay in the registry — done=true
			// rules are diff-checked again on subsequent cycles).
			//
			// We require the condition to hold for two consecutive
			// passes so in-flight PUTs from the previous pass have
			// time to settle in the GetAll view (otherwise a
			// mid-flight pass might falsely report "nothing
			// changed" while a PUT is still propagating through
			// the cache).
			inRangeAtTarget := alreadyAtTarget == size-skipped
			progressedThisPass := putsIssued > 0 || (prevAtTarget >= 0 && alreadyAtTarget > prevAtTarget)
			if skipped > 0 && inRangeAtTarget && !progressedThisPass && puterrs == 0 {
				stuckOnSkippedTicks++
				if stuckOnSkippedTicks >= maxStuckOnSkippedTicks {
					log.Printf("[INFO] reorderAll: %s — %d in-range rule(s) at target, %d deferred (declared orders > apiOrderable=%d). Returning so the next batch of Creates can extend the orderable range.", resourceType, alreadyAtTarget, skipped, count)
					return revision
				}
			} else {
				stuckOnSkippedTicks = 0
			}

			// Oscillation breaker (SUP-4231). A PUT ping-pong — every pass
			// issues PUTs for the same rules because each placement
			// displaces a neighbour — never produces a clean pass and
			// counts as "progress" for the no-progress net below, so
			// without this check the loop runs until terraform's context
			// deadline kills it. A pass "improves" only when its at-target
			// count exceeds the best seen for the current registration set;
			// after maxNonImprovingPutTicks consecutive PUT-issuing passes
			// with no improvement, log the contested rules and return so
			// the apply completes (the unplaced rules show as drift on the
			// next plan).
			improved := alreadyAtTarget > bestAtTarget
			if improved {
				bestAtTarget = alreadyAtTarget
			}
			switch {
			case size != lastSize:
				// New registrations legitimately reshuffle the target
				// set; restart the oscillation window.
				bestAtTarget = alreadyAtTarget
				nonImprovingPutTicks = 0
			case putsIssued > 0 && !improved && skipped == 0 && puterrs == 0:
				nonImprovingPutTicks++
				if nonImprovingPutTicks >= maxNonImprovingPutTicks {
					log.Printf("[WARN] reorderAll: %s — %d pass(es) kept re-issuing order updates without converging (best %d/%d rules at target). The declared orders appear unsatisfiable (duplicate order values in the tenant, or rules not managed by Terraform occupying positions in the declared range). Giving up so the apply can complete; re-run terraform plan to see the remaining differences. Contested in the last pass: %s", resourceType, nonImprovingPutTicks, bestAtTarget, size, strings.Join(putDetails, "; "))
					return revision
				}
			default:
				nonImprovingPutTicks = 0
			}

			// Slower no-progress safety net for genuinely
			// impossible configurations (e.g. declared order=100
			// with only 1 orderable rule). Real progress = size
			// grew, OR a PUT was issued, OR the deferred set
			// shrank.
			realProgress := size > lastSize || putsIssued > 0
			if !realProgress && prevSkipped >= 0 && skipped < prevSkipped {
				realProgress = true
			}
			if realProgress {
				noProgressTicks = 0
			} else if skipped > 0 || puterrs > 0 {
				noProgressTicks++
			} else {
				noProgressTicks = 0
			}

			if noProgressTicks >= maxNoProgressTicks {
				log.Printf("[WARN] reorderAll: no progress for %d ticks for %s — giving up. Likely the declared orders exceed the count of orderable rules (gaps in order numbers, or earlier Create failures). Re-run terraform apply to retry.", noProgressTicks, resourceType)
				return revision
			}

			lastSize = size
			prevSkipped = skipped
			prevAtTarget = alreadyAtTarget
		default:
			time.Sleep(reorderTickInterval / 2)
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

type OrderRule struct {
	Order int
	Rank  int
}

func reorderWithBeforeReorder(order OrderRule, id int, resourceType string, getCurrent func() (map[int]OrderRule, error), updateOrder func(id int, order OrderRule) error, beforeReorder func()) {
	rules.Lock()
	shouldCallReorder := false
	if rules.orderer == nil {
		rules.orderer = map[string]int{}
		rules.reorderDone = map[string]chan struct{}{}
		rules.revisions = map[string]uint64{}
	}
	if rules.orders == nil {
		rules.orders = map[string]map[int]orderWithState{}
	}
	if _, ok := rules.orderer[resourceType]; ok {
		// Check if the previous reorder goroutine has already finished
		select {
		case <-rules.reorderDone[resourceType]:
			// Channel is closed — previous reorder finished.
			// Start a new reorder cycle for late-arriving rules.
			log.Printf("[INFO] previous reorder for %s completed, starting new cycle for rule:%d", resourceType, id)
			rules.reorderDone[resourceType] = make(chan struct{})
			shouldCallReorder = true
		default:
			// Reorder goroutine still running — just register, don't start another
			shouldCallReorder = false
		}
	} else {
		rules.orderer[resourceType] = id
		shouldCallReorder = true
		rules.reorderDone[resourceType] = make(chan struct{})
	}
	if rules.orders[resourceType] == nil {
		rules.orders[resourceType] = map[int]orderWithState{}
	}
	rules.orders[resourceType][id] = orderWithState{order, false}
	rules.revisions[resourceType]++
	// Capture the done channel while still holding the lock: concurrent
	// registrations write rules.reorderDone under the lock, so reading the
	// map after Unlock races with them.
	doneCh := rules.reorderDone[resourceType]
	rules.Unlock()
	if shouldCallReorder {
		log.Printf("[INFO] starting to reorder the rules, delegating to rule:%d, order:%d", id, order)
		go func() {
			for {
				revision := reorderAll(resourceType, getCurrent, updateOrder, beforeReorder)
				rules.Lock()
				if revision != rules.revisions[resourceType] {
					rules.Unlock()
					continue
				}
				// Closing the cycle and checking for new registrations must
				// be atomic with registration's decision to start a new one.
				close(doneCh)
				rules.Unlock()
				return
			}
		}()
	}
}

// waitForReorder blocks until the reorder goroutine for the given resource type
// has completed. All rules should call this after markOrderRuleAsDone and before
// reading the resource, to ensure the final order is reflected in the API.
func waitForReorder(resourceType string) {
	rules.Lock()
	ch := rules.reorderDone[resourceType]
	rules.Unlock()
	if ch != nil {
		<-ch
	}
}

// crudFunc is the signature of a resource CRUD function.
type crudFunc = func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics

// recordRegisteredOrder writes the position registered with the reorder
// engine for the rule in d (its API id is d.Id()) into state. waitForReorder
// returns before the rule reaches that position whenever the engine releases
// a batch early — placement deferred until later creates extend the list, the
// deadlock-breaker, the oscillation breaker, failed reads — so the Read that
// follows records a transient position, and nothing rewrites state once a
// later cycle moves the rule. A registered target never changes, so it is the
// position the apply converges to; a refresh reads the API and reports any
// divergence that remains. Resources without a rank attribute ignore the rank
// write.
func recordRegisteredOrder(d *schema.ResourceData, resourceType string) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return
	}
	rules.Lock()
	reg, ok := rules.orders[resourceType][id]
	rules.Unlock()
	if !ok {
		return
	}
	_ = d.Set("order", reg.order.Order)
	if reg.order.Rank != 0 {
		_ = d.Set("rank", reg.order.Rank)
	}
}

// readOrderedRule runs the resource's Read after waitForReorder and then
// records the registered position, see recordRegisteredOrder.
func readOrderedRule(ctx context.Context, d *schema.ResourceData, meta interface{}, read crudFunc, resourceType string) diag.Diagnostics {
	diags := read(ctx, d, meta)
	recordRegisteredOrder(d, resourceType)
	return diags
}

func reorder(order OrderRule, id int, resourceType string, getCurrent func() (map[int]OrderRule, error), updateOrder func(id int, order OrderRule) error) {
	reorderWithBeforeReorder(order, id, resourceType, getCurrent, updateOrder, nil)
}

// placement is a registered rule and the position a reorder pass places it at.
type placement struct {
	id     int
	target OrderRule
}

func atTarget(cur, target OrderRule) bool {
	return cur.Order == target.Order && cur.Rank == target.Rank
}

// preferDirectPlacement skips the relative-order phase only when a simulation
// proves that the final placement phase alone reaches every declared target
// with fewer PUTs. Both trials use copies of the same snapshot and move only
// registered rules to their unchanged HCL targets. Actual writes and final
// verification still use the existing reorder loop.
func preferDirectPlacement(candidates []placement, current map[int]OrderRule, keep map[int]bool) bool {
	if len(candidates) < 2 {
		return false
	}
	// Restrict this optimization to an ordinary contiguous list with one
	// unchanged rank. Preserve the existing path for rank transitions,
	// ambiguous positions, and special API order data.
	rank := candidates[0].target.Rank
	positions := make(map[int]bool, len(current))
	count, maxOrder := 0, 0
	for _, cur := range current {
		if cur.Order < 1 {
			continue
		}
		if cur.Rank != rank || positions[cur.Order] {
			return false
		}
		positions[cur.Order] = true
		count++
		if cur.Order > maxOrder {
			maxOrder = cur.Order
		}
	}
	if maxOrder != count {
		return false
	}
	previous := 0
	for _, c := range candidates {
		cur, exists := current[c.id]
		if !exists || cur.Order < 1 || cur.Rank != c.target.Rank || c.target.Rank != rank || c.target.Order <= previous || c.target.Order > count {
			return false
		}
		previous = c.target.Order
	}

	direct := maps.Clone(current)
	directPuts := 0
	for _, c := range candidates {
		if !atTarget(direct[c.id], c.target) {
			applyMove(direct, c.id, c.target)
			directPuts++
		}
	}
	if directPuts == 0 {
		return false
	}
	for _, c := range candidates {
		if !atTarget(direct[c.id], c.target) {
			return false
		}
	}

	relative := maps.Clone(current)
	relativePuts := 0
	move := func(c placement) {
		if !atTarget(relative[c.id], c.target) {
			applyMove(relative, c.id, c.target)
			relativePuts++
		}
	}
	for _, c := range candidates {
		if !keep[c.id] {
			move(c)
		}
	}
	for _, c := range candidates {
		move(c)
	}
	return directPuts < relativePuts
}

// inOrderRun returns the rules (candidates are in target order) that form the
// longest run whose current positions are already in that order: the rules
// that do not need to move for the others to be placed around them.
func inOrderRun(candidates []placement, current map[int]OrderRule) map[int]bool {
	// Longest strictly increasing subsequence of the current positions.
	tails := []int{} // tails[k]: index of the smallest tail of a run of length k+1
	prev := make([]int, len(candidates))
	for i, c := range candidates {
		pos := current[c.id].Order
		lo, hi := 0, len(tails)
		for lo < hi {
			mid := (lo + hi) / 2
			if current[candidates[tails[mid]].id].Order < pos {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		prev[i] = -1
		if lo > 0 {
			prev[i] = tails[lo-1]
		}
		if lo == len(tails) {
			tails = append(tails, i)
		} else {
			tails[lo] = i
		}
	}
	keep := make(map[int]bool, len(tails))
	if len(tails) > 0 {
		for i := tails[len(tails)-1]; i >= 0; i = prev[i] {
			keep[candidates[i].id] = true
		}
	}
	return keep
}

// applyMove updates a snapshot of rule positions for a rule moved to order:
// like the API, the rules between its old and new position shift by one.
func applyMove(current map[int]OrderRule, id int, to OrderRule) {
	from := current[id].Order
	for other, c := range current {
		if other == id || c.Order < 1 {
			continue
		}
		switch {
		case from < to.Order && c.Order > from && c.Order <= to.Order:
			c.Order--
		case from > to.Order && c.Order >= to.Order && c.Order < from:
			c.Order++
		default:
			continue
		}
		current[other] = c
	}
	current[id] = to
}

// ruleUpdateMode is how an update of an ordered rule writes the rule before
// handing its declared position to the reorder engine.
type ruleUpdateMode int

const (
	// ruleUpdateOrderOnly: only `order` changed. Nothing is written; the
	// reorder engine moves the rule (and skips it when it is already in place,
	// for example after a rule above it was deleted).
	ruleUpdateOrderOnly ruleUpdateMode = iota
	// ruleUpdateInPlace: other attributes changed, rank did not. They are
	// written with the rule left at its current position.
	ruleUpdateInPlace
	// ruleUpdateMoveToEnd: rank changed. The rule is written at the end of the
	// list with rank 7, where any rank is allowed, and the reorder engine then
	// places it.
	ruleUpdateMoveToEnd
)

// orderedRuleUpdateMode picks the update mode from the planned changes.
func orderedRuleUpdateMode(d *schema.ResourceData) ruleUpdateMode {
	if !d.HasChangesExcept("order") {
		return ruleUpdateOrderOnly
	}
	if !d.HasChange("rank") {
		return ruleUpdateInPlace
	}
	return ruleUpdateMoveToEnd
}

// updateOrderedRule writes an update of an ordered rule according to its
// update mode, before the caller hands the declared position to the reorder
// engine. An order-only update writes nothing.
//
// known is the rule's position as already read by the caller (nil when not
// available); current reads it again, and is used when known is nil or the API
// rejected the position. end reads the end of the rule list (see endOfList).
// write sends the update at a position.
func updateOrderedRule(d *schema.ResourceData, known *OrderRule, current, end func() (OrderRule, error), write func(OrderRule) error) error {
	switch orderedRuleUpdateMode(d) {
	case ruleUpdateOrderOnly:
		log.Printf("[INFO] rule %s: only the order changed; leaving the move to the reorder engine", d.Id())
		return nil
	case ruleUpdateMoveToEnd:
		return writeAtFreshPosition(func(bool) (OrderRule, error) { return end() }, write)
	default:
		return writeAtFreshPosition(func(fresh bool) (OrderRule, error) {
			if known != nil && !fresh {
				return *known, nil
			}
			return current()
		}, write)
	}
}

// endOfList is the last position of a rule list, given the orders of its rules,
// with rank 7: a rule written there is allowed whatever the ranks of the rules
// above it.
func endOfList(orders []int) OrderRule {
	last := 1
	for _, o := range orders {
		if o > last {
			last = o
		}
	}
	return OrderRule{Order: last, Rank: 7}
}

// maxStaleOrderRetries bounds writeAtFreshPosition's retries.
const maxStaleOrderRetries = 3

var staleOrderPattern = regexp.MustCompile(`not allowed at order [0-9]+`)

// isStaleOrderError reports the API rejecting a position that does not exist
// in the rule list, e.g. "Rule is not allowed at order 20" after a rule of the
// same type was deleted while the request was in flight.
func isStaleOrderError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && staleOrderPattern.MatchString(err.Error())
}

// writeAtFreshPosition writes a rule at a position read just before the write.
// position is called with fresh=false first and may reuse a value it already
// has; when the API rejects the position because the list changed in between,
// it is called again with fresh=true and the write is retried, up to
// maxStaleOrderRetries times. Other errors are returned as they are.
func writeAtFreshPosition(position func(fresh bool) (OrderRule, error), write func(OrderRule) error) error {
	var err error
	for attempt := 0; attempt <= maxStaleOrderRetries; attempt++ {
		var pos OrderRule
		if pos, err = position(attempt > 0); err != nil {
			return err
		}
		if err = write(pos); !isStaleOrderError(err) {
			return err
		}
		log.Printf("[WARN] order %d is no longer valid (the rule list changed during the request); reading the position again (attempt %d of %d)", pos.Order, attempt+1, maxStaleOrderRetries)
	}
	return err
}
