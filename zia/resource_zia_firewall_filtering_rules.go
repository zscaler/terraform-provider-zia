package zia

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
)

var (
	firewallFilteringLock          sync.Mutex
	firewallFilteringStartingOrder int
)

func resourceFirewallFilteringRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallFilteringRulesCreate,
		Read:   resourceFirewallFilteringRulesRead,
		Update: resourceFirewallFilteringRulesUpdate,
		Delete: resourceFirewallFilteringRulesDelete,
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
					resp, err := zClient.filteringrules.GetByName(id)
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the Firewall Filtering policy rule",
				ValidateFunc: validation.StringLenBetween(0, 31),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional information about the rule",
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule order number of the Firewall Filtering policy rule",
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the Firewall Filtering policy rule",
			},
			"enable_full_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action the Firewall Filtering policy rule takes when packets match the rule",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK_DROP",
					"BLOCK_RESET",
					"BLOCK_ICMP",
					"EVAL_NWAPP",
				}, false),
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dest_ip_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"default_rule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, the default rule is applied",
			},
			"predefined": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, a predefined rule is applied",
			},
			"locations":             setIDsSchemaTypeCustom(intPtr(8), "list of locations for which rule must be applied"),
			"location_groups":       setIDsSchemaTypeCustom(intPtr(32), "list of locations groups"),
			"users":                 setIDsSchemaTypeCustom(intPtr(4), "list of users for which rule must be applied"),
			"groups":                setIDsSchemaTypeCustom(intPtr(8), "list of groups for which rule must be applied"),
			"departments":           setIDsSchemaTypeCustom(intPtr(140000), "list of departments for which rule must be applied"),
			"time_windows":          setIDsSchemaTypeCustom(intPtr(2), "The time interval in which the Firewall Filtering policy rule applies"),
			"labels":                setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"device_groups":         setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":               setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"src_ip_groups":         setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"dest_ip_groups":        setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"app_service_groups":    setIDsSchemaTypeCustom(nil, "list of application service groups"),
			"app_services":          setIDsSchemaTypeCustom(nil, "list of application services"),
			"nw_application_groups": setIDsSchemaTypeCustom(nil, "list of nw application groups"),
			"nw_service_groups":     setIDsSchemaTypeCustom(nil, "list of nw service groups"),
			"workload_groups":       setIdNameSchemaCustom(255, "The list of preconfigured workload groups to which the policy must be applied"),
			"nw_services":           setIDsSchemaTypeCustom(intPtr(1024), "list of nw services"),
			"zpa_app_segments":      setExtIDNameSchemaCustom(intPtr(255), "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method."),
			"dest_countries":        getDestinationCountries(),
			"nw_applications":       getCloudFirewallNwApplications(),
			"device_trust_levels":   getDeviceTrustLevels(),
		},
	}
}

func validateFirewallRule(req filteringrules.FirewallFilteringRules) error {
	if req.Name == "Office 365 One Click Rule" || req.Name == "UCaaS One Click Rule" {
		return fmt.Errorf("deletion of the predefined rule '%s' is not allowed", req.Name)
	}
	if req.Name == "Block All IPv6" || req.Name == "Block malicious IPs and domains" {
		return fmt.Errorf("deletion of the predefined rule '%s' is not allowed", req.Name)
	}
	if req.Name == "Default Firewall Filtering Rule" {
		return fmt.Errorf("deletion of the predefined rule '%s' is not allowed", req.Name)
	}
	return nil
}

func resourceFirewallFilteringRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	req := expandFirewallFilteringRules(d)
	log.Printf("[INFO] Creating zia firewall filtering rule\n%+v\n", req)

	if err := validateFirewallRule(req); err != nil {
		return err
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		firewallFilteringLock.Lock()
		if firewallFilteringStartingOrder == 0 {
			list, _ := zClient.filteringrules.GetAll()
			for _, r := range list {
				if r.Order > firewallFilteringStartingOrder {
					firewallFilteringStartingOrder = r.Order
				}
			}
			if firewallFilteringStartingOrder == 0 {
				firewallFilteringStartingOrder = 1
			}
		}
		firewallFilteringLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = firewallFilteringStartingOrder
		resp, err := zClient.filteringrules.Create(&req)
		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, order, req.Rank, currentOrderVsRankWording(zClient), err)
				}
				if time.Since(start) < timeout {
					log.Printf("[INFO] Creating firewall filtering rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return fmt.Errorf("error creating resource: %s", err)
		}

		log.Printf("[INFO] Created zia firewall filtering rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "firewall_filtering_rules", func() (int, error) {
			list, err := zClient.filteringrules.GetAll()
			return len(list), err
		}, func(id, order int) error {
			rule, err := zClient.filteringrules.Get(id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = zClient.filteringrules.Update(id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		err = resourceFirewallFilteringRulesRead(d, m)
		if err != nil {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return err
		}
		markOrderRuleAsDone(resp.ID, "firewall_filtering_rules")
		break
	}
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return nil
}

func resourceFirewallFilteringRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no zia firewall filtering rule id is set")
	}

	resp, err := zClient.filteringrules.Get(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	processedDestCountries := make([]string, len(resp.DestCountries))
	for i, country := range resp.DestCountries {
		processedDestCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	log.Printf("[INFO] Getting firewall filtering rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("enable_full_logging", resp.EnableFullLogging)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
	_ = d.Set("action", resp.Action)
	_ = d.Set("state", resp.State)
	_ = d.Set("description", resp.Description)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", processedDestCountries)
	_ = d.Set("nw_applications", resp.NwApplications)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationsGroups)); err != nil {
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

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return err
	}

	if err := d.Set("src_ip_groups", flattenIDs(resp.SrcIpGroups)); err != nil {
		return err
	}

	if err := d.Set("dest_ip_groups", flattenIDs(resp.DestIpGroups)); err != nil {
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

	if err := d.Set("app_services", flattenIDs(resp.AppServices)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
		return err
	}
	if err := d.Set("app_service_groups", flattenIDs(resp.AppServiceGroups)); err != nil {
		return err
	}

	if err := d.Set("device_groups", flattenIDs(resp.DeviceGroups)); err != nil {
		return err
	}

	if err := d.Set("devices", flattenIDs(resp.Devices)); err != nil {
		return err
	}
	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return fmt.Errorf("error setting workload_groups: %s", err)
	}
	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return err
	}
	return nil
}

func resourceFirewallFilteringRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule ID not set: %v\n", id)
		return fmt.Errorf("firewall filtering rule ID not set")
	}
	log.Printf("[INFO] Updating firewall filtering rule ID: %v\n", id)
	req := expandFirewallFilteringRules(d)
	if err := validateFirewallRule(req); err != nil {
		return err
	}
	if _, err := zClient.filteringrules.Get(id); err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := zClient.filteringrules.Update(id, &req)
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating firewall filtering rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return fmt.Errorf("error updating resource: %s", err)
		}

		reorder(req.Order, req.ID, "firewall_filtering_rules", func() (int, error) {
			list, err := zClient.filteringrules.GetAll()
			return len(list), err
		}, func(id, order int) error {
			rule, err := zClient.filteringrules.Get(id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = zClient.filteringrules.Update(id, rule)
			return err
		})

		err = resourceFirewallFilteringRulesRead(d, m)
		if err != nil {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return err
		}
		markOrderRuleAsDone(req.ID, "firewall_filtering_rules")
		break
	}
	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return nil
}

func resourceFirewallFilteringRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule not set: %v\n", id)
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := zClient.filteringrules.Get(id)
	if err != nil {
		return fmt.Errorf("error retrieving firewall filtering rule %d: %v", id, err)
	}

	// Validate if the rule can be deleted
	if err := validateFirewallRule(*rule); err != nil {
		return err
	}

	log.Printf("[INFO] Deleting firewall filtering rule ID: %v\n", (d.Id()))
	if _, err := zClient.filteringrules.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] firewall filtering rule deleted")

	if activationErr := triggerActivation(zClient); activationErr != nil {
		return activationErr
	}
	return nil
}

func expandFirewallFilteringRules(d *schema.ResourceData) filteringrules.FirewallFilteringRules {
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

	result := filteringrules.FirewallFilteringRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Order:               d.Get("order").(int),
		Rank:                d.Get("rank").(int),
		Action:              d.Get("action").(string),
		State:               d.Get("state").(string),
		Description:         d.Get("description").(string),
		SrcIps:              SetToStringList(d, "src_ips"),
		DestAddresses:       SetToStringList(d, "dest_addresses"),
		DestIpCategories:    SetToStringList(d, "dest_ip_categories"),
		DeviceTrustLevels:   SetToStringList(d, "device_trust_levels"),
		DestCountries:       processedDestCountries,
		NwApplications:      SetToStringList(d, "nw_applications"),
		EnableFullLogging:   d.Get("enable_full_logging").(bool),
		DefaultRule:         d.Get("default_rule").(bool),
		Predefined:          d.Get("predefined").(bool),
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
		AppServices:         expandIDNameExtensionsSet(d, "app_services"),
		AppServiceGroups:    expandIDNameExtensionsSet(d, "app_service_groups"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
		DeviceGroups:        expandIDNameExtensionsSet(d, "device_groups"),
		Devices:             expandIDNameExtensionsSet(d, "devices"),
		WorkloadGroups:      expandWorkloadGroups(d, "workload_groups"),
		ZPAAppSegments:      expandZPAAppSegmentSet(d, "zpa_app_segments"),
	}
	return result
}
