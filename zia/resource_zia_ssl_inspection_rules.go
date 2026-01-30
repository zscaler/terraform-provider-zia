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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sslinspection"
)

var (
	sslInspectionLock          sync.Mutex
	sslInspectionStartingOrder int
)

func resourceSSLInspectionRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSLInspectionRulesCreate,
		ReadContext:   resourceSSLInspectionRulesRead,
		UpdateContext: resourceSSLInspectionRulesUpdate,
		DeleteContext: resourceSSLInspectionRulesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// Extract the action block
			actionList := d.Get("action").([]interface{})
			if len(actionList) == 0 {
				return fmt.Errorf("action block must be set")
			}
			actionMap := actionList[0].(map[string]interface{})

			// Extract action type
			actionType := actionMap["type"].(string)

			// Validate based on action type
			switch actionType {
			case "DECRYPT":
				// Ensure decryptSubActions is set and not empty
				decryptSubActionsList := actionMap["decrypt_sub_actions"].([]interface{})
				if len(decryptSubActionsList) == 0 {
					return fmt.Errorf("when action.type is 'DECRYPT', decrypt_sub_actions block must be set")
				}
				decryptSubActionsMap := decryptSubActionsList[0].(map[string]interface{})

				// Ensure all required fields in decryptSubActions are set
				if decryptSubActionsMap["server_certificates"].(string) == "" ||
					decryptSubActionsMap["min_client_tls_version"].(string) == "" ||
					decryptSubActionsMap["min_server_tls_version"].(string) == "" {
					return fmt.Errorf("when action.type is 'DECRYPT', all required fields in decrypt_sub_actions must be set")
				}

				// ssl_interception_cert is mandatory when action.type is DECRYPT and override_default_certificate is true
				overrideDefaultCertificate := actionMap["override_default_certificate"].(bool)
				sslInterceptionCertList := actionMap["ssl_interception_cert"].([]interface{})

				if overrideDefaultCertificate {
					if len(sslInterceptionCertList) == 0 {
						return fmt.Errorf("when action.type is 'DECRYPT' and override_default_certificate is true, ssl_interception_cert must be set")
					}
					// Ensure the cert block has a valid id set
					if certMap, ok := sslInterceptionCertList[0].(map[string]interface{}); ok {
						certID, hasID := certMap["id"]
						if !hasID || certID == nil {
							return fmt.Errorf("when action.type is 'DECRYPT' and override_default_certificate is true, ssl_interception_cert must have an id set")
						}
						switch v := certID.(type) {
						case int:
							if v == 0 {
								return fmt.Errorf("when action.type is 'DECRYPT' and override_default_certificate is true, ssl_interception_cert must have a non-zero id")
							}
						case int64:
							if v == 0 {
								return fmt.Errorf("when action.type is 'DECRYPT' and override_default_certificate is true, ssl_interception_cert must have a non-zero id")
							}
						}
					}
				}

				if actionMap["show_eun"].(bool) || actionMap["show_eunatp"].(bool) {
					return fmt.Errorf("when action.type is 'DECRYPT', neither show_eun nor show_eunatp can be set")
				}

			case "DO_NOT_DECRYPT":
				// Ensure doNotDecryptSubActions is set and not empty
				doNotDecryptSubActionsList := actionMap["do_not_decrypt_sub_actions"].([]interface{})
				if len(doNotDecryptSubActionsList) == 0 {
					return fmt.Errorf("when action.type is 'DO_NOT_DECRYPT', do_not_decrypt_sub_actions block must be set")
				}
				doNotDecryptSubActionsMap := doNotDecryptSubActionsList[0].(map[string]interface{})

				// If bypassOtherPolicies is true, serverCertificates and minTLSVersion cannot be set
				bypassOtherPolicies := doNotDecryptSubActionsMap["bypass_other_policies"].(bool)
				if bypassOtherPolicies {
					if doNotDecryptSubActionsMap["server_certificates"].(string) != "" ||
						doNotDecryptSubActionsMap["min_tls_version"].(string) != "" {
						return fmt.Errorf("when action.type is 'DO_NOT_DECRYPT' and bypass_other_policies is true, serverCertificates and minTLSVersion cannot be set")
					}
				} else {
					// If bypassOtherPolicies is false, ensure serverCertificates and minTLSVersion are set
					if doNotDecryptSubActionsMap["server_certificates"].(string) == "" ||
						doNotDecryptSubActionsMap["min_tls_version"].(string) == "" {
						return fmt.Errorf("when action.type is 'DO_NOT_DECRYPT' and bypass_other_policies is false, serverCertificates and minTLSVersion must be set")
					}
				}

			case "BLOCK":
				// Ensure decryptSubActions and doNotDecryptSubActions are not set
				if len(actionMap["decrypt_sub_actions"].([]interface{})) > 0 ||
					len(actionMap["do_not_decrypt_sub_actions"].([]interface{})) > 0 {
					return fmt.Errorf("when action.type is 'BLOCK', neither decrypt_sub_actions nor do_not_decrypt_sub_actions can be set")
				}

				// When action.type is BLOCK and overrideDefaultCertificate is false,
				// sslInterceptionCert cannot be set
				overrideDefaultCertificate := actionMap["override_default_certificate"].(bool)
				sslInterceptionCertList := actionMap["ssl_interception_cert"].([]interface{})
				if !overrideDefaultCertificate && len(sslInterceptionCertList) > 0 {
					return fmt.Errorf("when action.type is 'BLOCK' and override_default_certificate is false, ssl_interception_cert cannot be set")
				}

				if actionMap["show_eunatp"].(bool) {
					return fmt.Errorf("when action.type is 'BLOCK', show_eunatp cannot be set")
				}

			default:
				return fmt.Errorf("invalid action type: %s", actionType)
			}

			return nil
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
					resp, err := sslinspection.GetByName(ctx, service, id)
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
				Description: "The name of the SSL Inspection rule",
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Additional information about the SSL Inspection rule",
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables the SSL Inspection rules.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank of the admin who creates this rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The rule order of execution for the  SSL Inspection rules with respect to other rules.",
			},
			"road_warrior_for_kerberos": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When set to true, the rule is applied to remote users that use PAC with Kerberos authentication.",
			},
			"cloud_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of cloud applications to which the SSL Inspection rule must be applied
				Use the data source zia_cloud_applications to get the list of available cloud applications:
				https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_applications
				`,
			},
			"url_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of URL Categories to which the SSL inspection rule must be applied.
				See the URL Categories API for the list of available categories:
				https://help.zscaler.com/zia/url-categories#/urlCategories-get`,
			},
			"action": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the action configuration for the SSL inspection rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"BLOCK",
								"DECRYPT",
								"DO_NOT_DECRYPT",
							}, false),
						},
						"show_eun": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"show_eunatp": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"override_default_certificate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"ssl_interception_cert": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "SSL interception certificate. Required when action.type is 'DECRYPT'.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"decrypt_sub_actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_certificates": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ocsp_check": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"block_ssl_traffic_with_no_sni_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"min_client_tls_version": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"CLIENT_TLS_1_0",
											"CLIENT_TLS_1_1",
											"CLIENT_TLS_1_2",
											"CLIENT_TLS_1_3",
										}, false),
									},
									"min_server_tls_version": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SERVER_TLS_1_0",
											"SERVER_TLS_1_1",
											"SERVER_TLS_1_2",
											"SERVER_TLS_1_3",
										}, false),
									},
									"block_undecrypt": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"http2_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"do_not_decrypt_sub_actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bypass_other_policies": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"server_certificates": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ocsp_check": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"block_ssl_traffic_with_no_sni_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"min_tls_version": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"SERVER_TLS_1_0",
											"SERVER_TLS_1_1",
											"SERVER_TLS_1_2",
											"SERVER_TLS_1_3",
										}, false),
									},
								},
							},
						},
					},
				},
			},
			"locations":           setIDsSchemaTypeCustom(intPtr(8), "list of locations for which rule must be applied"),
			"location_groups":     setIDsSchemaTypeCustom(intPtr(32), "list of locations groups"),
			"users":               setIDsSchemaTypeCustom(intPtr(4), "list of users for which rule must be applied"),
			"groups":              setIDsSchemaTypeCustom(intPtr(8), "list of groups for which rule must be applied"),
			"departments":         setIDsSchemaTypeCustom(intPtr(140000), "list of departments for which rule must be applied"),
			"time_windows":        setIDsSchemaTypeCustom(intPtr(2), "The time interval in which the Firewall Filtering policy rule applies"),
			"labels":              setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"device_groups":       setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":             setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"source_ip_groups":    setIDsSchemaTypeCustom(nil, "list of source ip groups"),
			"dest_ip_groups":      setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"workload_groups":     setIdNameSchemaCustom(255, "The list of preconfigured workload groups to which the policy must be applied"),
			"proxy_gateways":      setIDsSchemaTypeCustom(nil, "The proxy chaining gateway for which this rule is applicable. Ignore if the forwarding method is not Proxy Chaining."),
			"zpa_app_segments":    setExtIDNameSchemaCustom(intPtr(255), "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method."),
			"user_agent_types":    getUserAgentTypes(),
			"device_trust_levels": getDeviceTrustLevels(),
			"platforms":           getSSLInspectionPlatforms(),
		},
	}
}

func resourceSSLInspectionRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandSSLInspectionRules(d)
	log.Printf("[INFO] Creating zia ssl inspection rule\n%+v\n", req)

	start := time.Now()

	for {
		sslInspectionLock.Lock()
		if sslInspectionStartingOrder == 0 {
			list, _ := sslinspection.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > sslInspectionStartingOrder {
					sslInspectionStartingOrder = r.Order
				}
			}
			if sslInspectionStartingOrder == 0 {
				sslInspectionStartingOrder = 1
			}
		}
		sslInspectionLock.Unlock()
		startWithoutLocking := time.Now()

		intendedOrder := req.Order
		intendedRank := req.Rank
		if intendedRank < 7 {
			// always start rank 7 rules at the next available order after all ranked rules
			req.Rank = 7
		}
		req.Order = sslInspectionStartingOrder
		resp, err := sslinspection.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, intendedOrder, req.Rank, currentFirewallOrderVsRankWording(ctx, zClient), err))
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia ssl inspection rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
		resourceType := "ssl_inspection_rules"

		reorderWithBeforeReorder(
			OrderRule{Order: intendedOrder, Rank: intendedRank},
			resp.ID,
			resourceType,
			func() (int, error) {
				allRules, err := sslinspection.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				// Custom updateOrder that handles predefined rules
				rule, err := sslinspection.Get(ctx, service, id)
				if err != nil {
					return err
				}
				// to avoid the STALE_CONFIGURATION_ERROR
				rule.LastModifiedTime = 0
				rule.LastModifiedBy = nil
				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = sslinspection.Update(ctx, service, id, rule)
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

		return resourceSSLInspectionRulesRead(ctx, d, meta)
	}
}

func resourceSSLInspectionRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia ssl inspection rule id is set"))
	}
	resp, err := sslinspection.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing ssl inspection rule rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting ssl inspection rule rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("state", resp.State)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("platforms", resp.Platforms)
	_ = d.Set("cloud_applications", resp.CloudApplications)
	_ = d.Set("road_warrior_for_kerberos", resp.RoadWarriorForKerberos)
	_ = d.Set("url_categories", resp.URLCategories)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
	_ = d.Set("user_agent_types", resp.UserAgentTypes)

	if err := d.Set("action", flattenSSLInspectionAction(resp.Action)); err != nil {
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

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("source_ip_groups", flattenIDExtensionsListIDs(resp.SourceIPGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("dest_ip_groups", flattenIDExtensionsListIDs(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("proxy_gateways", flattenIDExtensionsListIDs(resp.ProxyGateways)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("zpa_app_segments", flattenZPAAppSegmentsSimple(resp.ZPAAppSegments)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
	}

	return nil
}

func resourceSSLInspectionRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] ssl inspection rule ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("ssl inspection rule ID not set"))
	}
	log.Printf("[INFO] Updating ssl inspection rule ID: %v\n", id)
	req := expandSSLInspectionRules(d)

	if _, err := sslinspection.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	existingRules, err := sslinspection.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all ssl inspection rules: %v", err)
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

	_, err = sslinspection.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating ssl inspection rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "ssl_inspection_rules",
		func() (int, error) {
			allRules, err := sslinspection.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := sslinspection.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// to avoid the STALE_CONFIGURATION_ERROR
			rule.LastModifiedTime = 0
			rule.LastModifiedBy = nil
			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = sslinspection.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceSSLInspectionRulesRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "ssl_inspection_rules")

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

func resourceSSLInspectionRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] ssl inspection rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting ssl inspection rule ID: %v\n", (d.Id()))

	if _, err := sslinspection.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] ssl inspection rule deleted")

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

func expandSSLInspectionRules(d *schema.ResourceData) sslinspection.SSLInspectionRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandSSLInspectionRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := sslinspection.SSLInspectionRules{
		ID:                     id,
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		Order:                  order,
		Rank:                   d.Get("rank").(int),
		State:                  d.Get("state").(string),
		RoadWarriorForKerberos: d.Get("road_warrior_for_kerberos").(bool),
		CloudApplications:      SetToStringList(d, "cloud_applications"),
		DeviceTrustLevels:      SetToStringList(d, "device_trust_levels"),
		Platforms:              SetToStringList(d, "platforms"),
		UserAgentTypes:         SetToStringList(d, "user_agent_types"),
		URLCategories:          SetToStringList(d, "url_categories"),
		DeviceGroups:           expandIDNameExtensionsSet(d, "device_groups"),
		Devices:                expandIDNameExtensionsSet(d, "devices"),
		Locations:              expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:         expandIDNameExtensionsSet(d, "location_groups"),
		Groups:                 expandIDNameExtensionsSet(d, "groups"),
		Departments:            expandIDNameExtensionsSet(d, "departments"),
		SourceIPGroups:         expandIDNameExtensionsSet(d, "source_ip_groups"),
		DestIpGroups:           expandIDNameExtensionsSet(d, "dest_ip_groups"),
		Users:                  expandIDNameExtensionsSet(d, "users"),
		TimeWindows:            expandIDNameExtensionsSet(d, "time_windows"),
		Labels:                 expandIDNameExtensionsSet(d, "labels"),
		ProxyGateways:          expandIDNameExtensionsSet(d, "proxy_gateways"),
		ZPAAppSegments:         expandZPAAppSegmentSet(d, "zpa_app_segments"),
		WorkloadGroups:         expandWorkloadGroupsIDName(d, "workload_groups"),
	}

	// ADD THIS:
	// Parse the 'action' block from Terraform into your Action struct
	result.Action = expandSSLInspectionAction(d.Get("action"))

	return result
}

