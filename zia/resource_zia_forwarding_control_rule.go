package zia

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
				Description: "Determines whether the Firewall Filtering policy rule is enabled or disabled",
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
			"app_service_groups":             setIDsSchemaTypeCustom(nil, "list of application service groups"),
			"proxy_gateway":                  setIdNameSchemaCustom(1, "The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method."),
			"zpa_gateway":                    setIdNameSchemaCustom(1, "The ZPA Server Group for which this rule is applicable. Only the Server Groups that are associated with the selected Application Segments are allowed. This field is applicable only for the ZPA forwarding method."),
			"zpa_app_segments":               setExtIDNameSchemaCustom(intPtr(255), "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method."),
			"zpa_application_segments":       setIDsSchemaTypeCustom(intPtr(255), "List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
			"zpa_application_segment_groups": setIDsSchemaTypeCustom(intPtr(255), "List of ZPA Application Segment Groups for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)."),
			"dest_countries":                 getDestinationCountries(),
		},
	}
}

func validateForwardingRuleConstraints(d *schema.ResourceData) error {
	forwardMethod := d.Get("forward_method").(string)
	ruleType := d.Get("type").(string)

	isSet := func(attr string) bool {
		_, ok := d.GetOk(attr)
		return ok
	}

	if ruleType == "FORWARDING" {
		switch forwardMethod {
		case "ZPA":
			requiredAttrs := []string{"zpa_app_segments", "zpa_gateway"}
			var missingAttrs []string
			for _, attr := range requiredAttrs {
				if !isSet(attr) {
					missingAttrs = append(missingAttrs, attr)
				}
			}
			if len(missingAttrs) > 0 {
				return fmt.Errorf("the following attributes are required for ZPA forwarding: %v", missingAttrs)
			}

		case "DIRECT":
			prohibitedAttrs := []string{"zpa_gateway", "proxy_gateway", "zpa_app_segments", "zpa_application_segments", "zpa_application_segment_groups"}
			for _, attr := range prohibitedAttrs {
				if isSet(attr) {
					return fmt.Errorf("%s attribute cannot be set when type is 'FORWARDING' and forward_method is 'DIRECT'", attr)
				}
			}

		case "PROXYCHAIN":
			if !isSet("proxy_gateway") {
				return fmt.Errorf("proxy gateway is mandatory for Proxy Chaining forwarding")
			}
			prohibitedAttrs := []string{"zpa_gateway", "zpa_app_segments", "zpa_application_segments", "zpa_application_segment_groups"}
			for _, attr := range prohibitedAttrs {
				if isSet(attr) {
					return fmt.Errorf("%s attribute cannot be set when type is 'FORWARDING' and forward_method is 'PROXYCHAIN'", attr)
				}
			}
		}
	}

	return nil
}

func validatePredefinedRules(req forwarding_rules.ForwardingRules) error {
	if req.Name == "Client Connector Traffic Direct" || req.Name == "ZPA Pool For Stray Traffic" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	if req.Name == "ZIA Inspected ZPA Apps" || req.Name == "Fallback mode of ZPA Forwarding" {
		return fmt.Errorf("predefined rule '%s' cannot be deleted", req.Name)
	}
	return nil
}

