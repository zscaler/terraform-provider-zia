package zia

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/forwarding_rules"
)

// var (
// 	forwardingControlLock          sync.Mutex
// 	forwardingControlStartingOrder int
// )

func resourceForwardingControlRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceForwardingControlRuleCreate,
		Read:   resourceForwardingControlRuleRead,
		Update: resourceForwardingControlRuleUpdate,
		Delete: resourceForwardingControlRuleDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := zClient.forwarding_rules.GetByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("rule_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique identifier assigned to the forwarding rule",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the forwarding rule",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the forwarding rule",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rule type selected from the available options",
				ValidateFunc: validation.StringInSlice([]string{
					"FIREWALL",
					"DNS",
					"DNAT",
					"SNAT",
					"FORWARDING",
					"INTRUSION_PREVENTION",
					"EC_DNS",
					"EC_RDR",
					"EC_SELF",
					"DNS_RESPONSE",
				}, false),
			},
			"forward_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of traffic forwarding method selected from the available options",
				ValidateFunc: validation.StringInSlice([]string{
					"INVALID",
					"DIRECT",
					"PROXYCHAIN",
					"ZIA",
					"ZPA",
					"ECZPA",
					"ECSELF",
					"DROP",
				}, false),
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The order of execution for the forwarding rule order",
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank assigned to the forwarding rule",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates whether the forwarding rule is enabled or disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"src_ips": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.",
			},
			"dest_addresses": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.",
			},
			"dest_ip_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.",
			},
			"res_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination domain categories to which the rule applies",
			},
			"zpa_gateway": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "ID of the ZPA Gateway.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the ZPA Gateway.",
						},
					},
				},
			},

			"zpa_app_segments": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of the application segments.",
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
			},
			"locations":                      setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations."),
			"location_groups":                setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the location groups to which the forwarding rule applies"),
			"ec_groups":                      setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies"),
			"departments":                    setIDsSchemaTypeCustom(intPtr(140000), "list of departments for which rule must be applied"),
			"groups":                         setIDsSchemaTypeCustom(intPtr(8), "list of groups for which rule must be applied"),
			"users":                          setIDsSchemaTypeCustom(intPtr(4), "list of users for which rule must be applied"),
			"src_ip_groups":                  setIDsSchemaTypeCustom(nil, "Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group"),
			"src_ipv6_groups":                setIDsSchemaTypeCustom(nil, "Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group"),
			"dest_ip_groups":                 setIDsSchemaTypeCustom(nil, "User-defined destination IP address groups to which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group"),
			"dest_ipv6_groups":               setIDsSchemaTypeCustom(nil, "Destination IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group"),
			"nw_services":                    setIDsSchemaTypeCustom(intPtr(1024), "User-defined network services to which the rule applies. If not set, the rule is not restricted to a specific network service."),
			"nw_service_groups":              setIDsSchemaTypeCustom(nil, "User-defined network service group to which the rule applies. If not set, the rule is not restricted to a specific network service group."),
			"labels":                         setIDsSchemaTypeCustom(intPtr(1), "Labels that are applicable to the rule"),
			"nw_application_groups":          setIDsSchemaTypeCustom(nil, "User-defined network service application groups to which the rule applied. If not set, the rule is not restricted to a specific network service application group."),
			"time_windows":                   setIDsSchemaTypeCustom(intPtr(2), "list of time interval during which rule must be enforced."),
			"device_groups":                  setIDsSchemaTypeCustom(nil, "Name-ID pairs of device groups for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation."),
			"devices":                        setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which the rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation."),
			"proxy_gateway":                  setIDsSchemaTypeCustom(intPtr(1), "The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method."),
			"zpa_application_segments":       setIDsSchemaTypeCustom(intPtr(255), "List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
			"zpa_application_segment_groups": setIDsSchemaTypeCustom(intPtr(255), "List of ZPA Application Segment Groups for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
			"nw_applications":                getNwApplications(),
			"dest_countries":                 getDestinationCountries(),
		},
	}
}

func resourceForwardingControlRuleCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandForwardingControlRule(d)
	log.Printf("[INFO] Creating zia forwarding control rule\n%+v\n", req)

	resp, err := zClient.forwarding_rules.Create(&req)
	if err != nil {
		return err
	}
	// if err := validateForwardMethodAttrs(d); err != nil {
	// 	return err
	// }

	log.Printf("[INFO] Created zia ip source groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)
	return resourceForwardingControlRuleRead(d, m)
}

func resourceForwardingControlRuleRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no zia forwarding control rule id is set")
	}
	resp, err := zClient.forwarding_rules.Get(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing forwarding control rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting forwarding control rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("forward_method", resp.ForwardMethod)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("state", resp.State)
	_ = d.Set("type", resp.Type)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", resp.DestCountries)
	_ = d.Set("res_categories", resp.ResCategories)
	_ = d.Set("nw_applications", resp.NwApplications)

	if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationsGroups)); err != nil {
		return err
	}

	if err := d.Set("ec_groups", flattenIDNameExtensions(resp.ECGroups)); err != nil {
		return err
	}

	if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
		return err
	}

	if err := d.Set("time_windows", flattenIDNameExtensions(resp.TimeWindows)); err != nil {
		return err
	}

	if err := d.Set("src_ip_groups", flattenIDNameExtensions(resp.SrcIpGroups)); err != nil {
		return err
	}

	if err := d.Set("src_ipv6_groups", flattenIDNameExtensions(resp.SrcIpv6Groups)); err != nil {
		return err
	}

	if err := d.Set("dest_ip_groups", flattenIDNameExtensions(resp.DestIpGroups)); err != nil {
		return err
	}

	if err := d.Set("dest_ipv6_groups", flattenIDNameExtensions(resp.DestIpv6Groups)); err != nil {
		return err
	}

	if err := d.Set("nw_services", flattenIDNameExtensions(resp.NwServices)); err != nil {
		return err
	}

	if err := d.Set("nw_service_groups", flattenIDNameExtensions(resp.NwServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_application_groups", flattenIDNameExtensions(resp.NwApplicationGroups)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
		return err
	}

	if err := d.Set("proxy_gateway", flattenIDNameExtensions(resp.ProxyGateway)); err != nil {
		return err
	}

	if err := d.Set("zpa_gateway", flattenZPAGateway(resp.ZPAGateway)); err != nil {
		return err
	}

	if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
		return err
	}

	if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
		return err
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return err
	}

	// if err := d.Set("zpa_application_segments", flattenZPAApplicationSegmentsSimple(resp.ZPAApplicationSegments)); err != nil {
	// 	return err
	// }

	// if err := d.Set("zpa_application_segment_groups", flattenZPAApplicationSegmentGroupsSimple(resp.ZPAApplicationSegmentGroups)); err != nil {
	// 	return err
	// }
	return nil
}

func resourceForwardingControlRuleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] forwarding control rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia forwarding control rule ID: %v\n", id)
	req := expandForwardingControlRule(d)
	if _, err := zClient.forwarding_rules.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := zClient.forwarding_rules.Update(id, &req); err != nil {
		return err
	}

	return resourceForwardingControlRuleRead(d, m)
}

func resourceForwardingControlRuleDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] forwarding control rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting forwarding control rule ID: %v\n", (d.Id()))

	if _, err := zClient.forwarding_rules.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] forwarding control rule deleted")
	return nil
}