func expandSSLInspectionAction(v interface{}) sslinspection.Action {
	var action sslinspection.Action

	list, ok := v.([]interface{})
	if !ok || len(list) < 1 {
		return action
	}

	raw := list[0].(map[string]interface{})

	// Simple fields
	if val, ok := raw["type"].(string); ok {
		action.Type = val
	}
	if val, ok := raw["show_eun"].(bool); ok {
		action.ShowEUN = val
	}
	if val, ok := raw["show_eunatp"].(bool); ok {
		action.ShowEUNATP = val
	}
	if val, ok := raw["override_default_certificate"].(bool); ok {
		action.OverrideDefaultCertificate = val
	}

	// ssl_interception_cert block
	if certList, ok := raw["ssl_interception_cert"].([]interface{}); ok && len(certList) > 0 {
		action.SSLInterceptionCert = expandSSLInterceptionCert(certList[0].(map[string]interface{}))
	}

	// do_not_decrypt_sub_actions block
	if doNotDecryptList, ok := raw["do_not_decrypt_sub_actions"].([]interface{}); ok && len(doNotDecryptList) > 0 {
		action.DoNotDecryptSubActions = expandSSLInspectionDoNotDecryptSubActions(doNotDecryptList[0].(map[string]interface{}))
	}

	// decrypt_sub_actions block (NEW)
	if decryptList, ok := raw["decrypt_sub_actions"].([]interface{}); ok && len(decryptList) > 0 {
		action.DecryptSubActions = expandSSLInspectionDecryptSubActions(decryptList[0].(map[string]interface{}))
	}

	return action
}

func expandSSLInterceptionCert(m map[string]interface{}) *sslinspection.SSLInterceptionCert {
	cert := &sslinspection.SSLInterceptionCert{}

	if v, ok := m["id"].(int); ok {
		cert.ID = v
	}
	// if v, ok := m["name"].(string); ok {
	// 	cert.Name = v
	// }
	// if v, ok := m["default_certificate"].(bool); ok {
	// 	cert.DefaultCertificate = v
	// }

	return cert
}

func expandSSLInspectionDecryptSubActions(m map[string]interface{}) *sslinspection.DecryptSubActions {
	sub := &sslinspection.DecryptSubActions{}

	if v, ok := m["server_certificates"].(string); ok {
		sub.ServerCertificates = v
	}
	if v, ok := m["ocsp_check"].(bool); ok {
		sub.OcspCheck = v
	}
	if v, ok := m["block_ssl_traffic_with_no_sni_enabled"].(bool); ok {
		sub.BlockSslTrafficWithNoSniEnabled = v
	}
	if v, ok := m["min_client_tls_version"].(string); ok {
		sub.MinClientTLSVersion = v
	}
	if v, ok := m["min_server_tls_version"].(string); ok {
		sub.MinServerTLSVersion = v
	}
	if v, ok := m["block_undecrypt"].(bool); ok {
		sub.BlockUndecrypt = v
	}
	if v, ok := m["http2_enabled"].(bool); ok {
		sub.HTTP2Enabled = v
	}

	return sub
}

func expandSSLInspectionDoNotDecryptSubActions(m map[string]interface{}) *sslinspection.DoNotDecryptSubActions {
	sub := &sslinspection.DoNotDecryptSubActions{}

	if v, ok := m["bypass_other_policies"].(bool); ok {
		sub.BypassOtherPolicies = v
	}
	if v, ok := m["server_certificates"].(string); ok {
		sub.ServerCertificates = v
	}
	if v, ok := m["ocsp_check"].(bool); ok {
		sub.OcspCheck = v
	}
	if v, ok := m["block_ssl_traffic_with_no_sni_enabled"].(bool); ok {
		sub.BlockSslTrafficWithNoSniEnabled = v
	}
	if v, ok := m["min_tls_version"].(string); ok {
		// Make sure your Go struct matches what the API expects (minTLSVersion vs minServerTLSVersion)
		sub.MinTLSVersion = v
	}

	return sub
}
