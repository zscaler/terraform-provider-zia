package zia

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlp_idmp_profile "github.com/willguibr/terraform-provider-zia/gozscaler/dlp_idm_profile"
)

func dataSourceDLPIDMProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDLPIDMProfileRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The identifier (1-64) for the IDM template (i.e., IDM profile) that is unique within the organization.",
			},
			"profile_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IDM template name, which is unique per Index Tool.",
			},
			"profile_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IDM template's description.",
			},
			"profile_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IDM template's type.",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully qualified domain name (FQDN) of the IDM template's host machine.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number of the IDM template's host machine",
			},
			"profile_dir_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IDM template's directory file path, where all files are present.",
			},
			"schedule_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The schedule type for the IDM template's schedule (i.e., Monthly, Weekly, Daily, or None). This attribute is required by PUT and POST requests.",
			},
			"schedule_day": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The day the IDM template is scheduled for. This attribute is required by PUT and POST requests.",
			},
			"schedule_day_of_month": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The day of the month that the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to MONTHLY.",
			},
			"schedule_day_of_week": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The day of the week the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to WEEKLY.",
			},
			"schedule_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number of the IDM template's host machine",
			},
			"schedule_disabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, the schedule for the IDM template is temporarily in a disabled state. This attribute is required by PUT requests in order to disable or enable a schedule.",
			},
			"upload_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the file uploaded to the Index Tool for the IDM template.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username to be used on the IDM template's host machine.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version number for the IDM template.",
			},
			"idm_client": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The unique identifer for the Index Tool that was used to create the IDM template. This attribute is required by POST requests, but ignored if provided in PUT requests.",
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
			"volume_of_documents": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total volume of all the documents associated to the IDM template.",
			},
			"num_documents": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of documents associated to the IDM template.",
			},
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The date and time the IDM template was last modified.",
			},
			"modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the IDM template last.",
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
	}
}

func dataSourceDLPIDMProfileRead(d *schema.ResourceData, m interface{}) error {
	zClient := m.(*Client)

	var resp *dlp_idmp_profile.DLPIDMProfile
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for dlp idm profile id: %d\n", id)
		res, err := zClient.dlp_idmp_profile.Get(id)
		if err != nil {
			return err
		}
		resp = res
	}

	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for dlp idm profile: %s\n", name)
		res, err := zClient.dlp_idmp_profile.GetByName(name)
		if err != nil {
			return err
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ProfileID))
		_ = d.Set("profile_id", resp.ProfileID)
		_ = d.Set("profile_name", resp.ProfileName)
		_ = d.Set("profile_desc", resp.ProfileDesc)
		_ = d.Set("profile_type", resp.ProfileType)
		_ = d.Set("host", resp.Host)
		_ = d.Set("port", resp.Port)
		_ = d.Set("profile_dir_path", resp.ProfileDirPath)
		_ = d.Set("schedule_type", resp.ScheduleType)
		_ = d.Set("Schedule_day", resp.ScheduleDay)
		_ = d.Set("schedule_day_of_month", resp.ScheduleDayOfMonth)
		_ = d.Set("schedule_day_of_week", resp.ScheduleDayOfWeek)
		_ = d.Set("schedule_time", resp.ScheduleTime)
		_ = d.Set("schedule_disabled", resp.ScheduleDisabled)
		_ = d.Set("upload_status", resp.UploadStatus)
		_ = d.Set("user_name", resp.UserName)
		_ = d.Set("version", resp.Version)
		_ = d.Set("volume_of_documents", resp.VolumeOfDocuments)
		_ = d.Set("num_documents", resp.NumDocuments)
		_ = d.Set("last_modified_time", resp.LastModifiedTime)

		if err := d.Set("idm_client", flattenIDExtensions(resp.IdmClient)); err != nil {
			return err
		}

		if err := d.Set("modified_by", flattenIDExtensions(resp.ModifiedBy)); err != nil {
			return err
		}

	} else {
		return fmt.Errorf("couldn't find any dlp idm profile with name '%s' or id '%d'", name, id)
	}

	return nil
}
