package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_ueba_alerts/alert_definitions"
)

var supportedUEBAAlertNames = []string{
	"LDAP_SUCCESS", "LDAP_FAILURE", "LDAP_CONNECTION_DOWN", "AUTH_BRIDGE_DOWN",
	"ADP_SCHEDULE_UPDATE_FAILURE", "IDM_SCHEDULE_UPDATE_FAILURE",
	"EDM_SCHEMA_INDEXING_FAILURE", "OUTGOING_VIRUSES", "OUTGOING_SPYWARE",
	"OUTGOING_MALWARE", "OUTGOING_UNSCANNABLE_FILES", "INCOMING_VIRUSES",
	"INCOMING_SPYWARE", "INCOMING_MALWARE", "INCOMING_UNSCANNABLE_FILES",
	"BOTNET", "MALICIOUS_CONTENT", "PHISHING", "URL_FILTERING_BLOCKED_SITES",
	"STREAMING_UPLOAD", "STREAMING_VIEW_LISTEN", "SOCIAL_NETWORK_POST",
	"CHAT_FILE_TRANSFER", "WEBMAIL_FILE_ATTACHMENT", "HIPAA_VIOLATION",
	"PCI_VIOLATION", "GLBA_VIOLATION", "CUSTOM_ENGINE_VIOLATION",
	"POLICY_VIOLATION", "BA_ADWARE", "BA_MALWARE", "BA_ANONYMIZER",
	"ADVANCED_SECURITY", "PEER_TO_PEER", "UNAUTH_COMM", "CROSS_SITE_SCRIPTING",
	"BROWSER_EXPLOIT", "SUSPICIOUS_DESTINATION", "ADWARE_SPYWARE", "WEB_SPAM",
	"PAGE_RISK", "ADSPYWARE_SITES", "CRYPTOMINING", "TRAFFIC_DECREASE",
	"TRAFFIC_INCREASE", "DGA_DOMAINS", "BA_PATIENT0",
}

func resourceUEBAAlertDefinitions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUEBAAlertDefinitionsCreate,
		ReadContext:   resourceUEBAAlertDefinitionsRead,
		UpdateContext: resourceUEBAAlertDefinitionsUpdate,
		DeleteContext: resourceUEBAAlertDefinitionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("alert_definition_id", idInt)
				} else {
					resp, err := alert_definitions.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("alert_definition_id", resp.ID)
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
			"alert_definition_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DISABLED"}, false),
				Description:  "The status of the alert rule.",
			},
			"alert_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(supportedUEBAAlertNames, false),
				Description:  "The alert name that identifies the threat or event type the alert is triggered for.",
			},
			"occurrence": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OCCURRENCE_1", "OCCURRENCE_5", "OCCURRENCE_10", "OCCURRENCE_100", "OCCURRENCE_1000"}, false),
				Description:  "Specifies the occurrence of an ongoing alert for a specific threat or event type.",
			},
			"traffic_change_percent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the percentage change in traffic. Applicable to traffic-based alerts such as TRAFFIC_INCREASE and TRAFFIC_DECREASE.",
			},
			"interval": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"INTERVAL_5_MINUTES", "INTERVAL_15_MINUTES", "INTERVAL_30_MINUTES", "INTERVAL_1_HOUR", "INTERVAL_1_DAY"}, false),
				Description:  "The time span within which an event's occurrence triggers an alert.",
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "LOCATION", "DEPARTMENT", "ORGANIZATION"}, false),
				Description:  "Specifies if the alert is triggered for a user, location, department, or organization.",
			},
			"severity": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"}, false),
				Description:  "The threat severity based on which the alert is triggered.",
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the triggered alert.",
			},
			"entity": setSingleIDSchemaTypeCustom("An immutable reference to the entity (user, location, or department) the alert is scoped to."),
		},
	}
}

func resourceUEBAAlertDefinitionsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandUEBAAlertDefinitions(d)
	log.Printf("[INFO] Creating ZIA UEBA alert definition\n%+v\n", req)

	resp, _, err := alert_definitions.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA UEBA alert definition. ID: %v\n", resp.ID)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("alert_definition_id", resp.ID)

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceUEBAAlertDefinitionsRead(ctx, d, meta)
}

func resourceUEBAAlertDefinitionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "alert_definition_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no UEBA alert definition id is set"))
	}
	resp, err := alert_definitions.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia UEBA alert definition %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia UEBA alert definition:\n%+v\n", resp)

	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("alert_definition_id", resp.ID)
	_ = d.Set("status", resp.Status)
	_ = d.Set("alert_name", resp.AlertName)
	_ = d.Set("occurrence", resp.Occurrence)
	_ = d.Set("traffic_change_percent", resp.TrafficChangePercent)
	_ = d.Set("interval", resp.Interval)
	_ = d.Set("scope", resp.Scope)
	_ = d.Set("severity", resp.Severity)
	_ = d.Set("comments", resp.Comments)
	if err := d.Set("entity", flattenSingleIDNameExtensions(resp.Entity)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUEBAAlertDefinitionsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "alert_definition_id")
	if !ok {
		log.Printf("[ERROR] UEBA alert definition ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia UEBA alert definition ID: %v\n", id)
	req := expandUEBAAlertDefinitions(d)

	if _, err := alert_definitions.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := alert_definitions.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	if shouldActivate() {
		time.Sleep(2 * time.Second)
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceUEBAAlertDefinitionsRead(ctx, d, meta)
}

func resourceUEBAAlertDefinitionsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "alert_definition_id")
	if !ok {
		log.Printf("[ERROR] UEBA alert definition ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia UEBA alert definition ID: %v\n", d.Id())

	if _, err := alert_definitions.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia UEBA alert definition deleted")

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

func expandUEBAAlertDefinitions(d *schema.ResourceData) alert_definitions.AlertDefinitions {
	id, _ := getIntFromResourceData(d, "alert_definition_id")
	result := alert_definitions.AlertDefinitions{
		ID:                   id,
		Status:               d.Get("status").(string),
		AlertName:            d.Get("alert_name").(string),
		Occurrence:           d.Get("occurrence").(string),
		TrafficChangePercent: d.Get("traffic_change_percent").(int),
		Interval:             d.Get("interval").(string),
		Scope:                d.Get("scope").(string),
		Severity:             d.Get("severity").(string),
		Comments:             d.Get("comments").(string),
		Entity:               expandSingleIDNameExtensions(d, "entity"),
	}
	return result
}
