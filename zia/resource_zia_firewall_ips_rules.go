package zia

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallipscontrolpolicies"
)

var (
	firewallIPSLock          sync.Mutex
	firewallIPSStartingOrder int
)

func resourceFirewallIPSRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallIPSRulesCreate,
		ReadContext:   resourceFirewallIPSRulesRead,
		UpdateContext: resourceFirewallIPSRulesUpdate,
		DeleteContext: resourceFirewallIPSRulesDelete,
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
					resp, err := firewallipscontrolpolicies.GetByName(ctx, service, id)
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
				Description:  "The name of the IPS Control rule",
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
				Computed:    true,
				Description: "Rule order number of the Firewall Filtering policy rule",
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "The admin rank specified for the rule based on your assigned admin rank. Admin rank determines the rule order that can be specified for the rule. ",
			},
			"enable_full_logging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "A Boolean value that indicates whether full logging is enabled. A true value indicates that full logging is enabled, whereas a false value indicates that aggregate logging is enabled.",
			},
			"capture_pcap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value that indicates whether packet capture (PCAP) is enabled or not",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action configured for the rule that must take place if the traffic matches the rule criteria, such as allowing or blocking the traffic or bypassing the rule.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK_DROP",
					"BLOCK_RESET",
					"BYPASS_IPS",
				}, false),
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The state of the rule indicating whether it is enabled or disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"res_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "URL categories associated with resolved IP addresses to which the rule applies. If not set, the rule is not restricted to a specific URL category.",
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
					Type:         schema.TypeString,
					ValidateFunc: validateDestAddress, // Apply the custom validation function here
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
			"locations":         setIDsSchemaTypeCustom(intPtr(8), "list of locations for which rule must be applied"),
			"location_groups":   setIDsSchemaTypeCustom(intPtr(32), "list of locations groups"),
			"users":             setIDsSchemaTypeCustom(intPtr(4), "list of users for which rule must be applied"),
			"groups":            setIDsSchemaTypeCustom(intPtr(8), "list of groups for which rule must be applied"),
			"departments":       setIDsSchemaTypeCustom(intPtr(140000), "list of departments for which rule must be applied"),
			"time_windows":      setIDsSchemaTypeCustom(intPtr(2), "The time interval in which the Firewall Filtering policy rule applies"),
			"labels":            setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"device_groups":     setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":           setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"src_ip_groups":     setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"src_ipv6_groups":   setIDsSchemaTypeCustom(nil, "list of Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group."),
			"dest_ip_groups":    setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"dest_ipv6_groups":  setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"nw_services":       setIDsSchemaTypeCustom(intPtr(1024), "list of nw services"),
			"nw_service_groups": setIDsSchemaTypeCustom(nil, "list of nw service groups"),
			"threat_categories": setIDsSchemaTypeCustom(nil, "list of Advanced threat categories to which the rule applies"),
			"zpa_app_segments":  setExtIDNameSchemaCustom(intPtr(255), "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method."),
			"dest_countries":    getISOCountryCodes(),
			"source_countries":  getISOCountryCodes(),
		},
	}
}

func resourceFirewallIPSRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFirewallIPSRules(d)
	log.Printf("[INFO] Creating zia firewall ips rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		firewallIPSLock.Lock()
		if firewallIPSStartingOrder == 0 {
			list, _ := firewallipscontrolpolicies.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > firewallIPSStartingOrder {
					firewallIPSStartingOrder = r.Order
				}
			}
			if firewallIPSStartingOrder == 0 {
				firewallIPSStartingOrder = 1
			}
		}
		firewallIPSLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = firewallIPSStartingOrder

		resp, err := firewallipscontrolpolicies.Create(ctx, service, &req)
		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, order, req.Rank, currentOrderVsRankWording(ctx, zClient), err))
				}
				if time.Since(start) < timeout {
					log.Printf("[INFO] Creating firewall ips rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia firewall ips rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "firewall_filtering_rules", func() (int, error) {
			list, err := firewallipscontrolpolicies.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := firewallipscontrolpolicies.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = firewallipscontrolpolicies.Update(ctx, service, id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceFirewallIPSRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "firewall_filtering_rules")
		break
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceFirewallIPSRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia firewall ips rule id is set"))
	}

	resp, err := firewallipscontrolpolicies.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall ips rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	processedDestCountries := make([]string, len(resp.DestCountries))
	for i, country := range resp.DestCountries {
		processedDestCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	processedSrcCountries := make([]string, len(resp.SourceCountries))
	for i, country := range resp.SourceCountries {
		processedSrcCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	log.Printf("[INFO] Getting firewall ips rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("enable_full_logging", resp.EnableFullLogging)
	_ = d.Set("action", resp.Action)
	_ = d.Set("state", resp.State)
	_ = d.Set("description", resp.Description)
	_ = d.Set("res_categories", resp.ResCategories)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", processedDestCountries)
	_ = d.Set("source_countries", processedSrcCountries)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationsGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", flattenIDs(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("src_ip_groups", flattenIDs(resp.SrcIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("src_ipv6_groups", flattenIDs(resp.SrcIpv6Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ip_groups", flattenIDs(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ipv6_groups", flattenIDs(resp.DestIpv6Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_services", flattenIDs(resp.NwServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_service_groups", flattenIDs(resp.NwServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("threat_categories", flattenIDs(resp.ThreatCategories)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFirewallIPSRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall ips rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("firewall ips rule ID not set"))
	}
	log.Printf("[INFO] Updating firewall ips rule ID: %v\n", id)
	req := expandFirewallIPSRules(d)

	if _, err := firewallipscontrolpolicies.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := firewallipscontrolpolicies.Update(ctx, service, id, &req)
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating firewall ips rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorder(req.Order, req.ID, "firewall_filtering_rules", func() (int, error) {
			list, err := firewallipscontrolpolicies.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := firewallipscontrolpolicies.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = firewallipscontrolpolicies.Update(ctx, service, id, rule)
			return err
		})

		if diags := resourceFirewallIPSRulesRead(ctx, d, meta); diags.HasError() {
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

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceFirewallIPSRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall ips rule not set: %v\n", id)
	}

	// Retrieve the rule to check if it's predefined
	rule, err := firewallipscontrolpolicies.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving firewall IPS rule %d: %v", id, err))
	}

	// Prevent deletion if the rule is predefined
	if rule.Predefined {
		return diag.FromErr(fmt.Errorf("deletion of predefined rule '%s' is not allowed", rule.Name))
	}

	log.Printf("[INFO] Deleting firewall ips rule ID: %v\n", (d.Id()))
	if _, err := firewallipscontrolpolicies.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] firewall ips rule deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandFirewallIPSRules(d *schema.ResourceData) firewallipscontrolpolicies.FirewallIPSRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Process DestCountries and SourceCountries using the helper function
	processedDestCountries := processCountries(SetToStringList(d, "dest_countries"))
	processedSourceCountries := processCountries(SetToStringList(d, "source_countries"))

	result := firewallipscontrolpolicies.FirewallIPSRules{
		ID:                id,
		Name:              d.Get("name").(string),
		Order:             d.Get("order").(int),
		Rank:              d.Get("rank").(int),
		Action:            d.Get("action").(string),
		State:             d.Get("state").(string),
		Description:       d.Get("description").(string),
		SrcIps:            SetToStringList(d, "src_ips"),
		DestAddresses:     SetToStringList(d, "dest_addresses"),
		DestIpCategories:  SetToStringList(d, "dest_ip_categories"),
		ResCategories:     SetToStringList(d, "res_categories"),
		DestCountries:     processedDestCountries,
		SourceCountries:   processedSourceCountries,
		EnableFullLogging: d.Get("enable_full_logging").(bool),
		DefaultRule:       d.Get("default_rule").(bool),
		Predefined:        d.Get("predefined").(bool),
		Locations:         expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups:   expandIDNameExtensionsSet(d, "location_groups"),
		Departments:       expandIDNameExtensionsSet(d, "departments"),
		Groups:            expandIDNameExtensionsSet(d, "groups"),
		Users:             expandIDNameExtensionsSet(d, "users"),
		TimeWindows:       expandIDNameExtensionsSet(d, "time_windows"),
		SrcIpGroups:       expandIDNameExtensionsSet(d, "src_ip_groups"),
		SrcIpv6Groups:     expandIDNameExtensionsSet(d, "src_ipv6_groups"),
		DestIpGroups:      expandIDNameExtensionsSet(d, "dest_ip_groups"),
		DestIpv6Groups:    expandIDNameExtensionsSet(d, "dest_ipv6_groups"),
		NwServices:        expandIDNameExtensionsSet(d, "nw_services"),
		NwServiceGroups:   expandIDNameExtensionsSet(d, "nw_service_groups"),
		Labels:            expandIDNameExtensionsSet(d, "labels"),
		DeviceGroups:      expandIDNameExtensionsSet(d, "device_groups"),
		Devices:           expandIDNameExtensionsSet(d, "devices"),
		ThreatCategories:  expandIDNameExtensionsSet(d, "threat_categories"),
		ZPAAppSegments:    expandZPAAppSegmentSet(d, "zpa_app_segments"),
	}
	return result
}
