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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/roles"
)

func resourceAdminRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminRolesCreate,
		ReadContext:   resourceAdminRolesRead,
		UpdateContext: resourceAdminRolesUpdate,
		DeleteContext: resourceAdminRolesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("role_id", idInt)
				} else {
					resp, err := roles.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("role_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		CustomizeDiff: adminRolesCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the admin role.",
			},
			"rank": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Admin rank of this admin role. This is applicable only when admin rank is enabled in the advanced settings. Default value is 7 (the lowest rank). The assigned admin rank determines the roles or admin users this user can manage, and which rule orders this admin can access.",
			},
			"policy_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Policy access permission.",
			},
			"alerting_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Alerting access permission",
			},
			"dashboard_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Dashboard access permission.",
			},
			"report_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Report access permission.",
			},
			"analysis_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Insights logs access permission.",
			},
			"username_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Username access permission. When set to NONE, the username will be obfuscated.",
			},
			"device_info_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Device information access permission. When set to NONE, device information is obfuscated.",
			},
			"admin_acct_access": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Admin and role management access permission.",
			},
			"is_auditor": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether this is an auditor role.",
			},
			"feature_permissions": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Feature access permission. Indicates which features an admin role can access and if the admin has both read and write access, or read-only access.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ext_feature_permissions": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "External feature access permission.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_non_editable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether or not this admin user is editable/deletable.",
			},
			"logs_limit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Log range limit.",
				ValidateFunc: validation.StringInSlice([]string{
					"UNRESTRICTED",
					"MONTH_1",
					"MONTH_2",
					"MONTH_3",
					"MONTH_4",
					"MONTH_5",
					"MONTH_6",
				}, false),
			},
			"role_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The admin role type. ()This attribute is subject to change.)",
				ValidateFunc: validation.StringInSlice([]string{
					"ORG_ADMIN",
					"EXEC_INSIGHT",
					"EXEC_INSIGHT_AND_ORG_ADMIN",
					"SDWAN",
				}, false),
			},
			"report_time_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Time duration allocated to the report dashboard. ",
			},
			"permissions": getAdminRolePermissions(),
		},
	}
}

func resourceAdminRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandAdminRoles(d)
	log.Printf("[INFO] Creating ZIA admin roles\n%+v\n", req)

	resp, _, err := roles.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA admin roles request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("role_id", resp.ID)

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

	return resourceAdminRolesRead(ctx, d, meta)
}

func resourceAdminRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "role_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no admin role id is set"))
	}
	resp, err := roles.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia admin roles %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia admin roles:\n%+v\n", resp)

	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("rank", resp.Rank)
	_ = d.Set("name", resp.Name)
	_ = d.Set("policy_access", resp.PolicyAccess)
	_ = d.Set("alerting_access", resp.AlertingAccess)
	_ = d.Set("dashboard_access", resp.DashboardAccess)
	_ = d.Set("report_access", resp.ReportAccess)
	_ = d.Set("analysis_access", resp.AnalysisAccess)
	_ = d.Set("username_access", resp.UsernameAccess)
	_ = d.Set("device_info_access", resp.DeviceInfoAccess)
	_ = d.Set("admin_acct_access", resp.AdminAcctAccess)
	_ = d.Set("is_auditor", resp.IsAuditor)
	_ = d.Set("permissions", resp.Permissions)
	_ = d.Set("feature_permissions", resp.FeaturePermissions)
	_ = d.Set("ext_feature_permissions", resp.ExtFeaturePermissions)
	_ = d.Set("is_non_editable", resp.IsNonEditable)
	_ = d.Set("logs_limit", resp.LogsLimit)
	_ = d.Set("role_type", resp.RoleType)
	_ = d.Set("report_time_duration", resp.ReportTimeDuration)

	return nil
}

func resourceAdminRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "role_id")
	if !ok {
		log.Printf("[ERROR] admin role ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia admin roleID: %v\n", id)
	req := expandAdminRoles(d)
	if _, err := roles.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := roles.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
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

	return resourceAdminRolesRead(ctx, d, meta)
}

func resourceAdminRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "role_id")
	if !ok {
		log.Printf("[ERROR] role ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia role ID: %v\n", (d.Id()))

	if _, err := roles.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia role deleted")

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

func expandAdminRoles(d *schema.ResourceData) roles.AdminRoles {
	id, _ := getIntFromResourceData(d, "role_id")

	featurePermissions := make(map[string]interface{})
	if v, ok := d.GetOk("feature_permissions"); ok {
		for key, val := range v.(map[string]interface{}) {
			featurePermissions[key] = val
		}
	}

	extFeaturePermissions := make(map[string]interface{})
	if v, ok := d.GetOk("ext_feature_permissions"); ok {
		for key, val := range v.(map[string]interface{}) {
			extFeaturePermissions[key] = val
		}
	}

	return roles.AdminRoles{
		ID:                    id,
		Name:                  d.Get("name").(string),
		Rank:                  d.Get("rank").(int),
		PolicyAccess:          d.Get("policy_access").(string),
		AlertingAccess:        d.Get("alerting_access").(string),
		DashboardAccess:       d.Get("dashboard_access").(string),
		ReportAccess:          d.Get("report_access").(string),
		AnalysisAccess:        d.Get("analysis_access").(string),
		UsernameAccess:        d.Get("username_access").(string),
		DeviceInfoAccess:      d.Get("device_info_access").(string),
		AdminAcctAccess:       d.Get("admin_acct_access").(string),
		IsAuditor:             d.Get("is_auditor").(bool),
		IsNonEditable:         d.Get("is_non_editable").(bool),
		LogsLimit:             d.Get("logs_limit").(string),
		RoleType:              d.Get("role_type").(string),
		ReportTimeDuration:    d.Get("report_time_duration").(int),
		Permissions:           SetToStringList(d, "permissions"),
		FeaturePermissions:    featurePermissions,
		ExtFeaturePermissions: extFeaturePermissions,
	}
}

// CustomizeDiff function enforces your special validations at plan-time
func adminRolesCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// 1) Once the resource is created (has an ID), the role_type cannot be changed
	//    i.e. if there's an existing ID and role_type is changing => error
	if d.Id() != "" && d.HasChange("role_type") {
		old, new := d.GetChange("role_type")
		if old != new {
			return fmt.Errorf("the role_type cannot be changed once the resource is created. old=%v new=%v", old, new)
		}
	}

	// 2) If role_type == "SDWAN", policy_access must be "READ_WRITE"
	// 3) If role_type == "SDWAN", alerting_access must be "NONE"
	roleType := d.Get("role_type").(string)
	if roleType == "SDWAN" {
		if d.Get("policy_access").(string) != "READ_WRITE" {
			return fmt.Errorf("for SDWAN roles, policy_access must be READ_WRITE")
		}
		if d.Get("alerting_access").(string) != "NONE" {
			return fmt.Errorf("for SDWAN roles, alerting_access must be NONE")
		}
	}

	return nil
}
