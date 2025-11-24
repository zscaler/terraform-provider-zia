package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/users"
)

func dataSourceUserManagement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
							Type:     schema.TypeInt,
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
				Type:     schema.TypeSet,
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
						"idp_id": {
							Type:     schema.TypeInt,
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
			"auth_methods": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Accepted Authentication Methods",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"BASIC",
						"DIGEST",
					}, false),
				},
			},
			"is_auditor": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceUserManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Always fetch all users and search locally
	log.Printf("[INFO] Fetching all users\n")
	allUsers, err := users.GetAllUsers(ctx, service, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all users: %s", err))
	}

	log.Printf("[DEBUG] Retrieved %d users\n", len(allUsers))

	var resp *users.Users
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	// Search by ID first if provided
	if idProvided {
		log.Printf("[INFO] Searching for user by ID: %d\n", id)
		for _, user := range allUsers {
			if user.ID == id {
				resp = &user
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting user by ID %d: user not found", id))
		}
	}

	// Search by name if not found by ID and name is provided
	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for user by name: %s\n", nameStr)
		for _, user := range allUsers {
			if user.Name == nameStr {
				resp = &user
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting user by name %s: user not found", nameStr))
		}
	}

	// If neither ID nor name provided, or no match found
	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("couldn't find any user with name '%s' or id '%d'", nameStr, id))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	// Set the resource data
	d.SetId(fmt.Sprintf("%d", resp.ID))
	err = d.Set("name", resp.Name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	err = d.Set("email", resp.Email)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting email: %s", err))
	}
	err = d.Set("comments", resp.Comments)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting comments: %s", err))
	}
	err = d.Set("temp_auth_email", resp.TempAuthEmail)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting temp_auth_email: %s", err))
	}
	err = d.Set("admin_user", resp.AdminUser)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting admin_user: %s", err))
	}
	err = d.Set("type", resp.Type)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting type: %s", err))
	}
	err = d.Set("auth_methods", resp.AuthMethods)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting auth_methods: %s", err))
	}

	if err := d.Set("department", flattenUserDepartment(resp.Department)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", flattenUserGroups(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] User found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}
