package zia

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudapplications/risk_profiles"
)

func dataSourceRiskProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRiskProfilesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Unique identifier for the risk profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Cloud application risk profile name.",
			},
			"profile_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Risk profile type. The default profile type is CLOUD_APPLICATIONS. This field is not applicable to the Lite API.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Risk profile status. The default profile type is CLOUD_APPLICATIONS. This field is not applicable to the Lite API.",
			},
			"exclude_certificates": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates if the certificates are included or not. This field is not applicable to the Lite API.",
			},
			"poor_items_of_service": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on the presence of questionable terms and conditions in their legal agreements, such as sharing customer data with third-party applications.",
			},
			"admin_audit_logs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for logging and tracking all administrative activities to identify potential security threats.",
			},
			"data_breach": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their history of reported data breaches in the last three years.",
			},
			"source_ip_restrictions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their ability to restrict access to specific IP addresses, reducing the attack surface.",
			},
			"mfa_support": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for multi-factor authentication to enhance user account security.",
			},
			"ssl_pinned": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their use of pinned SSL certificates, making it difficult for attackers to decrypt and validate traffic.",
			},
			"http_security_headers": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their implementation of all security headers (X-XSS-Protection, X-Frame-Options, Strict-Transport-Security, Content-Security-Policy, and X-Content-Type-Options) to protect against common web attacks. This field is not applicable to the Lite API.",
			},
			"evasive": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for anonymous access without requiring user authentication that can increase the risk of malicious activity.",
			},
			"dns_caa_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their implementation of the DNS Certification Authority Authorization (CAA) policy that helps prevent unauthorized SSL certificates.",
			},
			"weak_cipher_support": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for weak ciphers with key sizes less than 128 bits that can compromise SSL connections.",
			},
			"password_strength": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their password strength requirements under Hosting Info.",
			},
			"ssl_cert_validity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on the validity period of their SSL certificates.",
			},
			"vulnerability": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their published CVE vulnerabilities.",
			},
			"malware_scanning_for_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for content malware scanning.",
			},
			"file_sharing": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for file sharing features that can increase the risk of data exfiltration.",
			},
			"ssl_cert_key_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on the key size of their SSL certificates.",
			},
			"vulnerable_to_heart_bleed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their vulnerability to the Heartbleed attack.",
			},
			"vulnerable_to_log_jam": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their vulnerability to the Logjam attack.",
			},
			"vulnerable_to_poodle": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their vulnerability to the POODLE attack.",
			},
			"vulnerability_disclosure": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their policy for disclosing known vulnerabilities, allowing ethical hackers to report potential security threats.",
			},
			"support_for_waf": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for web application firewalls (WAFs) to protect against common web attacks.",
			},
			"remote_screen_sharing": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for remote access screen sharing, which can increase the risk of data exfiltration if not properly secured.",
			},
			"sender_policy_framework": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for sender policy framework (SPF) authentication, which helps prevent email spoofing.",
			},
			"domain_keys_identified_mail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for DomainKeys Identified Mail (DKIM) authentication, which helps prevent email tampering.",
			},
			"domain_based_message_auth": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filters applications based on their support for Domain-Based Message Authentication, Reporting, and Conformance (DMARC), which helps prevent email spoofing and phishing attacks.",
			},
			"last_mod_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp of when the profile was last modified.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Timestamp of when the profile was created.",
			},
			"certifications": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of certifications to be included or excluded for the profile.",
			},
			"data_encryption_in_transit": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filters applications based on their support for encrypting data in transit.",
			},
			"risk_index": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "The risk index number of the cloud applications. It represents the risk score assigned to each cloud application based on the risk attribute values.",
			},
			"custom_tags": {
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
					},
				},
			},
			"last_modified_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The admin that modified the rule label last. This is a read-only field. Ignored by PUT requests.",
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
		},
	}
}

func dataSourceRiskProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	var resp *risk_profiles.RiskProfiles
	id, ok := getIntFromResourceData(d, "id")
	if ok {
		log.Printf("[INFO] Getting data for risk profile id: %d\n", id)
		res, err := risk_profiles.Get(ctx, service, id)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}
	name, _ := d.Get("name").(string)
	if resp == nil && name != "" {
		log.Printf("[INFO] Getting data for risk profile name: %s\n", name)
		res, err := risk_profiles.GetByName(ctx, service, name)
		if err != nil {
			return diag.FromErr(err)
		}
		resp = res
	}

	if resp != nil {
		d.SetId(fmt.Sprintf("%d", resp.ID))
		_ = d.Set("name", resp.ProfileName)
		_ = d.Set("profile_type", resp.ProfileType)
		_ = d.Set("status", resp.Status)
		_ = d.Set("exclude_certificates", resp.ExcludeCertificates)
		_ = d.Set("poor_items_of_service", resp.PoorItemsOfService)
		_ = d.Set("admin_audit_logs", resp.AdminAuditLogs)
		_ = d.Set("data_breach", resp.DataBreach)
		_ = d.Set("source_ip_restrictions", resp.SourceIpRestrictions)
		_ = d.Set("mfa_support", resp.MfaSupport)
		_ = d.Set("ssl_pinned", resp.SslPinned)
		_ = d.Set("http_security_headers", resp.HttpSecurityHeaders)
		_ = d.Set("evasive", resp.Evasive)
		_ = d.Set("dns_caa_policy", resp.DnsCaaPolicy)
		_ = d.Set("weak_cipher_support", resp.WeakCipherSupport)
		_ = d.Set("password_strength", resp.PasswordStrength)
		_ = d.Set("ssl_cert_validity", resp.SslCertValidity)
		_ = d.Set("vulnerability", resp.Vulnerability)
		_ = d.Set("malware_scanning_for_content", resp.MalwareScanningForContent)
		_ = d.Set("file_sharing", resp.FileSharing)
		_ = d.Set("ssl_cert_key_size", resp.SslCertKeySize)
		_ = d.Set("vulnerable_to_heart_bleed", resp.VulnerableToHeartBleed)
		_ = d.Set("vulnerable_to_log_jam", resp.VulnerableToLogJam)
		_ = d.Set("vulnerable_to_poodle", resp.VulnerableToPoodle)
		_ = d.Set("vulnerability_disclosure", resp.VulnerabilityDisclosure)
		_ = d.Set("support_for_waf", resp.SupportForWaf)
		_ = d.Set("remote_screen_sharing", resp.RemoteScreenSharing)
		_ = d.Set("sender_policy_framework", resp.SenderPolicyFramework)
		_ = d.Set("domain_keys_identified_mail", resp.DomainKeysIdentifiedMail)
		_ = d.Set("domain_based_message_auth", resp.DomainBasedMessageAuth)
		_ = d.Set("last_mod_time", resp.LastModTime)
		_ = d.Set("create_time", resp.CreateTime)
		_ = d.Set("certifications", resp.Certifications)
		_ = d.Set("data_encryption_in_transit", resp.DataEncryptionInTransit)
		_ = d.Set("risk_index", resp.RiskIndex)

		if err := d.Set("custom_tags", flattenCommonIDNameExternalID(resp.CustomTags)); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("last_modified_by", flattenLastModifiedBy(resp.ModifiedBy)); err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("couldn't find any risk profile name '%s' or id '%d'", name, id))
	}

	return nil
}