func expandForwardingControlRule(d *schema.ResourceData) forwarding_rules.ForwardingRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Process the DestCountries to add the prefix where needed
	rawDestCountries := SetToStringList(d, "dest_countries")
	processedDestCountries := make([]string, len(rawDestCountries))
	for i, country := range rawDestCountries {
		if country != "ANY" && country != "NONE" && len(country) == 2 { // Assuming the 2 letter code is an ISO Alpha-2 Code
			processedDestCountries[i] = "COUNTRY_" + country
		} else {
			processedDestCountries[i] = country
		}
	}

	result := forwarding_rules.ForwardingRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Order:               d.Get("order").(int),
		Rank:                d.Get("rank").(int),
		Type:                d.Get("type").(string),
		State:               d.Get("state").(string),
		ForwardMethod:       d.Get("forward_method").(string),
		ResCategories:       SetToStringList(d, "res_categories"),
		SrcIps:              SetToStringList(d, "src_ips"),
		DestAddresses:       SetToStringList(d, "dest_addresses"),
		DestIpCategories:    SetToStringList(d, "dest_ip_categories"),
		NwApplications:      SetToStringList(d, "nw_applications"),
		DestCountries:       processedDestCountries,
		Locations:           expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups:     expandIDNameExtensionsSet(d, "location_groups"),
		Departments:         expandIDNameExtensionsSet(d, "departments"),
		Groups:              expandIDNameExtensionsSet(d, "groups"),
		Users:               expandIDNameExtensionsSet(d, "users"),
		TimeWindows:         expandIDNameExtensionsSet(d, "time_windows"),
		SrcIpGroups:         expandIDNameExtensionsSet(d, "src_ip_groups"),
		DestIpGroups:        expandIDNameExtensionsSet(d, "dest_ip_groups"),
		NwServices:          expandIDNameExtensionsSet(d, "nw_services"),
		NwServiceGroups:     expandIDNameExtensionsSet(d, "nw_service_groups"),
		NwApplicationGroups: expandIDNameExtensionsSet(d, "nw_application_groups"),
		DeviceGroups:        expandIDNameExtensionsSet(d, "device_groups"),
		Devices:             expandIDNameExtensionsSet(d, "devices"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
		ECGroups:            expandIDNameExtensionsSet(d, "ec_groups"),
		ProxyGateway:        expandIDNameExtensionsSet(d, "proxy_gateway"),
		ZPAGateway:          expandZPAGateway(d, "zpa_gateway"),
		ZPAAppSegments:      expandZPAAppSegmentSet(d, "zpa_app_segments"),
		//ZPAApplicationSegments:      expandZPAApplicationSegmentSet(d, "zpa_application_segments"),
		//ZPAApplicationSegmentGroups: expandZPApplicationSegmentGroupSet(d, "zpa_application_segment_groups"),
	}
	return result
}

func flattenZPAGateway(gateway forwarding_rules.ZPAGateway) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":   gateway.ID,
			"name": gateway.Name,
		},
	}
}

func expandZPAGateway(d *schema.ResourceData, key string) forwarding_rules.ZPAGateway {
	listInterface, exists := d.GetOk(key)
	if !exists || len(listInterface.([]interface{})) == 0 {
		return forwarding_rules.ZPAGateway{}
	}

	gatewayMap := listInterface.([]interface{})[0].(map[string]interface{})

	return forwarding_rules.ZPAGateway{
		ID:   gatewayMap["id"].(int),
		Name: gatewayMap["name"].(string),
	}
}

func expandZPAAppSegmentSet(d *schema.ResourceData, key string) []forwarding_rules.ZPAAppSegments {
	setInterface, exists := d.GetOk(key)
	if !exists {
		return nil
	}

	inputSet := setInterface.(*schema.Set).List()
	var result []forwarding_rules.ZPAAppSegments
	for _, item := range inputSet {
		itemMap := item.(map[string]interface{})
		name := itemMap["name"].(string)
		externalID := itemMap["external_id"].(string)

		segment := forwarding_rules.ZPAAppSegments{
			Name:       name,
			ExternalID: externalID,
		}
		result = append(result, segment)
	}
	return result
}

func flattenZPAAppSegmentsSimple(list []forwarding_rules.ZPAAppSegments) []interface{} {
	flattenedList := make([]interface{}, 0, len(list))
	for _, val := range list {
		r := map[string]interface{}{
			"name":        val.Name,
			"external_id": val.ExternalID,
		}
		flattenedList = append(flattenedList, r)
	}
	return flattenedList
}

/*
func flattenZPAApplicationSegmentsSimple(list []forwarding_rules.ZPAApplicationSegments) []interface{} {
	flattenedList := make([]interface{}, 0, len(list))
	for _, val := range list {
		r := map[string]interface{}{
			"name": val.Name,
			// "external_id": val.ExternalID,
		}
		flattenedList = append(flattenedList, r)
	}
	return flattenedList
}

func flattenZPAApplicationSegmentGroupsSimple(list []forwarding_rules.ZPAApplicationSegmentGroups) []interface{} {
	flattenedList := make([]interface{}, 0, len(list))
	for _, val := range list {
		r := map[string]interface{}{
			"name": val.Name,
			// "external_id": val.ExternalID,
		}
		flattenedList = append(flattenedList, r)
	}
	return flattenedList
}
*/