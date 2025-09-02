---
layout: "zscaler"
page_title: "Release Notes"
description: |-
  The Zscaler Internet Access (ZIA) provider Release Notes
---

# ZIA Provider: Release Notes

## USAGE

Track all ZIA Terraform provider's releases. New resources, features, and bug fixes will be tracked here.

---
``Last updated: v4.4.11``

---

## 4.4.11 (September, 2 2025)

### Notes

- Release date: **(September, 2 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #471](https://github.com/zscaler/terraform-provider-zia/pull/471) - Fixed Firewall Rules Description Heredoc handling

## 4.4.10 (August, 27 2025)

### Notes

- Release date: **(August, 27 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #470](https://github.com/zscaler/terraform-provider-zia/pull/470) - Fixed resource attribute `zpa_gateway` in `zia_forwarding_control_rule` due to missing `name` attribute

## 4.4.9 (August, 26 2025)

### Notes

- Release date: **(August, 26 2025)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #468](https://github.com/zscaler/terraform-provider-zia/pull/468) - Added data source `zia_cloud_to_cloud_ir` - Retrieves the Cloud-to-Cloud Incident Receiver (C2CIR) information configured in the ZIA Admin Portal. This data source can be used to set the corresponding receiver when configuring the resource `zia_dlp_web_rules` or `zia_casb_dlp_rules`

- [PR #468](https://github.com/zscaler/terraform-provider-zia/pull/468) - Added attribute `receiver` to `zia_dlp_web_rules` and `zia_casb_dlp_rules` resources to allow configuration of Cloud-to-Cloud Incident Receivers.

### Bug Fixes

- [PR #468](https://github.com/zscaler/terraform-provider-zia/pull/468) - Added `val` attribute to `zia_url_categories` resource to enable consistent referencing of URL categories in DLP web rules and other resources
- [PR #468](https://github.com/zscaler/terraform-provider-zia/pull/468) - Fixed performance issue in firewall filtering rules reordering by removing unnecessary predefined rule processing that was causing excessive wait times

### Documentation

- [PR #468](https://github.com/zscaler/terraform-provider-zia/pull/468) - Updated documentation for `zia_url_categories` resource to include new `val` attribute
- [PR #468](https://github.com/zscaler/terraform-provider-zia/pull/468) - Updated documentation for `zia_forwarding_control_rule` to remove unsupported attributes `devices` and `device_groups`

## 4.4.8 (August, 22 2025)

### Notes

- Release date: **(August, 22 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #466](https://github.com/zscaler/terraform-provider-zia/pull/466) - Enhanced `zia_device_groups` data source to support retrieving all device groups when no name is specified, in addition to existing single device group lookup by name. Added a new list field to return all device groups for bulk operations while maintaining backward compatibility.

## 4.4.7 (August, 22 2025)

### Notes

- Release date: **(August, 22 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #465](https://github.com/zscaler/terraform-provider-zia/pull/465) - Fixed ZPA Gateway app segments drift issue caused by API bug where individual gateway retrieval returns all possible app segments instead of only associated ones. Updated resource and data source to use GetAll() endpoint with local filtering to ensure correct app segments are returned, preventing Terraform drift during plan operations.

## 4.4.6 (August, 19 2025)

### Notes

- Release date: **(August, 19 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #463](https://github.com/zscaler/terraform-provider-zia/pull/463) - Fixed import in the resource `zia_url_filtering_rules` to ensure correct `cbi_profile` import due to API limitation.
- [PR #463](https://github.com/zscaler/terraform-provider-zia/pull/463) - Set attribute `state` to Computed in `zia_location_management` to handle odd API behavior.

## 4.4.5 (August, 13 2025)

### Notes

- Release date: **(August, 6 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #458](https://github.com/zscaler/terraform-provider-zia/pull/458) - Fixed drift in `zia_dlp_dictionaries` on attribute `hierarchical_identifiers` and updated examples

## 4.4.4 (August, 6 2025)

### Notes

- Release date: **(August, 6 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #457](https://github.com/zscaler/terraform-provider-zia/pull/457) - Added new `dictionary_type` value `MIP_TAG` to resource `zia_dlp_dictionaries`


## 4.4.3 (July, 31 2025)

### Notes

- Release date: **(July, 31 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #456](https://github.com/zscaler/terraform-provider-zia/pull/456) - Removed validation for attribute `nw_applications` from the resource `zia_firewall_filtering_rule`. See respective documentations for each resource for further instructions.
- [PR #456](https://github.com/zscaler/terraform-provider-zia/pull/456) - Applied heredoc formatting to support non-standard multi-line text on `description` attribute across supported resources.

## 4.4.2 (July, 29 2025)

### Notes

- Release date: **(July, 29 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #455](https://github.com/zscaler/terraform-provider-zia/pull/455) - Fixed action validation in `zia_cloud_app_control_rule`.

## 4.4.1 (July, 24 2025)

### Notes

- Release date: **(July, 24 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #453](https://github.com/zscaler/terraform-provider-zia/pull/453) - Upgraded to [Zscaler-SDK-GO v3.5.6](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.5.6)

## 4.4.0 (July, 24 2025)

### Notes

- Release date: **(July, 24 2025)**
- Supported Terraform version: **v1.x**

### NEW - RESOURCES

The following new resources have been introduced:

- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and resource``zia_cloud_nss_feed`` - Adds a new cloud NSS feed
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added the resource ``zia_bandwidth_classes`` - Bandwidth Classes
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and resource ``zia_bandwidth_control_rule`` - Bandwidth Control Rules
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and resource ``zia_bandwidth_classes_file_size`` - Bandwidth Classes File Size
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and resource``zia_bandwidth_classes_web_conferencing`` - Bandwidth Classes Web Conferencing

### NEW - DATA SOURCES

The following new data sources have been introduced:

- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and datasource``zia_cloud_nss_feed`` - Retrieves the cloud NSS feeds configured in the ZIA Admin Portal
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added the datasource ``zia_bandwidth_control_rule`` - Bandwidth Control Rules
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added and datasource ``zia_cloud_app_control_rule_actions`` - Retrieve all available actions for Cloud App Control Rules. This data source can be used to set the corresponding actions when configuring the resource `zia_cloud_app_control_rule`

### Enhancement

- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Added support to attribute `browser_eun_template_id` in the following resources `zia_cloud_app_control_rule`, `zia_url_filtering_rules` and `zia_file_type_control_rules`. The attribute allow for the configuration of [Browser End User Notification Templates](https://help.zscaler.com/zia/about-browser-eun-template)
- [PR #452](https://github.com/zscaler/terraform-provider-zia/pull/452) - Removed local validation for the attributes `url_categories`, `cloud_applications`, `applications`, `validity_time_zone_id`. See respective documentations for each resource for further instructions.

## 4.3.3 (June, 30 2025)

### Notes

- Release date: **(June, 30 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #449](https://github.com/zscaler/terraform-provider-zia/pull/449) - Fixed attribute `url_categories` in `zia_file_type_control_rules` expanding function.

## 4.3.2 (June, 23 2025)

### Notes

- Release date: **(June, 23 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #446](https://github.com/zscaler/terraform-provider-zia/pull/446) - Fixed `zia_dlp_web_rules` customizeDiff validation for attributes `external_auditor_email`, `auditor` and `notification_template`. Resoruce now allows for rule configuration when these attributes not not set.

## 4.3.1 (June, 23 2025)

### Notes

- Release date: **(June, 23 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #446](https://github.com/zscaler/terraform-provider-zia/pull/446) - Upgraded to [Zscaler-SDK-GO v3.5.1](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.5.1) to fix api error message parsing issues via legacy client.

## 4.3.0 (June, 19 2025)

### Notes

- Release date: **(June, 19 2025)**
- Supported Terraform version: **v1.x**

### NEW - RESOURCES, DATA SOURCES

- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - The following new resources and data sources have been introduced:

- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_browser_control_policy`` - Browser Control Policy
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_casb_dlp_rules`` - SaaS Security API (Casb DLP Rules)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_casb_malware_rules`` - SaaS Security API (Casb Malware Rules)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_cloud_application_instance`` - Cloud Application Instance
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_risk_profiles`` - Risk Profiles

### NEW DATA SOURCES

- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_casb_tenant`` - SaaS Security API (Casb Tenant)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_casb_email_label`` - SaaS Security API (Casb Email Label)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_casb_tombstone_template`` - SaaS Security API (Casb Quarantine Tombstone Template)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_casb_tombstone_template`` - SaaS Security API (Casb Quarantine Tombstone Template)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_domain_profiles`` - SaaS Security API (Casb Domain Profiles)
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added the datasource and resource ``zia_tenant_restriction_profile`` - Tenant Restriction Profile

### Bug Fixes

- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Added validation to ``zia_dlp_web_rules`` to prevent conflict between attributes: `auditor`, `external_auditor_email` and `notification_template`
- [PR #444](https://github.com/zscaler/terraform-provider-zia/pull/444) - Removed validation function `validateDestAddress` for the attribute `dest_addresses` to support both IPv4 Addresses and Wildcard FQDN.

## 4.2.0 (June, 11 2025)

### Notes

- Release date: **(June, 11 2025)**
- Supported Terraform version: **v1.x**

### NEW - RESOURCES, DATA SOURCES

- [PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) - The following new resources and data sources have been introduced:

- Added the datasource and resource ``zia_subscription_alert`` [PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Subscription Alerts
- Added the datasource and resource ``zia_forwarding_control_proxies``[PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Manage proxy for a third-party proxy service
- Added the datasource and resource ``zia_ftp_control_policy``[PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Manage FTP Control status and the list of URL categories for which FTP is allowed
- Added the datasource and resource ``zia_mobile_malware_protection_policy``[PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Manage Mobile Malware Protection rule
- Added the datasource and resource ``zia_nat_control_rules``[PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Manage DNAT Control policy rule
- Added the datasource and resource ``zia_nss_server``[PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Manage NSS server objects
- Added the datasource and resource ``zia_virtual_service_edge_cluster``[PR #439](https://github.com/zscaler/terraform-provider-zia/pull/439) :rocket: - Manage Virtual Service Edge cluster

## 4.1.5 (June, 5 2025)

### Notes

- Release date: **(June, 5 2025)**
- Supported Terraform version: **v1.x**

### Documentation

- [PR #435](https://github.com/zscaler/terraform-provider-zia/pull/435) - Fixed documentation spellings

## 4.1.4 (June, 5 2025)

### Notes

- Release date: **(June, 5 2025)**
- Supported Terraform version: **v1.x**

### Enhancement

- [PR #435](https://github.com/zscaler/terraform-provider-zia/pull/435) - Fixed `zia_firewall_filtering_rule` import issue with predefined rules.
- [PR #435](https://github.com/zscaler/terraform-provider-zia/pull/435) - Fixed country name and timezone validation for `zia_location_management` resource.
- [PR #435](https://github.com/zscaler/terraform-provider-zia/pull/435) - Fixed documentation spellings

## 4.1.3 (June, 5 2025)

### Notes

- Release date: **(June, 5 2025)**
- Supported Terraform version: **v1.x**

### Enhancement

- [PR #435](https://github.com/zscaler/terraform-provider-zia/pull/435) - Fixed `zia_firewall_filtering_rule` import issue with predefined rules.
- [PR #435](https://github.com/zscaler/terraform-provider-zia/pull/435) - Fixed country name and timezone validation for `zia_location_management` resource.

## 4.1.2 (May, 20 2025)

### Notes

- Release date: **(May, 13 2025)**
- Supported Terraform version: **v1.x**

### Enhancement

- [PR #431](https://github.com/zscaler/terraform-provider-zia/pull/431) - Upgraded to [Zscaler SDK GO v3.3.1] - Fixing ZIA User Pagination parameters.

## 4.1.1 (May, 13 2025)

### Notes

- Release date: **(May, 13 2025)**
- Supported Terraform version: **v1.x**

### Enhancement

- [PR #429](https://github.com/zscaler/terraform-provider-zia/pull/429) - Added new action to `CONFIRM` to resource `zia_dlp_web_rules`
- [PR #429](https://github.com/zscaler/terraform-provider-zia/pull/429) - Added new file_types to `FTCATEGORY_MS_PROJ` and `FTCATEGORY_APPINSTALLER` to resource `zia_dlp_web_rules`

## 4.1.0 (April, 18 2025)

### Notes

- Release date: **(April, 18 2025)**
- Supported Terraform version: **v1.x**

### Enhancement

- [PR #422](https://github.com/zscaler/terraform-provider-zia/pull/422) - Added new resource `resource_zia_admin_roles` to Admin Roles.

## 4.0.10 (April, 7 2025)

### Notes

- Release date: **(April, 7 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #416](https://github.com/zscaler/terraform-provider-zia/pull/416) - Fixed `zia_dlp_web_rules` sub rule reorder logic to ensure rules are ordered correctly.
- [PR #416](https://github.com/zscaler/terraform-provider-zia/pull/416) - Replaced attribute `malicious_urls` with `bypass_urls` in the resource `zia_atp_security_exceptions` documentation.
- [PR #416](https://github.com/zscaler/terraform-provider-zia/pull/416) - Fixed the flattening function `flattenIDExtensionsListIDs` and schema function `setIDsSchemaTypeCustom`. This will ensure Terraform identifies plan changes when block lists are removed from the configuration.
- [PR #416](https://github.com/zscaler/terraform-provider-zia/pull/416) - Fix to attribute the `order` attribute in all rule based resources to ensure consistency on ordering logic.
- [PR #416](https://github.com/zscaler/terraform-provider-zia/pull/416) - Fix  custom order logic on the resource `zia_firewall_filtering_rules` to ensure pre-defined rules are placed in the correct position to prevent drifts

## 4.0.9 (March, 14 2025)

### Notes

- Release date: **(March, 14 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #410](https://github.com/zscaler/terraform-provider-zia/pull/410) - Fixed `zia_dlp_web_rules` resource to fail fast during API errors.
- [PR #410](https://github.com/zscaler/terraform-provider-zia/pull/410) - Added fix to `zia_sandbox_rules` to ignore the `order` attribute for the default rule named: `Default BA Rule`. This will prevent potential drifts when rule is returned with a non default order number.

## 4.0.8 (February, 14 2025)

### Notes

- Release date: **(February, 14 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #402](https://github.com/zscaler/terraform-provider-zia/pull/402) - Fixed missing `url_category` attribute within the expand function for the resource `zia_ssl_inspection_rules`.
- [PR #402](https://github.com/zscaler/terraform-provider-zia/pull/402) - Updated provider to [zscaler-sdk-go v4.0.2](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v3.1.6)

## 4.0.7 (February, 13 2025)

### Notes

- Release date: **(February, 13  2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #397](https://github.com/zscaler/terraform-provider-zia/pull/397) - Fixed panic with `zia_ssl_inspection_rules` due to misconfigured flattening ID function within the read function.

## 4.0.6 (February, 12 2025)

### Notes

- Release date: **(February, 12  2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #396](https://github.com/zscaler/terraform-provider-zia/pull/396) - Fixed `zia_ssl_inspection_rules` validation error and panic issue.

## 4.0.5 (February, 10 2025)

### Notes

- Release date: **(February, 10  2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #393](https://github.com/zscaler/terraform-provider-zia/pull/393) - Fixed the custom ID for the following resources:
  - `zia_auth_settings_urls`
  - `zia_sandbox_behavioral_analysis`
  - `zia_security_settings`

## 4.0.4 (February, 6 2025)

### Notes

- Release date: **(February, 6  2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #392](https://github.com/zscaler/terraform-provider-zia/pull/392) - Improved the rule reorder logic to expedite reorder process for the following resources:
  - `zia_firewall_filtering_rule`
  - `zia_firewall_dns_rule`
  - `zia_firewall_ips_rule`
  - `zia_file_type_control_rules`
  - `zia_forwarding_control_rule`
  - `zia_ssl_inspection_rules`
  - `zia_sandbox_rules`

### Documentation

- [PR #392](https://github.com/zscaler/terraform-provider-zia/pull/392) - Updated documentation for tghe following resources describing reorder process and concept of predefined vs default rules
  - `zia_firewall_filtering_rule`
  - `zia_firewall_dns_rule`
  - `zia_ssl_inspection_rules`

## 4.0.3 (February, 5 2025)

### Notes

- Release date: **(February, 5  2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #391](https://github.com/zscaler/terraform-provider-zia/pull/391) - Added new url categories to validation function. The following new categories have been added:
  - `GLOBAL_INT_OFC365_ALLOW`
  - `GLOBAL_INT_OFC365_DEFAULT`
  - `GLOBAL_INT_OFC365_OPTIMIZE`

### IMPORTANT WARNING

- [PR #391](https://github.com/zscaler/terraform-provider-zia/pull/391) - For security reasons, authentication via configuration yaml file is not supported in this provider. Please use one of the documented authentication methods:
  - Environment Variables
  - Provider Block configuration

For information on the supported authentication methods please visit the Terraform Provider Registry [here](https://registry.terraform.io/providers/zscaler/zia/latest/docs)

## 4.0.2 (January, 31 2025)

### Notes

- Release date: **(January, 31 2025)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #388](https://github.com/zscaler/terraform-provider-zia/pull/388) - Fixed ZIA import resource for `zia_dlp_notification_templates` due to heredoc missformatting.
- [PR #388](https://github.com/zscaler/zscaler-terraformer/pull/257). Fixed ZIA import resource for `zia_end_user_notification` due to heredoc missformatting and attribute validation issue. - [Issue #387](https://github.com/zscaler/terraform-provider-zia/issues/387)
- [PR #388](https://github.com/zscaler/zscaler-terraformer/pull/388). Fixed ZIA import resources for: `zia_forwarding_control_zpa_gateway` due to missing attribute `type`.

## 4.0.1 (January, 29 2025)

### Notes

- Release date: **(January, 29 2025)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #384](https://github.com/zscaler/terraform-provider-zia/pull/384) - Fixed panic related to attribute `proxy_gateway` in the resource `zia_ssl_inspection_rules`.

## 4.0.0 (January, 22 2025) - BREAKING CHANGES

### Notes

- Release date: **(January, 22 2025)**
- Supported Terraform version: **v1.x**

#### Enhancements - Zscaler OneAPI Support

[PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383): The ZIA Terraform Provider now offers support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

⚠️ **WARNING**: Please refer to the [Index Page](https://github.com/zscaler/terraform-provider-zia/blob/master/docs/index.md) page for details on authentication requirements prior to upgrading your provider configuration.

⚠️ **WARNING**: Attention Government customers. OneAPI and Zidentity is not currently supported for the following clouds: `zscalergov` and `zscalerten`. Refer to the [Legacy API Framework](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index) section for more information on how authenticate to these environments using the legacy method.

### NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS

#### ENV VARS: ZIA Sandbox Submission - BREAKING CHANGES

[PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383): Authentication to Zscaler Sandbox service now use the following attributes.

- `sandboxToken` - Can also be sourced from the `ZSCALER_SANDBOX_TOKEN` environment variable.
- `sandboxCloud` - Can also be sourced from the `ZSCALER_SANDBOX_CLOUD` environment variable.

The use of the previous envioronment variables combination `ZIA_SANDBOX_TOKEN` and `ZIA_CLOUD` is now deprecated.

### NEW - RESOURCES, DATA SOURCES

[PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383): The following new resources and data sources have been introduced:

- Added the datasource and resource ``zia_sandbox_rules`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manage Sandbox Rules
- Added the datasource and resource ``zia_firewall_dns_rule``[PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manage Cloud Firewall DNS Rules
- Added the datasource and resource ``zia_firewall_ips_rule`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manage Cloud Firewall IPS Rules
- Added the datasource and resource ``zia_file_type_control_rules`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manage File Type Control Rules
- Added the datasource and resource ``zia_advanced_threat_settings`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages advanced threat configuration settings
- Added the datasource and resource ``zia_atp_malicious_urls`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages malicious URLs added to the denylist in ATP policy
- Added the datasource and resource ``zia_atp_security_exceptions`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages Security Exceptions (URL Bypass List) for the ATP policy
- Added the datasource and resource ``zia_advanced_settings`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages Advanced Settings configuration. [Configuring Advanced Settings](https://help.zscaler.com/zia/configuring-advanced-settings)
- Added the datasource and resource ``zia_atp_malware_inspection`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages Advanced Threat Protection Malware Inspection configuration. [Malware Protection](https://help.zscaler.com/zia/policies/malware-protection)
- Added the datasource and resource ``zia_atp_malware_protocols`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages Advanced Threat Protection Malware Protocols configuration. [Malware Protection](https://help.zscaler.com/zia/policies/malware-protection)
- Added the datasource and resource ``zia_atp_malware_settings`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages Advanced Threat Protection Malware Settings. [Malware Protection](https://help.zscaler.com/zia/policies/malware-protection)
- Added the datasource and resource ``zia_atp_malware_policy`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages Advanced Threat Protection Malware Policy. [Malware Protection](https://help.zscaler.com/zia/policies/malware-protection)
- Added the datasource and resource ``zia_end_user_notification`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Retrieves information of browser-based end user notification (EUN) configuration details.[Understanding Browser-Based End User Notifications](https://help.zscaler.com/zia/understanding-browser-based-end-user-notifications)
- Added the datasource and resource ``zia_url_filtering_and_cloud_app_settings`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages the URL and Cloud App Control advanced policy settings.[Configuring Advanced Policy Settings](https://help.zscaler.com/zia/configuring-advanced-policy-settings)
- Added the datasource ``zia_cloud_applications`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Retrieves Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, File Type Control rules, and SSL Inspection rules.
- Added the datasource ``zia_forwarding_control_proxy_gateway`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Retrieves information of existing Proxy Gateway configuration.
- Added the datasource and resource ``zia_ssl_inspection_rules`` [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) :rocket: - Manages SSL Inspection Rules.

#### NEW ATTRIBUTES

- [PR #383](https://github.com/zscaler/terraform-provider-zia/pull/383) - Added new `actions` values to resource `zia_cloud_app_control_rule`.
Please refer to the [Cloud Application Control - Rule Types vs Actions Matrix](https://github.com/zscaler/terraform-provider-zia/blob/master/docs/resources/zia_cloud_app_control_rule.md#cloud-application-control---rule-types-vs-actions-matrix) page for details each action per `rule_type`

## 3.0.7 (November, 17 2024)

### Notes

- Release date: **(November, 17  2024)**
- Supported Terraform version: **v1.x**

### Internal Fixes

- [PR #374](https://github.com/zscaler/terraform-provider-zia/pull/374) - Added new `file_types` supported values in the `zia_dlp_web_rules` resource. See the [zia_dlp_web_rules](https://registry.terraform.io/providers/zscaler/zia/latest/docs/resources/zia_dlp_web_rules) documentation.

## 3.0.6 (October, 8 2024)

### Notes

- Release date: **(October, 8  2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #374](https://github.com/zscaler/terraform-provider-zia/pull/374) - Added missing attribute `sourceCountries` to ZIA `firewallfilteringrule`

## 3.0.5 (October, 4 2024)

### Notes

- Release date: **(October, 4  2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #373](https://github.com/zscaler/terraform-provider-zia/pull/373) - The resource `zia_forwarding_control_rule` now pauses for 60 seconds before proceeding with the create or update process whenever the `forward_method` attribute is set to `ZPA`. In case of a failure related to resource synchronization, the provider will retry the resource creation or update up to 3 times, waiting 30 seconds between each retry. This behavior ensures that ZIA and ZPA have sufficient time to synchronize and replicate the necessary resource IDs, reducing the risk of transient errors during provisioning.
  **NOTE** This retry mechanism helps to automatically overcome temporary latency without manual intervention. This behavior does not affect forwarding rules configured with other forward_methods such as `DIRECT`.

## 3.0.4 (September, 6 2024)

### Notes

- Release date: **(September, 6 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #369](https://github.com/zscaler/terraform-provider-zia/pull/369) - Fixed `zia_dlp_web_rules` validation function for the attribute `file_types`.

## 3.0.3 (August, 26 2024)

### Notes

- Release date: **(August, 26 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #368](https://github.com/zscaler/terraform-provider-zia/pull/368) - Implemented runtime validation for the attribute `dest_addresses` in the resource: `zia_firewall_filtering_rule`. The provider now validates if the IP address provided is an IPv4.

## 3.0.2 (August, 19 2024)

### Notes

- Release date: **(August, 19 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #366](https://github.com/zscaler/terraform-provider-zia/pull/366) - Implemented runtime validation for resource: `zia_forwarding_control_rule`. The provider now validates incompatible attributes during the plan and apply stages at the schema level.
- [PR #366](https://github.com/zscaler/terraform-provider-zia/pull/366) - Fixed the datasource `zia_traffic_forwarding_gre_vip_recommended_list` to allow Geo location  information override when needed. The datasource now supports the following optional attributes:
  - `routable_ip` - (Boolean) The routable IP address.
  - `within_country_only` - (Boolean) Search within country only.
  - `include_private_service_edge` - (Boolean) Include ZIA Private Service Edge VIPs.
  - `include_current_vips` - (Boolean) Include currently assigned VIPs.
  - `latitude` - (Number) The latitude coordinate of the GRE tunnel source.
  - `longitude` - (Number) The longitude coordinate of the GRE tunnel source.
  - `subcloud` - (String) The longitude coordinate of the GRE tunnel source.

- [PR #366](https://github.com/zscaler/terraform-provider-zia/pull/366) - Added centralized semaphore functionality to manipulate concurrent request limitations.

## 3.0.1 (August, 13 2024)

### Notes

- Release date: **(August, 13 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #365](https://github.com/zscaler/terraform-provider-zia/pull/365) - Fixed `ports` attribute in `zia_location_management` resource to support `TypeSet` with elements of `TypeInt`.

### Documentation

- [PR #365](https://github.com/zscaler/terraform-provider-zia/pull/365) - Updated documentation for resources: `zia_location_management` and `zia_cloud_app_control_rule`

## 3.0.0 (August, 12 2024)

### Notes

- Release date: **(August, 12 2024)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #361](https://github.com/zscaler/terraform-provider-zia/pull/361) - Added new resource and datasource `zia_cloud_app_control_rule` for Cloud Application Control rule management.
- [PR #361](https://github.com/zscaler/terraform-provider-zia/pull/361) - Added new datasource `zia_dlp_dictionary_predefined_identifiers` to retrieve DLP Dictionary Hierarchical Identifiers. The information can be used when configuring DLP Dictionary resource attribute `hierarchical_identifiers` to clone predefined dictionaries.
- [PR #361](https://github.com/zscaler/terraform-provider-zia/pull/361) - Added new attribute `hierarchical_identifiers` to `zia_dlp_dictionaries` resource.
- [PR #361](https://github.com/zscaler/terraform-provider-zia/pull/361) - Enhanced `zia_security_settings` to support maximum number of blacklist urls.

### Bug Fixes

- [PR #361](https://github.com/zscaler/terraform-provider-zia/pull/361) - Added Semaphore retry logic to resource ``zia_url_categories`` to assist with rate limiting management.
- [PR #361](https://github.com/zscaler/terraform-provider-zia/pull/361) - Fixed `ports` attribute in `zia_location_management` resource to support `TypeList`.

## 2.91.4 (July, 3 2024)

### Notes

- Release date: **(July, 3  2024)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #357](https://github.com/zscaler/terraform-provider-zia/pull/357) - Fixed ``zia_url_filtering_rules`` drift due to attribute conversion ``validatidy_start_time`` and ``validity_end_time``.

## 2.91.3 (July, 2 2024)

### Notes

- Release date: **(July, 2  2024)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #356](https://github.com/zscaler/terraform-provider-zia/pull/356) - Fixed ``zia_url_filtering_rules`` schema validation to ensure proper validation during plan and apply stages.
- [PR #356](https://github.com/zscaler/terraform-provider-zia/pull/356) - Fixed ``zia_location_management`` drift due to missing `state` attribute in the READ function.

## 2.91.2 (July, 2 2024)

### Notes

- Release date: **(July, 2  2024)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #356](https://github.com/zscaler/terraform-provider-zia/pull/356) - Fixed ``zia_url_filtering_rules`` schema validation to ensure proper validation during plan and apply stages.

## 2.91.1 (June, 29 2024)

### Notes

- Release date: **(June, 29  2024)**
- Supported Terraform version: **v1.x**

### Bug Fix

- [PR #354](https://github.com/zscaler/terraform-provider-zia/pull/354) - Fixed go.mod and go.sum
- [PR #354](https://github.com/zscaler/terraform-provider-zia/pull/354) - Fixed computed attributes in the schema

## 2.91.0 (June, 19 2024)

### Notes

- Release date: **(June, 19  2024)**
- Supported Terraform version: **v1.x**

### BREAKING CHANGES and ENHACEMENTS

- [PR #350](https://github.com/zscaler/terraform-provider-zia/pull/350)
  - `zia_url_filtering_rules` - The provider now explicitly validates during the plan and apply stages which attributes can be set based on the `action` value.
  - `zia_url_filtering_rules` - The provider now allows for the use of `RFC1123` date and time format i.e `Sun, 16 Jun 2024 15:04:05 UTC` when setting the attributes `validity_start_time` and `validity_end_time` instead of the native epoch unix format.

    ~> **NOTE** This change is not backwards compatible.
  - `zia_url_filtering_rules` - The provider now explicitly validates the attribute `validity_time_zone_id` against the official [IANA List](https://nodatime.org/TimeZones). The supported format is: `"US/Pacific"`

    ~> **NOTE** This change is not backwards compatible.

  - `ziaActivator` - The Out-of-band ZIA Activator has been updated to directly leverage the [Zscaler-SDK-GO](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.61.0).
    ~> **NOTE** If you plan to update your provider installation to the latest v2.91.0, you must re-compile the utility program.
    ~> **NOTE** Note that as of release [v2.8.2](https://github.com/zscaler/terraform-provider-zia/releases/tag/v2.8.2) the provider offers the option to trigger activation by setting the `ZIA_ACTIVATION` environment variable. With this enhancement the activation occurs only when this environment variable is set to `true`.

### Internal Changes

- [PR #350](https://github.com/zscaler/terraform-provider-zia/pull/350) - Upgraded to [Zscaler-SDK-GO](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.61.0). The upgrade supports easier ZIA API Client instantiation for existing and new resources.
- [PR #350](https://github.com/zscaler/terraform-provider-zpa/pull/350) Upgraded ``releaser.yml`` to [GoReleaser v6](https://github.com/goreleaser/goreleaser-action/releases/tag/v6.0.0)

## 2.9.1 (June, 14 2024)

### Notes

- Release date: **(June, 14  2024)**
- Supported Terraform version: **v1.x**

### Internal Changes

- [PR #350](https://github.com/zscaler/terraform-provider-zia/pull/350) - Upgraded to [Zscaler-SDK-GO](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.61.0). The upgrade supports easier ZIA API Client instantiation for existing and new resources.
- [PR #350](https://github.com/zscaler/terraform-provider-zpa/pull/pull/350) Upgraded ``releaser.yml`` to [GoReleaser v6](https://github.com/goreleaser/goreleaser-action/releases/tag/v6.0.0)

## 2.9.0 (May, 22 2024) - BREAKING CHANGE

### Notes

- Release date: **(May, 22 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes - BREAKING CHANGE

- [PR #345](https://github.com/zscaler/terraform-provider-zia/pull/345) - The attribute `ocr_enabled` has been deprecated at the upstream API and is no longer accepted. The OCR feature must be enabled via the [DLP Advanced Settings](https://help.zscaler.com/zia/configuring-dlp-advanced-settings).
  **NOTE** DLP engines support OCR scanning of `PNG`, `JPEG`, `TIFF`, and `BMP` files.

- [PR #345](https://github.com/zscaler/terraform-provider-zia/pull/345) - Implemented Fix for `zia_dlp_web_rules` for new attributes `parent_rule` and `sub_rules`. A parent rule must be configured with rank 0 and prior to any potential subrule. It is not possible to add existing rules as as subrules under the parent rule.

## 2.8.31 (May, 21 2024)

### Notes

- Release date: **(May, 21 2024)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- [PR #344](https://github.com/zscaler/terraform-provider-zia/344) - Fixed `id` conversion for the resource `zia_traffic_forwarding_vpn_credentials` to ensure proper state file setting.

- [PR #344](https://github.com/zscaler/terraform-provider-zia/pull/344) - Upgraded to [Zscaler-SDK-GO v2.5.2](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.5.2)

## 2.8.3 (May, 7 2024)

### Notes

- Release date: **(May, 7 2024)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #340](https://github.com/zscaler/terraform-provider-zia/pull/340) - Added new ZIA URL Filtering Rule attribute `source_ip_groups` to resources: `zia_url_filtering_rules` and `zia_dlp_web_rules`
- [PR #340](https://github.com/zscaler/terraform-provider-zia/pull/340) - Upgraded to [Zscaler-GO-SDK v2.5.0](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.5.0)

## 2.8.21 (April, 8 2024)

### Notes

- Release date: **(April, 8 2024)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #336](https://github.com/zscaler/terraform-provider-zia/pull/336) - Upgraded provider to [Zscaler-SDK-GO v2.4.35](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v2.4.35)

## 2.8.2 (April, 8 2024)

### Notes

- Release date: **(April, 8 2024)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #332](https://github.com/zscaler/terraform-provider-zia/pull/332) - Implemented optional environment variable `ZIA_ACTIVATION` for optional configuration activation. This is an improved version of the initial release [v2.8.0](https://github.com/zscaler/terraform-provider-zia/releases/tag/v2.8.0) where activations were done implicitly for every resource. With this enhancement the activation will only occur when this environment variable is set to true.

## 2.8.1 (March, 27 2024)

### Notes

- Release date: **(March, 27 2024)**
- Supported Terraform version: **v1.x**

### Documentation

- Redacted several password creation examples to prevent GitGuardian false positives. A header comment has also been added to advise.

## 2.8.0 (March, 27 2024)

### Notes

- Release date: **(March, 27 2024)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #330](https://github.com/zscaler/terraform-provider-zia/pull/330) - Implemented auto activation functionality to all supported resources. Configurations will now be activated during `CREATE`, `UPDATE` AND `DELETE` actions when executing `terraform apply` or `terraform destroy`, which removes the need of out of band activation or the use of the resource: `zia_activation_status`.

### Fixes

- [PR #330](https://github.com/zscaler/terraform-provider-zia/pull/330) - Fixed `zia_user_management` resource to support activation pre and post user enrolment using `BASIC` authentication method.

## 2.7.33 (March, 6 2024)

### Notes

- Release date: **(March, 6 2024)**
- Supported Terraform version: **v1.x**

### Enhacements

- [PR #325](https://github.com/zscaler/terraform-provider-zia/pull/325) Updated [support guide](/docs/guides/support.md) with new Zscaler support model.
- [PR #325](https://github.com/zscaler/terraform-provider-zia/pull/325) - Added support to import of the following resources:
- ``zia_auth_settings_urls``
- ``zia_sandbox_behavioral_analysis``
- ``zia_security_settings``

## 2.7.32 (February, 28 2024)

### Notes

- Release date: **(February, 28 2024)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #322](https://github.com/zscaler/terraform-provider-zia/pull/322) - Fixed validation `zia_url_filtering_rules` resource to validate `protocols` attribute to accept `HTTP_RULE` and `HTTPS_RULE`.
- [PR #322](https://github.com/zscaler/terraform-provider-zia/pull/322) - Fixed validation `zia_url_filtering_rules` validations for rules with `action` configured as `ISOLATE`.
- [PR #322](https://github.com/zscaler/terraform-provider-zia/pull/322) - Fixed linter issues across several acceptance tests resources and data sources.

## 2.7.31 (February, 28 2024)

### Notes

- Release date: **(February, 28 2024)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #321](https://github.com/zscaler/terraform-provider-zia/pull/321) - Fixed validation function in the resource `zia_url_filtering_rules` for the attribute `protocols`. The provider now validates the following API supported values: `SMRULEF_ZPA_BROKERS_RULE`, `ANY_RULE`, `TCP_RULE`, `UDP_RULE`, `DOHTTPS_RULE`, `TUNNELSSL_RULE`, `HTTP_PROXY`, `FOHTTP_RULE`, `FTP_RULE`, `HTTPS_RULE`, `HTTP_RULE`, `SSL_RULE`, `TUNNEL_RULE`, `WEBSOCKETSSL_RULE`, `WEBSOCKET_RULE`,

# 2.7.3 (February 14, 2024)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #319](https://github.com/zscaler/terraform-provider-zia/pull/319) - Implemented validation to the following resources:
  - `zia_firewall_filtering_destination_groups`
  - `zia_firewall_filtering_rule`
  - `zia_forwarding_control_zpa_gateway`
  - `zia_forwarding_control_policy`

# 2.7.2 (January 31, 2024)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #315](https://github.com/zscaler/terraform-provider-zia/pull/315) - Added support to new `workload_groups` attributes to the following resources:
  - ``zia_firewall_filtering_rule``
  - ``zia_url_filtering_rules``
  - ``zia_dlp_web_rules``

### Fixes

- [PR #315](https://github.com/zscaler/terraform-provider-zia/pull/315) - Fixed panic within the resource ``zia_location_management`` when setting the attribute ``ip_addresses`` in a sub-location. The provider now supports and validates the following ``ip_addresses`` formats:
  - `10.0.0.0-10.0.0.255`
  - `10.0.0.1`

  ~> **NOTE** CIDR notation is currently not supported due to API response incosistencies that may introduce drifts in the Terraform execution. This issue will be addressed in the future.

# 2.7.1 (January 26, 2024)

## Notes
- Golang: **v1.19**

### Enhacements

- [PR #313](https://github.com/zscaler/terraform-provider-zia/pull/313) - Added support for ZIA Workload Groups Tagging

## 2.7.0 (January, 15 2023)

### Notes

- Release date: **(January, 15 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

NEW - RESOURCES, DATA SOURCES

- [PR #293](https://github.com/zscaler/terraform-provider-zia/pull/293) - ✨ Added support for ZIA 🆕 Custom ZPA Gateway for use with Forwarding Control policy to forward traffic to ZPA for Source IP Anchoring.
- [PR #294](https://github.com/zscaler/terraform-provider-zia/pull/294) - ✨ Added support for ZIA 🆕 Forwarding Control Rule configuration.

- [PR #295](https://github.com/zscaler/terraform-provider-zia/pull/295) - ✨ Added ZIA Sandbox MD5 Hash and verdict report submission Resources:
  - **Sandbox Advanced Settings** - `zia_sandbox_behavioral_analysis` Gets and Upddates the custom list of MD5 file hashes that are blocked by Sandbox.
  - **Sandbox Report** - `zia_sandbox_report` Gets a full (i.e., complete) or summary detail report for an MD5 hash of a file that was analyzed by Sandbox.

- [PR #295](https://github.com/zscaler/terraform-provider-zia/pull/295) - ✨ Added ZIA Sandbox raw and archive file submission:
  - **Sandbox Submission** - `zia_sandbox_file_submission` - Submits raw or archive files (e.g., ZIP) to Sandbox for analysis. You can submit up to 100 files per day and it supports all file types that are currently supported by Sandbox.
  - **Sandbox Submission** - `zia_sandbox_file_submission` -  Submits raw or archive files (e.g., ZIP) to the Zscaler service for out-of-band file inspection to generate real-time verdicts for known and unknown files. It leverages capabilities such as Malware Prevention, Advanced Threat Prevention, Sandbox cloud effect, AI/ML-driven file analysis, and integrated third-party threat intelligence feeds to inspect files and classify them as benign or malicious instantaneously.
    ⚠️ **Note:**: The ZIA Terraform provider requires both the `ZIA_CLOUD` and `ZIA_SANDBOX_TOKEN` in order to authenticate to the Zscaler Cloud Sandbox environment. For details on how obtain the API Token visit the Zscaler help portal [About Sandbox API Token](https://help.zscaler.com/zia/about-sandbox-api-token)

- [PR #302](https://github.com/zscaler/terraform-provider-zia/pull/302) - Added new `zia_dlp_web_rules` attributes:
  - `severity` - Supported values: `RULE_SEVERITY_HIGH`, `RULE_SEVERITY_MEDIUM`, `RULE_SEVERITY_LOW`, `RULE_SEVERITY_INFO`
  - `user_risk_score_levels` - Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`
  - `parent_rule`
  - `sub_rules`

- [PR #308](https://github.com/zscaler/terraform-provider-zia/pull/308) - ✨ Added 🆕 Cloud Browser Isolation Profile data source. The data source can be used to associate a CBI profile with the `zia_url_filtering_rules` resource when the action is set to `ISOLATE`

- [PR #309](https://github.com/zscaler/terraform-provider-zia/pull/309) - ✨ Added 🆕 support to the following attributes within the `zia_firewall_filtering_rule`:
  - `device_trust_levels` - Supported values: `UNKNOWN_DEVICETRUSTLEVEL`, `LOW_TRUST`, `MEDIUM_TRUST`, `HIGH_TRUST`
  - `user_risk_score_levels` - Supported values: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`
  - `devices`
  - `device_groups`

- [PR #309](https://github.com/zscaler/terraform-provider-zia/pull/309) - ✨ Added new attribute `zpa_app_segments` to `zia_firewall_filtering_rule` to support ZPA Application Segments. Only ZPA application segments that have the Source IP Anchor option enabled are supported.

### Fixes

- [PR #299](https://github.com/zscaler/terraform-provider-zia/pull/299) - Fixed panic with ``zia_url_categories``.
- [PR #302](https://github.com/zscaler/terraform-provider-zia/pull/302) - Fixed `zia_dlp_web_rules` File Types validation function.

## 2.6.6 (November, 23 2023)

### Notes

- Release date: **(November, 23 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #291](https://github.com/zscaler/terraform-provider-zia/pull/291) - Fixed panic with resource `zia_admin_users` due to API changes.

## 2.6.5 (November, 5 2023)

### Notes

- Release date: **(November, 5 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #285](https://github.com/zscaler/terraform-provider-zia/pull/285) - Fixed drift within `zia_firewall_filtering_rule` for the attribute `dest_countries`.

## 2.6.3 (October, 18 2023)

### Notes

- Release date: **(October, 18 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #278](https://github.com/zscaler/terraform-provider-zia/pull/278) - Provider HTTP Header now includes enhanced ``User-Agent`` information for troubleshooting assistance.
  - i.e ``User-Agent: (darwin arm64) Terraform/1.5.5 Version/2.6.3``
- [PR #283](https://github.com/zscaler/terraform-provider-zia/pull/283) - Upgrade to Zscaler-SDK-GO v2.1.4

## 2.6.2 (September, 19 2023)

### Notes

- Release date: **(September, 19 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #276](https://github.com/zscaler/terraform-provider-zia/pull/276) - Added Country code validation for attribute `dest_countries` in the resource `zia_firewall_filtering_rule`. The provider validates the use of proper 2 letter country codes [ISO3166 By Alpha2Code](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)

- [PR #276](https://github.com/zscaler/terraform-provider-zia/pull/276) - Added Country name validation for attribute `country` in the resource `zia_location_management`. The provider validates the use uppercase country codes using [ISO-3166-1](https://en.wikipedia.org/wiki/ISO_3166-1)

## 2.6.1 (August, 29 2023)

### Notes

- Release date: **(August, 29 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #258](https://github.com/zscaler/terraform-provider-zia/pull/258) Improved geographical coordinates for attributes `latitude` and `longitude` in the resource `zia_traffic_forwarding_static_ip` to ensures that the state always mirrors the backend system's values.

### Fixes

- [PR #259](https://github.com/zscaler/terraform-provider-zia/pull/259) Fixed drift problem within the resource `zia_firewall_filtering_network_service_groups`.
- [PR #266](https://github.com/zscaler/terraform-provider-zia/pull/266) Fixed drift problem within the resource `zia_url_filtering_rules` order attribute

- [PR #260](https://github.com/zscaler/terraform-provider-zia/pull/260) Updated `zia_firewall_filtering_network_service` resource documentation.
!> **NOTE:** Resources of type `PREDEFINED` are built-in resources within the ZIA cloud and must be imported before the Terraform execution. Attempting to update the resource directly will return `DUPLICATE_ITEM` error message. To import a predefined built-in resource use the following command for example: `terraform import zia_firewall_filtering_network_service.this "DHCP"`

## 2.6.0 (August, 1 2023)

### Notes

- Release date: **(August, 1 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #257](https://github.com/zscaler/terraform-provider-zia/pull/257) Added New Public ZIA DLP Engine Endpoints (POST/PUT/DELETE)
⚠️ **WARNING:** "Before using the new ``zia_dlp_engines`` resource contact [Zscaler Support](https://help.zscaler.com/login-tickets)." and request the following API methods ``POST``, ``PUT``, and ``DELETE`` to be enabled for your organization.

### Fixes

- [PR #251](https://github.com/zscaler/terraform-provider-zia/pull/251) Added new predefied URL Category ``AI_ML_APPS`` to resource ``resource_zia_url_categories``.
- [PR #253](https://github.com/zscaler/terraform-provider-zia/pull/253) Fixed documentation for resource ``zia_firewall_filtering_destination_groups``

## 2.5.6 (June, 10 2023)

### Notes

- Release date: **(June, 10 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- Updated to Zscaler-SDK-GO v1.5.5. The update improves search mechanisms for ZIA resources, to ensure streamline upstream GET API requests and responses using ``search`` parameter. Notice that not all current API endpoints support the search parameter, in which case, all resources will be returned.

## 2.5.5 (May, 29 2023)

### Notes

- Release date: **(May, 29 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #244](https://github.com/zscaler/terraform-provider-zia/pull/244) Fix ``zia_user_management`` to ensure when the ``auth_methods``attribute is set, and user password is changed, the provide will re-enroll the user to update the password.

## 2.5.4 (May, 25 2023)

### Notes

- Release date: **(May, 25 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #234](https://github.com/zscaler/terraform-provider-zia/pull/234) Fix expand functions to ensure correct API response processing across all resource rule creation.

## 2.5.3 (May, 13 2023)

### Notes

- Release date: **(May, 13 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #231](https://github.com/zscaler/terraform-provider-zia/pull/219) ``zia_dlp_web_rules``: Fixed panic with ``zia_web_dlp_rules`` due to ``dlp_engines`` attribute expand function

## 2.5.2 (May, 1 2023)

### Notes

- Release date: **(May, 1 2023)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #224](https://github.com/zscaler/terraform-provider-zia/pull/224) ``zia_dlp_web_rule``: Reduced TimeTicker for faster rule order processing during creation and modifications.
- [PR #224](https://github.com/zscaler/terraform-provider-zia/pull/224) ``zia_dlp_web_rule``: Updated DLP Web Rule documentation with more examples
- [PR #226](https://github.com/zscaler/terraform-provider-zia/pull/226) Expanded ZIA search criteria to include auditor users.
- [PR #227](https://github.com/zscaler/terraform-provider-zia/pull/227) Introduced new attribute ``parent_name`` to the resource ``zia_location_management``. The attribute will allow the ability to search for sublocation resources across multiple parent locations specially when overlapping names are in use. Issue [#223](https://github.com/zscaler/terraform-provider-zia/issues/223)

### Fixes

- [PR #219](https://github.com/zscaler/terraform-provider-zia/pull/219) ``zia_dlp_web_rules``: Fixed drift issues with attributes ``url_categories`` and ``dlp_engines``
- [PR #221](https://github.com/zscaler/terraform-provider-zia/pull/221) ``zia_dlp_dictionary``: Fix DLP dictionary resource when ``phrase`` attribute is not provided
- [PR #228](https://github.com/zscaler/terraform-provider-zia/pull/228) ``zia_dlp_dictionary``: Fixed ``idm_profile_match_accuracy`` attribute to prevent drifts, plus accept ``zia_dlp_idm_profile_lite`` template_id when selecting ``dictionary_type`` INDEXED_DATA_MATCH

## 2.5.1 (April, 12 2023)

### Notes

- Release date: **(April, 12 2023)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #213](https://github.com/zscaler/terraform-provider-zia/pull/213) ``zia_location_management``: Added to support to sub-location search within data source. Issue [#209](https://github.com/zscaler/terraform-provider-zia/issues/209)

### Fixes

- [PR #217](https://github.com/zscaler/terraform-provider-zia/pull/217) ``zia_dlp_engines``: Fixed DLP Engine data source to allow search for predefined engines. Issue [#216](https://github.com/zscaler/terraform-provider-zia/issues/216)
- [PR #219](https://github.com/zscaler/terraform-provider-zia/pull/219) ``zia_dlp_web_rules``: DLP Web rule configuration drift for certain attributes when not set in order.

## 2.5.0 (March, 20 2023)

### Notes

- Release date: **(March, 20 2023)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #202](https://github.com/zscaler/terraform-provider-zia/pull/202) ``zia_user_management``: Implemented new attribute ``auth_methods``. The attribute supports the following values: ``BASIC`` and/or ``DIGEST``.
- ``zia_location_management``: Implemented new attribute ``basic_auth_enabled``. The supported values are: ``true`` or ``false``

- [PR #202](https://github.com/zscaler/terraform-provider-zia/pull/202) The provider now supports authentication to Zscaler ``preview`` and ``zscalerten`` clouds.

- [PR #211](https://github.com/zscaler/terraform-provider-zia/pull/211) Added new datasource ``zia_location_lite``. This data source can be used to return the "Road Warrior" location, which can then be used in the following resources: ``zia_url_filtering_rules``, ``zia_firewall_filtering_rule`` and ``zia_dlp_web_rules``

- [PR #213](https://github.com/zscaler/terraform-provider-zia/pull/213) Added support to search for sub-location within the resource ``zia_location_management``

### Fixes

- [PR #212](https://github.com/zscaler/terraform-provider-zia/pull/212) ``zia_user_management``: Fixed flattening function to expand group attribute values. Issue [#205](https://github.com/zscaler/terraform-provider-zia/issues/205)

## 2.4.6 (March, 6 2023)

### Notes

- Release date: **(March, 6 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- ``zia_location_management``: Fixed IPv4 Address and IPv4Address range validation.
- ``zia_traffic_forwarding_static_ip``: Fixed Longitude and Latitude computed attributes.
- ``zia_url_categories``: Removed ``Default: false`` attribute to prevent drifts.

## 2.4.5 (March, 2 2023)

### Notes

- Release date: **(March, 2 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #199](https://github.com/zscaler/terraform-provider-zia/pull/199) Improved ``Timeout`` reorder functions to ensure the rules across the below resources are organized correctly.
  - ``zia_firewall_filtering_rule`

- [PR #200](https://github.com/zscaler/terraform-provider-zia/pull/200) Improved ``Timeout`` reorder functions to ensure the rules across the below resources are organized correctly.
  - ``zia_dlp_web_rules`
  - ``zia_url_filtering_rules`

## 2.4.4 (March, 1 2023)

### Notes

- Release date: **(March, 1 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #193](https://github.com/zscaler/terraform-provider-zia/pull/193) Added new following new datasources:
  - ``zia_firewall_filtering_application_services`` The returned values are:
    - ``SKYPEFORBUSINESS``, ``FILE_SHAREPT_ONEDRIVE``, ``EXCHANGEONLINE``, ``M365COMMON``, ``ZOOMMEETING``, ``WEBEXMEETING``, ``WEBEXTEAMS``, ``WEBEXCALLING``, ``RINGCENTRALMEETING``, ``GOTOMEETING``, ``GOTOMEETING_INROOM``, ``LOGMEINMEETING``, ``LOGMEINRESCUE``

  - ``zia_firewall_filtering_application_services_group`` The returned values are:
    - ``OFFICE365``, ``ZOOM``, ``WEBEX``, ``RINGCENTRAL``, ``LOGMEIN``

### Fixes

- [PR #194](https://github.com/zscaler/terraform-provider-zia/pull/194) Improved ``Timeout`` reorder functions to ensure the rules across the below resources are organized correctly.
  - ``zia_dlp_web_rules``
  - ``zia_url_filtering_rules``
  - ``zia_firewall_filtering_rule`

⚠️ **WARNING:** Due to API limitations, we recommend to limit the number of requests to ONE, when configuring the above resources.

  This will allow the API to settle these resources in the correct order. Pushing large batches of security rules at once, may incur in Terraform to Timeout after 20 mins, as it will try to place the rules in the incorrect order. This issue will be addressed in future versions.

In order to accomplish this, make sure you set the
[parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at or
below this limit to prevent performance impacts.

- [PR #195](https://github.com/zscaler/terraform-provider-zia/pull/195) Fixed ``zia_traffic_forwarding_gre_tunnel`` by removing unecessary computed values to prevent drifts.

## 2.4.3 (February, 28 2023)

### Notes

- Release date: **(February, 28 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #193](https://github.com/zscaler/terraform-provider-zia/pull/193) Added new following new datasources:
  - ``zia_firewall_filtering_application_services`` The returned values are:
    - ``SKYPEFORBUSINESS``, ``FILE_SHAREPT_ONEDRIVE``, ``EXCHANGEONLINE``, ``M365COMMON``, ``ZOOMMEETING``, ``WEBEXMEETING``, ``WEBEXTEAMS``, ``WEBEXCALLING``, ``RINGCENTRALMEETING``, ``GOTOMEETING``, ``GOTOMEETING_INROOM``, ``LOGMEINMEETING``, ``LOGMEINRESCUE``

  - ``zia_firewall_filtering_application_services_group`` The returned values are:
    - ``OFFICE365``, ``ZOOM``, ``WEBEX``, ``RINGCENTRAL``, ``LOGMEIN``

### Fixes

- [PR #194](https://github.com/zscaler/terraform-provider-zia/pull/194) Improved ``Timeout`` reorder functions to ensure the rules across the below resources are organized correctly.
  - ``zia_dlp_web_rules``
  - ``zia_url_filtering_rules``
  - ``zia_firewall_filtering_rule`

⚠️ **WARNING:** Due to API limitations, we recommend to limit the number of requests to ONE, when configuring the above resources.

  This will allow the API to settle these resources in the correct order. Pushing large batches of security rules at once, may incur in Terraform to Timeout after 20 mins, as it will try to place the rules in the incorrect order. This issue will be addressed in future versions.

In order to accomplish this, make sure you set the
[parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at or
below this limit to prevent performance impacts.

- [PR #195](https://github.com/zscaler/terraform-provider-zia/pull/195) Fixed ``zia_traffic_forwarding_gre_tunnel`` by removing unecessary computed values to prevent drifts.

## 2.4.2 (February, 13 2023)

### Notes

- Release date: **(February, 13 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #180](https://github.com/zscaler/terraform-provider-zia/pull/180) Implemented customizable ``Timeouts`` for Create and Update functions to help with rule reorder across the following resources:
  - ``zia_dlp_web_rules``
  - ``zia_url_filtering_rules``
  - ``zia_firewall_filtering_rule``

- [PR #182](https://github.com/zscaler/terraform-provider-zia/pull/182) Implemented validation for ``ocr_enabled`` attribute validation for ``zia_dlp_web_rules``

## 2.4.1 (February, 10 2023)

### Notes

- Release date: **(February, 10 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #181](https://github.com/zscaler/terraform-provider-zia/pull/181) Implemented customizable ``Timeouts`` for Create and Update functions to help with rule reorder across the following resources:
  - ``zia_dlp_web_rules``
  - ``zia_url_filtering_rules``
  - ``zia_firewall_filtering_rule``

## 2.4.0 (January, 31 2023)

### Notes

- Release date: **(January, 31 2023)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #176](https://github.com/zscaler/terraform-provider-zia/pull/176) Added the following ZIA data sources
  - ``zia_dlp_icap_servers`` - Gets a the list of DLP servers using ICAP
  - ``zia_dlp_incident_receiver_servers`` - Gets a list of DLP Incident Receivers
  - ``zia_dlp_idm_profiles`` - Indexed Document Match (IDM) template (or profile) information.

## 2.3.6 (January, 25 2023)

### Notes

- Release date: **(January, 25 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #171](https://github.com/zscaler/terraform-provider-zia/pull/171) - Update to Zscaler-Go-SDK to fix bool parameter ``enable_full_logging`` in the ZIA Firewall Filtering resource.
- [PR #174](https://github.com/zscaler/terraform-provider-zia/pull/174) - Fix ``zia_web_rules`` file_types attribute to accept empty values. Also, added new supported file types to the validation fuction.

## 2.3.5 (January, 12 2023)

### Notes

- Release date: **(January, 12 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #160](https://github.com/zscaler/terraform-provider-zia/pull/160) - Fixed Pagination Issues across all resources

## 2.3.4 (January, 4 2023)

### Notes

- Release date: **(January, 4 2023)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #168](https://github.com/zscaler/terraform-provider-zia/pull/168) ``zia_firewall_filtering_rule`` Added the following new network applications to validation function
  - ``VMWARE_HORIZON_VIEW``,``ADOBE_CREATIVE_CLOUD``, ``ZOOMINFO``, ``SERVICE_NOW``, ``MS_SSAS``, ``GOOGLE_DNS``, ``CLOUDFLARE_DNS``, ``ADGUARD``, ``QUAD9``, ``OPENDNS``, ``CLEANBROWSING``, ``COMCAST_DNS``, ``NEXTDNS``, ``POWERDNS``,``BLAHDNS``,``SECUREDNS``,``RUBYFISH``,``DOH_UNKNOWN``,``GOOGLE_KEEP``,``AMAZON_CHIME``,``WORKDAY``,``FIFA``,``ROBLOX``,``WANGWANG``,``S7COMM_PLUS``,``DOH``,``AGORA_IO``,``MS_DFSR``,``WS_DISCOVERY``,``STUN``,``FOLDINGATHOME``,``GE_PROCIFY``,``MOXA_ASPP``,``APP_CH``,``GLASSDOOR``,``TINDER``,``BAIDU_TIEBA``,``MIMEDIA``,``FILESANYWHERE``,``HOUSEPARTY``,``GBRIDGE``,``HAMACHI``,``HEXATECH``,``HOTSPOT_SHIELD``,``MEGAPROXY``,``OPERA_VPN``,``SPOTFLUX``,``TUNNELBEAR``,``ZENMATE``, ``OPENGW``, ``VPNOVERDNS``, ``HOXX_VPN``, ``VPN1_COM``, ``SPRINGTECH_VPN``, ``BARRACUDA_VPN``, ``HIDEMAN_VPN``, ``WINDSCRIBE``, ``BROWSEC_VPN``, ``EPIC_BROWSER_VPN``, ``SKYVPN``, ``KPN_TUNNEL``, ``ERSPAN``,``EVASIVE_PROTOCOL``, ``DOTDASH``, ``ADOBE_DOCUMENT_CLOUD``, ``FLIPKART_BOOKS``

- [PR #165](https://github.com/zscaler/terraform-provider-zia/pull/162) ``zia_url_filtering_rules`` Added new URL Categories

## 2.3.3 (January, 1 2023)

### Notes

- Release date: **(January, 1 2023)**
- Supported Terraform version: **v1.x**
## 2.3.2 (December, 30 2022)

### Notes

- Release date: **(December, 30 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #164](https://github.com/zscaler/terraform-provider-zia/pull/164) Added missing URL Category resource parameters
- [PR #165](https://github.com/zscaler/terraform-provider-zia/pull/162) Added missing new URL Category pre-validation to ``zia_url_filtering_rule`` The new categories are: `DYNAMIC_DNS` and `NEWLY_REVIVED_DOMAINS`

## 2.3.1 (December, 3 2022)

### Notes

- Release date: **(December, 3 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #150](https://github.com/zscaler/terraform-provider-zia/pull/150) Fixed DLP Web rule resource panic due to incorrect assignment
- [PR #150](https://github.com/zscaler/terraform-provider-zia/pull/150) Fixed DLP Notification Template resource panic due to incorrect assignment
- [PR #151](https://github.com/zscaler/terraform-provider-zia/pull/151) Fixed DLP Dictionary panic due to incorrect assignment

## 2.3.0 (November, 25 2022)

### Notes

- Release date: **(November, 25 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #147](https://github.com/zscaler/terraform-provider-zia/pull/147) Fixed Read/Update/Delete functions to allow automatic recreation of resources, that have been manually deleted via the UI.
- [PR #147](https://github.com/zscaler/terraform-provider-zia/pull/147) Removed ``deprecated`` helper from ``zia_location_management`` resource.

## 2.2.3 (October, 20 2022)

### Notes

- Release date: **(October, 20 2022)**
- Supported Terraform version: **v1.x**

### Enhancements

- [PR #137](https://github.com/zscaler/terraform-provider-zia/pull/137) Added Customizable Timeouts to zia_activation_status resource.
- [PR #138](https://github.com/zscaler/terraform-provider-zia/pull/138) Added acceptance test to ``zia_activation_status`` data source.

### Fixes

- [PR #134](https://github.com/zscaler/terraform-provider-zia/pull/134) Update to zscaler-sdk-go v0.1.1
- [PR #135](https://github.com/zscaler/terraform-provider-zia/pull/135) Update to zscaler-sdk-go v0.1.2
- [PR #135](https://github.com/zscaler/terraform-provider-zia/pull/135) Added missing parameter ``comment`` to ``zia_traffic_forwarding_static_ips``
- [PR #136](https://github.com/zscaler/terraform-provider-zia/pull/136) Updated Documentation for zia_activation_status resource and data source.

## 2.2.2 (September, 25 2022)

### Notes

- Release date: **(September, 25 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #130](https://github.com/zscaler/terraform-provider-zia/pull/130) Fix Import Resource By ID

## 2.2.1 (September, 21 2022)

### Notes

- Release date: **(September, 21 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #127](https://github.com/zscaler/terraform-provider-zia/pull/127) Updated provider to [zscaler-sdk-go v0.0.10](https://github.com/zscaler/zscaler-sdk-go/releases/tag/v0.0.10)
- [PR #127](https://github.com/zscaler/terraform-provider-zia/pull/127) zia_user_management group attribute to hold a list of group IDs as a typeList instead of typeSet.

## 2.2.0

### Notes

- Release date: **(August 19 2022)**
- Supported Terraform version: **v1.x**

- [PR #113](https://github.com/zscaler/terraform-provider-zia/pull/113) Integrated newly created Zscaler GO SDK. Models are now centralized in the repository [zscaler-sdk-go](https://github.com/zscaler/zscaler-sdk-go)

### Fixes

- Terraform import failing for zia_traffic_forwarding_static_ip resource. Search by IP criteria was not implemented.

## 2.1.2

### Notes

- Release date: **(July 19 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #110](https://github.com/zscaler/terraform-provider-zia/pull/110) Added Terraform UserAgent for Backend API tracking

### Fixes

- [PR #111](https://github.com/zscaler/terraform-provider-zia/pull/111) Updated Import GPG key in goreleaser to [paultyng/ghaction-import-gpg](https://github.com/paultyng/ghaction-import-gpg)
- [PR #111](https://github.com/zscaler/terraform-provider-zia/pull/111) Updated golangci-lint to use golang 18

## 2.1.1

### Notes

- Release date: **(June, 7 2022)**
- Supported Terraform version: **v1.x**

- Fix: Fixed provider file to include resource and datasource hooks.

### New Features

- `zia_auth_settings_urls` Added new resource to support adding and removing URLs to ZIA exemption list.
- `zia_security_policy_settings` Added new resource to support adding and removing whitelisted and blacklisted URLs to the Advanced Threat Protection feature in ZIA.
  - Important: [API](https://community.zscaler.com/tags/api) limits apply based on the type of URLs being added. The ZIA API today allows: for 25K URL into the denylist and 255 into the allowlist. Refer to the [API](https://community.zscaler.com/tags/api) documentation [Here](https://help.zscaler.com/zia/api)

## 2.1.0

### Notes

- Release date: **(June, 7 2022)**
- Supported Terraform version: **v1.x**

### New Features

- `zia_auth_settings_urls` Added new resource to support adding and removing URLs to ZIA exemption list.
- `zia_security_policy_settings` Added new resource to support adding and removing whitelisted and blacklisted URLs to the Advanced Threat Protection feature in ZIA.

## 2.0.3

### Notes

- Release date: **(May, 18 2022)**
- Supported Terraform version: **v1.x**

### Announcement

The Terraform Provider for Zscaler Internet Access (ZIA) is now officially hosted under Zscaler's GitHub account and published in the Terraform Registry. For more details, visit the Zscaler Community Article [Here](https://community.zscaler.com/t/zpa-and-zia-terraform-providers-now-verified/16675)
Administrators who used previous versions of the provider, and followed instructions to install the binary as a custom provider, must update their provider block as such:

```hcl
terraform {
  required_providers {
    zia = {
      source = "zscaler/zia"
      version = "2.0.3"
    }
  }
}
provider "zia" {}

```

### New Data Sources

- ``zia_dlp_engines`` - [PR#91](https://github.com/zscaler/terraform-provider-zia/pull/91) 🔧

## 2.0.2

### Notes

- Release date: **(May, 17 2022)**
- Supported Terraform version: **v1.x**

### Announcement

The Terraform Provider for Zscaler Internet Access (ZIA) is now officially hosted under Zscaler's GitHub account and published in the Terraform Registry.
Administrators who used previous versions of the provider, and followed instructions to install the binary as a custom provider, must update their provider block as such:

```hcl
terraform {
  required_providers {
    zia = {
      source = "zscaler/zia"
      version = "2.0.3"
    }
  }
}
provider "zia" {}

```

### New Data Sources

- ``zia_dlp_engines`` - [PR#91](https://github.com/zscaler/terraform-provider-zia/pull/91) 🔧

## 2.0.1

### Notes

- Release date: **(April 17, 2022)**
- Supported Terraform version: **v1.x**

### Bug Fixes

Several schema type, expand and flattening function fixes were implemented to prevent undesired plan refresh updates and further provider optimization.

- ``zia_dlp_dictionaries`` - [PR#61](https://github.com/zscaler/terraform-provider-zia/pull/61) 🔧
- ``zia_dlp_web_rules`` - [PR#62](https://github.com/zscaler/terraform-provider-zia/pull/62) 🔧
- ``zia_firewall_filtering_rule`` - Added schema validation ``order`` parameter to ensure value is at least 1. [PR#63](https://github.com/zscaler/terraform-provider-zia/pull/63) 🔧
- ``zia_url_filtering_rules`` - [PR#66](https://github.com/zscaler/terraform-provider-zia/pull/66) 🔧
- ``zia_admin_users`` - [PR#67](https://github.com/zscaler/terraform-provider-zia/pull/67) 🔧
- ``zia_user_management`` - [PR#67](https://github.com/zscaler/terraform-provider-zia/pull/67) 🔧

### Enhancements

1. Updated ZIA API client to validate the corresponding Zscaler cloud name. The previous environment variable ``ZIA_BASE_URL`` was replaced with ``ZIA_CLOUD``. [PR#58](https://github.com/zscaler/terraform-provider-zia/pull/58)

2. The provider now validates the proper Zscaler cloud name. [PR#58](https://github.com/zscaler/terraform-provider-zia/pull/58) For instructions on how to find your Zscaler cloud name, refer to the following help article [Here](https://help.zscaler.com/zia/getting-started-zia-api#RetrieveAPIKey)

3. Added and fixed multiple acceptance tests to easily and routinely verify that Terraform Plugins produce the expected outcome

4. Updated GitHub Actions CI to include both build and acceptance test workflow

5. Added new optimized acceptance tests - [PR#71](https://github.com/zscaler/terraform-provider-zia/pull/71) 🔧

## 2.0.0

### Notes

- Release date: **(February 9, 2022)**
- Supported Terraform version: **v1.x**

### New Resources and DataSources

The ZIA cloud service API  now includes new endpoints in order to fully support Data Loss Prevention (DLP) rule creation and updates. The following Terraform resources and data source have been added:

DATA SOURCES:

- ``data_source_zia_device_group`` [PR#50](https://github.com/zscaler/terraform-provider-zia/pull/50) :rocket:
- ``data_source_zia_dlp_notification_templates``.[PR#53](https://github.com/zscaler/terraform-provider-zia/pull/53) :rocket:
- ``data_source_zia_dlp_web_rules``.[PR#53](https://github.com/zscaler/terraform-provider-zia/pull/53) :rocket:
- ``data_source_zia_dlp_engines``.[PR#53](https://github.com/zscaler/terraform-provider-zia/pull/53) :rocket:

RESOURCES:

- ``resource_zia_dlp_notification_templates``.[PR#53](https://github.com/zscaler/terraform-provider-zia/pull/53):rocket:
- ``resource_zia_dlp_web_rules``.[PR#53](https://github.com/zscaler/terraform-provider-zia/pull/53) :rocket:
- ``resource_zia_dlp_engines``.[PR#53](https://github.com/zscaler/terraform-provider-zia/pull/53) :rocket:

UPDATES:

- Added ``zia_device_groups`` to ``resource_zia_url_filtering_rules``.[PR#51](https://github.com/zscaler/terraform-provider-zia/pull/51) :rocket:

### New Acceptance Tests

- Added multiple acceptance tests to easily and routinely verify that Terraform Plugins produce the expected outcome. [PR#54](https://github.com/zscaler/terraform-provider-zia/pull/51)
- Added GoRelease workflow to GitHub Actions CI/CD for automatic software release.

## 1.0.3

### Notes

- Release date: **(December 28, 2021)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- Fixed issue where Terraform showed that resources had been modified even though nothing had been changed in the upstream resources. [PR#45](https://github.com/zscaler/terraform-provider-zia/pull/45) 🔧

### Enhacements

- Added multiple validators across several resources for better API abstraction and mistake prevention during `terraform apply` [PR#46](https://github.com/zscaler/terraform-provider-zia/pull/46) :rocket:

- The provider now supports the ability to import resources via its `name` and/or `id` property to support easier migration of existing ZIA resources via `terraform import` command.
The  following resources are supported:
  - resource_zia_admin_users - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47)] :rocket:
  - resource_zia_dlp_dictionaries - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_firewall_filtering_rules - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_fw_filtering_ip_destination_groups - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_fw_filtering_ip_source_groups - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_fw_filtering_network_application_groups - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_fw_filtering_network_services_groups - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_fw_filtering_network_services - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_location_management - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_url_categories - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_url_filtering_rules - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:
  - resource_zia_user_management_users - [PR#47](https://github.com/zscaler/terraform-provider-zia/pull/47) :rocket:

## 1.0.2

### Notes

- Release date: **(November 29, 2021)**
- Supported Terraform version: **v1.x**

### Bug Fixes

- VPN Credentials: Fixed issue where when creating a VPN credential and `type` was set to `IP`, the field `ip_address` was being returned as a non-expected argument. The issue was addressed on [PR#36](https://github.com/zscaler/terraform-provider-zia/pull/36)

- VPN Credentials: Fixed issue where when creating VPN credential and `type` was set to `UFQDN`, the parameter was not being validated if it was empty. The issue was addressed on [PR#36](https://github.com/zscaler/terraform-provider-zia/pull/36)

- VPN Credentials: Removed unsupported VPN Credential types `CN` and `XAUTH`. The issue was addressed on [PR#36](https://github.com/zscaler/terraform-provider-zia/pull/36)

- Location Management: Fixed issue where when creating a sub-location and the `ip_addresses` field was empty or the value was not a valid IPv4 address r IPv4 range, the provider pushed partial configuration and then exited with failure. The new validation function, will check if the `parent_id` has been set to a value greater than `0` and if the `ip_addresses` parameter has been fullfilled. The issue was addressed on [PR#37](https://github.com/zscaler/terraform-provider-zia/pull/37)

### Enhacements

- Static IP: Added ``ForceNew`` option to ``ip_address`` in the schema, so the resource will be destroyed and recreated [PR#40](https://github.com/zscaler/terraform-provider-zia/pull/40)

- VPN Credentials: Added ``ForceNew`` option to ``type`` in the schema, so the resource will be destroyed and recreated if the type of the VPN resource needs to be changed from ``IP`` to ``UFQDN`` and vice-versa [PR#41](https://github.com/zscaler/terraform-provider-zia/pull/41)

# 1.0.0

### Notes

- Release date: **(November, 12 2021)**
- Supported Terraform version: **v1.x**

### Initial Release

### RESOURCE FEATURES

- New Resource: resource_zia_admin_users 🆕
- New Resource: resource_zia_dlp_dictionaries 🆕
- New Resource: resource_zia_firewall_filtering_rules 🆕
- New Resource: resource_zia_fw_filtering_ip_destination_groups 🆕
- New Resource: resource_zia_fw_filtering_ip_source_groups 🆕
- New Resource: resource_zia_fw_filtering_network_application_groups 🆕
- New Resource: resource_zia_fw_filtering_network_services_groups 🆕
- New Resource: resource_zia_fw_filtering_network_services 🆕
- New Resource: resource_zia_location_management 🆕
- New Resource: resource_zia_traffic_forwarding_gre_tunnels 🆕
- New Resource: resource_zia_traffic_forwarding_static_ips 🆕
- New Resource: resource_zia_traffic_forwarding_vpn_credentials 🆕
- New Resource: resource_zia_url_categories 🆕
- New Resource: resource_zia_url_filtering_rules 🆕
- New Resource: resource_zia_user_management_users 🆕

### DATA SOURCE FEATURES

- New Data Source: data_source_zia_admin_roles 🆕
- New Data Source: data_source_zia_admin_users 🆕
- New Data Source: data_source_zia_dlp_dictionaries 🆕
- New Data Source: data_source_zia_firewall_filtering_rules 🆕
- New Data Source: data_source_zia_fw_filtering_ip_destination_groups 🆕
- New Data Source: data_source_zia_fw_filtering_ip_source_groups 🆕
- New Data Source: data_source_zia_fw_filtering_network_application_groups 🆕
- New Data Source: data_source_zia_fw_filtering_network_application 🆕
- New Data Source: data_source_zia_fw_filtering_network_service_groups 🆕
- New Data Source: data_source_zia_fw_filtering_network_services 🆕
- New Data Source: data_source_zia_fw_filtering_time_window 🆕
- New Data Source: data_source_zia_location_groups 🆕
- New Data Source: data_source_zia_location_management 🆕
- New Data Source: data_source_zia_traffic_forwarding_gre_internal_ranges 🆕
- New Data Source: data_source_zia_traffic_forwarding_gre_tunnels 🆕
- New Data Source: data_source_zia_traffic_forwarding_gre_vip_recommended_list 🆕
- New Data Source: data_source_zia_traffic_forwarding_ip_gre_tunnel_info 🆕
- New Data Source: data_source_zia_traffic_forwarding_public_nodes_vips 🆕
- New Data Source: data_source_zia_traffic_forwarding_static_ips 🆕
- New Data Source: data_source_zia_traffic_forwarding_vpn_credentials 🆕
- New Data Source: data_source_zia_url_categories 🆕
- New Data Source: data_source_zia_url_filtering_rules 🆕
- New Data Source: data_source_zia_user_management_departments 🆕
- New Data Source: data_source_zia_user_management_groups 🆕
- New Data Source: data_source_zia_user_management_users 🆕