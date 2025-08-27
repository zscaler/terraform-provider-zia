package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/c2c_incident_receiver"
)

func dataSourceDLPCloudToCloudIR() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudToCloudIRRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "System-generated unique ID of the Cloud-to-Cloud DLP Incident Forwarding tenant",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "User-provided name for the SaaS application tenant",
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The status of the tenant",
			},
			"modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the tenant was last modified",
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the user who last modified the tenant",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier for the user",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the user",
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "External identifier for the user",
						},
						"extensions": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Additional properties for the user",
						},
					},
				},
			},
			"last_tenant_validation_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp when the tenant was last validated",
			},

			"last_validation_msg": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Last validation message information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error message from validation",
						},
						"error_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Error code from validation",
						},
					},
				},
			},
			"onboardable_entity": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the onboardable entity",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique identifier for the onboardable entity",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the onboardable entity",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the onboardable entity",
						},
						"enterprise_tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise tenant ID",
						},
						"application": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application name",
						},
						"last_validation_msg": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Last validation message for the onboardable entity",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error message from validation",
									},
									"error_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Error code from validation",
									},
								},
							},
						},
						"tenant_authorization_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tenant authorization information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_token": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Access token for authorization",
									},
									"bot_token": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bot token for authorization",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Redirect URL for authorization",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of authorization",
									},
									"env": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Environment type",
									},
									"temp_auth_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Temporary authorization code",
									},
									"subdomain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subdomain for the tenant",
									},
									"apicp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API CP value",
									},
									"client_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client ID for authorization",
									},
									"client_secret": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client secret for authorization",
									},
									"secret_token": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Secret token for authorization",
									},
									"user_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username for authorization",
									},
									"user_pwd": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "User password for authorization",
									},
									"instance_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance URL",
									},
									"role_arn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Role ARN for AWS integration",
									},
									"quarantine_bucket_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Quarantine bucket name",
									},
									"cloud_trail_bucket_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud trail bucket name",
									},
									"bot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bot ID",
									},
									"org_api_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Organization API key",
									},
									"external_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "External ID",
									},
									"enterprise_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enterprise ID",
									},
									"cred_json": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Credentials JSON",
									},
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Role for the tenant",
									},
									"organization_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Organization ID",
									},
									"workspace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace name",
									},
									"workspace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace ID",
									},
									"qtn_channel_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Quarantine channel URL",
									},
									"features_supported": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "List of supported features",
									},
									"mal_qtn_lib_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Malware quarantine library name",
									},
									"dlp_qtn_lib_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DLP quarantine library name",
									},
									"credentials": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Credentials information",
									},
									"token_endpoint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Token endpoint URL",
									},
									"rest_api_endpoint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "REST API endpoint URL",
									},
									"smir_bucket_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SMIR bucket configuration",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Unique identifier for the bucket config",
												},
												"config_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Configuration name for the bucket",
												},
												"metadata_bucket_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metadata bucket name URL",
												},
												"data_bucket_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Data bucket name URL",
												},
											},
										},
									},
									"qtn_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Quarantine information",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"admin_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Admin ID",
												},
												"qtn_folder_path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Quarantine folder path",
												},
												"mod_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Modification time",
												},
											},
										},
									},
									"qtn_info_cleared": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether quarantine info is cleared",
									},
								},
							},
						},
						"zscaler_app_tenant_id": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Zscaler app tenant ID information",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unique identifier for the Zscaler app tenant",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the Zscaler app tenant",
									},
									"external_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "External identifier for the Zscaler app tenant",
									},
									"extensions": {
										Type:        schema.TypeMap,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Additional properties for the Zscaler app tenant",
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

func dataSourceCloudToCloudIRRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *c2c_incident_receiver.C2CIncidentReceiver
	var err error

	// Check if ID is provided
	if id, ok := d.GetOk("id"); ok {
		resp, err = c2c_incident_receiver.Get(ctx, service, id.(int))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to get cloud-to-cloud IR with ID %d: %v", id, err))
		}
	} else if name, ok := d.GetOk("name"); ok {
		resp, err = c2c_incident_receiver.GetC2CIRName(ctx, service, name.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to get cloud-to-cloud IR with name %s: %v", name, err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	log.Printf("[INFO] Retrieved cloud-to-cloud IR: %+v", resp)

	// Check if response is nil
	if resp == nil {
		return diag.FromErr(fmt.Errorf("received nil response from cloud-to-cloud IR API"))
	}

	// Debug logging
	log.Printf("[DEBUG] Response ID: %d", resp.ID)
	log.Printf("[DEBUG] Response Name: %s", resp.Name)
	log.Printf("[DEBUG] LastValidationMsg: %+v", resp.LastValidationMsg)
	log.Printf("[DEBUG] OnboardableEntity: %+v", resp.OnboardableEntity)

	// Set the ID
	d.SetId(fmt.Sprintf("%d", resp.ID))

	// Set basic attributes
	_ = d.Set("id", resp.ID)
	_ = d.Set("name", resp.Name)
	_ = d.Set("status", resp.Status)
	_ = d.Set("modified_time", resp.ModifiedTime)
	_ = d.Set("last_tenant_validation_time", resp.LastTenantValidationTime)

	// Set last_modified_by with safety check
	if resp.LastModifiedBy != nil {
		if err := d.Set("last_modified_by", flattenIDExtensionsList(resp.LastModifiedBy)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set last_validation_msg
	if resp.LastValidationMsg != nil && (resp.LastValidationMsg.ErrorMsg != "" || resp.LastValidationMsg.ErrorCode != "") {
		if err := d.Set("last_validation_msg", flattenLastValidationMsg(resp.LastValidationMsg)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set onboardable_entity
	if resp.OnboardableEntity != nil {
		if err := d.Set("onboardable_entity", flattenOnboardableEntity(resp.OnboardableEntity)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func flattenLastValidationMsg(msg *c2c_incident_receiver.LastValidationMsg) []map[string]interface{} {
	if msg == nil || (msg.ErrorMsg == "" && msg.ErrorCode == "") {
		return nil
	}

	return []map[string]interface{}{
		{
			"error_msg":  msg.ErrorMsg,
			"error_code": msg.ErrorCode,
		},
	}
}

func flattenOnboardableEntity(entity *c2c_incident_receiver.OnboardableEntity) []map[string]interface{} {
	if entity == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":                   entity.ID,
		"name":                 entity.Name,
		"type":                 entity.Type,
		"enterprise_tenant_id": entity.EnterpriseTenantID,
		"application":          entity.Application,
	}

	// Set last_validation_msg
	if entity.LastValidationMsg.ErrorMsg != "" || entity.LastValidationMsg.ErrorCode != "" {
		result["last_validation_msg"] = flattenLastValidationMsg(&entity.LastValidationMsg)
	}

	// Set tenant_authorization_info
	result["tenant_authorization_info"] = flattenTenantAuthorizationInfo(&entity.TenantAuthorizationInfo)

	// Set zscaler_app_tenant_id
	if entity.ZscalerAppTenantID != nil {
		result["zscaler_app_tenant_id"] = flattenIDExtensionsList(entity.ZscalerAppTenantID)
	}

	return []map[string]interface{}{result}
}

func flattenTenantAuthorizationInfo(info *c2c_incident_receiver.TenantAuthorizationInfo) []map[string]interface{} {
	if info == nil {
		return nil
	}

	result := map[string]interface{}{
		"access_token":            info.AccessToken,
		"bot_token":               info.BotToken,
		"redirect_url":            info.RedirectUrl,
		"type":                    info.Type,
		"env":                     info.Env,
		"temp_auth_code":          info.TempAuthCode,
		"subdomain":               info.Subdomain,
		"apicp":                   info.Apicp,
		"client_id":               info.ClientID,
		"client_secret":           info.ClientSecret,
		"secret_token":            info.SecretToken,
		"user_name":               info.UserName,
		"user_pwd":                info.UserPwd,
		"instance_url":            info.InstanceUrl,
		"role_arn":                info.RoleArn,
		"quarantine_bucket_name":  info.QuarantineBucketName,
		"cloud_trail_bucket_name": info.CloudTrailBucketName,
		"bot_id":                  info.BotID,
		"org_api_key":             info.OrgApiKey,
		"external_id":             info.ExternalID,
		"enterprise_id":           info.EnterpriseID,
		"cred_json":               info.CredJson,
		"role":                    info.Role,
		"organization_id":         info.OrganizationID,
		"workspace_name":          info.WorkspaceName,
		"workspace_id":            info.WorkspaceID,
		"qtn_channel_url":         info.QtnChannelUrl,
		"features_supported":      info.FeaturesSupported,
		"mal_qtn_lib_name":        info.MalQtnLibName,
		"dlp_qtn_lib_name":        info.DlpQtnLibName,
		"credentials":             info.Credentials,
		"token_endpoint":          info.TokenEndpoint,
		"rest_api_endpoint":       info.RestApiEndpoint,
		"qtn_info_cleared":        info.QtnInfoCleared,
	}

	// Set smir_bucket_config
	if len(info.SmirBucketConfig) > 0 {
		result["smir_bucket_config"] = flattenSmirBucketConfig(info.SmirBucketConfig)
	}

	// Set qtn_info
	if len(info.QtnInfo) > 0 {
		result["qtn_info"] = flattenQtnInfo(info.QtnInfo)
	}

	return []map[string]interface{}{result}
}

func flattenSmirBucketConfig(configs []c2c_incident_receiver.SmirBucketConfig) []map[string]interface{} {
	var result []map[string]interface{}
	for _, config := range configs {
		result = append(result, map[string]interface{}{
			"id":                   config.ID,
			"config_name":          config.ConfigName,
			"metadata_bucket_name": config.MetadataBucketName,
			"data_bucket_name":     config.DataBucketName,
		})
	}
	return result
}

func flattenQtnInfo(qtnInfo []interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, info := range qtnInfo {
		if infoMap, ok := info.(map[string]interface{}); ok {
			result = append(result, map[string]interface{}{
				"admin_id":        infoMap["adminId"],
				"qtn_folder_path": infoMap["qtnFolderPath"],
				"mod_time":        infoMap["modTime"],
			})
		}
	}
	return result
}
