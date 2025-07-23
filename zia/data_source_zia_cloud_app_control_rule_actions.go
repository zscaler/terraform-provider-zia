package zia

import (
	"context"
	"fmt"
	"log"

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
			"available_actions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of available actions for the specified cloud applications and rule type.",
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

	// Set ID using a synthetic ID based on hash of inputs
	d.SetId(fmt.Sprintf("%s-%d", ruleType, len(cloudApps)))
	_ = d.Set("available_actions", actions)

	return nil
}
