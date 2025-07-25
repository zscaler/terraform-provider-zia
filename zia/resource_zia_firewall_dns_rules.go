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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewalldnscontrolpolicies"
)

var (
	firewallDNSLock          sync.Mutex
	firewallDNSStartingOrder int
)

func resourceFirewallDNSRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallDNSRulesCreate,
		ReadContext:   resourceFirewallDNSRulesRead,
		UpdateContext: resourceFirewallDNSRulesUpdate,
		DeleteContext: resourceFirewallDNSRulesDelete,

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
					resp, err := firewalldnscontrolpolicies.GetByName(ctx, service, id)
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
				Description: "The name of the IPS Control rule",
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
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "Rule order number. If omitted, the rule will be added to the end of the rule set.",
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "The admin rank specified for the rule based on your assigned admin rank. Admin rank determines the rule order that can be specified for the rule. ",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action configured for the rule that must take place if the traffic matches the rule criteria, such as allowing, blocking, or redirecting the traffic.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK",
					"REDIR_REQ",
					"REDIR_RES",
					"REDIR_ZPA",
					"REDIR_REQ_DOH",
					"REDIR_REQ_KEEP_SENDER",
					"REDIR_REQ_TCP",
					"REDIR_REQ_UDP",
					"BLOCK_WITH_RESPONSE",
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
			"redirect_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address to which the traffic will be redirected to when the DNAT rule is triggered. If not set, no redirection is done to specific IP addresses.",
			},
			"block_response_code": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The IP address to which the traffic will be redirected to when the DNAT rule is triggered. If not set, no redirection is done to specific IP addresses.",
				ValidateFunc: validateBlockResponseCode,
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
					Type: schema.TypeString,
				},
				Description: "Destination addresses. Supports IPv4, FQDNs, or wildcard FQDNs",
			},
			"dest_ip_categories": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Destination IP categories to which the rule applies. If not set, the rule is not restricted to specific categories.",
			},
			"capture_pcap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value that indicates whether packet capture (PCAP) is enabled or not",
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
			"applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `User-defined network service applications on which the rule is applied.
				If not set, the rule is not restricted to a specific network service application
				Use the data source zia_cloud_applications to get the list of available cloud applications:
				https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_applications
				`,
			},
			"locations":              setIDsSchemaTypeCustom(intPtr(8), "list of locations for which rule must be applied"),
			"location_groups":        setIDsSchemaTypeCustom(intPtr(32), "list of locations groups"),
			"users":                  setIDsSchemaTypeCustom(intPtr(4), "list of users for which rule must be applied"),
			"groups":                 setIDsSchemaTypeCustom(intPtr(8), "list of groups for which rule must be applied"),
			"departments":            setIDsSchemaTypeCustom(intPtr(140000), "list of departments for which rule must be applied"),
			"time_windows":           setIDsSchemaTypeCustom(intPtr(2), "The time interval in which the Firewall Filtering policy rule applies"),
			"labels":                 setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"device_groups":          setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":                setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"src_ip_groups":          setIDsSchemaTypeCustom(nil, "list of Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group."),
			"src_ipv6_groups":        setIDsSchemaTypeCustom(nil, "list of Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group."),
			"dest_ip_groups":         setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"dest_ipv6_groups":       setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"application_groups":     setIDsSchemaTypeCustom(nil, "list of nw application groups"),
			"dns_gateway":            setIdNameSchemaCustom(1, "The DNS gateway used to redirect traffic, specified when the rule action is to redirect DNS request to an external DNS service"),
			"zpa_ip_group":           setIdNameSchemaCustom(1, "The ZPA IP pool specified when the rule action is to resolve domain names of ZPA applications to an ephemeral IP address from a preconfigured IP pool"),
			"edns_ecs_object":        setIdNameSchemaCustom(1, "The EDNS ECS object which resolves DNS request"),
			"dest_countries":         getISOCountryCodes(),
			"source_countries":       getISOCountryCodes(),
			"dns_rule_request_types": getDnsRuleRequestTypes(),
			// "applications":           getCloudApplications(),
			"protocols": getDNSRuleProtocols(),
		},
	}
}

func resourceFirewallDNSRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFirewallDNSRules(d)
	log.Printf("[INFO] Creating zia firewall dns rule\n%+v\n", req)

	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		firewallDNSLock.Lock()
		if firewallDNSStartingOrder == 0 {
			list, _ := firewalldnscontrolpolicies.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > firewallDNSStartingOrder {
					firewallDNSStartingOrder = r.Order
				}
			}
			if firewallDNSStartingOrder == 0 {
				firewallDNSStartingOrder = 1
			}
		}
		firewallDNSLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = firewallDNSStartingOrder

		resp, err := firewalldnscontrolpolicies.Create(ctx, service, &req)

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
					log.Printf("[INFO] Creating firewall dns rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia firewall dns rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "firewall_dns_rule", func() (int, error) {
			list, err := firewalldnscontrolpolicies.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := firewalldnscontrolpolicies.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = firewalldnscontrolpolicies.Update(ctx, service, id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceFirewallDNSRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "firewall_dns_rule")
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

func resourceFirewallDNSRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia firewall dns rule id is set"))
	}

	resp, err := firewalldnscontrolpolicies.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing firewall dns rule %s from state because it no longer exists in ZIA", d.Id())
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

	log.Printf("[INFO] Getting firewall dns rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("action", resp.Action)
	_ = d.Set("state", resp.State)
	_ = d.Set("block_response_code", resp.BlockResponseCode)
	_ = d.Set("dns_rule_request_types", resp.DNSRuleRequestTypes)
	_ = d.Set("res_categories", resp.ResCategories)
	_ = d.Set("redirect_ip", resp.RedirectIP)
	_ = d.Set("applications", resp.Applications)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", processedDestCountries)
	_ = d.Set("source_countries", processedSrcCountries)
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("capture_pcap", resp.CapturePCAP)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)

	if err := d.Set("application_groups", flattenIDExtensionsListIDs(resp.ApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDExtensionsListIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

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

	if err := d.Set("src_ipv6_groups", flattenIDExtensionsListIDs(resp.SrcIpv6Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ip_groups", flattenIDExtensionsListIDs(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ipv6_groups", flattenIDExtensionsListIDs(resp.DestIpv6Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dns_gateway", flattenIDNameSet(resp.DNSGateway)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("edns_ecs_object", flattenIDNameSet(resp.EDNSEcsObject)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_ip_group", flattenIDNameSet(resp.ZPAIPGroup)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFirewallDNSRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall dns rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("firewall dns rule ID not set"))
	}
	log.Printf("[INFO] Updating firewall dns rule ID: %v\n", id)
	req := expandFirewallDNSRules(d)

	if _, err := firewalldnscontrolpolicies.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := firewalldnscontrolpolicies.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating firewall dns rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorder(req.Order, req.ID, "firewall_dns_rule", func() (int, error) {
			list, err := firewalldnscontrolpolicies.GetAll(ctx, service)
			return len(list), err
		}, func(id, order int) error {
			rule, err := firewalldnscontrolpolicies.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = firewalldnscontrolpolicies.Update(ctx, service, id, rule)
			return err
		})

		if diags := resourceFirewallDNSRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "firewall_dns_rule")
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

func resourceFirewallDNSRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] firewall dns rule not set: %v\n", id)
	}

	// Retrieve the rule to check if it's predefined
	rule, err := firewalldnscontrolpolicies.Get(ctx, service, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving firewall DNS rule %d: %v", id, err))
	}

	// Prevent deletion if the rule is predefined
	if rule.Predefined {
		return diag.FromErr(fmt.Errorf("deletion of predefined rule '%s' is not allowed", rule.Name))
	}

	log.Printf("[INFO] Deleting firewall dns rule ID: %v\n", (d.Id()))
	if _, err := firewalldnscontrolpolicies.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] firewall dns rule deleted")

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

func expandFirewallDNSRules(d *schema.ResourceData) firewalldnscontrolpolicies.FirewallDNSRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandFirewallDNSRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	// Process DestCountries and SourceCountries using the helper function
	processedDestCountries := processCountries(SetToStringList(d, "dest_countries"))
	processedSourceCountries := processCountries(SetToStringList(d, "source_countries"))

	result := firewalldnscontrolpolicies.FirewallDNSRules{
		ID:                  id,
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Order:               order,
		Rank:                d.Get("rank").(int),
		Action:              d.Get("action").(string),
		State:               d.Get("state").(string),
		BlockResponseCode:   d.Get("block_response_code").(string),
		RedirectIP:          d.Get("redirect_ip").(string),
		DefaultRule:         d.Get("default_rule").(bool),
		Predefined:          d.Get("predefined").(bool),
		Applications:        SetToStringList(d, "applications"),
		DNSRuleRequestTypes: SetToStringList(d, "dns_rule_request_types"),
		DestAddresses:       SetToStringList(d, "dest_addresses"),
		DestIpCategories:    SetToStringList(d, "dest_ip_categories"),
		Protocols:           SetToStringList(d, "protocols"),
		ResCategories:       SetToStringList(d, "res_categories"),
		SrcIps:              SetToStringList(d, "src_ips"),
		DestCountries:       processedDestCountries,
		SourceCountries:     processedSourceCountries,
		ApplicationGroups:   expandIDNameExtensionsSet(d, "application_groups"),
		Locations:           expandIDNameExtensionsSet(d, "locations"),
		LocationsGroups:     expandIDNameExtensionsSet(d, "location_groups"),
		Departments:         expandIDNameExtensionsSet(d, "departments"),
		Groups:              expandIDNameExtensionsSet(d, "groups"),
		Users:               expandIDNameExtensionsSet(d, "users"),
		TimeWindows:         expandIDNameExtensionsSet(d, "time_windows"),
		SrcIpGroups:         expandIDNameExtensionsSet(d, "src_ip_groups"),
		SrcIpv6Groups:       expandIDNameExtensionsSet(d, "src_ipv6_groups"),
		DestIpGroups:        expandIDNameExtensionsSet(d, "dest_ip_groups"),
		DestIpv6Groups:      expandIDNameExtensionsSet(d, "dest_ipv6_groups"),
		Labels:              expandIDNameExtensionsSet(d, "labels"),
		DeviceGroups:        expandIDNameExtensionsSet(d, "device_groups"),
		Devices:             expandIDNameExtensionsSet(d, "devices"),
		DNSGateway:          expandIDNameSet(d, "dns_gateway"),
		EDNSEcsObject:       expandIDNameSet(d, "edns_ecs_object"),
		ZPAIPGroup:          expandIDNameSet(d, "zpa_ip_group"),
	}
	return result
}
