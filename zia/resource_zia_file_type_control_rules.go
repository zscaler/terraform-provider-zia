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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol"
)

var (
	fileTypeLock          sync.Mutex
	fileTypeStartingOrder int
)

func resourceFileTypeControlRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileTypeControlRulesCreate,
		ReadContext:   resourceFileTypeControlRulesRead,
		UpdateContext: resourceFileTypeControlRulesUpdate,
		DeleteContext: resourceFileTypeControlRulesDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// Check if active_content is enabled
			if activeContent, ok := d.GetOk("active_content"); ok && activeContent.(bool) {
				// Validate file_types
				if fileTypes, ok := d.GetOk("file_types"); ok {
					allowedFileTypes := []string{
						"FTCATEGORY_MS_WORD",
						"FTCATEGORY_MS_POWERPOINT",
						"FTCATEGORY_PDF_DOCUMENT",
						"FTCATEGORY_MS_EXCEL",
					}
					fileTypeSet := fileTypes.(*schema.Set).List()
					for _, fileType := range fileTypeSet {
						// Ensure each fileType is one of the allowedFileTypes
						fileTypeStr := fileType.(string)
						if !contains(allowedFileTypes, fileTypeStr) {
							return fmt.Errorf(
								"attribute 'active_content' can only be enabled when 'file_types' contains only the following: %v. Invalid value: %s",
								allowedFileTypes, fileTypeStr,
							)
						}
					}
				} else {
					return fmt.Errorf(
						"attribute 'active_content' requires 'file_types' to be set with one of the following values: %v",
						[]string{"FTCATEGORY_MS_WORD", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_PDF_DOCUMENT", "FTCATEGORY_MS_EXCEL"},
					)
				}

				// Validate that unscanned is not enabled
				if unscannable, ok := d.GetOk("unscannable"); ok && unscannable.(bool) {
					return fmt.Errorf(
						"attribute 'unscannable' cannot be enabled when 'active_content' is enabled",
					)
				}
			}

			return nil
		},

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
					resp, err := filetypecontrol.GetByName(ctx, service, id)
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
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringLenBetween(0, 10240),
				Description:      "The description of the File Type Control rule.",
				StateFunc:        normalizeMultiLineString, // Ensures correct format before storing in Terraform state
				DiffSuppressFunc: noChangeInMultiLineText,  // Prevents unnecessary Terraform diffs
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables the File Type Control rule.",
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
				Description:  "The rule order of execution for the  File Type Control rule with respect to other rules.",
			},
			"filtering_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Action taken when traffic matches policy. This field is not applicable to the Lite API.",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW",
					"BLOCK",
					"CAUTION",
				}, false),
			},
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "File operation performed. This field is not applicable to the Lite API.",
				ValidateFunc: validation.StringInSlice([]string{
					"UPLOAD",
					"DOWNLOAD",
					"UPLOAD_DOWNLOAD",
				}, false),
			},
			"min_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 409600),
				Description:  "Minimum file size (in KB) used for evaluation of the FTP rule",
			},
			"max_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 409600),
				Description:  "Maximum file size (in KB) used for evaluation of the FTP rule",
			},
			"capture_pcap": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "A Boolean value that indicates whether packet capture (PCAP) is enabled or not",
			},
			"active_content": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Flag to check whether a file has active content or not",
			},
			"unscannable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Flag to check whether a file has active content or not",
			},
			"browser_eun_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "",
			},
			"cloud_applications": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of cloud applications to which the File Type Control rule must be applied
				Use the data source zia_cloud_applications to get the list of available cloud applications:
				https://registry.terraform.io/providers/zscaler/zia/latest/docs/data-sources/zia_cloud_applications
				`,
			},
			"url_categories": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: `The list of URL Categories to which the file type control rule rule must be applied.
				See the URL Categories API for the list of available categories:
				https://help.zscaler.com/zia/url-categories#/urlCategories-get`,
			},
			"file_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: `File type categories for which the policy is applied. If not set, the rule is applied across all file types
				See the available file types API in:
				https://help.zscaler.com/zia/file-type-control-policy#/fileTypeRules-get`,
			},
			"device_groups":       setIDsSchemaTypeCustom(nil, "This field is applicable for devices that are managed using Zscaler Client Connector."),
			"devices":             setIDsSchemaTypeCustom(nil, "Name-ID pairs of devices for which rule must be applied."),
			"locations":           setIDsSchemaTypeCustom(nil, "Name-ID pairs of locations for the which policy must be applied. If not set, policy is applied for all locations."),
			"location_groups":     setIDsSchemaTypeCustom(nil, "Name-ID pairs of locations groups for which rule must be applied."),
			"departments":         setIDsSchemaTypeCustom(nil, "The Name-ID pairs of departments to which the File Type Control rule must be applied."),
			"groups":              setIDsSchemaTypeCustom(nil, "The Name-ID pairs of groups to which the File Type Control rule must be applied."),
			"users":               setIDsSchemaTypeCustom(nil, "The Name-ID pairs of users to which the File Type Control rule must be applied."),
			"time_windows":        setIDsSchemaTypeCustom(intPtr(2), "list of time interval during which rule must be enforced."),
			"labels":              setIDsSchemaTypeCustom(intPtr(1), "list of Labels that are applicable to the rule."),
			"zpa_app_segments":    setExtIDNameSchemaCustom(intPtr(255), "List of Source IP Anchoring-enabled ZPA Application Segments for which this rule is applicable"),
			"device_trust_levels": getDeviceTrustLevels(),
			"protocols":           getFileTypeProtocols(),
		},
	}
}

func resourceFileTypeControlRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandFileTypeControlRules(d)
	log.Printf("[INFO] Creating zia file type control rule rule\n%+v\n", req)

	start := time.Now()

	for {
		fileTypeLock.Lock()
		if fileTypeStartingOrder == 0 {
			list, _ := filetypecontrol.GetAll(ctx, service)
			for _, r := range list {
				if r.Order > fileTypeStartingOrder {
					fileTypeStartingOrder = r.Order
				}
			}
			if fileTypeStartingOrder == 0 {
				fileTypeStartingOrder = 1
			}
		}
		fileTypeLock.Unlock()
		startWithoutLocking := time.Now()

		intendedOrder := req.Order
		intendedRank := req.Rank
		if intendedRank < 7 {
			// always start rank 7 rules at the next available order after all ranked rules
			req.Rank = 7
		}
		req.Order = fileTypeStartingOrder
		resp, err := filetypecontrol.Create(ctx, service, &req)

		// Fail immediately if INVALID_INPUT_ARGUMENT is detected
		if customErr := failFastOnErrorCodes(err); customErr != nil {
			return diag.Errorf("%v", customErr)
		}

		if err != nil {
			reg := regexp.MustCompile("Rule with rank [0-9]+ is not allowed at order [0-9]+")
			if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
				if reg.MatchString(err.Error()) {
					return diag.FromErr(fmt.Errorf("error creating resource: %s, please check the order %d vs rank %d, current rules:%s , err:%s", req.Name, intendedOrder, req.Rank, currentFileTypeControlOrderVsRankWording(ctx, zClient), err))
				}
			}
			return diag.FromErr(fmt.Errorf("error creating resource: %s", err))
		}

		log.Printf("[INFO] Created zia file type control rule rule request. Took: %s, without locking: %s, ID: %v\n", time.Since(start), time.Since(startWithoutLocking), resp)
		// Use separate resource type for rank 7 rules to avoid mixing with ranked rules
		resourceType := "file_type_control_rules"

		reorderWithBeforeReorder(
			OrderRule{Order: intendedOrder, Rank: intendedRank},
			resp.ID,
			resourceType,
			func() (int, error) {
				allRules, err := filetypecontrol.GetAll(ctx, service)
				if err != nil {
					return 0, err
				}
				// Count all rules including predefined ones for proper ordering
				return len(allRules), nil
			},
			func(id int, order OrderRule) error {
				// Custom updateOrder that handles predefined rules
				rule, err := filetypecontrol.Get(ctx, service, id)
				if err != nil {
					return err
				}
				// to avoid the STALE_CONFIGURATION_ERROR
				rule.LastModifiedTime = 0
				rule.LastModifiedBy = nil
				rule.Order = order.Order
				rule.Rank = order.Rank
				_, err = filetypecontrol.Update(ctx, service, id, rule)
				return err
			},
			nil, // Remove beforeReorder function to avoid adding too many rules to the map
		)

		d.SetId(strconv.Itoa(resp.ID))
		_ = d.Set("rule_id", resp.ID)

		markOrderRuleAsDone(resp.ID, resourceType)
		waitForReorder(resourceType)

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

		return resourceFileTypeControlRulesRead(ctx, d, meta)
	}
}

func resourceFileTypeControlRulesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no zia file type control rule id is set"))
	}
	resp, err := filetypecontrol.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing file type control rule %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting file type control rule:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rule_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("description", resp.Description)
	_ = d.Set("state", resp.State)
	_ = d.Set("order", resp.Order)
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("filtering_action", resp.FilteringAction)
	_ = d.Set("operation", resp.Operation)
	_ = d.Set("active_content", resp.ActiveContent)
	_ = d.Set("unscannable", resp.Unscannable)
	_ = d.Set("browser_eun_template_id", resp.BrowserEunTemplateID)
	_ = d.Set("capture_pcap", resp.CapturePCAP)
	_ = d.Set("protocols", resp.Protocols)
	_ = d.Set("file_types", resp.FileTypes)
	_ = d.Set("cloud_applications", resp.CloudApplications)
	_ = d.Set("url_categories", resp.URLCategories)
	_ = d.Set("device_trust_levels", resp.DeviceTrustLevels)
	_ = d.Set("max_size", resp.MaxSize)
	_ = d.Set("min_size", resp.MinSize)

	if err := d.Set("device_groups", flattenIDExtensionsListIDs(resp.DeviceGroups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("devices", flattenIDExtensionsListIDs(resp.Devices)); err != nil {
		return diag.FromErr(err)
	}

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

	if err := d.Set("time_windows", flattenIDExtensionsListIDs(resp.TimeWindows)); err != nil {
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

func resourceFileTypeControlRulesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] file type control rule ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating file type control rule ID: %v\n", id)
	req := expandFileTypeControlRules(d)

	if _, err := filetypecontrol.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}

	existingRules, err := filetypecontrol.GetAll(ctx, service)
	if err != nil {
		log.Printf("[ERROR] error getting all file type control rules: %v", err)
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

	_, err = filetypecontrol.Update(ctx, service, id, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	// Fail immediately if INVALID_INPUT_ARGUMENT is detected
	if customErr := failFastOnErrorCodes(err); customErr != nil {
		return diag.Errorf("%v", customErr)
	}

	if err != nil {
		if strings.Contains(err.Error(), "INVALID_INPUT_ARGUMENT") {
			log.Printf("[INFO] Updating file type control rule ID: %v, got INVALID_INPUT_ARGUMENT\n", id)
		}
		return diag.FromErr(fmt.Errorf("error updating resource: %s", err))
	}

	reorderWithBeforeReorder(OrderRule{Order: intendedOrder, Rank: intendedRank}, req.ID, "file_type_control_rules",
		func() (int, error) {
			allRules, err := filetypecontrol.GetAll(ctx, service)
			if err != nil {
				return 0, err
			}
			// Count all rules including predefined ones for proper ordering
			return len(allRules), nil
		},
		func(id int, order OrderRule) error {
			rule, err := filetypecontrol.Get(ctx, service, id)
			if err != nil {
				return err
			}
			// to avoid the STALE_CONFIGURATION_ERROR
			rule.LastModifiedTime = 0
			rule.LastModifiedBy = nil
			rule.Order = order.Order
			rule.Rank = order.Rank
			_, err = filetypecontrol.Update(ctx, service, id, rule)
			return err
		},
		nil, // Remove beforeReorder function to avoid adding too many rules to the map
	)

	if diags := resourceFileTypeControlRulesRead(ctx, d, meta); diags.HasError() {
		return diags
	}
	markOrderRuleAsDone(req.ID, "file_type_control_rules")
	waitForReorder("file_type_control_rules")

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

func resourceFileTypeControlRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "rule_id")
	if !ok {
		log.Printf("[ERROR] file type control rule not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting file type control rule ID: %v\n", (d.Id()))

	if _, err := filetypecontrol.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] file type control rule deleted")

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

func expandFileTypeControlRules(d *schema.ResourceData) filetypecontrol.FileTypeRules {
	id, _ := getIntFromResourceData(d, "rule_id")

	// Retrieve the order and fallback to 1 if it's 0
	order := d.Get("order").(int)
	if order == 0 {
		log.Printf("[WARN] expandFileTypeControlRules: Rule ID %d has order=0. Falling back to order=1", id)
		order = 1
	}

	result := filetypecontrol.FileTypeRules{
		ID:                   id,
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Order:                order,
		Rank:                 d.Get("rank").(int),
		State:                d.Get("state").(string),
		FilteringAction:      d.Get("filtering_action").(string),
		Operation:            d.Get("operation").(string),
		ActiveContent:        d.Get("active_content").(bool),
		Unscannable:          d.Get("unscannable").(bool),
		BrowserEunTemplateID: d.Get("browser_eun_template_id").(int),
		CapturePCAP:          d.Get("capture_pcap").(bool),
		MinSize:              d.Get("min_size").(int),
		MaxSize:              d.Get("max_size").(int),
		DeviceTrustLevels:    SetToStringList(d, "device_trust_levels"),
		Protocols:            SetToStringList(d, "protocols"),
		FileTypes:            SetToStringList(d, "file_types"),
		URLCategories:        SetToStringList(d, "url_categories"),
		CloudApplications:    SetToStringList(d, "cloud_applications"),
		DeviceGroups:         expandIDNameExtensionsSet(d, "device_groups"),
		Devices:              expandIDNameExtensionsSet(d, "devices"),
		Locations:            expandIDNameExtensionsSet(d, "locations"),
		LocationGroups:       expandIDNameExtensionsSet(d, "location_groups"),
		Groups:               expandIDNameExtensionsSet(d, "groups"),
		Departments:          expandIDNameExtensionsSet(d, "departments"),
		Users:                expandIDNameExtensionsSet(d, "users"),
		TimeWindows:          expandIDNameExtensionsSet(d, "time_windows"),
		Labels:               expandIDNameExtensionsSet(d, "labels"),
		ZPAAppSegments:       expandZPAAppSegmentSet(d, "zpa_app_segments"),
	}
	return result
}

func currentFileTypeControlOrderVsRankWording(ctx context.Context, zClient *Client) string {
	service := zClient.Service

	list, err := filetypecontrol.GetAll(ctx, service)
	if err != nil {
		return ""
	}
	result := ""
	for i, r := range list {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("Rank %d VS Order %d", r.Rank, r.Order)

	}
	return result
}
