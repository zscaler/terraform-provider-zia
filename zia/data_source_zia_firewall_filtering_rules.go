package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
)

// firewallFilteringRuleElementSchema returns the per-rule schema map used both
// as the top-level attributes (single-rule lookup mode) and as the Elem of the
// computed "rules" list (collection mode). Every field is Computed because the
// data source is read-only. The top-level builder layers Optional onto "id"
// and "name" for the lookup path and adds filter / search inputs.
//
// NOTE: this MUST return a fresh map on every call. terraform-plugin-sdk/v2
// mutates schema entries during validation, so sharing pointer-equal *schema.Schema
// values between the top level and the Elem map produces non-deterministic
// "field already set" panics.
func firewallFilteringRuleElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"order": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rank": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"access_control": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enable_full_logging": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"locations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"location_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"departments": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"users": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"time_windows": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"action": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_modified_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"last_modified_by": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"src_ips": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"src_ip_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"dest_addresses": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dest_ip_categories": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dest_countries": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dest_ip_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"nw_services": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"nw_service_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"nw_applications": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"nw_application_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"app_services": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"app_service_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"labels": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"device_groups": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"devices": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     idNameExtensionsElem(),
		},
		"device_trust_levels": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"zpa_app_segments": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "A unique identifier assigned to the Application Segment",
					},
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The name of the Application Segment",
					},
					"external_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Indicates the external ID. Applicable only when this reference is of an external entity.",
					},
				},
			},
		},
		"workload_groups": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The list of preconfigured workload groups to which the policy must be applied",
			Elem:        workloadGroupsElem(),
		},
		"default_rule": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"predefined": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

// idNameExtensionsElem returns a fresh Elem for the common
// {id, name, extensions} nested object used across many rule fields.
func idNameExtensionsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extensions": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// workloadGroupsElem returns a fresh Elem for the workload_groups field. Kept
// as a helper because the schema is deeply nested and would otherwise bloat
// firewallFilteringRuleElementSchema.
func workloadGroupsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "A unique identifier assigned to the workload group",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the workload group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the workload group",
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_modified_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     idNameExtensionsElem(),
			},
			"expression_json": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_container": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tags": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"value": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// firewallFilteringRuleFilterArgs returns the top-level filter inputs that map
// 1:1 to filteringrules.GetAllFilterOptions. These are forwarded to the API as
// query parameters so the result set is narrowed server-side BEFORE the
// optional JMESPath (search) is applied client-side.
//
// Five inputs carry a "filter_" prefix because their natural snake_case names
// (src_ips, dest_addresses, src_ip_groups, dest_ip_groups, nw_services,
// dest_ip_categories) collide with the per-rule OUTPUT attributes of the same
// names — which are list-of-objects, not strings. Every other filter uses the
// SDK field name verbatim.
func firewallFilteringRuleFilterArgs() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"predefined_rule_count": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If true, the API includes the predefined rule count in the response.",
		},
		"rule_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's ruleName query parameter (partial match).",
		},
		"rule_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's ruleLabel query parameter.",
		},
		"rule_label_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Server-side filter: matches the API's ruleLabelId query parameter.",
		},
		"rule_order": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's ruleOrder query parameter.",
		},
		"rule_description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's ruleDescription query parameter.",
		},
		"rule_action": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's ruleAction query parameter (e.g. ALLOW, BLOCK_DROP).",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's location query parameter.",
		},
		"department": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's department query parameter.",
		},
		"group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's group query parameter.",
		},
		"user": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's user query parameter.",
		},
		"device": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's device query parameter.",
		},
		"device_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's deviceGroup query parameter.",
		},
		"device_trust_level": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's deviceTrustLevel query parameter.",
		},
		"nw_application": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's nwApplication query parameter.",
		},
		"filter_src_ips": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's srcIps query parameter. Prefixed with filter_ to avoid colliding with the per-rule src_ips output attribute.",
		},
		"filter_dest_addresses": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's destAddresses query parameter. Prefixed with filter_ to avoid colliding with the per-rule dest_addresses output attribute.",
		},
		"filter_src_ip_groups": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's srcIpGroups query parameter. Prefixed with filter_ to avoid colliding with the per-rule src_ip_groups output attribute.",
		},
		"filter_dest_ip_groups": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's destIpGroups query parameter. Prefixed with filter_ to avoid colliding with the per-rule dest_ip_groups output attribute.",
		},
		"filter_nw_services": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's nwServices query parameter. Prefixed with filter_ to avoid colliding with the per-rule nw_services output attribute.",
		},
		"filter_dest_ip_categories": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Server-side filter: matches the API's destIpCategories query parameter. Prefixed with filter_ to avoid colliding with the per-rule dest_ip_categories output attribute.",
		},
	}
}

