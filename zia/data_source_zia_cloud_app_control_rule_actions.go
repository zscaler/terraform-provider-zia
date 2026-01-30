package zia

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudappcontrol"
)

func dataSourceCloudAppControlRuleActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudAppControlRuleActionsRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule type for the Cloud App Control policy (e.g., 'web', 'email').",
			},
			"cloud_apps": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of cloud applications to retrieve available actions for.",
			},
			"action_prefixes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Optional list of action prefixes to filter results. Valid values: ALLOW, DENY, BLOCK, CAUTION, ISOLATE, ESC. The underscore is automatically added. If specified, only actions starting with these prefixes will be included in filtered_actions.",
			},
			"available_actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of all available actions for the specified cloud applications and rule type (includes ISOLATE actions).",
			},
			"available_actions_without_isolate": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of available actions excluding ISOLATE actions. Use this for standard rules. ISOLATE actions cannot be mixed with other actions.",
			},
			"isolate_actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of only ISOLATE actions. Use this for Cloud Browser Isolation rules. ISOLATE actions require cbi_profile configuration and cannot be mixed with other actions.",
			},
			"filtered_actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of actions filtered by action_prefixes parameter. Only populated if action_prefixes is specified. Use this for custom filtering by action type (ALLOW_, DENY_, BLOCK_, etc.).",
			},
		},
	}
}

func dataSourceCloudAppControlRuleActionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	ruleType := d.Get("type").(string)
	rawCloudApps := d.Get("cloud_apps").([]interface{})
	cloudApps := make([]string, len(rawCloudApps))
	for i, app := range rawCloudApps {
		cloudApps[i] = app.(string)
	}

	log.Printf("[INFO] Calling All Available Actions for ruleType %q with apps: %v", ruleType, cloudApps)

	payload := cloudappcontrol.AvailableActionsRequest{
		CloudApps: cloudApps,
		Type:      ruleType,
	}

	actions, err := cloudappcontrol.AllAvailableActions(ctx, service, ruleType, payload)
	if err != nil {
		return diag.FromErr(err)
	}

	// Separate ISOLATE actions from non-ISOLATE actions
	var actionsWithoutIsolate []string
	var isolateActions []string

	for _, action := range actions {
		if strings.HasPrefix(action, "ISOLATE_") {
			isolateActions = append(isolateActions, action)
		} else {
			actionsWithoutIsolate = append(actionsWithoutIsolate, action)
		}
	}

	// Handle custom prefix filtering if specified
	var filteredActions []string
	if prefixesInterface, ok := d.GetOk("action_prefixes"); ok {
		prefixesList := prefixesInterface.([]interface{})
		prefixes := make([]string, len(prefixesList))
		for i, p := range prefixesList {
			prefix := p.(string)
			// Automatically add underscore if not present
			if !strings.HasSuffix(prefix, "_") {
				prefix = prefix + "_"
			}
			prefixes[i] = prefix
		}

		// Filter actions by specified prefixes
		for _, action := range actions {
			for _, prefix := range prefixes {
				if strings.HasPrefix(action, prefix) {
					filteredActions = append(filteredActions, action)
					break
				}
			}
		}
		log.Printf("[DEBUG] Filtered actions by prefixes %v: %d actions matched", prefixes, len(filteredActions))
	}

	// Set ID using a synthetic ID based on hash of inputs
	d.SetId(fmt.Sprintf("%s-%d", ruleType, len(cloudApps)))
	_ = d.Set("available_actions", actions)
	_ = d.Set("available_actions_without_isolate", actionsWithoutIsolate)
	_ = d.Set("isolate_actions", isolateActions)
	_ = d.Set("filtered_actions", filteredActions)

	return nil
}
