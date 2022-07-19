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
``Last updated: v2.1.2``

---

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

- ``data_source_zia_device_group`` [PR#50](https://github.com/zscaler/terraform-provider-zpa/pull/50) :rocket:
- ``data_source_zia_dlp_notification_templates``.[PR#53](https://github.com/zscaler/terraform-provider-zpa/pull/53) :rocket:
- ``data_source_zia_dlp_web_rules``.[PR#53](https://github.com/zscaler/terraform-provider-zpa/pull/53) :rocket:
- ``data_source_zia_dlp_engines``.[PR#53](https://github.com/zscaler/terraform-provider-zpa/pull/53) :rocket:

RESOURCES:

- ``resource_zia_dlp_notification_templates``.[PR#53](https://github.com/zscaler/terraform-provider-zpa/pull/53):rocket:
- ``resource_zia_dlp_web_rules``.[PR#53](https://github.com/zscaler/terraform-provider-zpa/pull/53) :rocket:
- ``resource_zia_dlp_engines``.[PR#53](https://github.com/zscaler/terraform-provider-zpa/pull/53) :rocket:

UPDATES:

- Added ``zia_device_groups`` to ``resource_zia_url_filtering_rules``.[PR#51](https://github.com/zscaler/terraform-provider-zpa/pull/51) :rocket:

### New Acceptance Tests

- Added multiple acceptance tests to easily and routinely verify that Terraform Plugins produce the expected outcome. [PR#54](https://github.com/zscaler/terraform-provider-zpa/pull/51)
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
