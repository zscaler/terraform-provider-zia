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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
)

var (
	sandboxLock          sync.Mutex
	sandboxStartingOrder int
)

func resourceSandboxRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSandboxRulesCreate,
		ReadContext:   resourceSandboxRulesRead,
		UpdateContext: resourceSandboxRulesUpdate,
		DeleteContext: resourceSandboxRulesDelete,

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
					resp, err := sandbox_rules.GetByName(ctx, service, id)
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
				Description: "The File Type Control policy rule name.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 10240),
				Description:  "Additional information about the Sandbox rule",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables the sandbox rules.",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 7),
				Description:  "Admin rank of the admin who creates this rule",
			},
			"order": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The rule order of execution for the  sandbox rules with respect to other rules.",
			},
			"ba_rule_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The action configured for the rule that must take place if the traffic matches the rule criteria.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK",
				}, false),
			},
			"first_time_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value indicating whether a First-Time Action is specifically configured for the rule",
			},
			"first_time_operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The action that must take place when users download unknown files for the first time",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW_SCAN",
					"QUARANTINE",
					"ALLOW_NOSCAN",
					"QUARANTINE_ISOLATE",
				}, false),
			},
			"ml_action_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When set to true, this indicates that 'Machine Learning Intelligence Action' checkbox has been checked on",
			},
			"by_threat_score": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "",
			},
			"locations":            setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of locations for the which policy must be applied. If not set, policy is applied for all locations."),
			"location_groups":      setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of locations groups for which rule must be applied."),
			"departments":          setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of departments to which the sandbox rules must be applied."),
			"groups":               setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of groups to which the sandbox rules must be applied."),
			"users":                setIDsSchemaTypeCustom(intPtr(4), "The Name-ID pairs of users to which the sandbox rules must be applied."),
			"labels":               setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"zpa_app_segments":     setExtIDNameSchemaCustom(intPtr(255), "List of Source IP Anchoring-enabled ZPA Application Segments for which this rule is applicable"),
			"url_categories":       getURLCategories(),
			"ba_policy_categories": getBaPolicyCategories(),
			"file_types":           getSandboxFileTypes(),
			"protocols":            getSandboxRuleProtocols(),
		},
	}
}

func resourceSandboxRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandSandboxRules(d)
	log.Printf("[INFO] Creating ZIA sandbox rule\n%+v\n", req)

	// Create timeout for the operation
	timeout := d.Timeout(schema.TimeoutCreate)
	start := time.Now()

	for {
		sandboxLock.Lock()
		if sandboxStartingOrder == 0 {
			list, _ := sandbox_rules.GetAll(ctx, service)
			for _, r := range list {
				// Ignore default rule
				if r.Order == 127 || r.Name == "Default BA Rule" {
					continue
				}
				if r.Order > sandboxStartingOrder {
					sandboxStartingOrder = r.Order
				}
			}
			if sandboxStartingOrder == 0 {
				sandboxStartingOrder = 1
			}
		}
		sandboxLock.Unlock()
		startWithoutLocking := time.Now()

		order := req.Order
		req.Order = sandboxStartingOrder

		resp, err := sandbox_rules.Create(ctx, service, &req)

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
					log.Printf("[INFO] Creating sandbox rule name: %v, got INVALID_INPUT_ARGUMENT\n", req.Name)
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia sandbox rule request. took:%s, without locking:%s,  ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		reorder(order, resp.ID, "sandbox_rules", func() (int, error) {
			list, err := sandbox_rules.GetAll(ctx, service)
			filteredList := filterOutDefaultRule(list)
			return len(filteredList), err
		}, func(id, order int) error {
			rule, err := sandbox_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = sandbox_rules.Update(ctx, service, id, rule)
			return err
		})

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		if diags := resourceSandboxRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(resp.ID, "sandbox_rules")
		break
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceSandboxRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia sandbox rules id is set"))
	}
	resp, err := sandbox_rules.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing sandbox rules %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	// Ignore Default BA Rule to prevent drift
	if resp.Order == 127 || resp.Name == "Default BA Rule" {
		log.Printf("[INFO] Skipping default rule '%s' with order %d to prevent drift", resp.Name, resp.Order)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Getting sandbox rules:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("state", resp.State)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("ba_rule_action", resp.BaRuleAction)
	_ = d.Set("first_time_enable", resp.FirstTimeEnable)
	_ = d.Set("first_time_operation", resp.FirstTimeOperation)
	_ = d.Set("ml_action_enabled", resp.MLActionEnabled)
	_ = d.Set("by_threat_score", resp.ByThreatScore)
	_ = d.Set("ba_policy_categories", resp.BaPolicyCategories)
	_ = d.Set("url_categories", resp.URLCategories)
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("file_types", resp.FileTypes)

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", flattenIDExtensionsListIDs(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenIDExtensionsListIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDExtensionsListIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSandboxRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] sandbox rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("sandbox rule ID not set"))
	}
	log.Printf("[INFO] Updating sandbox rule ID: %v\n", id)
	req := expandSandboxRules(d)

	if _, err := sandbox_rules.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	start := time.Now()

	for {
		_, err := sandbox_rules.Update(ctx, service, id, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				log.Printf("[INFO] Updating sandbox rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
				if time.Since(start) < timeout {
					time.Sleep(10 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorder(req.Order, req.ID, "sandbox_rules", func() (int, error) {
			list, err := sandbox_rules.GetAll(ctx, service)
			filteredList := filterOutDefaultRule(list)
			return len(filteredList), err
		}, func(id, order int) error {
			rule, err := sandbox_rules.Get(ctx, service, id)
			if err != nil {
				return err
			}
			rule.Order = order
			_, err = sandbox_rules.Update(ctx, service, id, rule)
			return err
		})

		if diags := resourceSandboxRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}
		markOrderRuleAsDone(req.ID, "sandbox_rules")
		break
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func resourceSandboxRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] sandbox rules not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting sandbox rules ID: %v\n", (d.Id()))

	if _, err := sandbox_rules.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] sandbox rules deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandSandboxRules(d *schema.ResourceData) sandbox_rules.SandboxRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandSandboxRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := sandbox_rules.SandboxRules{
		ID:                 id,
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Order:              order,
		Rank:               d.Get("rank").(int),
		State:              d.Get("state").(string),
		BaRuleAction:       d.Get("ba_rule_action").(string),
		FirstTimeEnable:    d.Get("first_time_enable").(bool),
		FirstTimeOperation: d.Get("first_time_operation").(string),
		MLActionEnabled:    d.Get("ml_action_enabled").(bool),
		ByThreatScore:      d.Get("by_threat_score").(int),
		Protocols:          SetToStringList(d, "protocols"),
		BaPolicyCategories: SetToStringList(d, "ba_policy_categories"),
		URLCategories:      SetToStringList(d, "url_categories"),
		FileTypes:          SetToStringList(d, "file_types"),
		Locations:          expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:     expandIDNameExtensionsSet(d, "location_groups"),
		Groups:             expandIDNameExtensionsSet(d, "groups"),
		Departments:        expandIDNameExtensionsSet(d, "departments"),
		Users:              expandIDNameExtensionsSet(d, "users"),
		Labels:             expandIDNameExtensionsSet(d, "labels"),
		ZPAAppSegments:     expandZPAAppSegmentSet(d, "zpa_app_segments"),
	}
	return result
}

func filterOutDefaultRule(rules []sandbox_rules.SandboxRules) []sandbox_rules.SandboxRules {
	var filteredRules []sandbox_rules.SandboxRules
	for _, rule := range rules {
		if rule.Order != 127 && rule.Name != "Default BA Rule" {
			filteredRules = append(filteredRules, rule)
		} else {
			log.Printf("[INFO] Ignoring default rule '%s' with order %d", rule.Name, rule.Order)
		}
	}
	return filteredRules
}
