package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_web_rules"
)

func resourceDlpWebRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceDlpWebRulesCreate,
		Read:   resourceDlpWebRulesRead,
		Update: resourceDlpWebRulesUpdate,
		Delete: resourceDlpWebRulesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					_ = d.Set("rule_id", id)
				} else {
					resp, err := zClient.dlp_web_rules.GetByName(id)
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
				Computed:    true,
				Description: "The DLP policy rule name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of the DLP policy rule.",
			},
			"access_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The access privilege for this DLP policy rule based on the admin's state.",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"READ_ONLY",
					"READ_WRITE",
				}, false),
			},
			"protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The protocol criteria specified for the DLP policy rule.",
			},
			"rank": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(1, 7),
				Description:  "Admin rank of the admin who creates this rule",
			},
			"order": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Order of execution of rule with respect to other URL Filtering rules",
			},
			"cloud_applications": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of cloud applications to which the DLP policy rule must be applied.",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 96000),
				Description:  "The minimum file size (in KB) used for evaluation of the DLP policy rule.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The action taken when traffic matches the DLP policy rule criteria.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"NONE",
					"BLOCK",
					"ALLOW",
					"ICAP_RESPONSE",
				}, false),
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables the DLP policy rule.",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"DISABLED",
				}, false),
			},
			"external_auditor_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The email address of an external auditor to whom DLP email notifications are sent.",
			},
			"match_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The match only criteria for DLP engines.",
			},
			"without_content_inspection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates a DLP policy rule without content inspection, when the value is set to true.",
			},
			"ocr_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables image file scanning.",
			},
			"zscaler_incident_reciever": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.",
			},
			"file_types":            getDLPRuleFileTypes("The list of file types to which the DLP policy rule must be applied."),
			"locations":             listIDsSchemaTypeCustom(8, "The Name-ID pairs of locations to which the DLP policy rule must be applied."),
			"location_groups":       listIDsSchemaTypeCustom(32, "The Name-ID pairs of locations groups to which the DLP policy rule must be applied."),
			"users":                 listIDsSchemaTypeCustom(4, "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"excluded_users":        listIDsSchemaTypeCustom(4, "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"groups":                listIDsSchemaTypeCustom(8, "The Name-ID pairs of groups to which the DLP policy rule must be applied."),
			"excluded_groups":       listIDsSchemaTypeCustom(4, "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"departments":           listIDsSchemaType("The Name-ID pairs of departments to which the DLP policy rule must be applied."),
			"excluded_departments":  listIDsSchemaTypeCustom(4, "The Name-ID pairs of users to which the DLP policy rule must be applied."),
			"dlp_engines":           listIDsSchemaTypeCustom(4, "The list of DLP engines to which the DLP policy rule must be applied."),
			"time_windows":          listIDsSchemaType("list of time interval during which rule must be enforced."),
			"labels":                listIDsSchemaType("list of Labels that are applicable to the rule."),
			"url_categories":        listIDsSchemaType("The list of URL categories to which the DLP policy rule must be applied."),
			"auditor":               listIDsSchemaTypeCustom(1, "The auditor to which the DLP policy rule must be applied."),
			"notification_template": listIDsSchemaTypeCustom(1, "The template used for DLP notification emails."),
			"icap_server":           listIDsSchemaTypeCustom(1, "The DLP server, using ICAP, to which the transaction content is forwarded."),
		},
	}
}

func resourceDlpWebRulesCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	req := expandDlpWebRules(d)
	errValidation := validateDlpWebRules(req)
	if errValidation != nil {
		return errValidation
	}
	log.Printf("[INFO] Creating zia web dlp rule\n%+v\n", req)

	resp, _, err := zClient.dlp_web_rules.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia web dlp rule request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("rule_id", resp.ID)

	return resourceDlpWebRulesRead(d, m)
}

func resourceDlpWebRulesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return fmt.Errorf("no zia web dlp rule id is set")
	}
	resp, err := zClient.dlp_web_rules.Get(id)

	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing web dlp rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting web dlp rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("access_control", resp.AccessControl)
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("description", resp.Description)
	_ = d.Set("file_types", resp.FileTypes)
	_ = d.Set("cloud_applications", resp.CloudApplications)
	_ = d.Set("state", resp.State)
	_ = d.Set("min_size", resp.MinSize)
	_ = d.Set("action", resp.Action)
	_ = d.Set("match_only", resp.MatchOnly)
	_ = d.Set("external_auditor_email", resp.ExternalAuditorEmail)
	_ = d.Set("without_content_inspection", resp.WithoutContentInspection)
	_ = d.Set("ocr_enabled", resp.OcrEnabled)
	_ = d.Set("zscaler_incident_reciever", resp.ZscalerIncidentReciever)

	if err := d.Set("locations", flattenIDExtensionsListIDs(resp.Locations)); err != nil {
		return err
	}

	if err := d.Set("location_groups", flattenIDExtensionsListIDs(resp.LocationGroups)); err != nil {
		return err
	}

	if err := d.Set("groups", flattenIDExtensionsListIDs(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("departments", flattenIDExtensionsListIDs(resp.Departments)); err != nil {
		return err
	}

	if err := d.Set("users", flattenIDExtensionsListIDs(resp.Users)); err != nil {
		return err
	}

	if err := d.Set("url_categories", flattenIDExtensionsListIDs(resp.URLCategories)); err != nil {
		return err
	}

	if err := d.Set("dlp_engines", flattenIDExtensionsListIDs(resp.DLPEngines)); err != nil {
		return err
	}

	if err := d.Set("time_windows", flattenIDExtensionsListIDs(resp.TimeWindows)); err != nil {
		return err
	}

	if err := d.Set("auditor", flattenIDExtensionListIDs(resp.Auditor)); err != nil {
		return err
	}

	if err := d.Set("notification_template", flattenIDExtensionListIDs(resp.NotificationTemplate)); err != nil {
		return err
	}

	if err := d.Set("icap_server", flattenIDExtensionListIDs(resp.IcapServer)); err != nil {
		return err
	}

	if err := d.Set("labels", flattenIDExtensionsListIDs(resp.Labels)); err != nil {
		return err
	}
	if err := d.Set("excluded_groups", flattenIDExtensions(resp.ExcludedGroups)); err != nil {
		return err
	}
	if err := d.Set("excluded_departments", flattenIDExtensions(resp.ExcludedDepartments)); err != nil {
		return err
	}
	if err := d.Set("excluded_users", flattenIDExtensions(resp.ExcludedUsers)); err != nil {
		return err
	}

	return nil
}

func resourceDlpWebRulesUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] web dlp rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating web dlp rule ID: %v\n", id)
	req := expandDlpWebRules(d)
	errValidation := validateDlpWebRules(req)
	if errValidation != nil {
		return errValidation
	}
	if _, _, err := zClient.dlp_web_rules.Update(id, &req); err != nil {
		return err
	}

	return resourceDlpWebRulesRead(d, m)
}

func resourceDlpWebRulesDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] web dlp rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting fweb dlp rule ID: %v\n", (d.Id()))

	if _, err := zClient.dlp_web_rules.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] web dlp rule deleted")
	return nil
}

func validateDlpWebRules(dlp dlp_web_rules.WebDLPRules) error {
	fileTypes := []string{"JPEG", "PNG", "TIFF", "BITMAP"}
	if dlp.OcrEnabled {
		// dlp.FileTypes must be a subset of fileTypes
		for _, t1 := range dlp.FileTypes {
			found := false
			for _, t2 := range fileTypes {
				if t1 == t2 {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("web dlp rule file types must be a subset of %v when OcrEnabled is disabled", fileTypes)
			}
		}
	} else if len(dlp.FileTypes) > 0 {
		return fmt.Errorf("web dlp rule file types must not be set when OcrEnabled is enabled")
	}
	return nil
}

func expandDlpWebRules(d *schema.ResourceData) dlp_web_rules.WebDLPRules {
	id, _ := getIntFromResourceData(d, "rule_id")
	result := dlp_web_rules.WebDLPRules{
		ID:                       id,
		Name:                     d.Get("name").(string),
		Order:                    d.Get("order").(int),
		Rank:                     d.Get("rank").(int),
		AccessControl:            d.Get("access_control").(string),
		Description:              d.Get("description").(string),
		Action:                   d.Get("action").(string),
		State:                    d.Get("state").(string),
		ExternalAuditorEmail:     d.Get("external_auditor_email").(string),
		MatchOnly:                d.Get("match_only").(bool),
		WithoutContentInspection: d.Get("without_content_inspection").(bool),
		OcrEnabled:               d.Get("ocr_enabled").(bool),
		ZscalerIncidentReciever:  d.Get("zscaler_incident_reciever").(bool),
		MinSize:                  d.Get("min_size").(int),
		Protocols:                SetToStringList(d, "protocols"),
		FileTypes:                SetToStringList(d, "file_types"),
		CloudApplications:        SetToStringList(d, "cloud_applications"),
		Auditor:                  expandIDNameExtensions(d, "auditor"),
		NotificationTemplate:     expandIDNameExtensions(d, "notification_template"),
		IcapServer:               expandIDNameExtensions(d, "icap_server"),
		Locations:                expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:           expandIDNameExtensionsSet(d, "location_groups"),
		Groups:                   expandIDNameExtensionsSet(d, "groups"),
		Departments:              expandIDNameExtensionsSet(d, "departments"),
		Users:                    expandIDNameExtensionsSet(d, "users"),
		URLCategories:            expandIDNameExtensionsSet(d, "url_categories"),
		DLPEngines:               expandIDNameExtensionsSet(d, "dlp_engines"),
		TimeWindows:              expandIDNameExtensionsSet(d, "time_windows"),
		Labels:                   expandIDNameExtensionsSet(d, "labels"),
		ExcludedUsers:            expandIDNameExtensionsSet(d, "excluded_groups"),
		ExcludedGroups:           expandIDNameExtensionsSet(d, "excluded_departments"),
		ExcludedDepartments:      expandIDNameExtensionsSet(d, "excluded_users"),
	}
	return result
}
