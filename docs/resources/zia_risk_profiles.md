---
subcategory: "Cloud Application Risk Profile"
layout: "zscaler"
page_title: "ZIA: risk_profiles"
description: |-
  Official documentation https://help.zscaler.com/zia/about-cloud-application-risk-profile
  API documentation https://help.zscaler.com/zia/cloud-applications#/riskProfiles-get
  Adds a new cloud application risk profile
---

# zia_risk_profiles (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-cloud-application-risk-profile)
* [API documentation](https://help.zscaler.com/zia/cloud-applications#/riskProfiles-get)

Use the **zia_risk_profiles** resource allows the creation and management of cloud application risk profile in the Zscaler Internet Access cloud or via the API.
See [About Cloud Application Risk Profile](https://help.zscaler.com/zia/about-cloud-application-risk-profile) for more details.

## Example Usage - Create a Risk Profile

```hcl
resource "zia_risk_profiles" "this" {
    profile_name = "RiskProfile_12346"
    status="SANCTIONED"
    risk_index=[1, 2, 3, 4, 5]
    certifications=["AICPA", "CCPA", "CISP"]
    password_strength="GOOD"
    poor_items_of_service="YES"
    admin_audit_logs="YES"
    data_breach="YES"
    source_ip_restrictions="YES"
    file_sharing="YES"
    mfa_support="YES"
    ssl_pinned="YES"
    data_encryption_in_transit=[
        "SSLV2", "SSLV3", "TLSV1_0", "TLSV1_1", "TLSV1_2", "TLSV1_3", "UN_KNOWN"
    ]
    http_security_headers="YES"
    evasive="YES"
    dns_caa_policy="YES"
    ssl_cert_validity="YES"
    weak_cipher_support="YES"
    vulnerability="YES"
    vulnerable_to_heart_bleed="YES"
    ssl_cert_key_size="BITS_2048"
    vulnerable_to_poodle="YES"
    support_for_waf="YES"
    vulnerability_disclosure="YES"
    domain_keys_identified_mail="YES"
    malware_scanning_for_content="YES"
    domain_based_message_auth="YES"
    sender_policy_framework="YES"
    remote_screen_sharing="YES"
    vulnerable_to_log_jam="YES"
    profile_type="CLOUD_APPLICATIONS"
    custom_tags {
        id = [1, 2]
    }
}
```

## Argument Reference

The following arguments are supported:

* `profile_name` - (Required, String) Cloud application risk profile name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `profile_type` - (String) Risk profile type. Supported value: `CLOUD_APPLICATIONS`. Default is `CLOUD_APPLICATIONS`.
* `status` - (String) Status of the applications. Supported values: `UN_SANCTIONED`, `SANCTIONED`, `ANY`.
* `exclude_certificates` - (Int) Indicates if the certificates are included or not.
* `poor_items_of_service` - (String) Filters applications based on questionable legal terms. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `admin_audit_logs` - (String) Filters based on support for administrative logging. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `data_breach` - (String) Filters based on history of data breaches. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `source_ip_restrictions` - (String) Filters based on IP restriction support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `mfa_support` - (String) Filters based on multi-factor authentication support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `ssl_pinned` - (String) Filters based on use of pinned SSL certificates. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `http_security_headers` - (String) Filters based on HTTP security headers support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `evasive` - (String) Filters based on anonymous access support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `dns_caa_policy` - (String) Filters based on DNS CAA policy implementation.
* `weak_cipher_support` - (String) Filters based on weak cipher usage. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `password_strength` - (String) Filters based on password strength policy. Supported values: `ANY`, `GOOD`, `POOR`, `UN_KNOWN`.
* `ssl_cert_validity` - (String) Filters based on SSL certificate validity period. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `vulnerability` - (String) Filters based on published CVE vulnerabilities. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `malware_scanning_for_content` - (String) Filters based on content malware scanning. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `file_sharing` - (String) Filters based on file sharing capability. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `ssl_cert_key_size` - (String) Filters based on SSL certificate key size. Supported values: `ANY`, `UN_KNOWN`, `BITS_1024`, `BITS_2048`, `BITS_256`, `BITS_3072`, `BITS_384`, `BITS_4096`, `BITS_8192`.
* `vulnerable_to_heart_bleed` - (String) Filters based on Heartbleed vulnerability. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `vulnerable_to_log_jam` - (String) Filters based on Logjam vulnerability. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `vulnerable_to_poodle` - (String) Filters based on POODLE vulnerability. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `vulnerability_disclosure` - (String) Filters based on vulnerability disclosure policy. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `support_for_waf` - (String) Filters based on Web Application Firewall (WAF) support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `remote_screen_sharing` - (String) Filters based on remote screen sharing support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `sender_policy_framework` - (String) Filters based on SPF authentication support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `domain_keys_identified_mail` - (String) Filters based on DKIM authentication support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `domain_based_message_auth` - (String) Filters based on DMARC support. Supported values: `ANY`, `YES`, `NO`, `UN_KNOWN`.
* `custom_tags` - (Set) List of custom tags to be included or excluded for the profile.
* `data_encryption_in_transit` - (Optional) Filters based on encryption of data in transit.
* `risk_index` - (Optional) Filters based on risk index thresholds.
* `certifications` - (Optional) Filters based on supported certifications.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_risk_profiles** can be imported by using `<PROFILE_ID>` or `<PROFILE_NAME>` as the import ID.

For example:

```shell
terraform import zia_risk_profiles.example <profile_id>
```

or

```shell
terraform import zia_risk_profiles.example <profile_name>
```
