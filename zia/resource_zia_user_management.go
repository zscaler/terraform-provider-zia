package zia

import (
	"fmt"
	"log"

	"github.com/dome9/dome9-sdk-go/dome9/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func resourceUserManagement() *schema.Resource {
	return &schema.Resource{
		Create:   resourceUserManagementCreate,
		Read:     resourceUserManagementRead,
		Update:   resourceUserManagementUpdate,
		Delete:   resourceUserManagementDelete,
		Importer: &schema.ResourceImporter{},

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

func resourceUserManagementCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandApplicationSegmentRequest(d)
	log.Printf("[INFO] Creating application segment request\n%+v\n", req)
	if req.SegmentGroupID == "" {
		log.Println("[ERROR] Please provde a valid segment group for the application segment")
		return fmt.Errorf("please provde a valid segment group for the application segment")
	}
	resp, _, err := zClient.usermanagement.CreateUsers(req)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Created application segment request. ID: %v\n", resp.ID)
	d.SetId(resp.ID)

	return resourceUserManagementRead(d, m)
}

func resourceUserManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	resp, _, err := zClient.usermanagement.GetUsers(d.Id())

	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing user %s from state because it no longer exists in ZPA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Reading user: %+v\n", resp)
	_ = d.Set("name", resp.Name)
	_ = d.Set("email", resp.Email)
	_ = d.Set("comments", resp.Comments)
	_ = d.Set("temp_auth_email", resp.TempAuthEmail)
	_ = d.Set("password", resp.Password)
	_ = d.Set("adminUser", resp.AdminUser)
	_ = d.Set("type", resp.Type)
	_ = d.Set("id", resp.ID)
	_ = d.Set("groups", flattenGroups(resp.Groups))
	_ = d.Set("department", flattenDepartment(resp.Department))

	return nil
}

func resourceUserManagementUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id := d.Id()
	log.Printf("[INFO] Updating segment group ID: %v\n", id)
	req := expandUserManagement(d)

	if _, err := zClient.usermanagement.UpdateUsers(id, &req); err != nil {
		return err
	}

	return resourceUserManagementRead(d, m)
}

func resourceUserManagementDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	log.Printf("[INFO] Deleting user ID: %v\n", d.Id())

	if _, err := zClient.usermanagement.DeleteUsers(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	log.Printf("[INFO] user deleted")
	return nil
}

func expandApplicationSegmentRequest(d *schema.ResourceData) usermanagement.Users {
	return usermanagement.Users{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		Email:         d.Get("email").(string),
		Comments:      d.Get("comments").(string),
		TempAuthEmail: d.Get("temp_auth_email").(string),
		Password:      d.Get("password").(string),
		AdminUser:     d.Get("admin_user").(string),
		Type:          d.Get("type").(string),
		Groups:        expandGroups(d.Get("groups").([]interface{})),
		Department:    expandDepartment(d.Get("department").([]interface{})),
	}
}

func expandGroups(d *schema.ResourceData) usermanagement.Groups {
	Groups := usermanagement.Groups{
		ID:       d.Get("id").(string),
		Name:     d.Get("name").(string),
		IdpID:    d.Get("idp_id").(string),
		Comments: d.Get("comments").(string),
	}
	return Groups
}

func expandDepartment(department []interface{}) []usermanagement.Department {
	departments := make([]usermanagement.Department, len(department))

	for i, department := range department {
		departmentItem := department.(map[string]interface{})
		departments[i] = usermanagement.Department{
			ID:       departmentItem["id"].(string),
			Name:     departmentItem["name"].(string),
			IdpID:    departmentItem["idp_id"].(string),
			Comments: departmentItem["comments"].(string),
			Deleted:  departmentItem["deleted"].(bool),
		}

	}

	return departments
}
