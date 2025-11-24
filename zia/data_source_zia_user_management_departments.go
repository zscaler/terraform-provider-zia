package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/departments"
)

func dataSourceDepartmentManagement() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDepartmentManagementRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
	}
}

func dataSourceDepartmentManagementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	// Always fetch all departments and search locally
	log.Printf("[INFO] Fetching all departments\n")
	allDepartments, err := departments.GetAll(ctx, service, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting all departments: %s", err))
	}

	log.Printf("[DEBUG] Retrieved %d departments\n", len(allDepartments))

	var resp *departments.Department
	id, idProvided := getIntFromResourceData(d, "id")
	nameObj, nameProvided := d.GetOk("name")
	nameStr := ""
	if nameProvided {
		nameStr = nameObj.(string)
	}

	// Search by ID first if provided
	if idProvided {
		log.Printf("[INFO] Searching for department by ID: %d\n", id)
		for _, dept := range allDepartments {
			if dept.ID == id {
				resp = &dept
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting department by ID %d: department not found", id))
		}
	}

	// Search by name if not found by ID and name is provided
	if resp == nil && nameProvided && nameStr != "" {
		log.Printf("[INFO] Searching for department by name: %s\n", nameStr)
		for _, dept := range allDepartments {
			if dept.Name == nameStr {
				resp = &dept
				break
			}
		}
		if resp == nil {
			return diag.FromErr(fmt.Errorf("error getting department by name %s: department not found", nameStr))
		}
	}

	// If neither ID nor name provided, or no match found
	if resp == nil {
		if idProvided || (nameProvided && nameStr != "") {
			return diag.FromErr(fmt.Errorf("couldn't find any department with name '%s' or id '%d'", nameStr, id))
		}
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	// Set the resource data
	d.SetId(fmt.Sprintf("%d", resp.ID))
	err = d.Set("name", resp.Name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	err = d.Set("idp_id", resp.IdpID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting idp_id: %s", err))
	}
	err = d.Set("comments", resp.Comments)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting comments: %s", err))
	}
	err = d.Set("deleted", resp.Deleted)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting deleted: %s", err))
	}

	log.Printf("[DEBUG] Department found: ID=%d, Name=%s\n", resp.ID, resp.Name)
	return nil
}
