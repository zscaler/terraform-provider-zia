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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/traffic_capture"
)

var (
	trafficCaptureLock          sync.Mutex
	trafficCaptureStartingOrder int
)

func resourceTrafficCaptureRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFiresourceTrafficCaptureRulesCreate,
		ReadContext:   resourceFiresourceTrafficCaptureRulesRead,
		UpdateContext: resourceFiresourceTrafficCaptureRulesUpdate,
		DeleteContext: resourceFiresourceTrafficCaptureRulesDelete,
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
					resp, err := traffic_capture.GetByName(ctx, service, id)
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional information about the rule",
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule order number. If omitted, the rule will be added to the end of the rule set.",
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the Firewall Filtering policy rule",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action to be enforced when the traffic matches the rule criteria",
				ValidateFunc: validation.StringInSlice([]string{
					"CAPTURE",
					"SKIP",
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
			"nw_applications": {
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
			"txn_size_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The maximum size of traffic to capture per connection",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"UNLIMITED",
					"THIRTY_TWO_KB",
					"TWO_FIFTY_SIX_KB",
					"TWO_MB",
					"FOUR_MB",
					"THIRTY_TWO_MB",
					"SIXTY_FOUR_MB",
				}, false),
			},
			"txn_sampling": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The percentage of connections sampled for capturing each time the rule is triggered",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"ONE_PERCENT",
					"TWO_PERCENT",
					"FIVE_PERCENT",
					"TEN_PERCENT",
					"TWENTY_FIVE_PERCENT",
					"HUNDRED_PERCENT",
				}, false),
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
			"nw_application_groups": setIDsSchemaTypeCustom(nil, "list of nw application groups"),
			"nw_service_groups":     setIDsSchemaTypeCustom(nil, "list of nw service groups"),
			"nw_services":           setIDsSchemaTypeCustom(intPtr(1024), "list of nw services"),
			"workload_groups":       setIdNameSchemaCustom(255, "The list of preconfigured workload groups to which the policy must be applied"),
			"dest_countries":        getISOCountryCodes(),
			"source_countries":      getISOCountryCodes(),
			"device_trust_levels":   getDeviceTrustLevels(),
		},
	}
}

func resourceFiresourceTrafficCaptureRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFiresourceTrafficCaptureRules(d)
	log.Printf("[INFO] Creating zia firewall filtering rule\n%+v\n", req)

	start := time.Now()

	trafficCaptureLock.Lock()
	if trafficCaptureStartingOrder == 0 {
		list, _ := traffic_capture.GetAll(ctx, service, nil)
		for _, r := range list {
			if r.Order > trafficCaptureStartingOrder {
				trafficCaptureStartingOrder = r.Order
			}
		}
		if trafficCaptureStartingOrder == 0 {
			trafficCaptureStartingOrder = 1
		} else {
			trafficCaptureStartingOrder++
		}
	}
	trafficCaptureLock.Unlock()
	startWithoutLocking := time.Now()

	// Store the intended order from HCL
	intendedOrder := req.Order
	intendedRank := req.Rank
	if intendedRank < 7 {
		// always start rank 7 rules at the next available order after all ranked rules
		req.Rank = 7
	}
	req.Order = trafficCaptureStartingOrder
	resp, err := traffic_capture.Create(ctx, service, &req)

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			if reg.MatchString(err.Error()) {
				return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, intendedOrder, req.Rank, currentTrafficCaptureOrderVsRankWording(ctx, zClient), err))
			}
		}
		return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
	}

	log.Printf("[INFO] Created zia firewall filtering rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
	// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
	resourceType := "firewall_filtering_rules"

	reorderWithBeforeReorder(
		OrderRule{Order: intendedOrder, Rank: intendedRank},
		resp.ID,
		resourceType,
		func() (int, error) {
			allRules, err := traffic_capture.GetAll(ctx, service, nil)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			// Custom updateOrder that handles predefined rules
			rule, err := traffic_capture.Get(ctx, service, id)
			if err != nil {
				return err
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = traffic_capture.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)

	markOrderRuleAsDone(resp.ID, resourceType)

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

	return resourceFiresourceTrafficCaptureRulesRead(ctx, d, meta)
}

func resourceFiresourceTrafficCaptureRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia firewall filtering rule id is set"))
	}

	// Use GetAll() instead of Get() to reduce API calls during terraform refresh
	allRules, err := traffic_capture.GetAll(ctx, service, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the specific rule by ID
	var resp *traffic_capture.TrafficCaptureRules
	for i := range allRules {
		if allRules[i].ID == id {
			resp = &allRules[i]
			break
		}
	}

	// Rule not found
	if resp == nil {
		log.Printf("[WARN] Removing firewall filtering rule %s from state because it no longer exists in ZIA", d.Id())
		d.SetId("")
		return nil
	}
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
	_ = d.Set("txn_size_limit", resp.TxnSizeLimit)
	_ = d.Set("txn_sampling", resp.TxnSampling)

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
	return nil
}

func resourceFiresourceTrafficCaptureRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("firewall filtering rule ID not set"))
	}
	log.Printf("[INFO] Updating firewall filtering rule ID: %v\n", id)
	req := expandFiresourceTrafficCaptureRules(d)

	if _, err := traffic_capture.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	existingRules, err := traffic_capture.GetAll(ctx, service, nil)
	if err != nil {
		log.Printf("[ERROR] error getting all filtering rules: %v", err)
	}
	sort.Slice(existingRules, func(i, j int) bool {
		return existingRules[i].Rank < existingRules[j].Rank || (existingRules[i].Rank == existingRules[j].Rank && existingRules[i].Order < existingRules[j].Order)
	})
	intendedOrder := req.Order
	intendedRank := req.Rank
	nextAvailableOrder := existingRules[len(existingRules)-1].Order
	// always start rank 7 rules at the next available order after all ranked rules
	req.Rank = 7

	req.Order = nextAvailableOrder

	_, err = traffic_capture.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating firewall filtering rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "firewall_filtering_rules",
		func() (int, error) {
			allRules, err := traffic_capture.GetAll(ctx, service, nil)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := traffic_capture.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// Optional: avoid unnecessary updates if the current order is already correct
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = traffic_capture.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceFiresourceTrafficCaptureRulesRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "firewall_filtering_rules")

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

func resourceFiresourceTrafficCaptureRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall filtering rule not set: %v\n", id)
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := traffic_capture.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving firewall filtering rule %d: %v", id, err))
	}

	// Additional check for any predefined rule (backup validation)
	if rule.Predefined {
		return diag.FromErr(fmt.Errorf("deletion of predefined rule '%s' is not allowed", rule.Name))
	}

	log.Printf("[INFO] Deleting firewall filtering rule ID: %v\n", (d.Id()))
	if _, err := traffic_capture.Delete(ctx, service, id); err != nil {
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

func expandFiresourceTrafficCaptureRules(d *schema.ResourceData) traffic_capture.TrafficCaptureRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandFiresourceTrafficCaptureRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	// Process DestCountries and SourceCountries using the helper function
	processedDestCountries := processCountries(SetToStringList(d, "dest_countries"))
	processedSourceCountries := processCountries(SetToStringList(d, "source_countries"))

	result := traffic_capture.TrafficCaptureRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Order:               order,
		Rank:                d.Get("rank").(int),
		Action:              d.Get("action").(string),
		State:               d.Get("state").(string),
		Description:         d.Get("description").(string),
		TxnSizeLimit:        d.Get("txn_size_limit").(string),
		TxnSampling:         d.Get("txn_sampling").(string),
		SrcIps:              SetToStringList(d, "src_ips"),
		DestAddresses:       SetToStringList(d, "dest_addresses"),
		DestIpCategories:    SetToStringList(d, "dest_ip_categories"),
		DeviceTrustLevels:   SetToStringList(d, "device_trust_levels"),
		DestCountries:       processedDestCountries,
		SourceCountries:     processedSourceCountries,
		NwApplications:      SetToStringList(d, "nw_applications"),
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
		AppServiceGroups:    expandIDNameExtensionsSet(d, "app_service_groups"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
		DeviceGroups:        expandIDNameExtensionsSet(d, "device_groups"),
		Devices:             expandIDNameExtensionsSet(d, "devices"),
		WorkloadGroups:      expandWorkloadGroupsIDName(d, "workload_groups"),
	}
	return result
}

func currentTrafficCaptureOrderVsRankWording(ctx context.Context, zClient *Client) string {
	service := zClient.Service

	list, err := traffic_capture.GetAll(ctx, service, nil)
	if err != nil {
		return ""
	}
	result := ""
	for i, r := range list {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("Rank %d VS Order %d", r.Rank, r.Order)

	}
	return result
}
