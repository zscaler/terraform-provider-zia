package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
)

func dataSourceAdminRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAdminRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rank": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dashboard_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"report_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"analysis_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin_acct_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_auditor": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_non_editable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"logs_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAdminRolesRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *adminuserrolemgmt.AdminRoles
	idObj, idSet := d.GetOk("id")
	id, idIsInt := idObj.(int)
	if idSet && idIsInt && id > 0 {
		log.Printf("[INFO] Getting data for admin role id: %d\n", id)
		res, err := zClient.adminuserrolemgmt.GetAdminRoles(id)
		if err != nil {
			return err
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for admin role name: %s\n", name)
		res, err := zClient.adminuserrolemgmt.GetAdminRolesByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("rank", resp.Rank)
		_ = d.Set("name", resp.Name)
		_ = d.Set("policy_access", resp.PolicyAccess)
		_ = d.Set("dashboard_access", resp.DashboardAccess)
		_ = d.Set("report_access", resp.ReportAccess)
		_ = d.Set("is_auditor", resp.IsAuditor)
		_ = d.Set("analysis_access", resp.AnalysisAccess)
		_ = d.Set("username_access", resp.UsernameAccess)
		_ = d.Set("admin_acct_access", resp.AdminAcctAccess)
		_ = d.Set("is_auditor", resp.IsAuditor)
		_ = d.Set("permissions", resp.Permissions)
		_ = d.Set("is_non_editable", resp.IsNonEditable)
		_ = d.Set("logs_limit", resp.LogsLimit)
		_ = d.Set("role_type", resp.RoleType)

	} else {
		return fmt.Errorf("couldn't find any admin role name '%s' or id '%d'", name, id)
	}

	return nil
}
