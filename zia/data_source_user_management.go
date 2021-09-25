package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func dataSourceUserManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"department": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comments": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"temp_auth_email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_auditor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"admin_user": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *usermanagement.Users
	id, ok := d.Get("id").(string)
	if ok && id != "" {
		log.Printf("[INFO] Getting data for admin role management  %s\n", id)
		res, _, err := zClient.usermanagement.GetUsers(id)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(resp.ID)
		_ = d.Set("name", resp.Name)
		_ = d.Set("email", resp.Email)
		_ = d.Set("comments", resp.Comments)
		_ = d.Set("temp_auth_email", resp.TempAuthEmail)
		_ = d.Set("password", resp.Password)
		_ = d.Set("adminUser", resp.AdminUser)
		_ = d.Set("type", resp.Type)

		if err := d.Set("groups", flattenGroups(resp.Groups)); err != nil {
			return err
		}

		if err := d.Set("department", flattenDepartment(resp.Department)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("couldn't find any user with name '%s' or id '%s'", id)
	}

	return nil
}

func flattenGroups(groups []usermanagement.Groups) []interface{} {
	group := make([]interface{}, len(groups))
	for i, val := range groups {
		group[i] = map[string]interface{}{
			"id":       val.ID,
			"name":     val.Name,
			"idp_id":   val.IdpID,
			"comments": val.Comments,
		}
	}

	return group
}

func flattenDepartment(department []usermanagement.Department) []interface{} {
	departments := make([]interface{}, len(department))
	for i, val := range department {
		departments[i] = map[string]interface{}{
			"id":       val.ID,
			"name":     val.Name,
			"idp_id":   val.IdpID,
			"comments": val.Comments,
			"deleted":  val.Deleted,
		}
	}

	return departments
}
