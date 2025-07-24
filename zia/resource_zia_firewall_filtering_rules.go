package zia

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
)

var (
	firewallFilteringLock          sync.Mutex
	firewallFilteringStartingOrder int
)

func resourceFirewallFilteringRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallFilteringRulesCreate,
		ReadContext:   resourceFirewallFilteringRulesRead,
		UpdateContext: resourceFirewallFilteringRulesUpdate,
		DeleteContext: resourceFirewallFilteringRulesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := filteringrules.GetByName(ctx, service, id)
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Firewall Filtering policy rule",
				// ValidateFunc: validation.StringLenBetween(0, 31),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional information about the rule",
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"order": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Rule order number. If omitted, the rule will be added to the end of the rule set.",
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Destination addresses. Supports IPv4, FQDNs, or wildcard FQDNs",
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
			"dest_countries":        getISOCountryCodes(),
			"source_countries":      getISOCountryCodes(),
			"nw_applications":       getCloudApplications(),
			"device_trust_levels":   getDeviceTrustLevels(),
		},
	}
}

/*
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
*/

func beforeReorderFirewallFilteringRules(ctx context.Context, service *zscaler.Service) func() {
	return func() {
		log.Printf("[INFO] beforeReorderFirewallFilteringRules")
		// get all predefined rules and set their order to come first
		rules, err := filteringrules.GetAll(ctx, service)
		if err != nil {
			log.Printf("[ERROR] beforeReorderFirewallFilteringRules: %v", err)
		}
		// first order predefined rules by their order
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Order < rules[j].Order
		})
		order := 1
		for _, r := range rules {
			if r.Predefined {
				r.Order = order
				_, err = filteringrules.Update(ctx, service, r.ID, &r)
				if err != nil {
					log.Printf("[ERROR] beforeReorderFirewallFilteringRules: %v", err)
				}
				order++
			}
		}
	}
}

func resourceFirewallFilteringRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFirewallFilteringRules(d)
	log.Printf("[INFO] Creating zia firewall filtering rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		firewallFilteringLock.Lock()
		if firewallFilteringStartingOrder == 0 {
			list, _ := filteringrules.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > firewallFilteringStartingOrder {
					firewallFilteringStartingOrder = r.Order
				}
			}
			if firewallFilteringStartingOrder == 0 {
				firewallFilteringStartingOrder = 1
			} else {
				firewallFilteringStartingOrder++
			}
		}
		firewallFilteringLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = firewallFilteringStartingOrder

		resp, err := filteringrules.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, order, req.Rank, currentOrderVsRankWording(ctx, zClient), err))
				}
				if time.Since(start) < timeout {
					log.Printf("[INFO] Creating firewall filtering rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia firewall filtering rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorderWithBeforeReorder(order, resp.ID, "firewall_filtering_rules",
			func() (int, error) {
				allRules, err := filteringrules.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id, order int) error {
				// Custom updateOrder that handles predefined rules
				rule, err := filteringrules.Get(ctx, service, id)
				if err != nil {
					return err
				}
				if rule.Predefined {
					log.Printf("[INFO] Skipping reorder update for predefined rule ID %d (order: %d)", id, rule.Order)
					return nil
				}

				rule.Order = order
				_, err = filteringrules.Update(ctx, service, id, rule)
				return err
			},
			beforeReorderFirewallFilteringRules(ctx, service),
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceFirewallFilteringRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "firewall_filtering_rules")
		break
	}

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceFirewallFilteringRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia firewall filtering rule id is set"))
	}

	resp, err := filteringrules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall filtering rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	// if resp.Predefined {
	// 	log.Printf("[INFO] Rule ID %d is predefined — ignoring from Terraform state", resp.ID)
	// 	d.SetId("") // clear from Terraform state
	// 	return nil
	// }
	processedDestCountries := make([]string, len(resp.DestCountries))
	for i, country := range resp.DestCountries {
		processedDestCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	processedSrcCountries := make([]string, len(resp.SourceCountries))
	for i, country := range resp.SourceCountries {
		processedSrcCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
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
	_ = d.Set("source_countries", processedSrcCountries)
	_ = d.Set("nw_applications", resp.NwApplications)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationsGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenIDExtensionsListIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", flattenIDExtensionsListIDs(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDExtensionsListIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_windows", flattenIDExtensionsListIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("src_ip_groups", flattenIDExtensionsListIDs(resp.SrcIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ip_groups", flattenIDExtensionsListIDs(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_services", flattenIDExtensionsListIDs(resp.NwServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_service_groups", flattenIDExtensionsListIDs(resp.NwServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_application_groups", flattenIDExtensionsListIDs(resp.NwApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("app_services", flattenIDExtensionsListIDs(resp.AppServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_service_groups", flattenIDExtensionsListIDs(resp.AppServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDExtensionsListIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFirewallFilteringRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("firewall filtering rule ID not set"))
	}
	log.Printf("[INFO] Updating firewall filtering rule ID: %v\n", id)
	req := expandFirewallFilteringRules(d)

	if _, err := filteringrules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := filteringrules.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating firewall filtering rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorderWithBeforeReorder(req.Order, req.ID, "firewall_filtering_rules",
			func() (int, error) {
				allRules, err := filteringrules.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id, order int) error {
				rule, err := filteringrules.Get(ctx, service, id)
				if err != nil {
					return err
				}
				if rule.Predefined {
					log.Printf("[INFO] Skipping reorder update for predefined rule ID %d (order: %d)", id, rule.Order)
					return nil
				}

				// Optional: avoid unnecessary updates if the current order is already correct
				if rule.Order == order {
					return nil
				}

				rule.Order = order
				_, err = filteringrules.Update(ctx, service, id, rule)
				return err
			},
			beforeReorderFirewallFilteringRules(ctx, service),
		)

		if diags := resourceFirewallFilteringRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "firewall_filtering_rules")
		break
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceFirewallFilteringRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule not set: %v\n", id)
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := filteringrules.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving firewall filtering rule %d: %v", id, err))
	}

	// Validate if the rule can be deleted
	// Prevent deletion if the rule is predefined
	if rule.Predefined {
		return diag.FromErr(fmt.Errorf("deletion of predefined rule '%s' is not allowed", rule.Name))
	}

	log.Printf("[INFO] Deleting firewall filtering rule ID: %v\n", (d.Id()))
	if _, err := filteringrules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] firewall filtering rule deleted")

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		// Sleep for 2 seconds before potentially triggering the activation
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandFirewallFilteringRules(d *schema.ResourceData) filteringrules.FirewallFilteringRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandFirewallFilteringRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	// Process DestCountries and SourceCountries using the helper function
	processedDestCountries := processCountries(SetToStringList(d, "dest_countries"))
	processedSourceCountries := processCountries(SetToStringList(d, "source_countries"))

	result := filteringrules.FirewallFilteringRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Order:               order,
		Rank:                d.Get("rank").(int),
		Action:              d.Get("action").(string),
		State:               d.Get("state").(string),
		Description:         d.Get("description").(string),
		SrcIps:              SetToStringList(d, "src_ips"),
		DestAddresses:       SetToStringList(d, "dest_addresses"),
		DestIpCategories:    SetToStringList(d, "dest_ip_categories"),
		DeviceTrustLevels:   SetToStringList(d, "device_trust_levels"),
		DestCountries:       processedDestCountries,
		SourceCountries:     processedSourceCountries,
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
