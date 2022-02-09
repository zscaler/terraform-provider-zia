package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
)

func dataSourceAdminUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAdminUsersRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_non_editable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_auditor": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_password_login_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"pwd_last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_security_report_comm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_service_update_comm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_product_update_comm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_password_expired": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_exec_mobile_app_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"role": {
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
							Computed: true,
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
			"exec_mobile_app_tokens": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"org_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token_expiry": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAdminUsersRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *adminuserrolemgmt.AdminUsers
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for location id: %d\n", id)
		res, err := zClient.adminuserrolemgmt.GetAdminUsers(id)
		if err != nil {
			return err
		}
		resp = res
	}
	loginName, _ := d.Get("login_name").(string)
	if resp == nil && loginName != "" {
		log.Printf("[INFO] Getting data for location name: %s\n", loginName)
		res, err := zClient.adminuserrolemgmt.GetAdminUsersByName(loginName)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
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
		_ = d.Set("is_exec_mobile_app_enabled", resp.IsExecMobileAppEnabled)

		if err := d.Set("role", flattenAdminUserRole(resp.Role)); err != nil {
			return fmt.Errorf("failed to read admin user role %s", err)
		}

		if err := d.Set("admin_scope", flattenAdminScope(resp)); err != nil {
			return fmt.Errorf("failed to read admin scope %s", err)
		}

		if err := d.Set("exec_mobile_app_tokens", flattenExecMobileAppTokens(resp)); err != nil {
			return fmt.Errorf("failed to read mobile app tokens %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any admin user login name '%s' or id '%d'", loginName, id)
	}

	return nil
}

func flattenAdminUserRole(role *adminuserrolemgmt.Role) interface{} {
	return []map[string]interface{}{
		{
			"id":         role.ID,
			"name":       role.Name,
			"extensions": role.Extensions,
		},
	}
}

func flattenAdminScope(scopes *adminuserrolemgmt.AdminUsers) []interface{} {
	scope := make([]interface{}, 1)
	scope[0] = map[string]interface{}{
		"type":                        scopes.AdminScopeType,
		"scope_group_member_entities": flattenIDNameExtensions(scopes.AdminScopeGroupMemberEntities),
		"scope_entities":              flattenIDNameExtensions(scopes.AdminScopeEntities),
	}
	return scope
}

func flattenExecMobileAppTokens(mobileAppTokens *adminuserrolemgmt.AdminUsers) []interface{} {
	execMobileAppTokens := make([]interface{}, len(mobileAppTokens.ExecMobileAppTokens))
	for i, execMobileApp := range mobileAppTokens.ExecMobileAppTokens {
		execMobileAppTokens[i] = map[string]interface{}{
			"cloud":        execMobileApp.Cloud,
			"org_id":       execMobileApp.OrgId,
			"name":         execMobileApp.Name,
			"token_id":     execMobileApp.TokenId,
			"token":        execMobileApp.Token,
			"token_expiry": execMobileApp.TokenExpiry,
			"create_time":  execMobileApp.CreateTime,
			"device_id":    execMobileApp.DeviceId,
			"device_name":  execMobileApp.DeviceName,
		}
	}

	return execMobileAppTokens
}