func dataSourceFirewallFilteringRule() *schema.Resource {
	// Build the top-level schema from three layers:
	//   1) every rule attribute (Computed) for the single-rule lookup path
	//   2) lookup inputs (rule_id, name) + the JMESPath search input
	//   3) the GetAllFilterOptions inputs + the rules computed list output
	//
	// IMPORTANT: we DROP the "id" entry from the element schema at the top
	// level. terraform-plugin-sdk/v2 reflects d.SetId(...) back into a
	// top-level schema field named "id" during State(); when that field is
	// declared as TypeInt the framework runs strconv.ParseInt on the resource
	// id and panics on values that aren't valid int64 (e.g. uint64s above
	// MaxInt64, or non-numeric collection-mode fingerprints). Renaming the
	// per-rule numeric id to "rule_id" at the top level lets us put any
	// string we want into d.SetId() without the framework trying to parse it.
	//
	// The "id" field stays inside the rules[] Elem schema unchanged, so users
	// can still do `for r in data...rules : r.id => r`.
	top := firewallFilteringRuleElementSchema()
	delete(top, "id")

	top["rule_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "Numeric id of a specific rule to look up. When set, the data source returns that rule's attributes at the top level and leaves rules[] empty. The matching rule's id is also exposed on the Terraform resource id (data.zia_firewall_filtering_rule.foo.id).",
	}
	top["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Name of a specific rule to look up (case-insensitive). When set, the data source returns that rule's attributes at the top level and leaves rules[] empty.",
	}

	top["search"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "JMESPath expression applied client-side after the SDK aggregates all pages. Use camelCase JSON field names (e.g. \"[?action=='ALLOW' && enableFullLogging]\").",
	}

	top["rules"] = &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of rules matching the filters. Populated when the data source is invoked without an explicit rule_id/name, or when the filters yield more than one result.",
		Elem: &schema.Resource{
			Schema: firewallFilteringRuleElementSchema(),
		},
	}

	for k, v := range firewallFilteringRuleFilterArgs() {
		top[k] = v
	}

	return &schema.Resource{
		ReadContext: dataSourceFirewallFilteringRuleRead,
		Schema:      top,
	}
}

func dataSourceFirewallFilteringRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Attach JMESPath expression (if any) BEFORE the SDK call so the
	// pagination engine applies it after aggregating pages.
	if searchExpr, ok := d.GetOk("search"); ok {
		ctx = zscaler.ContextWithJMESPath(ctx, searchExpr.(string))
		log.Printf("[INFO] zia_firewall_filtering_rule JMESPath filter set: %s\n", searchExpr.(string))
	}

	opts := buildFirewallFilteringFilterOptions(d)
	log.Printf("[INFO] zia_firewall_filtering_rule listing rules with opts=%+v\n", opts)

	allRules, err := filteringrules.GetAll(ctx, service, opts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing firewall filtering rules: %w", err))
	}
	log.Printf("[DEBUG] zia_firewall_filtering_rule retrieved %d rule(s) after server-side filters + JMESPath\n", len(allRules))

	// Single-rule lookup path: rule_id wins over name, both are exact matches.
	id, idProvided := getIntFromResourceData(d, "rule_id")
	nameStr, _ := d.Get("name").(string)
	nameProvided := nameStr != ""

	if idProvided {
		for i := range allRules {
			if allRules[i].ID == id {
				return setFirewallFilteringRuleTopLevel(d, &allRules[i])
			}
		}
		return diag.FromErr(fmt.Errorf("no firewall filtering rule found with rule_id %d (after filters)", id))
	}

	if nameProvided {
		for i := range allRules {
			if strings.EqualFold(allRules[i].Name, nameStr) {
				return setFirewallFilteringRuleTopLevel(d, &allRules[i])
			}
		}
		return diag.FromErr(fmt.Errorf("no firewall filtering rule found with name %q (after filters)", nameStr))
	}

	// Collection path: populate the rules list and leave the top-level
	// per-rule attributes at their zero values. The id of the data source is
	// a deterministic hash of the filter inputs so repeated plans converge.
	rulesList := make([]interface{}, 0, len(allRules))
	for i := range allRules {
		rulesList = append(rulesList, flattenFirewallFilteringRuleToMap(&allRules[i]))
	}
	if err := d.Set("rules", rulesList); err != nil {
		return diag.FromErr(fmt.Errorf("error setting rules: %w", err))
	}
	// The data source's Terraform-level id is metadata only; the real ids
	// live on each rules[*].id. Using a constant keeps it stable across
	// plans regardless of the filter inputs.
	d.SetId("firewall_filtering_rules")
	log.Printf("[DEBUG] zia_firewall_filtering_rule collection populated with %d rule(s)\n", len(rulesList))
	return nil
}

