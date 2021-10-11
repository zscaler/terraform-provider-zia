package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
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
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User name. This appears when choosing users for policies.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization.",
			},
			"groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Groups a user belongs to. Groups are used in policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Group name",
						},
						"idp_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Unique identfier for the identity provider (IdP)",
						},
						"comments": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Additional information about the group",
						},
					},
				},
			},
			"department": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Department name",
						},
						"idp_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Identity provider (IdP) ID",
						},
						"comments": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Additional information about this department",
						},
						"deleted": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"comments": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional information about this user.",
			},
			"temp_auth_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Temporary Authentication Email. If you enabled one-time tokens or links, enter the email address to which the Zscaler service sends the tokens or links. If this is empty, the service will send the email to the User email.",
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				// Sensitive: true,
				Description: "User's password. Applicable only when authentication type is Hosted DB. Password strength must follow what is defined in the auth settings.",
			},
		},
	}
}

func resourceUserManagementCreate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	req := expandUsers(d)
	log.Printf("[INFO] Creating zia user with request\n%+v\n", req)

	resp, err := zClient.usermanagement.Create(&req)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Created zia user request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("user_id", resp.ID)
	return resourceUserManagementRead(d, m)
}

func resourceUserManagementRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "user_id")
	if !ok {
		return fmt.Errorf("no users id is set")
	}
	resp, err := zClient.usermanagement.Get(id)
	if err != nil {
		if err.(*client.ErrorResponse).IsObjectNotFound() {
			log.Printf("[WARN] Removing user %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	log.Printf("[INFO] Getting user:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("user_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("email", resp.Email)
	_ = d.Set("comments", resp.Comments)
	_ = d.Set("temp_auth_email", resp.TempAuthEmail)
	_ = d.Set("password", resp.Password)

	if err := d.Set("groups", flattenGroupsSimple(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("department", flattenDepartmentsSimple(resp.Departments)); err != nil {
		return err
	}

	return nil
}

func resourceUserManagementUpdate(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "user_id")
	if !ok {
		log.Printf("[ERROR] user ID not set: %v\n", id)
	}

	log.Printf("[INFO] Updating users ID: %v\n", id)
	req := expandUsers(d)

	if _, _, err := zClient.usermanagement.Update(id, &req); err != nil {
		return err
	}

	return resourceUserManagementRead(d, m)
}

func resourceUserManagementDelete(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	id, ok := getIntFromResourceData(d, "user_id")
	if !ok {
		log.Printf("[ERROR] user ID not set: %v\n", id)
	}

	log.Printf("[INFO] Deleting user ID: %v\n", id)

	if _, err := zClient.usermanagement.Delete(id); err != nil {
		return err
	}

	d.SetId("")
	log.Printf("[INFO] user deleted")
	return nil
}

func expandUsers(d *schema.ResourceData) usermanagement.User {
	id, _ := getIntFromResourceData(d, "user_id")
	result := usermanagement.User{
		ID:            id,
		Name:          d.Get("name").(string),
		Email:         d.Get("email").(string),
		Comments:      d.Get("comments").(string),
		TempAuthEmail: d.Get("temp_auth_email").(string),
		Password:      d.Get("password").(string),
		Groups:        expandGroups(d),
		Departments:   expandDepartments(d),
	}
	groups := expandGroups(d)
	if groups != nil {
		result.Groups = groups
	}
	departments := expandDepartments(d)
	if departments != nil {
		result.Departments = departments
	}
	return result
}

func expandGroups(d *schema.ResourceData) []usermanagement.Groups {
	var groups []usermanagement.Groups
	if groupsInterface, ok := d.GetOk("groups"); ok {
		group := groupsInterface.([]interface{})
		groups = make([]usermanagement.Groups, len(group))
		for i, val := range group {
			groupItem := val.(map[string]interface{})
			groups[i] = usermanagement.Groups{
				ID:       groupItem["id"].(int),
				Name:     groupItem["name"].(string),
				IdpID:    groupItem["idp_id"].(int),
				Comments: groupItem["comments"].(string),
			}
		}
	}

	return groups
}

func expandDepartments(d *schema.ResourceData) *usermanagement.Departments {
	departmentObj, ok := d.GetOk("department")
	if !ok {
		return nil
	}
	departments, ok := departmentObj.(*schema.Set)
	if !ok {
		return nil
	}
	if len(departments.List()) > 0 {
		departmentObj := departments.List()[0]
		department, ok := departmentObj.(map[string]interface{})
		if !ok {
			return nil
		}
		return &usermanagement.Departments{
			ID:       department["id"].(int),
			Name:     department["name"].(string),
			IdpID:    department["idp_id"].(int),
			Comments: department["comments"].(string),
			Deleted:  department["deleted"].(bool),
		}
	}
	return nil
}

func flattenGroupsSimple(group []usermanagement.Groups) []interface{} {
	groups := make([]interface{}, len(group))
	for i, val := range group {
		groups[i] = map[string]interface{}{
			"id":       val.ID,
			"name":     val.Name,
			"idp_id":   val.IdpID,
			"comments": val.Comments,
		}
	}

	return groups
}

func flattenDepartmentsSimple(departments *usermanagement.Departments) interface{} {
	return []map[string]interface{}{
		{
			"id":       departments.ID,
			"name":     departments.Name,
			"idp_id":   departments.IdpID,
			"comments": departments.Comments,
			"deleted":  departments.Deleted,
		},
	}
}
