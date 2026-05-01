package zia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
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
				Optional: true,
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
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JMESPath expression to filter results client-side. Applied after pagination completes. Example: \"[?contains(name, 'Admin')]\"",
			},
		},
	}
}

func dataSourceUserManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	searchExprRaw, searchProvided := d.GetOk("search")
	searchExpr := ""
	if searchProvided {
		searchExpr = searchExprRaw.(string)
		ctx = zscaler.ContextWithJMESPath(ctx, searchExpr)
		log.Printf("[INFO] JMESPath filter set: %s\n", searchExpr)
	}

	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}
	emailObj, emailProvided := d.GetOk("email")
	emailStr := ""
	if emailProvided {
		emailStr = emailObj.(string)
	}

	if !idProvided && nameStr == "" && emailStr == "" {
		return diag.FromErr(fmt.Errorf("one of 'id', 'name', or 'email' must be provided"))
	}

	// Build server-side filter options.
	//
	// The /users endpoint exposes a `name` query parameter that performs a
	// partial match against the display name AND email fields. We use it to
	// narrow the API-side pool when looking up by `name` or `email`, EXCEPT
	// when a `search` (JMESPath) expression is also provided: in that case we
	// must fetch the unfiltered pool so the JMESPath can evaluate against the
	// full set of users — otherwise the expression is effectively scoped to
	// "users whose name/email contains <lookup>", which produces confusing
	// "user not found" errors when the JMESPath expression itself is
	// misconfigured.
	filterOpts := &users.GetAllUsersFilterOptions{}
	if !searchProvided {
		switch {
		case nameStr != "":
			filterOpts.Name = nameStr
		case emailStr != "":
			filterOpts.Name = emailStr
		}
	}

	log.Printf("[INFO] Fetching users (server-side filter: name=%q, jmespath=%q)\n", filterOpts.Name, searchExpr)
	allUsers, err := users.GetAllUsers(ctx, service, filterOpts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all users: %s", err))
	}

	log.Printf("[DEBUG] Retrieved %d users (after server-side filter and JMESPath)\n", len(allUsers))

	// Helper to surface a clearer error when the candidate pool is empty
	// because of an over-restrictive or malformed `search` expression. This
	// is the most common cause of "user not found" when `search` is used.
	notFoundf := func(key, value string) error {
		base := fmt.Sprintf("error getting user by %s %s: user not found", key, value)
		if searchProvided {
			return fmt.Errorf("%s. The `search` (JMESPath) expression %q filtered the candidate pool down to %d user(s) before the %s match. Verify the expression references valid camelCase fields (e.g. `department.name`, not `department.email`) and that the target user satisfies the predicate", base, searchExpr, len(allUsers), key)
		}
		return errors.New(base)
	}

	var resp *users.Users

	if idProvided {
		log.Printf("[INFO] Searching for user by ID: %d\n", id)
		for _, user := range allUsers {
			if user.ID == id {
				u := user
				resp = &u
				break
			}
		}
		if resp == nil {
			return diag.FromErr(notFoundf("ID", strconv.Itoa(id)))
		}
	}

	if resp == nil && emailStr != "" {
		log.Printf("[INFO] Searching for user by email: %s\n", emailStr)
		for _, user := range allUsers {
			if strings.EqualFold(user.Email, emailStr) {
				u := user
				resp = &u
				break
			}
		}
		if resp == nil {
			return diag.FromErr(notFoundf("email", emailStr))
		}
	}

	if resp == nil && nameStr != "" {
		log.Printf("[INFO] Searching for user by name: %s\n", nameStr)
		for _, user := range allUsers {
			if user.Name == nameStr {
				u := user
				resp = &u
				break
			}
		}
		if resp == nil {
			return diag.FromErr(notFoundf("name", nameStr))
		}
	}

	if resp == nil {
		return diag.FromErr(fmt.Errorf("couldn't find any user matching the provided id/name/email"))
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
