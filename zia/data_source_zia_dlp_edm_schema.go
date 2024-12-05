package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_exact_data_match"
)

func dataSourceDLPEDMSchema() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDLPEDMSchemaRead,
		Schema: map[string]*schema.Schema{
			"schema_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The identifier (1-65519) for the EDM schema (i.e., EDM template) that is unique within the organization.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The EDM schema (i.e., EDM template) name. This attribute is ignored by PUT requests, but required for POST requests.",
			},
			"revision": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The revision number of the CSV file upload to the Index Tool. This attribute is required by PUT requests.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated filename, excluding the extention",
			},
			"original_file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated filename, excluding the extention.",
			},
			"file_upload_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EDM template's CSV file upload to the Index Tool. This attribute is required by PUT and POST requests.",
			},
			"orig_col_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of actual columns selected from the CSV file. This attribute is required by PUT and POST requests.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The date and time the IDM template was last modified.",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the DLP policy rule last.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
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
			"created_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the DLP policy rule last.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
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
			"edm_client": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the DLP policy rule last.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifier that uniquely identifies an entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configured name of the entity",
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
			"cells_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of cells used by the EDM schema (i.e., EDM template).",
			},
			"schema_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates the status of a specified EDM schema (i.e., EDM template). If this value is set to true, the schema is active and can be used by DLP dictionaries.",
			},
			"schedule_present": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The total number of cells used by the EDM schema (i.e., EDM template).",
			},
			"token_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates the status of a specified EDM schema (i.e., EDM template). If this value is set to true, the schema is active and can be used by DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The token (i.e., criteria) name. This attribute is required by PUT and POST requests.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The token (i.e., criteria) name. This attribute is required by PUT and POST requests.",
						},
						"primary_key": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the token is a primary key.",
						},
						"original_column": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The column position for the token in the original CSV file uploaded to the Index Tool, starting from 1. This attribue required by PUT and POST requests.",
						},
						"hash_file_column_order": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The column position for the token in the hashed file, starting from 1.",
						},
						"col_length_bitmap": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The length of the column bitmap in the hashed file.",
						},
					},
				},
			},
			"schedule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates the status of a specified EDM schema (i.e., EDM template). If this value is set to true, the schema is active and can be used by DLP dictionaries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The schedule type for the IDM template's schedule (i.e., Monthly, Weekly, Daily, or None). This attribute is required by PUT and POST requests.",
						},
						"schedule_day_of_month": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The day of the month that the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to MONTHLY.",
						},
						"schedule_day_of_week": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The day of the week the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to WEEKLY.",
						},
						"schedule_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time of the day (in minutes) that the IDM template is scheduled for. For example: at 3am= 180 mins. This attribute is required by PUT and POST requests.",
						},
						"schedule_disabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to true, the schedule for the IDM template is temporarily in a disabled state. This attribute is required by PUT requests in order to disable or enable a schedule.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDLPEDMSchemaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *dlp_exact_data_match.DLPEDMSchema
	schemaID, ok := getIntFromResourceData(d, "profile_id")
	if ok {
		log.Printf("[INFO] Getting data for dlp edm schema id: %d\n", schemaID)
		res, err := dlp_exact_data_match.GetDLPEDMSchemaID(ctx, service, schemaID)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	projectName, _ := d.Get("project_name").(string)
	if resp == nil && projectName != "" {
		log.Printf("[INFO] Getting data for dlp edm schema name: %s\n", projectName)
		res, err := dlp_exact_data_match.GetDLPEDMByName(ctx, service, projectName)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.SchemaID))
		_ = d.Set("schema_id", resp.SchemaID)
		_ = d.Set("project_name", resp.ProjectName)
		_ = d.Set("revision", resp.Revision)
		_ = d.Set("file_name", resp.Filename)
		_ = d.Set("original_file_name", resp.OriginalFileName)
		_ = d.Set("file_upload_status", resp.FileUploadStatus)
		_ = d.Set("orig_col_count", resp.OrigColCount)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)
		_ = d.Set("cells_used", resp.CellsUsed)
		_ = d.Set("schema_active", resp.SchemaActive)
		_ = d.Set("schedule_present", resp.SchedulePresent)

		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.ModifiedBy)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("created_by", flattenIDExtensionsList(resp.CreatedBy)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("edm_client", flattenIDExtensionsList(resp.EDMClient)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("token_list", flattenEDMTokenList(resp)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("schedule", flattenEDMSchedule(&resp.Schedule)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any dlp edn schema name '%s' or id '%d'", projectName, schemaID))
	}

	return nil
}

func flattenEDMTokenList(edm *dlp_exact_data_match.DLPEDMSchema) []interface{} {
	tokenList := make([]interface{}, len(edm.TokenList))
	for i, val := range edm.TokenList {
		tokenList[i] = map[string]interface{}{
			"name":                   val.Name,
			"type":                   val.Type,
			"primary_key":            val.PrimaryKey,
			"original_column":        val.OriginalColumn,
			"hash_file_column_order": val.HashfileColumnOrder,
			"col_length_bitmap":      val.ColLengthBitmap,
		}
	}

	return tokenList
}

func flattenEDMSchedule(edmSchedule *dlp_exact_data_match.Schedule) []map[string]interface{} {
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["schedule_type"] = edmSchedule.ScheduleType
	result[0]["schedule_day_of_month"] = edmSchedule.ScheduleDayOfMonth
	result[0]["schedule_day_of_week"] = edmSchedule.ScheduleDayOfWeek
	result[0]["schedule_time"] = edmSchedule.ScheduleTime
	result[0]["schedule_disabled"] = edmSchedule.ScheduleDisabled
	return result
}
