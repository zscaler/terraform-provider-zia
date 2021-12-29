package zia

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func resourceUserManagement() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserManagementCreate,
		Read:   resourceUserManagementRead,
		Update: resourceUserManagementUpdate,
		Delete: resourceUserManagementDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				zClient := m.(*Client)

				id := d.Id()
				_, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					// assume if the passed value is an int
					d.Set("user_id", id)
				} else {
					resp, err := zClient.usermanagement.GetUserByName(id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						d.Set("user_id", resp.ID)
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
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name. This appears when choosing users for policies.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization.",
			},
			"department": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Additional information about this user.",
				ValidateFunc: validation.StringLenBetween(0, 10240),
			},
			"temp_auth_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Temporary Authentication Email. If you enabled one-time tokens or links, enter the email address to which the Zscaler service sends the tokens or links. If this is empty, the service will send the email to the User email.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "User's password. Applicable only when authentication type is Hosted DB. Password strength must follow what is defined in the auth settings.",
			},
			"groups": listIDsSchemaTypeCustom(8, "List of Groups a user belongs to. Groups are used in policies."),
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
		if respErr, ok := err.(*client.ErrorResponse); ok && respErr.IsObjectNotFound() {
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

	if err := d.Set("groups", flattenUserGroupSet(resp.Groups)); err != nil {
		return err
	}

	if err := d.Set("department", flattenUserDepartment(resp.Department)); err != nil {
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

func expandUsers(d *schema.ResourceData) usermanagement.Users {
	id, _ := getIntFromResourceData(d, "user_id")
	result := usermanagement.Users{
		ID:            id,
		Name:          d.Get("name").(string),
		Email:         d.Get("email").(string),
		Comments:      d.Get("comments").(string),
		TempAuthEmail: d.Get("temp_auth_email").(string),
		Password:      d.Get("password").(string),
	}
	groups := expandUserGroups(d, "groups")
	if groups != nil {
		result.Groups = groups
	}

	department := expandUserDepartment(d)
	if department != nil {
		result.Department = department
	}
	return result
}
