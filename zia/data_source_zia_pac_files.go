package zia

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/pacfiles"
)

func dataSourcePacFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePacFilesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the PAC file to retrieve. When set, only the matching PAC file is returned.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the PAC file to retrieve. When set, only the matching PAC file is returned.",
			},
			"filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"pac_content"}, false),
				Description:  "When set to `pac_content`, the PAC file content is omitted from the results.",
			},
			"pac_files": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of PAC files in deployed state, including default and custom PAC files.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier for the PAC file.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the PAC file.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the PAC file.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain of your organization to which the PAC file applies.",
						},
						"pac_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL location of the PAC file, auto-generated when the PAC file is first added.",
						},
						"pac_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The content of the PAC file. Empty when the `filter` attribute is set to `pac_content`.",
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the PAC file is editable.",
						},
						"pac_sub_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The obfuscated URL of the PAC file. Returned when `pac_url_obfuscated` is true.",
						},
						"pac_url_obfuscated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the PAC file URL is obfuscated.",
						},
						"pac_verification_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The verification status of the PAC file, indicating whether any syntax errors were identified. Supported values: `VERIFY_NOERR`, `VERIFY_ERR`, `NOVERIFY`.",
						},
						"pac_version_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the PAC file version: deployed, staged for deployment, or marked as the last known good version. Supported values: `DEPLOYED`, `STAGE`, `LKG`.",
						},
						"pac_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number of the PAC file.",
						},
						"pac_commit_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The commit message entered when the PAC file version was saved.",
						},
						"total_hits": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of times the PAC file was used during the last 30 days.",
						},
						"last_modified_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp when the PAC file was last modified, in Unix time.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timestamp when the PAC file was created, in Unix time.",
						},
						"last_modified_by": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The admin that last modified the PAC file.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identifier that uniquely identifies an entity.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configured name of the entity.",
									},
									"extensions": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourcePacFilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	filter := ""
	if f, ok := d.GetOk("filter"); ok {
		filter = f.(string)
	}

	var pacFiles []pacfiles.PACFileConfig

	if id, ok := getIntFromResourceData(d, "id"); ok {
		allFiles, err := pacfiles.GetPacFiles(ctx, service, filter)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, pacFile := range allFiles {
			if pacFile.ID == id {
				pacFiles = append(pacFiles, pacFile)
				break
			}
		}
		if len(pacFiles) == 0 {
			return diag.FromErr(fmt.Errorf("no PAC file found with ID: %d", id))
		}
	} else if name, ok := d.GetOk("name"); ok {
		foundFile, err := pacfiles.GetPacFileByName(ctx, service, name.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		if filter == "pac_content" {
			// The name lookup always returns the full object; honour the
			// filter by omitting the content, matching the list behaviour.
			foundFile.PACContent = ""
		}
		pacFiles = append(pacFiles, *foundFile)
	} else {
		var err error
		pacFiles, err = pacfiles.GetPacFiles(ctx, service, filter)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(pacFiles) == 0 {
			return diag.FromErr(fmt.Errorf("no PAC files found"))
		}
	}

	var result []map[string]interface{}
	for _, pacFile := range pacFiles {
		pacData := map[string]interface{}{
			"id":                      pacFile.ID,
			"name":                    pacFile.Name,
			"description":             pacFile.Description,
			"domain":                  pacFile.Domain,
			"pac_url":                 pacFile.PACUrl,
			"pac_content":             pacFile.PACContent,
			"editable":                pacFile.Editable,
			"pac_sub_url":             pacFile.PACSubURL,
			"pac_url_obfuscated":      pacFile.PACUrlObfuscated,
			"pac_verification_status": pacFile.PACVerificationStatus,
			"pac_version_status":      pacFile.PACVersionStatus,
			"pac_version":             pacFile.PACVersion,
			"pac_commit_message":      pacFile.PACCommitMessage,
			"total_hits":              pacFile.TotalHits,
			"last_modified_time":      pacFile.LastModificationTime,
			"create_time":             pacFile.CreateTime,
		}

		if pacFile.LastModifiedBy != nil {
			pacData["last_modified_by"] = flattenLastModifiedBy(pacFile.LastModifiedBy)
		} else {
			pacData["last_modified_by"] = []interface{}{}
		}

		result = append(result, pacData)
	}

	if err := d.Set("pac_files", result); err != nil {
		return diag.FromErr(fmt.Errorf("error setting PAC files: %w", err))
	}

	// Single-file lookups adopt that file's ID; the all-files form uses a
	// stable synthetic ID. The internal ID must stay numeric in every path:
	// the top-level `id` attribute is an integer, and the SDK parses the
	// internal ID into it when refreshing state.
	if len(pacFiles) == 1 {
		d.SetId(strconv.Itoa(pacFiles[0].ID))
		_ = d.Set("id", pacFiles[0].ID)
		_ = d.Set("name", pacFiles[0].Name)
	} else {
		d.SetId(strconv.Itoa(schema.HashString("all_pac_files")))
	}

	return nil
}
