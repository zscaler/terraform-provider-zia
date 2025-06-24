package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/users"
)

func resourceUserManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserManagementCreate,
		ReadContext:   resourceUserManagementRead,
		UpdateContext: resourceUserManagementUpdate,
		DeleteContext: resourceUserManagementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("user_id", idInt)
				} else {
					resp, err := users.GetUserByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("user_id", resp.ID)
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 127),
				Description:  "User name. This appears when choosing users for policies.",
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 127),
				Description:  "User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization.",
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
			"auth_methods": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Accepted Authentication Methods",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"BASIC",
					}, false),
				},
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "User's password. Applicable only when authentication type is Hosted DB. Password strength must follow what is defined in the auth settings.",
			},
			"groups": setIDsSchemaTypeCustom(nil, "List of Groups a user belongs to. Groups are used in policies."),
			"department": {
				Type:     schema.TypeSet,
				Required: true,
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
							Computed:    true,
							Description: "Department name",
						},
						"idp_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Identity provider (IdP) ID",
						},
						"comments": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Additional information about this department",
						},
						"deleted": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceUserManagementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	req := expandUsers(d)
	log.Printf("[INFO] Creating zia user with request\n%+v\n", req)

	resp, err := users.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	// Sleep for 5 seconds before triggering the activation
	time.Sleep(5 * time.Second)

	// Trigger activation after creating the rule label
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	log.Printf("[INFO] Created zia user request. ID: %v\n", resp)
	authMethods := SetToStringList(d, "auth_methods")
	if len(authMethods) > 0 {
		_, err = users.EnrollUser(ctx, service, resp.ID, users.EnrollUserRequest{
			AuthMethods: authMethods,
			Password:    resp.Password,
		})
		if err != nil {
			log.Printf("[ERROR] enrolling user failed: %v\n", err)
		}
	}
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("user_id", resp.ID)

	time.Sleep(5 * time.Second)

	// Trigger activation after creating the rule label
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}
	return resourceUserManagementRead(ctx, d, meta)
}

func resourceUserManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "user_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no users id is set"))
	}
	resp, err := users.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing user %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting user:\n%+v\n", resp)
	d.SetId(fmt.Sprintf("%d", resp.ID))
	_ = d.Set("user_id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("email", resp.Email)
	_ = d.Set("comments", resp.Comments)
	_ = d.Set("temp_auth_email", resp.TempAuthEmail)

	if err := d.Set("groups", flattenUserGroups(resp.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("department", flattenUserDepartment(resp.Department)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceUserManagementUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "user_id")
	if !ok {
		log.Printf("[ERROR] user ID not set: %v\n", id)
	}

	log.Printf("[INFO] Updating users ID: %v\n", id)
	req := expandUsers(d)
	if _, err := users.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := users.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Sleep for 5 seconds before triggering the activation
	time.Sleep(5 * time.Second)

	// Trigger activation after creating the rule label
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	authMethods := SetToStringList(d, "auth_methods")
	if (d.HasChange("password") || d.HasChange("auth_methods")) && len(authMethods) > 0 {
		_, err := users.EnrollUser(ctx, service, id, users.EnrollUserRequest{
			AuthMethods: authMethods,
			Password:    req.Password,
		})
		if err != nil {
			log.Printf("[ERROR] enrolling user failed: %v\n", err)
		}
	}

	// Sleep for 5 seconds before triggering the activation
	time.Sleep(5 * time.Second)

	// Trigger activation after creating the rule label
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceUserManagementRead(ctx, d, meta)
}

func resourceUserManagementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "user_id")
	if !ok {
		log.Printf("[ERROR] user ID not set: %v\n", id)
	}

	log.Printf("[INFO] Deleting user ID: %v\n", id)
	err := DetachRuleIDNameExtensions(
		zClient,
		id,
		"Users",
		func(r *filteringrules.FirewallFilteringRules) []common.IDNameExtensions {
			return r.Users
		},
		func(r *filteringrules.FirewallFilteringRules, ids []common.IDNameExtensions) {
			r.Users = ids
		},
	)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := users.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] user deleted")

	// Sleep for 5 seconds before triggering the activation
	time.Sleep(5 * time.Second)

	// Trigger activation after creating the rule label
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandUsers(d *schema.ResourceData) users.Users {
	id, _ := getIntFromResourceData(d, "user_id")
	result := users.Users{
		ID:            id,
		Name:          d.Get("name").(string),
		Email:         d.Get("email").(string),
		Comments:      d.Get("comments").(string),
		TempAuthEmail: d.Get("temp_auth_email").(string),
		Password:      d.Get("password").(string),
		Groups:        expandUserGroups(d, "groups"),
	}

	department := expandUserDepartment(d)
	if department != nil {
		result.Department = department
	}
	return result
}
