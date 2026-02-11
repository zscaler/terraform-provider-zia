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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/forwarding_rules"
)

var (
	forwardingControlLock          sync.Mutex
	forwardingControlStartingOrder int
)

func resourceForwardingControlRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardingControlRuleCreate,
		ReadContext:   resourceForwardingControlRuleRead,
		UpdateContext: resourceForwardingControlRuleUpdate,
		DeleteContext: resourceForwardingControlRuleDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			forwardMethod := d.Get("forward_method").(string)
			ruleType := d.Get("type").(string)

			// Function to check if an attribute is set
			isSet := func(attr string) bool {
				_, ok := d.GetOk(attr)
				return ok
			}

			// Validate the constraints based on the rule type and forward method
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

			// Combined validation: `dest_addresses` and `dest_countries` can only be set when `forward_method` is either `PROXYCHAIN` or `DIRECT`
			if (isSet("dest_addresses") || isSet("dest_countries") || isSet("dest_ip_categories")) && forwardMethod != "PROXYCHAIN" && forwardMethod != "DIRECT" {
				return fmt.Errorf("dest_addresses, dest_countries and dest_ip_categories can only be set when forward_method is either 'PROXYCHAIN' or 'DIRECT'")
			}

			return nil
		},
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
					resp, err := forwarding_rules.GetByName(ctx, service, id)
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional information about the forwarding rule",
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
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
					"ENATDEDIP",
				}, false),
			},
			"order": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The order of execution for the forwarding rule order",
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
			"device_groups":                  setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
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
			"dest_countries":                 getISOCountryCodes(),
		},
	}
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

func resourceForwardingControlRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandForwardingControlRule(d)
	log.Printf("[INFO] Creating zia forwarding control rule\n%+v\n", req)

	if err := validatePredefinedRules(req); err != nil {
		return diag.FromErr(err)
	}

	start := time.Now()

	forwardMethod := d.Get("forward_method").(string)
	if forwardMethod == "ZPA" {
		// Sleep for 60 seconds before invoking Create
		time.Sleep(60 * time.Second)
	}

	for {
		forwardingControlLock.Lock()
		if forwardingControlStartingOrder == 0 {
			list, _ := forwarding_rules.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > forwardingControlStartingOrder {
					forwardingControlStartingOrder = r.Order
				}
			}
			if forwardingControlStartingOrder == 0 {
				forwardingControlStartingOrder = 1
			}
		}
		forwardingControlLock.Unlock()
		startWithoutLocking := time.Now()

		intendedOrder := req.Order
		intendedRank := req.Rank
		if intendedRank < 7 {
			// always start rank 7 rules at the next available order after all ranked rules
			req.Rank = 7
		}
		req.Order = forwardingControlStartingOrder

		// Retry logic in case of specific error
		var resp *forwarding_rules.ForwardingRules
		var err error
		for i := 0; i < 3; i++ {
			resp, err = forwarding_rules.Create(ctx, service, &req)
			if err == nil {
				break
			}

			if forwardMethod == "ZPA" {
				if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.Response.StatusCode == 400 &&
					strings.Contains(respErr.Message, "is no longer an active Source IP Anchored App Segment") {
					log.Printf("[WARN] Received error indicating resource is no longer active. Retrying...\n")
					time.Sleep(30 * time.Second) // Wait for 30 seconds before retrying
					continue
				}
			}
			return diag.FromErr(err)
		}

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, intendedOrder, req.Rank, currentForwardingControlOrderVsRankWording(ctx, zClient), err))
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia forwarding control rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
		resourceType := "forwarding_control_rule"

		reorderWithBeforeReorder(
			OrderRule{Order: intendedOrder, Rank: intendedRank},
			resp.ID,
			resourceType,
			func() (int, error) {
				allRules, err := forwarding_rules.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				// Custom updateOrder that handles predefined rules
				rule, err := forwarding_rules.Get(ctx, service, id)
				if err != nil {
					return err
				}

				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = forwarding_rules.Update(ctx, service, id, rule)
				return err
			},
			nil, // Remove beforeReorder function to avoid adding too many rules to the map
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		markOrderRuleAsDone(resp.ID, resourceType)
		waitForReorder(resourceType)

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

		return resourceForwardingControlRuleRead(ctx, d, meta)
	}
}

func resourceForwardingControlRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia forwarding control rule id is set"))
	}
	resp, err := forwarding_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing forwarding control rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
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

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationsGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ec_groups", flattenIDExtensionsListIDs(resp.ECGroups)); err != nil {
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

	if err := d.Set("nw_application_groups", flattenIDExtensionsListIDs(resp.NwApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nw_application_groups", flattenIDExtensionsListIDs(resp.NwApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("app_service_groups", flattenIDExtensionsListIDs(resp.AppServiceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("proxy_gateway", flattenIDNameSet(resp.ProxyGateway)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_gateway", flattenIDNameSet(resp.ZPAGateway)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceForwardingControlRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] forwarding control rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("forwarding control rule ID not set"))
	}
	log.Printf("[INFO] Updating forwarding control rule ID: %v\n", id)
	req := expandForwardingControlRule(d)

	if err := validatePredefinedRules(req); err != nil {
		return diag.FromErr(err)
	}

	if _, err := forwarding_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	existingRules, err := forwarding_rules.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all forwarding control rules: %v", err)
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

	forwardMethod := d.Get("forward_method").(string)
	if forwardMethod == "ZPA" {
		// Sleep for 60 seconds before invoking Update
		time.Sleep(60 * time.Second)
	}

	// Retry logic in case of specific error
	for i := 0; i < 3; i++ {
		_, err = forwarding_rules.Update(ctx, service, id, &req)
		if err == nil {
			break
		}

		if forwardMethod == "ZPA" {
			if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.Response.StatusCode == 400 &&
				strings.Contains(respErr.Message, "is no longer an active Source IP Anchored App Segment") {
				log.Printf("[WARN] Received error indicating resource is no longer active. Retrying...\n")
				time.Sleep(30 * time.Second) // Wait for 30 seconds before retrying
				continue
			}
		}
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating forwarding control rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "forwarding_control_rule",
		func() (int, error) {
			allRules, err := forwarding_rules.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := forwarding_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// Optional: avoid unnecessary updates if the current order is already correct
			if rule.Order == order.Order && rule.Rank == order.Rank {
				return nil
			}

			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = forwarding_rules.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceForwardingControlRuleRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "forwarding_control_rule")
	waitForReorder("forwarding_control_rule")

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

func resourceForwardingControlRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("forwarding control rule ID not set: %v", id))
	}

	// Retrieve the rule to check if it's a predefined one
	rule, err := forwarding_rules.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving forwarding control rule %d: %v", id, err))
	}

	// Validate if the rule can be deleted
	if err := validatePredefinedRules(*rule); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting forwarding control rule ID: %v", id)
	if _, err := forwarding_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting forwarding control rule %d: %v", id, err))
	}

	d.SetId("")
	log.Printf("[INFO] Forwarding control rule deleted")
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

func expandForwardingControlRule(d *schema.ResourceData) forwarding_rules.ForwardingRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandFirewallIPSRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

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
		Order:               order,
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
		DeviceGroups:        expandIDNameExtensionsSet(d, "device_groups"),
	}

	return result
}

func currentForwardingControlOrderVsRankWording(ctx context.Context, zClient *Client) string {
	service := zClient.Service

	list, err := forwarding_rules.GetAll(ctx, service)
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
