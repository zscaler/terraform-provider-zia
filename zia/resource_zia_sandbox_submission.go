package zia

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_submission"
)

func resourceSandboxSubmission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSandboxSubmissionCreate,
		ReadContext:   resourceSandboxSubmissionRead,
		UpdateContext: resourceSandboxSubmissionUpdate,
		DeleteContext: resourceSandboxSubmissionDelete,

		Schema: map[string]*schema.Schema{
			"file_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"submission_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"submit",
					"discan",
				}, false),
			},
			"code": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"md5": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sandbox_submission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virus_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virus_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSandboxSubmissionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	filePath := d.Get("file_path").(string)
	force := d.Get("force").(bool)
	submissionMethod := d.Get("submission_method").(string)

	// Validation: If submission method is "discan", the "force" attribute should not be set
	if submissionMethod == "discan" && force {
		return diag.FromErr(fmt.Errorf("'force' attribute is not applicable for 'discan' submission method"))
	}

	file, err := os.Open(filePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to open file: %s", err))
	}
	defer file.Close()

	var result *sandbox_submission.ScanResult
	if submissionMethod == "submit" {
		forceStr := boolToString(force)
		result, err = sandbox_submission.SubmitFile(ctx, service, filePath, file, forceStr)
	} else if submissionMethod == "discan" {
		result, err = sandbox_submission.Discan(ctx, service, filePath, file)
	} else {
		return diag.FromErr(fmt.Errorf("invalid submission method: %s", submissionMethod))
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error submitting file to Sandbox: %s", err))
	}

	// Set Terraform resource attributes based on the response
	d.SetId(result.Md5)
	d.Set("code", result.Code)
	d.Set("message", result.Message)
	d.Set("file_type", result.FileType)
	d.Set("sandbox_submission", result.SandboxSubmission)
	d.Set("virus_name", result.VirusName)
	d.Set("virus_type", result.VirusType)
	d.Set("md5", result.Md5)

	return nil
}

func resourceSandboxSubmissionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Only POST methods are available, we can't fetch data again

	return nil
}

func resourceSandboxSubmissionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Trigger a re-creation of the resource if either 'file_path' or 'force' attributes change
	if d.HasChange("file_path") || d.HasChange("force") {
		// If there's a change, re-submit the file by calling the Create function
		return resourceSandboxSubmissionCreate(ctx, d, meta)
	}
	return nil
}

func resourceSandboxSubmissionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Since there is no DELETE method for this API, simply remove it from state
	d.SetId("")
	return nil
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
