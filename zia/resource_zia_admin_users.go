package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
)

func resourceAdminUsers() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAdminUsersCreate,
		Read:     resourceAdminUsersRead,
		Update:   resourceAdminUsersUpdate,
		Delete:   resourceAdminUsersDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
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
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Role of the admin. This is not required for an auditor.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_group_member_entities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ORGANIZATION",
								"DEPARTMENT",
								"LOCATION",
								"LOCATION_GROUP",
							}, false),
						},
						"scope_entities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"is_non_editable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_auditor": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password": {
				Type:        schema.TypeBool,
				Description: "The admin's password. If admin single sign-on (SSO) is disabled, then this field is mandatory for POST requests. This information is not provided in a GET response.",
				Optional:    true,
				Sensitive:   true,
			},
			"is_password_login_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_security_report_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_service_update_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_product_update_comm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_password_expired": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_exec_mobile_app_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAdminUsersCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandAdminUsers(d)
	log.Printf("[INFO] Creating zia admin user with request\n%+v\n", req)

	resp, err := zClient.adminuserrolemgmt.CreateAdminUser(req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia admin user request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))

	return resourceAdminUsersRead(d, m)
}

func resourceAdminUsersRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if !(idSet && idIsInt && id > 0) {
		return fmt.Errorf("no admin users id is set")
	}
	resp, err := zClient.adminuserrolemgmt.GetAdminUsers(id)
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing admin user %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting admin user:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("login_name", resp.LoginName)
	_ = d.Set("user_name", resp.UserName)
	_ = d.Set("email", resp.Email)
	_ = d.Set("comments", resp.Comments)
	_ = d.Set("is_non_editable", resp.IsNonEditable)
	_ = d.Set("disabled", resp.Disabled)
	_ = d.Set("is_auditor", resp.IsAuditor)
	_ = d.Set("is_password_login_allowed", resp.IsPasswordLoginAllowed)
	_ = d.Set("pwd_last_modified_time", resp.PasswordLastModifiedTime)
	_ = d.Set("is_security_report_comm_enabled", resp.IsSecurityReportCommEnabled)
	_ = d.Set("is_service_update_comm_enabled", resp.IsServiceUpdateCommEnabled)
	_ = d.Set("is_product_update_comm_enabled", resp.IsProductUpdateCommEnabled)
	_ = d.Set("is_password_expired", resp.IsPasswordExpired)
	_ = d.Set("admin_scope_type", resp.AdminScopeType)

	return nil
}

func resourceAdminUsersUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating admin users ID: %v\n", id)
	req := expandAdminUsers(d)

	if _, err := zClient.adminuserrolemgmt.UpdateAdminUser(id, req); err != nil {
		return err
	}

	return resourceAdminUsersRead(d, m)
}

func resourceAdminUsersDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	// Need to pass the ID (int) of the resource for deletion
	log.Printf("[INFO] Deleting admin user ID: %v\n", (d.Id()))

	if _, err := zClient.adminuserrolemgmt.DeleteAdminUser(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] admin user deleted")
	return nil
}

// Review expand functions to ensure it is set up correctly
func expandAdminUsers(d *schema.ResourceData) adminuserrolemgmt.AdminUsers {
	return adminuserrolemgmt.AdminUsers{
		LoginName:                   d.Get("login_name").(string),
		UserName:                    d.Get("user_name").(string),
		Email:                       d.Get("email").(string),
		Comments:                    d.Get("comments").(string),
		IsNonEditable:               d.Get("is_non_editable").(bool),
		Disabled:                    d.Get("disabled").(bool),
		IsAuditor:                   d.Get("is_auditor").(bool),
		Password:                    d.Get("password").(string),
		IsPasswordLoginAllowed:      d.Get("is_password_login_allowed").(bool),
		IsSecurityReportCommEnabled: d.Get("is_security_report_comm_enabled").(bool),
		IsServiceUpdateCommEnabled:  d.Get("is_service_update_comm_enabled").(bool),
		IsProductUpdateCommEnabled:  d.Get("is_product_update_comm_enabled").(bool),
		IsPasswordExpired:           d.Get("is_password_expired").(bool),
		IsExecMobileAppEnabled:      d.Get("is_exec_mobile_app_enabled").(bool),
		Role:                        expandAdminUserRoles(d),
		AdminScope:                  expandAdminScope(d),
	}
}

func expandAdminUserRoles(d *schema.ResourceData) adminuserrolemgmt.Role {
	return adminuserrolemgmt.Role{
		ID:         d.Get("id").(int),
		Name:       d.Get("name").(string),
		Extensions: d.Get("extensions").(map[string]interface{}),
	}
}

func expandAdminScope(d *schema.ResourceData) adminuserrolemgmt.AdminScope {
	adminScope := adminuserrolemgmt.AdminScope{
		ScopeGroupMemberEntities: expandScopeGroupMemberEntities(d),
		ScopeEntities:            expandScopeEntities(d),
	}

	return adminScope
}

func expandScopeGroupMemberEntities(d *schema.ResourceData) []adminuserrolemgmt.ScopeGroupMemberEntities {
	var scopeGroupMemberEntities []adminuserrolemgmt.ScopeGroupMemberEntities
	if scopeGroupInterface, ok := d.GetOk("scope_group_member_entities"); ok {
		scopes := scopeGroupInterface.([]interface{})
		scopeGroupMemberEntities = make([]adminuserrolemgmt.ScopeGroupMemberEntities, len(scopes))
		for i, scope := range scopes {
			scopeItem := scope.(map[string]interface{})
			scopeGroupMemberEntities[i] = adminuserrolemgmt.ScopeGroupMemberEntities{
				ID:         scopeItem["id"].(int),
				Name:       scopeItem["name"].(string),
				Extensions: scopeItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return scopeGroupMemberEntities
}

func expandScopeEntities(d *schema.ResourceData) []adminuserrolemgmt.ScopeEntities {
	var scopeEntities []adminuserrolemgmt.ScopeEntities
	if scopeGroupInterface, ok := d.GetOk("scope_entities"); ok {
		scopes := scopeGroupInterface.([]interface{})
		scopeEntities = make([]adminuserrolemgmt.ScopeEntities, len(scopeEntities))
		for i, scope := range scopes {
			scopeItem := scope.(map[string]interface{})
			scopeEntities[i] = adminuserrolemgmt.ScopeEntities{
				ID:         scopeItem["id"].(int),
				Name:       scopeItem["name"].(string),
				Extensions: scopeItem["extensions"].(map[string]interface{}),
			}
		}
	}

	return scopeEntities
}
