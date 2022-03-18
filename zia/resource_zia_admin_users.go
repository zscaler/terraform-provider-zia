package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/common"
)

func resourceAdminUsers() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminUsersCreate,
		Read:   resourceAdminUsersRead,
		Update: resourceAdminUsersUpdate,
		Delete: resourceAdminUsersDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("admin_id", id)
				} else {
					resp, err := zClient.adminuserrolemgmt.GetAdminUsersByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						d.Set("admin_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
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
				Computed:    true,
				Description: "Role of the admin. This is not required for an auditor.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_name_l10n_tag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extensions": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about the admin or auditor.",
			},
			"admin_scope": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_group_member_entities": listIDsSchemaType("list of scope group member IDs"),
						"type": {
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
						"scope_entities": listIDsSchemaType("list of scope IDs"),
					},
				},
			},
			"is_non_editable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_auditor": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},
			"is_security_report_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_service_update_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_product_update_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_password_expired": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_exec_mobile_app_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAdminUsersCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandAdminUsers(d)
	log.Printf("[INFO] Creating zia admin user with request\n%+v\n", req)
	if err := checkPasswordAllowed(req); err != nil {
		return err
	}
	if err := checkAdminScopeType(req); err != nil {
		return err
	}
	resp, err := zClient.adminuserrolemgmt.CreateAdminUser(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia admin user request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("admin_id", resp.ID)
	return resourceAdminUsersRead(d, m)
}

func checkPasswordAllowed(pass adminuserrolemgmt.AdminUsers) error {
	if pass.IsPasswordLoginAllowed && pass.Password == "" {
		return fmt.Errorf("enter a password for the admin. It can be 8 to 100 characters and must contain at least one number, one special character, and one upper-case letter")
	}
	return nil
}

func checkAdminScopeType(scopeType adminuserrolemgmt.AdminUsers) error {
	if scopeType.IsExecMobileAppEnabled && scopeType.AdminScopeType != "ORGANIZATION" {
		return fmt.Errorf("mobile app access can only be enabled for an admin with organization scope")
	}
	return nil
}

func resourceAdminUsersRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "admin_id")
	if !ok {
		return fmt.Errorf("no admin users id is set")
	}
	resp, err := zClient.adminuserrolemgmt.GetAdminUsers(id)
	if err != nil {
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing admin user %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
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
	_ = d.Set("is_password_login_allowed", resp.IsPasswordLoginAllowed)
	_ = d.Set("is_security_report_comm_enabled", resp.IsSecurityReportCommEnabled)
	_ = d.Set("is_service_update_comm_enabled", resp.IsServiceUpdateCommEnabled)
	_ = d.Set("is_product_update_comm_enabled", resp.IsProductUpdateCommEnabled)
	_ = d.Set("is_password_expired", resp.IsPasswordExpired)

	if err := d.Set("role", flattenAdminUserRole(resp.Role)); err != nil {
		return err
	}

	if err := d.Set("admin_scope", flattenAdminUsersScopesLite(resp)); err != nil {
		return err
	}

	return nil
}

func flattenAdminUsersScopesLite(resp *adminuserrolemgmt.AdminUsers) []interface{} {
	scope := make([]interface{}, 1)
	scope[0] = map[string]interface{}{
		"type":                        resp.AdminScopeType,
		"scope_group_member_entities": flattenIDs(resp.AdminScopeGroupMemberEntities),
		"scope_entities":              flattenIDs(resp.AdminScopeEntities),
	}
	return scope
}

func resourceAdminUsersUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	req := expandAdminUsers(d)
	log.Printf("[INFO] Updating admin users ID: %v\n", req.ID)
	if _, err := zClient.adminuserrolemgmt.UpdateAdminUser(req.ID, req); err != nil {
		return err
	}

	return resourceAdminUsersRead(d, m)
}

func resourceAdminUsersDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "admin_id")
	if !ok {
		return fmt.Errorf("cannot delete the resource admin users, no id found")
	}

	log.Printf("[INFO] Deleting admin user ID: %v\n", id)

	if _, err := zClient.adminuserrolemgmt.DeleteAdminUser(id); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] admin user deleted")
	return nil
}

func expandAdminUsers(d *schema.ResourceData) adminuserrolemgmt.AdminUsers {
	adminScopeType, AdminScopeGroupMemberEntities, AdminScopeEntities := expandAdminUsersScopes(d)
	result := adminuserrolemgmt.AdminUsers{
		ID:                            d.Get("admin_id").(int),
		LoginName:                     d.Get("login_name").(string),
		UserName:                      d.Get("username").(string),
		Email:                         d.Get("email").(string),
		Comments:                      d.Get("comments").(string),
		IsNonEditable:                 d.Get("is_non_editable").(bool),
		Disabled:                      d.Get("disabled").(bool),
		IsAuditor:                     d.Get("is_auditor").(bool),
		Password:                      d.Get("password").(string),
		IsPasswordLoginAllowed:        d.Get("is_password_login_allowed").(bool),
		IsSecurityReportCommEnabled:   d.Get("is_security_report_comm_enabled").(bool),
		IsServiceUpdateCommEnabled:    d.Get("is_service_update_comm_enabled").(bool),
		IsProductUpdateCommEnabled:    d.Get("is_product_update_comm_enabled").(bool),
		IsPasswordExpired:             d.Get("is_password_expired").(bool),
		IsExecMobileAppEnabled:        d.Get("is_exec_mobile_app_enabled").(bool),
		AdminScopeGroupMemberEntities: AdminScopeGroupMemberEntities,
		AdminScopeEntities:            AdminScopeEntities,
		AdminScopeType:                adminScopeType,
	}
	role := expandAdminUserRoles(d)
	if role != nil {
		result.Role = role
	}
	return result
}

func expandAdminUserRoles(d *schema.ResourceData) *adminuserrolemgmt.Role {
	rolesObj, ok := d.GetOk("role")
	if !ok {
		return nil
	}
	roles, ok := rolesObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(roles.List()) > 0 {
		rolesObj := roles.List()[0]
		role, ok := rolesObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &adminuserrolemgmt.Role{
			ID:         role["id"].(int),
			Name:       role["name"].(string),
			Extensions: role["extensions"].(map[string]interface{}),
		}
	}
	return nil
}
func expandAdminUsersScopes(d *schema.ResourceData) (string, []common.IDNameExtensions, []common.IDNameExtensions) {
	if scopeInterface, ok := d.GetOk("admin_scope"); ok {
		scopesSet, ok := scopeInterface.(*schema.Set)
		if !ok {
			return "", []common.IDNameExtensions{}, []common.IDNameExtensions{}
		}
		for _, val := range scopesSet.List() {
			scopeItem := val.(map[string]interface{})
			return scopeItem["type"].(string), expandIDNameExtensionsMap(scopeItem, "scope_group_member_entities"), expandIDNameExtensionsMap(scopeItem, "scope_entities")
		}
	}
	return "", []common.IDNameExtensions{}, []common.IDNameExtensions{}
}