// buildFirewallFilteringFilterOptions translates the optional filter_ args
// into the SDK's GetAllFilterOptions. Returns nil if no filters are set.
func buildFirewallFilteringFilterOptions(d *schema.ResourceData) *filteringrules.GetAllFilterOptions {
	opts := &filteringrules.GetAllFilterOptions{}
	any := false

	str := func(key string, dst *string) {
		if v, ok := d.GetOk(key); ok {
			*dst = v.(string)
			any = true
		}
	}
	intv := func(key string, dst *int) {
		if v, ok := d.GetOk(key); ok {
			*dst = v.(int)
			any = true
		}
	}
	boolv := func(key string, dst *bool) {
		// GetOkExists is deprecated in plugin-sdk/v2 but still the only way to
		// disambiguate an explicit false from an unset bool. We avoid it by
		// only treating "true" as a meaningful filter request; if the user
		// needs to force false, that's the server's default and a no-op.
		if v, ok := d.GetOk(key); ok && v.(bool) {
			*dst = true
			any = true
		}
	}

	boolv("predefined_rule_count", &opts.PredefinedRuleCount)
	str("rule_name", &opts.RuleName)
	str("rule_label", &opts.RuleLabel)
	intv("rule_label_id", &opts.RuleLabelId)
	str("rule_order", &opts.RuleOrder)
	str("rule_description", &opts.RuleDescription)
	str("rule_action", &opts.RuleAction)
	str("location", &opts.Location)
	str("department", &opts.Department)
	str("group", &opts.Group)
	str("user", &opts.User)
	str("device", &opts.Device)
	str("device_group", &opts.DeviceGroup)
	str("device_trust_level", &opts.DeviceTrustLevel)
	str("nw_application", &opts.NwApplication)
	str("filter_src_ips", &opts.SrcIps)
	str("filter_dest_addresses", &opts.DestAddresses)
	str("filter_src_ip_groups", &opts.SrcIpGroups)
	str("filter_dest_ip_groups", &opts.DestIpGroups)
	str("filter_nw_services", &opts.NwServices)
	str("filter_dest_ip_categories", &opts.DestIpCategories)

	if !any {
		return nil
	}
	return opts
}