func resourceForwardingControlRuleCreate(d *schema.ResourceData, m interface{}) error {

	if err := validateForwardingRuleConstraints(d); err != nil {
		return err
	}

	zClient := m.(*Client)
	req := expandForwardingControlRule(d)
	log.Printf("[INFO] Creating zia forwarding control rule\n%+v\n", req)

	if err := validatePredefinedRules(req); err != nil {
		return err
	}

	resp, err := zClient.forwarding_rules.Create(&req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created zia ip source groups request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)
	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
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

	processedDestCountries := make([]string, len(resp.DestCountries))
	for i, country := range resp.DestCountries {
		processedDestCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
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
	_ = d.Set("dest_countries", processedDestCountries)
	_ = d.Set("res_categories", resp.ResCategories)

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationsGroups)); err != nil {
		return err
	}

	if err := d.Set("ec_groups", flattenIDs(resp.ECGroups)); err != nil {
		return err
	}

	if err := d.Set("departments", flattenIDs(resp.Departments)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenIDs(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("users", flattenIDs(resp.Users)); err != nil {
		return err
	}

	if err := d.Set("src_ip_groups", flattenIDs(resp.SrcIpGroups)); err != nil {
		return err
	}

	if err := d.Set("src_ipv6_groups", flattenIDs(resp.SrcIpv6Groups)); err != nil {
		return err
	}

	if err := d.Set("dest_ip_groups", flattenIDs(resp.DestIpGroups)); err != nil {
		return err
	}

	if err := d.Set("dest_ipv6_groups", flattenIDs(resp.DestIpv6Groups)); err != nil {
		return err
	}

	if err := d.Set("nw_services", flattenIDs(resp.NwServices)); err != nil {
		return err
	}

	if err := d.Set("nw_service_groups", flattenIDs(resp.NwServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_application_groups", flattenIDs(resp.NwApplicationGroups)); err != nil {
		return err
	}

	if err := d.Set("nw_application_groups", flattenIDs(resp.NwApplicationGroups)); err != nil {
		return err
	}

	if err := d.Set("app_service_groups", flattenIDs(resp.AppServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
		return err
	}

	if err := d.Set("proxy_gateway", flattenIDNameSet(resp.ProxyGateway)); err != nil {
		return err
	}

	if err := d.Set("zpa_gateway", flattenIDNameSet(resp.ZPAGateway)); err != nil {
		return err
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return err
	}

	return nil
}

func resourceForwardingControlRuleUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] forwarding control rule ID not set: %v\n", id)
	}

	if err := validateForwardingRuleConstraints(d); err != nil {
		return err
	}

	log.Printf("[INFO] Updating zia forwarding control rule ID: %v\n", id)
	req := expandForwardingControlRule(d)

	if err := validatePredefinedRules(req); err != nil {
		return err
	}

	if _, err := zClient.forwarding_rules.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, err := zClient.forwarding_rules.Update(id, &req); err != nil {
		return err
	}
	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return resourceForwardingControlRuleRead(d, m)
}

func resourceForwardingControlRuleDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("forwarding control rule ID not set: %v", id)
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := zClient.forwarding_rules.Get(id)
	if err != nil {
		return fmt.Errorf("error retrieving forwarding control rule %d: %v", id, err)
	}

	// Validate if the rule can be deleted
	if err := validatePredefinedRules(*rule); err != nil {
		return err
	}

	log.Printf("[INFO] Deleting forwarding control rule ID: %v", id)
	if _, err := zClient.forwarding_rules.Delete(id); err != nil {
		return fmt.Errorf("error deleting forwarding control rule %d: %v", id, err)
	}

	d.SetId("")
	log.Printf("[INFO] Forwarding control rule deleted")
	// Trigger activation after creating the rule label
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
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
		DestCountries:       processedDestCountries,
		Locations:           expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups:     expandIDNameExtensionsSet(d, "location_groups"),
		Departments:         expandIDNameExtensionsSet(d, "departments"),
		Groups:              expandIDNameExtensionsSet(d, "groups"),
		Users:               expandIDNameExtensionsSet(d, "users"),
		SrcIpGroups:         expandIDNameExtensionsSet(d, "src_ip_groups"),
		DestIpGroups:        expandIDNameExtensionsSet(d, "dest_ip_groups"),
		NwServices:          expandIDNameExtensionsSet(d, "nw_services"),
		AppServiceGroups:    expandIDNameExtensionsSet(d, "app_service_groups"),
		NwServiceGroups:     expandIDNameExtensionsSet(d, "nw_service_groups"),
		NwApplicationGroups: expandIDNameExtensionsSet(d, "nw_application_groups"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
		ECGroups:            expandIDNameExtensionsSet(d, "ec_groups"),
		ProxyGateway:        expandIDNameSet(d, "proxy_gateway"),
		ZPAGateway:          expandIDNameSet(d, "zpa_gateway"),
		ZPAAppSegments:      expandZPAAppSegmentSet(d, "zpa_app_segments"),
	}

	return result
}
