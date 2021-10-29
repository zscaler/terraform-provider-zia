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
			"admin_id": {
				Type:     schema.TypeInt,
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
				Type:        schema.TypeSet,
				Optional:    true,
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_group_member_entities": {
							Type:     schema.TypeList,
							Computed: true,
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
										Optional: true,
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
						},
						"scope_entities": {
							Type:     schema.TypeList,
							Optional: true,
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
										Optional: true,
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
				Type:        schema.TypeString,
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
			"admin_scope_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ORGANIZATION",
					"DEPARTMENT",
					"LOCATION",
					"LOCATION_GROUP",
				}, false),
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
	_ = d.Set("admin_id", resp.ID)
	return resourceAdminUsersRead(d, m)
}

func resourceAdminUsersRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)
	id, ok := getIntFromResourceData(d, "admin_id")
	if !ok {
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
	_ = d.Set("admin_id", resp.ID)
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

	if err := d.Set("role", flattenAdminUserRole(resp.Role)); err != nil {
		return err
	}

	return nil
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
	id, _ := getIntFromResourceData(d, "admin_id")
	result := adminuserrolemgmt.AdminUsers{
		ID:                          id,
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
