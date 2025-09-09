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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_control_rules"
)

var (
	bandwidthControlLock          sync.Mutex
	bandwidthControlStartingOrder int
)

func resourceBandwdithControlRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwdithControlRulesCreate,
		ReadContext:   resourceBandwdithControlRulesRead,
		UpdateContext: resourceBandwdithControlRulesUpdate,
		DeleteContext: resourceBandwdithControlRulesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("rule_id", idInt)
				} else {
					resp, err := bandwidth_control_rules.GetByName(ctx, service, id)
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
				Optional:    true,
				Description: "The bandwidth control rule name",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The description of the bandwidth control rule",
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The order of the bandwidth control rule",
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank of the Bandwidth Control policy rule",
			},
			"min_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The minimum percentage of a location's bandwidth you want to be guaranteed for each selected bandwidth control rule",
			},
			"max_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum percentage of a location's bandwidth you want to be guaranteed for each selected bandwidth control rule",
			},
			"bandwidth_classes": setIDsSchemaTypeCustom(nil, "The bandwidth control rulees to which you want to apply this rule"),
			"locations":         setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of locations for which rule must be applied"),
			"location_groups":   setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of the location groups to which the rule must be applied."),
			"labels":            setIDsSchemaTypeCustom(nil, "Labels that are applicable to the rule"),
			"time_windows":      setIDsSchemaTypeCustom(nil, "The Name-ID pairs of time windows to which the bandwidth control rule must be applied"),
			"protocols":         getURLProtocols(),
		},
	}
}

func resourceBandwdithControlRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandBandwdithControlRules(d)
	log.Printf("[INFO] Creating ZIA bandwitdh control rule\n%+v\n", req)

	// Create timeout for the operation
	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		bandwidthControlLock.Lock()
		if bandwidthControlStartingOrder == 0 {
			list, _ := bandwidth_control_rules.GetAll(ctx, service)
			for _, r := range list {
				// Ignore default rule
				if r.Order == 125 || r.Name == "Default Bandwidth Control" {
					continue
				}
				if r.Order > bandwidthControlStartingOrder {
					bandwidthControlStartingOrder = r.Order
				}
			}
			if bandwidthControlStartingOrder == 0 {
				bandwidthControlStartingOrder = 1
			}
		}
		bandwidthControlLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = bandwidthControlStartingOrder

		resp, err := bandwidth_control_rules.Create(ctx, service, &req)

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
					log.Printf("[INFO] Creating bandwidth control rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia bandwidth control rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorderWithBeforeReorder(
			OrderRule{Order: order, Rank: req.Rank},
			resp.ID,
			"bandwidth_control_rule",
			func() (int, error) {
				list, err := bandwidth_control_rules.GetAll(ctx, service)
				filteredList := filterOutBandwidthDefaultRule(list)
				return len(filteredList), err
			},
			func(id int, order OrderRule) error {
				rule, err := bandwidth_control_rules.Get(ctx, service, id)
				if err != nil {
					return err
				}
				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = bandwidth_control_rules.Update(ctx, service, id, rule)
				return err
			},
			nil,
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceBandwdithControlRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "bandwidth_control_rule")
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

func resourceBandwdithControlRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no bandwidth control rules id is set"))
	}
	resp, err := bandwidth_control_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia bandwidth control rules %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia bandwidth control rules:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("state", resp.State)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("min_bandwidth", resp.MinBandwidth)
	_ = d.Set("max_bandwidth", resp.MaxBandwidth)
	_ = d.Set("protocols", resp.Protocols)

	if err := d.Set("bandwidth_classes", flattenIDExtensionsListIDs(resp.BandwidthClasses)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_windows", flattenIDExtensionsListIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceBandwdithControlRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] bandwidth control rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("bandwidth control rule ID not set"))
	}
	log.Printf("[INFO] Updating bandwidth control rule ID: %v\n", id)
	req := expandBandwdithControlRules(d)

	if _, err := bandwidth_control_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := bandwidth_control_rules.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating bandwidth control rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorderWithBeforeReorder(
			OrderRule{Order: req.Order, Rank: req.Rank},
			req.ID,
			"bandwidth_control_rule",
			func() (int, error) {
				list, err := bandwidth_control_rules.GetAll(ctx, service)
				filteredList := filterOutBandwidthDefaultRule(list)
				return len(filteredList), err
			},
			func(id int, order OrderRule) error {
				rule, err := bandwidth_control_rules.Get(ctx, service, id)
				if err != nil {
					return err
				}
				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = bandwidth_control_rules.Update(ctx, service, id, rule)
				return err
			},
			nil,
		)

		if diags := resourceBandwdithControlRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "bandwidth_control_rule")
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

func resourceBandwdithControlRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] bandwidth control rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia bandwidth control rule ID: %v\n", (d.Id()))

	if _, err := bandwidth_control_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia bandwidth control rule deleted")

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandBandwdithControlRules(d *schema.ResourceData) bandwidth_control_rules.BandwidthControlRules {
	id, _ := getIntFromResourceData(d, "rule_id")
	result := bandwidth_control_rules.BandwidthControlRules{
		ID:               id,
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		State:            d.Get("state").(string),
		Order:            d.Get("order").(int),
		Rank:             d.Get("rank").(int),
		MinBandwidth:     d.Get("min_bandwidth").(int),
		MaxBandwidth:     d.Get("max_bandwidth").(int),
		Protocols:        SetToStringList(d, "protocols"),
		BandwidthClasses: expandIDNameExtensionsSet(d, "bandwidth_classes"),
		Labels:           expandIDNameExtensionsSet(d, "labels"),
		Locations:        expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:   expandIDNameExtensionsSet(d, "location_groups"),
		TimeWindows:      expandIDNameExtensionsSet(d, "time_windows"),
	}
	return result
}

func filterOutBandwidthDefaultRule(rules []bandwidth_control_rules.BandwidthControlRules) []bandwidth_control_rules.BandwidthControlRules {
	var filteredRules []bandwidth_control_rules.BandwidthControlRules
	for _, rule := range rules {
		if rule.Order != 125 && rule.Name != "Default Bandwidth Control" {
			filteredRules = append(filteredRules, rule)
		} else {
			log.Printf("[INFO] Ignoring default rule '%s' with order %d", rule.Name, rule.Order)
		}
	}
	return filteredRules
}
