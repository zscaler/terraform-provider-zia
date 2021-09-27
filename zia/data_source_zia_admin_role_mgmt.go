package zia

/*
import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminrolemgmt"
)

func dataSourceAdminRoleMgmt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAdminRoleMgmtRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
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
				Default:  false,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_auditor": {
				Type:     schema.TypeBool,
				Computed: true,
				Default:  false,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"is_password_login_allowed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_security_report_comm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Default:  false,
			},
			"is_service_update_comm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Default:  false,
			},
			"is_product_update_comm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Default:  false,
			},
			"is_password_expired": {
				Type:     schema.TypeBool,
				Computed: true,
				Default:  false,
			},
			"is_exec_mobile_app_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Default:  false,
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
						"Scope_entities": {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceAdminRoleMgmtRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*gozscaler.Client)

	var resp *adminrolemgmt.AdminUsers
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for admin role management  %s\n", id)
		res, _, err := zClient.adminrolemgmt.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("login_name", resp.LoginName)
		_ = d.Set("user_name", resp.UserName)
		_ = d.Set("email", resp.Email)
		_ = d.Set("comments", resp.Comments)
		_ = d.Set("is_non_editable", resp.IsNonEditable)
		_ = d.Set("disabled", resp.Disabled)
		_ = d.Set("is_auditor", resp.IsAuditor)
		_ = d.Set("password", resp.Password)
		_ = d.Set("is_password_login_allowed", resp.IsPasswordLoginAllowed)
		_ = d.Set("is_security_report_comm_enabled", resp.IsSecurityReportCommEnabled)
		_ = d.Set("is_service_update_comm_enabled", resp.IsServiceUpdateCommEnabled)
		_ = d.Set("is_product_update_comm_enabled	", resp.IsProductUpdateCommEnabled)
		_ = d.Set("is_password_expired", resp.IsPasswordExpired)
		_ = d.Set("is_exec_mobile_app_enabled", resp.IsExecMobileAppEnabled)

		if err := d.Set("exec_mobile_app_tokens", flattenExecMobileAppTokens(resp)); err != nil {
			return fmt.Errorf("failed to read mobile app tokens %s", err)
		}
	} else {
		return fmt.Errorf("couldn't find any app connector group with name '%s' or id '%s'", name, id)
	}

	return nil
}

func flattenExecMobileAppTokens(mobileAppTokens *adminrolemgmt.AdminUsers) []interface{} {
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
*/
