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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/nat_control_policies"
)

var (
	natControlRuleLock          sync.Mutex
	natControlRuleStartingOrder int
)

func resourceNatControlRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNatControlRulesCreate,
		ReadContext:   resourceNatControlRulesRead,
		UpdateContext: resourceNatControlRulesUpdate,
		DeleteContext: resourceNatControlRulesDelete,
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
					resp, err := nat_control_policies.GetByName(ctx, service, id)
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
				Description: "Name of the nat control policy rule",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional information about the rule",
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
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
				Description:  "Admin rank of the nat control policy rule",
			},
			"enable_full_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines whether the nat control policy rule is enabled or disabled",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"redirect_fqdn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action the nat control policy rule takes when packets match the rule",
			},
			"redirect_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action the nat control policy rule takes when packets match the rule",
			},
			"redirect_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The action the nat control policy rule takes when packets match the rule",
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
			"res_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of destination domain categories to which the rule applies",
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
			"time_windows":      setIDsSchemaTypeCustom(intPtr(2), "The time interval in which the nat control policy rule applies"),
			"labels":            setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"device_groups":     setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":           setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"src_ip_groups":     setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"src_ipv6_groups":   setIDsSchemaTypeCustom(nil, "list of source ipv6 groups"),
			"dest_ip_groups":    setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"dest_ipv6_groups":  setIDsSchemaTypeCustom(nil, "list of destination ipv6 groups"),
			"nw_service_groups": setIDsSchemaTypeCustom(nil, "list of nw service groups"),
			"nw_services":       setIDsSchemaTypeCustom(intPtr(1024), "list of nw services"),
			"dest_countries":    getISOCountryCodes(),
		},
	}
}

func beforeReorderNatControlRules(ctx context.Context, service *zscaler.Service) func() {
	return func() {
		log.Printf("[INFO] beforeReorderNatControlRules")
		// get all predefined rules and set their order to come first
		rules, err := nat_control_policies.GetAll(ctx, service)
		if err != nil {
			log.Printf("[ERROR] beforeReorderNatControlRules: %v", err)
		}
		// first order predefined rules by their order
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Order < rules[j].Order
		})
		order := 1
		for _, r := range rules {
			if r.Predefined {
				r.Order = order
				_, err = nat_control_policies.Update(ctx, service, r.ID, &r)
				if err != nil {
					log.Printf("[ERROR] beforeReorderNatControlRules: %v", err)
				}
				order++
			}
		}
	}
}

func resourceNatControlRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandNatControlRules(d)
	log.Printf("[INFO] Creating zia nat control rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		natControlRuleLock.Lock()
		if natControlRuleStartingOrder == 0 {
			list, _ := nat_control_policies.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > natControlRuleStartingOrder {
					natControlRuleStartingOrder = r.Order
				}
			}
			if natControlRuleStartingOrder == 0 {
				natControlRuleStartingOrder = 1
			} else {
				natControlRuleStartingOrder++
			}
		}
		natControlRuleLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = natControlRuleStartingOrder

		resp, err := nat_control_policies.Create(ctx, service, &req)

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
					log.Printf("[INFO] Creating nat control rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia nat control rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: req.Rank},
			resp.ID,
			"nat_control_rules",
			func() (int, error) {
				allRules, err := nat_control_policies.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				// Custom updateOrder that handles predefined rules
				rule, err := nat_control_policies.Get(ctx, service, id)
				if err != nil {
					return err
				}
				if rule.Predefined {
					log.Printf("[INFO] Skipping reorder update for predefined rule ID %d (order: %d)", id, rule.Order)
					return nil
				}

				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = nat_control_policies.Update(ctx, service, id, rule)
				return err
			},
			beforeReorderNatControlRules(ctx, service),
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceNatControlRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "nat_control_rules")
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

func resourceNatControlRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia nat control rule id is set"))
	}

	resp, err := nat_control_policies.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing nat control rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	processedDestCountries := make([]string, len(resp.DestCountries))
	for i, country := range resp.DestCountries {
		processedDestCountries[i] = strings.TrimPrefix(country, "COUNTRY_")
	}

	log.Printf("[INFO] Getting nat control rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("enable_full_logging", resp.EnableFullLogging)
	_ = d.Set("state", resp.State)
	_ = d.Set("redirect_fqdn", resp.RedirectFqdn)
	_ = d.Set("redirect_ip", resp.RedirectIp)
	_ = d.Set("redirect_port", resp.RedirectPort)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", processedDestCountries)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)
	_ = d.Set("res_categories", resp.ResCategories)

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationGroups)); err != nil {
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

	if err := d.Set("src_ipv6_groups", flattenIDExtensionsListIDs(resp.SrcIpv6Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ip_groups", flattenIDExtensionsListIDs(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ipv6_groups", flattenIDExtensionsListIDs(resp.DestIpv6Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_services", flattenIDExtensionsListIDs(resp.NwServices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_service_groups", flattenIDExtensionsListIDs(resp.NwServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDExtensionsListIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNatControlRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] nat control rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("nat control rule ID not set"))
	}
	log.Printf("[INFO] Updating nat control rule ID: %v\n", id)
	req := expandNatControlRules(d)

	if _, err := nat_control_policies.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := nat_control_policies.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating nat control rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorderWithBeforeReorder(OrderRule{Order: req.Order, Rank: req.Rank}, req.ID, "nat_control_rules",
			func() (int, error) {
				allRules, err := nat_control_policies.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				rule, err := nat_control_policies.Get(ctx, service, id)
				if err != nil {
					return err
				}
				if rule.Predefined {
					log.Printf("[INFO] Skipping reorder update for predefined rule ID %d (order: %d)", id, rule.Order)
					return nil
				}

				// Optional: avoid unnecessary updates if the current order is already correct
				if rule.Order == order.Order && rule.Rank == order.Rank {
					return nil
				}

				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = nat_control_policies.Update(ctx, service, id, rule)
				return err
			},
			beforeReorderNatControlRules(ctx, service),
		)

		if diags := resourceNatControlRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "nat_control_rules")
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

func resourceNatControlRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] nat control rule not set: %v\n", id)
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := nat_control_policies.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving nat control rule %d: %v", id, err))
	}

	// Validate if the rule can be deleted
	// Prevent deletion if the rule is predefined
	if rule.Predefined {
		return diag.FromErr(fmt.Errorf("deletion of predefined rule '%s' is not allowed", rule.Name))
	}

	log.Printf("[INFO] Deleting nat control rule ID: %v\n", (d.Id()))
	if _, err := nat_control_policies.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] nat control rule deleted")

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

func expandNatControlRules(d *schema.ResourceData) nat_control_policies.NatControlPolicies {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandNatControlRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	processedDestCountries := processCountries(SetToStringList(d, "dest_countries"))

	result := nat_control_policies.NatControlPolicies{
		ID:                id,
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Order:             order,
		Rank:              d.Get("rank").(int),
		State:             d.Get("state").(string),
		RedirectFqdn:      d.Get("redirect_fqdn").(string),
		RedirectIp:        d.Get("redirect_ip").(string),
		RedirectPort:      d.Get("redirect_port").(int),
		SrcIps:            SetToStringList(d, "src_ips"),
		DestAddresses:     SetToStringList(d, "dest_addresses"),
		DestIpCategories:  SetToStringList(d, "dest_ip_categories"),
		ResCategories:     SetToStringList(d, "res_categories"),
		DestCountries:     processedDestCountries,
		EnableFullLogging: d.Get("enable_full_logging").(bool),
		DefaultRule:       d.Get("default_rule").(bool),
		Predefined:        d.Get("predefined").(bool),
		Locations:         expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:    expandIDNameExtensionsSet(d, "location_groups"),
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
	}
	return result
}
