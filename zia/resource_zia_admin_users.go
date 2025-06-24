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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
)

func resourceAdminUsers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdminUsersCreate,
		ReadContext:   resourceAdminUsersRead,
		UpdateContext: resourceAdminUsersUpdate,
		DeleteContext: resourceAdminUsersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("admin_id", idInt)
				} else {
					resp, err := admins.GetAdminUsersByLoginName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("admin_id", resp.ID)
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
			"admin_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Admin or auditor's username.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Admin or auditor's email address.",
			},
			"role": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Role of the admin. This is not required for an auditor.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the admin or auditor.",
			},
			"admin_scope_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ORGANIZATION",
					"DEPARTMENT",
					"LOCATION",
					"LOCATION_GROUP",
				}, false),
			},
			"admin_scope_entities": setIDsSchemaTypeCustom(nil, "list of destination ip groups"),
			"is_non_editable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_auditor": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password": {
				Type:         schema.TypeString,
				Description:  "The admin's password. If admin single sign-on (SSO) is disabled, then this field is mandatory for POST requests. This information is not provided in a GET response.",
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(8, 100),
			},
			"is_password_login_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_security_report_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_service_update_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_product_update_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_password_expired": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_exec_mobile_app_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAdminUsersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandAdminUsers(d)
	log.Printf("[INFO] Creating zia admin user with request\n%+v\n", req)
	if err := checkPasswordAllowed(req); err != nil {
		return diag.FromErr(err)
	}
	if err := checkAdminScopeType(req); err != nil {
		return diag.FromErr(err)
	}
	resp, err := admins.CreateAdminUser(ctx, service, req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created zia admin user request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("admin_id", resp.ID)

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

	return resourceAdminUsersRead(ctx, d, meta)
}

func checkPasswordAllowed(pass admins.AdminUsers) error {
	if pass.IsPasswordLoginAllowed && pass.Password == "" {
		return fmt.Errorf("enter a password for the admin. It can be 8 to 100 characters and must contain at least one number, one special character, and one upper-case letter")
	}
	return nil
}

func checkAdminScopeType(scopeType admins.AdminUsers) error {
	if scopeType.IsExecMobileAppEnabled && scopeType.AdminScopeType != "ORGANIZATION" {
		return fmt.Errorf("mobile app access can only be enabled for an admin with organization scope")
	}
	return nil
}

func resourceAdminUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "admin_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no admin users id is set"))
	}
	resp, err := admins.GetAdminUsers(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing admin user %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting admin user:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("admin_id", resp.ID)
	_ = d.Set("login_name", resp.LoginName)
	_ = d.Set("username", resp.UserName)
	_ = d.Set("email", resp.Email)
	_ = d.Set("comments", resp.Comments)
	_ = d.Set("is_non_editable", resp.IsNonEditable)
	_ = d.Set("disabled", resp.Disabled)
	_ = d.Set("is_auditor", resp.IsAuditor)
	_ = d.Set("admin_scope_type", resp.AdminScopeType)
	_ = d.Set("is_password_login_allowed", resp.IsPasswordLoginAllowed)
	_ = d.Set("is_security_report_comm_enabled", resp.IsSecurityReportCommEnabled)
	_ = d.Set("is_service_update_comm_enabled", resp.IsServiceUpdateCommEnabled)
	_ = d.Set("is_product_update_comm_enabled", resp.IsProductUpdateCommEnabled)
	_ = d.Set("is_password_expired", resp.IsPasswordExpired)

	if err := d.Set("role", flattenAdminUserRoleSimple(resp.Role)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("admin_scope_entities", flattenIDs(resp.AdminScopeEntities)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAdminUsersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "admin_id")
	if !ok {
		log.Printf("[ERROR] admin user ID not set: %v\n", id)
		return diag.FromErr(fmt.Errorf("admin user ID not set"))
	}

	log.Printf("[DEBUG] Updating admin user with ID: %d", id)

	req := expandAdminUsers(d)
	log.Printf("[DEBUG] Update request data: %+v", req)

	if _, err := admins.GetAdminUsers(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[INFO] Admin user %d not found. Removing from state", id)
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Error retrieving admin user before UpdateContext: %s", err)
		return diag.FromErr(err)
	}

	if _, err := admins.UpdateAdminUser(ctx, service, id, req); err != nil {
		log.Printf("[ERROR] Error updating admin user: %s", err)
		return diag.FromErr(err)
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

	return resourceAdminUsersRead(ctx, d, meta)
}

func resourceAdminUsersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "admin_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("cannot delete the resource admin users, no id found"))
	}

	log.Printf("[INFO] Deleting admin user ID: %v\n", id)

	if _, err := admins.DeleteAdminUser(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] admin user deleted")

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

func expandAdminUsers(d *schema.ResourceData) admins.AdminUsers {
	id, _ := getIntFromResourceData(d, "admin_id")
	result := admins.AdminUsers{
		ID:                          id,
		LoginName:                   d.Get("login_name").(string),
		UserName:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Comments:                    d.Get("comments").(string),
		IsNonEditable:               d.Get("is_non_editable").(bool),
		Disabled:                    d.Get("disabled").(bool),
		IsAuditor:                   d.Get("is_auditor").(bool),
		Password:                    d.Get("password").(string),
		AdminScopeType:              d.Get("admin_scope_type").(string),
		IsPasswordLoginAllowed:      d.Get("is_password_login_allowed").(bool),
		IsSecurityReportCommEnabled: d.Get("is_security_report_comm_enabled").(bool),
		IsServiceUpdateCommEnabled:  d.Get("is_service_update_comm_enabled").(bool),
		IsProductUpdateCommEnabled:  d.Get("is_product_update_comm_enabled").(bool),
		IsPasswordExpired:           d.Get("is_password_expired").(bool),
		IsExecMobileAppEnabled:      d.Get("is_exec_mobile_app_enabled").(bool),
		AdminScopeEntities:          expandIDNameExtensionsSet(d, "admin_scope_entities"),
	}
	role := expandAdminUserRoles(d)
	if role != nil {
		result.Role = role
	}
	return result
}

func flattenAdminUserRoleSimple(role *admins.Role) []interface{} {
	if role == nil {
		return []interface{}{}
	}
	roleMap := make(map[string]interface{})
	roleMap["id"] = role.ID

	return []interface{}{roleMap}
}

func expandAdminUserRoles(d *schema.ResourceData) *admins.Role {
	if v, ok := d.GetOk("role"); ok {
		roles := v.(*schema.Set).List()
		if len(roles) > 0 {
			roleMap := roles[0].(map[string]interface{})
			return &admins.Role{
				ID: roleMap["id"].(int),
			}
		}
	}
	return nil
}
