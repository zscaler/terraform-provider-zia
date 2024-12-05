package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
)

// var (
// 	sandboxLock          sync.Mutex
// 	sandboxStartingOrder int
// )

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
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The rule order of execution for the  sandbox rules with respect to other rules.",
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
			"device_groups":        setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":              setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"locations":            setIDsSchemaTypeCustom(intPtr(8), "Name-ID pairs of locations for the which policy must be applied. If not set, policy is applied for all locations."),
			"location_groups":      setIDsSchemaTypeCustom(intPtr(32), "Name-ID pairs of locations groups for which rule must be applied."),
			"departments":          setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of departments to which the sandbox rules must be applied."),
			"groups":               setIDsSchemaTypeCustom(intPtr(8), "The Name-ID pairs of groups to which the sandbox rules must be applied."),
			"users":                setIDsSchemaTypeCustom(intPtr(4), "The Name-ID pairs of users to which the sandbox rules must be applied."),
			"time_windows":         setIDsSchemaTypeCustom(intPtr(2), "list of time interval during which rule must be enforced."),
			"labels":               setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"zpa_app_segments":     setExtIDNameSchemaCustom(intPtr(255), "List of Source IP Anchoring-enabled ZPA Application Segments for which this rule is applicable"),
			"device_trust_levels":  getDeviceTrustLevels(),
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
		// Attempt to create the sandbox rule
		resp, err := sandbox_rules.Create(ctx, service, &req)
		if err != nil {
			// Handle specific retry scenarios
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") && time.Since(start) < timeout {
				log.Printf("[WARN] Retrying sandbox rule creation due to: %s", err)
				time.Sleep(5 * time.Second) // Wait before retrying
				continue
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		// Log successful creation and set resource ID
		log.Printf("[INFO] Created ZIA sandbox rule successfully. ID: %v", resp.ID)
		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		// Validate by reading the created resource
		if diags := resourceSandboxRulesRead(ctx, d, meta); diags.HasError() {
			if time.Since(start) < timeout {
				log.Printf("[WARN] Retrying sandbox rule read after creation.")
				time.Sleep(10 * time.Second) // Wait before retrying
				continue
			}
			return diags
		}

		break
	}

	// Optional: Sleep briefly before triggering activation
	time.Sleep(2 * time.Second)

	// Trigger activation if required
	if shouldActivate() {
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(fmt.Errorf("error triggering activation: %s", activationErr))
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation as ZIA_ACTIVATION is not set to true.")
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
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("file_types", resp.FileTypes)

	if err := d.Set("devices", flattenIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("device_groups", flattenIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", flattenIDs(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location_groups", flattenIDs(resp.LocationGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", flattenIDs(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("departments", flattenIDs(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("users", flattenIDs(resp.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("time_windows", flattenIDs(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", flattenIDs(resp.Labels)); err != nil {
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
		if err != nil {
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if time.Since(start) < timeout {
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}
			}
			return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
		}

		reorder(req.Order, req.ID, "sandbox_rules", func() (int, error) {
			list, err := sandbox_rules.GetAll(ctx, service)
			return len(list), err
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
		if activationErr := triggerActivation(zClient); activationErr != nil {
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
		if activationErr := triggerActivation(zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandSandboxRules(d *schema.ResourceData) sandbox_rules.SandboxRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	result := sandbox_rules.SandboxRules{
		ID:                 id,
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Order:              d.Get("order").(int),
		Rank:               d.Get("rank").(int),
		State:              d.Get("state").(string),
		BaRuleAction:       d.Get("ba_rule_action").(string),
		FirstTimeEnable:    d.Get("first_time_enable").(bool),
		FirstTimeOperation: d.Get("first_time_operation").(string),
		MLActionEnabled:    d.Get("ml_action_enabled").(bool),
		ByThreatScore:      d.Get("by_threat_score").(int),
		Protocols:          SetToStringList(d, "protocols"),
		BaPolicyCategories: SetToStringList(d, "ba_policy_categories"),
		FileTypes:          SetToStringList(d, "file_types"),
		DeviceGroups:       expandIDNameExtensionsSet(d, "device_groups"),
		Devices:            expandIDNameExtensionsSet(d, "devices"),
		Locations:          expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:     expandIDNameExtensionsSet(d, "location_groups"),
		Groups:             expandIDNameExtensionsSet(d, "groups"),
		Departments:        expandIDNameExtensionsSet(d, "departments"),
		Users:              expandIDNameExtensionsSet(d, "users"),
		TimeWindows:        expandIDNameExtensionsSet(d, "time_windows"),
		Labels:             expandIDNameExtensionsSet(d, "labels"),
		ZPAAppSegments:     expandZPAAppSegmentSet(d, "zpa_app_segments"),
	}
	return result
}