// setFirewallFilteringRuleTopLevel populates the top-level attributes from a
// single rule (single-rule lookup path). The rules list stays empty in this
// mode; consumers expecting a list should omit id/name and let the read
// function take the collection path.
func setFirewallFilteringRuleTopLevel(d *schema.ResourceData, resp *filteringrules.FirewallFilteringRules) diag.Diagnostics {
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("access_control", resp.AccessControl)
	_ = d.Set("enable_full_logging", resp.EnableFullLogging)
	_ = d.Set("action", resp.Action)
	_ = d.Set("state", resp.State)
	_ = d.Set("description", resp.Description)
	_ = d.Set("last_modified_time", resp.LastModifiedTime)
	_ = d.Set("src_ips", resp.SrcIps)
	_ = d.Set("dest_addresses", resp.DestAddresses)
	_ = d.Set("dest_ip_categories", resp.DestIpCategories)
	_ = d.Set("dest_countries", resp.DestCountries)
	_ = d.Set("nw_applications", resp.NwApplications)
	_ = d.Set("default_rule", resp.DefaultRule)
	_ = d.Set("predefined", resp.Predefined)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)

	if err := d.Set("locations", flattenIDNameExtensions(resp.Locations)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_groups", flattenIDNameExtensions(resp.LocationsGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("departments", flattenIDNameExtensions(resp.Departments)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("groups", flattenIDNameExtensions(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("users", flattenIDNameExtensions(resp.Users)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("time_windows", flattenIDNameExtensions(resp.TimeWindows)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.LastModifiedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("src_ip_groups", flattenIDNameExtensions(resp.SrcIpGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dest_ip_groups", flattenIDNameExtensions(resp.DestIpGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nw_services", flattenIDNameExtensions(resp.NwServices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nw_service_groups", flattenIDNameExtensions(resp.NwServiceGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nw_application_groups", flattenIDNameExtensions(resp.NwApplicationGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_services", flattenIDNameExtensions(resp.AppServices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("app_service_groups", flattenIDNameExtensions(resp.AppServiceGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("labels", flattenIDNameExtensions(resp.Labels)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("device_groups", flattenIDNameExtensions(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("devices", flattenIDNameExtensions(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workload_groups", flattenWorkloadGroups(resp.WorkloadGroups)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting workload_groups: %s", err))
	}
	if err := d.Set("zpa_app_segments", flattenZPAAppSegments(resp.ZPAAppSegments)); err != nil {
		return diag.FromErr(err)
	}
	// Always reset rules to empty in the single-rule path so plan output is
	// unambiguous about which mode the data source ran in.
	if err := d.Set("rules", []interface{}{}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// flattenFirewallFilteringRuleToMap projects a single API rule into the
// map[string]interface{} shape expected by the rules list Elem. Keys mirror
// firewallFilteringRuleElementSchema().
func flattenFirewallFilteringRuleToMap(r *filteringrules.FirewallFilteringRules) map[string]interface{} {
	return map[string]interface{}{
		"id":                    r.ID,
		"name":                  r.Name,
		"order":                 r.Order,
		"rank":                  r.Rank,
		"access_control":        r.AccessControl,
		"enable_full_logging":   r.EnableFullLogging,
		"action":                r.Action,
		"state":                 r.State,
		"description":           r.Description,
		"last_modified_time":    r.LastModifiedTime,
		"src_ips":               r.SrcIps,
		"dest_addresses":        r.DestAddresses,
		"dest_ip_categories":    r.DestIpCategories,
		"dest_countries":        r.DestCountries,
		"nw_applications":       r.NwApplications,
		"default_rule":          r.DefaultRule,
		"predefined":            r.Predefined,
		"device_trust_levels":   r.DeviceTrustLevels,
		"locations":             flattenIDNameExtensions(r.Locations),
		"location_groups":       flattenIDNameExtensions(r.LocationsGroups),
		"departments":           flattenIDNameExtensions(r.Departments),
		"groups":                flattenIDNameExtensions(r.Groups),
		"users":                 flattenIDNameExtensions(r.Users),
		"time_windows":          flattenIDNameExtensions(r.TimeWindows),
		"last_modified_by":      flattenLastModifiedBy(r.LastModifiedBy),
		"src_ip_groups":         flattenIDNameExtensions(r.SrcIpGroups),
		"dest_ip_groups":        flattenIDNameExtensions(r.DestIpGroups),
		"nw_services":           flattenIDNameExtensions(r.NwServices),
		"nw_service_groups":     flattenIDNameExtensions(r.NwServiceGroups),
		"nw_application_groups": flattenIDNameExtensions(r.NwApplicationGroups),
		"app_services":          flattenIDNameExtensions(r.AppServices),
		"app_service_groups":    flattenIDNameExtensions(r.AppServiceGroups),
		"labels":                flattenIDNameExtensions(r.Labels),
		"device_groups":         flattenIDNameExtensions(r.DeviceGroups),
		"devices":               flattenIDNameExtensions(r.Devices),
		"workload_groups":       flattenWorkloadGroups(r.WorkloadGroups),
		"zpa_app_segments":      flattenZPAAppSegments(r.ZPAAppSegments),
	}
}
