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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudapplications/risk_profiles"
)

func resourceRiskProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRiskProfilesCreate,
		ReadContext:   resourceRiskProfilesRead,
		UpdateContext: resourceRiskProfilesUpdate,
		DeleteContext: resourceRiskProfilesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				zClient := meta.(*Client)
				service := zClient.Service

				id := d.Id()
				idInt, parseIDErr := strconv.ParseInt(id, 10, 64)
				if parseIDErr == nil {
					_ = d.Set("profile_id", idInt)
				} else {
					resp, err := risk_profiles.GetByName(ctx, service, id)
					if err == nil {
						d.SetId(strconv.Itoa(resp.ID))
						_ = d.Set("profile_id", resp.ID)
					} else {
						return []*schema.ResourceData{d}, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the risk profile",
			},
			"profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier for the risk profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud application risk profile name",
			},
			"profile_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Risk profile type. The default profile type is CLOUD_APPLICATIONS",
				ValidateFunc: validation.StringInSlice([]string{
					"CLOUD_APPLICATIONS",
				}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the applications",
				ValidateFunc: validation.StringInSlice([]string{
					"UN_SANCTIONED",
					"SANCTIONED",
					"ANY",
				}, false),
			},
			"exclude_certificates": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Indicates if the certificates are included or not",
			},
			"poor_items_of_service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on the presence of questionable terms and conditions in their legal agreements, such as sharing customer data with third-party applications.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"admin_audit_logs": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for logging and tracking all administrative activities to identify potential security threats.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"data_breach": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their history of reported data breaches in the last three years.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"source_ip_restrictions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their ability to restrict access to specific IP addresses, reducing the attack surface.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"mfa_support": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for multi-factor authentication to enhance user account security.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"ssl_pinned": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their use of pinned SSL certificates, making it difficult for attackers to decrypt and validate traffic.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"http_security_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their implementation of all security headers (X-XSS-Protection, X-Frame-Options, Strict-Transport-Security, Content-Security-Policy, and X-Content-Type-Options) to protect against common web attacks. This field is not applicable to the Lite API.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"evasive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for anonymous access without requiring user authentication that can increase the risk of malicious activity",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"dns_caa_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their implementation of the DNS Certification Authority Authorization (CAA) policy that helps prevent unauthorized SSL certificates.",
			},
			"weak_cipher_support": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for weak ciphers with key sizes less than 128 bits that can compromise SSL connections.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"password_strength": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their password strength requirements under Hosting Info.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"GOOD",
					"POOR",
					"UN_KNOWN",
				}, false),
			},
			"ssl_cert_validity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on the validity period of their SSL certificates.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"vulnerability": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their published CVE vulnerabilities",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"malware_scanning_for_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for content malware scanning",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"file_sharing": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for file sharing features that can increase the risk of data exfiltration",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"ssl_cert_key_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on the key size of their SSL certificates",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"UN_KNOWN",
					"BITS_2048",
					"BITS_256",
					"BITS_3072",
					"BITS_384",
					"BITS_4096",
					"BITS_1024",
					"BITS_8192",
				}, false),
			},
			"vulnerable_to_heart_bleed": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their vulnerability to the Heartbleed attack",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"vulnerable_to_log_jam": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their vulnerability to the Logjam attack",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"vulnerable_to_poodle": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their vulnerability to the POODLE attack",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"vulnerability_disclosure": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their policy for disclosing known vulnerabilities, allowing ethical hackers to report potential security threats",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"support_for_waf": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for web application firewalls (WAFs) to protect against common web attacks.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"remote_screen_sharing": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for remote access screen sharing, which can increase the risk of data exfiltration if not properly secured.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"sender_policy_framework": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for sender policy framework (SPF) authentication, which helps prevent email spoofing.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"domain_keys_identified_mail": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for DomainKeys Identified Mail (DKIM) authentication, which helps prevent email tampering.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"domain_based_message_auth": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters applications based on their support for Domain-Based Message Authentication, Reporting, and Conformance (DMARC), which helps prevent email spoofing and phishing attacks.",
				ValidateFunc: validation.StringInSlice([]string{
					"ANY",
					"YES",
					"NO",
					"UN_KNOWN",
				}, false),
			},
			"custom_tags":                setIDExternalIDCustom(intPtr(255), "List of custom tags to be included or excluded for the profile."),
			"data_encryption_in_transit": getRiskProfileEncryptionInTransit(),
			"risk_index":                 getRiskProfileIndex(),
			"certifications":             getRiskProfileCertifications(),
		},
	}
}

func resourceRiskProfilesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient, ok := meta.(*Client)
	if !ok {
		return diag.Errorf("unexpected meta type: expected *Client, got %T", meta)
	}

	service := zClient.Service

	req := expandRiskProfiles(d)
	log.Printf("[INFO] Creating ZIA risk profiles\n%+v\n", req)

	resp, _, err := risk_profiles.Create(ctx, service, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ZIA risk profiles request. ID: %v\n", resp)
	d.SetId(strconv.Itoa(resp.ID))
	_ = d.Set("profile_id", resp.ID)

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceRiskProfilesRead(ctx, d, meta)
}

func resourceRiskProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "profile_id")
	if !ok {
		return diag.FromErr(fmt.Errorf("no risk profiles id is set"))
	}
	resp, err := risk_profiles.Get(ctx, service, id)
	if err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			log.Printf("[WARN] Removing zia risk profiles %s from state because it no longer exists in ZIA", d.Id())
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	log.Printf("[INFO] Getting zia risk profiles:\n%+v\n", resp)

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

	if err := d.Set("custom_tags", flattenCommonIDNameExternalIDSimple(resp.CustomTags)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRiskProfilesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "profile_id")
	if !ok {
		log.Printf("[ERROR] risk profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Updating zia risk profile ID: %v\n", id)
	req := expandRiskProfiles(d)
	if _, err := risk_profiles.Get(ctx, service, id); err != nil {
		if respErr, ok := err.(*errorx.ErrorResponse); ok && respErr.IsObjectNotFound() {
			d.SetId("")
			return nil
		}
	}
	if _, _, err := risk_profiles.Update(ctx, service, id, &req); err != nil {
		return diag.FromErr(err)
	}

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return resourceRiskProfilesRead(ctx, d, meta)
}

func resourceRiskProfilesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	id, ok := getIntFromResourceData(d, "profile_id")
	if !ok {
		log.Printf("[ERROR] risk profile ID not set: %v\n", id)
	}
	log.Printf("[INFO] Deleting zia risk profile ID: %v\n", (d.Id()))

	if _, err := risk_profiles.Delete(ctx, service, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	log.Printf("[INFO] zia risk profile deleted")

	// Sleep for 2 seconds before potentially triggering the activation
	time.Sleep(2 * time.Second)

	// Check if ZIA_ACTIVATION is set to a truthy value before triggering activation
	if shouldActivate() {
		if activationErr := triggerActivation(ctx, zClient); activationErr != nil {
			return diag.FromErr(activationErr)
		}
	} else {
		log.Printf("[INFO] Skipping configuration activation due to ZIA_ACTIVATION env var not being set to true.")
	}

	return nil
}

func expandRiskProfiles(d *schema.ResourceData) risk_profiles.RiskProfiles {
	id, _ := getIntFromResourceData(d, "profile_id")
	result := risk_profiles.RiskProfiles{
		ID: id,

		ProfileName:               d.Get("name").(string),
		ProfileType:               d.Get("profile_type").(string),
		Status:                    d.Get("status").(string),
		PoorItemsOfService:        d.Get("poor_items_of_service").(string),
		AdminAuditLogs:            d.Get("admin_audit_logs").(string),
		DataBreach:                d.Get("data_breach").(string),
		SourceIpRestrictions:      d.Get("source_ip_restrictions").(string),
		MfaSupport:                d.Get("mfa_support").(string),
		SslPinned:                 d.Get("ssl_pinned").(string),
		HttpSecurityHeaders:       d.Get("http_security_headers").(string),
		Evasive:                   d.Get("evasive").(string),
		DnsCaaPolicy:              d.Get("dns_caa_policy").(string),
		WeakCipherSupport:         d.Get("weak_cipher_support").(string),
		PasswordStrength:          d.Get("password_strength").(string),
		SslCertValidity:           d.Get("ssl_cert_validity").(string),
		Vulnerability:             d.Get("vulnerability").(string),
		MalwareScanningForContent: d.Get("malware_scanning_for_content").(string),
		FileSharing:               d.Get("file_sharing").(string),
		SslCertKeySize:            d.Get("ssl_cert_key_size").(string),
		VulnerableToHeartBleed:    d.Get("vulnerable_to_heart_bleed").(string),
		VulnerableToLogJam:        d.Get("vulnerable_to_log_jam").(string),
		VulnerableToPoodle:        d.Get("vulnerable_to_poodle").(string),
		VulnerabilityDisclosure:   d.Get("vulnerability_disclosure").(string),
		SupportForWaf:             d.Get("support_for_waf").(string),
		RemoteScreenSharing:       d.Get("remote_screen_sharing").(string),
		SenderPolicyFramework:     d.Get("sender_policy_framework").(string),
		DomainKeysIdentifiedMail:  d.Get("domain_keys_identified_mail").(string),
		DomainBasedMessageAuth:    d.Get("domain_based_message_auth").(string),
		ExcludeCertificates:       d.Get("exclude_certificates").(int),
		RiskIndex:                 SetToIntList(d, "risk_index"),
		Certifications:            SetToStringList(d, "certifications"),
		DataEncryptionInTransit:   SetToStringList(d, "data_encryption_in_transit"),
		CustomTags:                expandCommonIDNameExternalID(d, "custom_tags"),
	}
	return result
}
