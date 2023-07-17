# Changelog

## 2.6.0 (July, xx 2023)

### Notes

- Release date: **(July, xx 2023)**
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

- Updated to Zscaler-SDK-GO v1.5.5. The update improves search mechanisms for both ZIA and ZPA resources, to ensure streamline upstream GET API requests and responses using ``search`` parameter. Notice that not all current API endpoints support the search parameter, in which case, all resources will be returned.

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

## 2.5.0 (March, 27 2023)

### Notes

- Release date: **(March, 27 2023)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #202](https://github.com/zscaler/terraform-provider-zia/pull/202) ``zia_user_management``: Implemented new attribute ``auth_methods``. The attribute supports the following values: ``BASIC`` and/or ``DIGEST``.
- ``zia_location_management``: Implemented new attribute ``basic_auth_enabled``. The supported values are: ``true`` or ``false``

- [PR #202](https://github.com/zscaler/terraform-provider-zia/pull/202) The provider now supports authentication to Zscaler ``preview`` and ``zscalerten`` clouds.

- [PR #211](https://github.com/zscaler/terraform-provider-zia/pull/211) Added new datasource ``zia_location_lite``. This data source can be used to return the "Road Warrior" location, which can then be used in the following resources: ``zia_url_filtering_rules``, ``zia_firewall_filtering_rule`` and ``zia_dlp_web_rules``

- [PR #213](https://github.com/zscaler/terraform-provider-zia/pull/213) Added support to search for sub-location within the resource ``zia_location_management``

### Fixes

- [PR #212](https://github.com/zscaler/terraform-provider-zia/pull/212) ``zia_user_management``: Fixed flattening function to expand group attribute values. Issue [#205](https://github.com/zscaler/terraform-provider-zia/issues/205)

- [PR #214](https://github.com/zscaler/terraform-provider-zia/pull/214) ``zia_traffic_forwarding_gre_tunnel``: Fixed issue while creating GRE Tunnels. Issue #208

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

### Fixes

- [PR #167](https://github.com/zscaler/terraform-provider-zia/pull/167) Published provider as v2 go-module

## 2.3.2 (December, 30 2022)

### Notes

- Release date: **(December, 30 2022)**
- Supported Terraform version: **v1.x**

### Fixes

- [PR #164](https://github.com/zscaler/terraform-provider-zia/pull/164) Added missing URL Category resource parameters
- [PR #165](https://github.com/zscaler/terraform-provider-zia/pull/162) Added missing URL Category to ``zia_url_filtering_rule``

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

- [PR #127](https://github.com/zscaler/terraform-provider-zia/pull/127) Updated provider to zscaler-go-sdk v0.0.10
- [PR #127](https://github.com/zscaler/terraform-provider-zia/pull/127) zia_user_management group attribute to hold a list of group IDs as a typeList instead of typeSet.

## 2.2.0 (August, 19 2022)

### Notes

- Release date: **(August 19 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #113](https://github.com/zscaler/terraform-provider-zia/pull/113) Integrated newly created Zscaler GO SDK. Models are now centralized in the repository [zscaler-sdk-go](https://github.com/zscaler/zscaler-sdk-go)

### Fixes

- Terraform import failing for zia_traffic_forwarding_static_ip resource. Search by IP criteria was not implemented.

## 2.1.2 (June, 19 2022)

### Notes

- Release date: **(July 19 2022)**
- Supported Terraform version: **v1.x**

### Ehancements

- [PR #110](https://github.com/zscaler/terraform-provider-zia/pull/110) Added Terraform UserAgent for Backend API tracking

### Fixes

- [PR #111](https://github.com/zscaler/terraform-provider-zia/pull/111) Updated Import GPG key in goreleaser to [paultyng/ghaction-import-gpg](https://github.com/paultyng/ghaction-import-gpg)
- [PR #111](https://github.com/zscaler/terraform-provider-zia/pull/111) Updated golangci-lint to use golang 18

## 2.1.1 (June, 7 2022)

### Notes

- Supported Terraform version: **v1.x**

- Fix: Fixed provider file to include resource and datasource hooks.

## New Features

- `zia_auth_settings_urls` Added new resource to support adding and removing URLs to ZIA exemption list.
- `zia_security_policy_settings` Added new resource to support adding and removing whitelisted and blacklisted URLs to the Advanced Threat Protection feature in ZIA.

## 2.1.0 (June, 7 2022)

### Notes

- Supported Terraform version: **v1.x**

## New Features

- `zia_auth_settings_urls` Added new resource to support adding and removing URLs to ZIA exemption list.
- `zia_security_policy_settings` Added new resource to support adding and removing whitelisted and blacklisted URLs to the Advanced Threat Protection feature in ZIA.

# 2.0.3 (May, 18 2022)

## Notes

- Supported Terraform version: **v1.x**

## Announcement

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

## New Data Sources

- ``zia_dlp_engines`` - [PR#91](https://github.com/zscaler/terraform-provider-zia/pull/91) 🔧

## 2.0.1 (April 17, 2022)

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

## 2.0.0 (February 9, 2022)

## New Resources and DataSources

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

## New Acceptance Tests

- Added multiple acceptance tests to easily and routinely verify that Terraform Plugins produce the expected outcome. [PR#54](https://github.com/zscaler/terraform-provider-zpa/pull/51)
- Added GoRelease workflow to GitHub Actions CI/CD for automatic software release.

## 1.0.3 (December 28, 2021)

## Bug Fixes

- Fixed issue where Terraform showed that resources had been modified even though nothing had been changed in the upstream resources. [PR#45](https://github.com/zscaler/terraform-provider-zia/pull/45) 🔧

## Enhacements

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

## 1.0.2 (November 29, 2021)

## Bug Fixes

- VPN Credentials: Fixed issue where when creating a VPN credential and `type` was set to `IP`, the field `ip_address` was being returned as a non-expected argument. The issue was addressed on [PR#36](https://github.com/zscaler/terraform-provider-zia/pull/36)

- VPN Credentials: Fixed issue where when creating VPN credential and `type` was set to `UFQDN`, the parameter was not being validated if it was empty. The issue was addressed on [PR#36](https://github.com/zscaler/terraform-provider-zia/pull/36)

- VPN Credentials: Removed unsupported VPN Credential types `CN` and `XAUTH`. The issue was addressed on [PR#36](https://github.com/zscaler/terraform-provider-zia/pull/36)

- Location Management: Fixed issue where when creating a sub-location and the `ip_addresses` field was empty or the value was not a valid IPv4 address r IPv4 range, the provider pushed partial configuration and then exited with failure. The new validation function, will check if the `parent_id` has been set to a value greater than `0` and if the `ip_addresses` parameter has been fullfilled. The issue was addressed on [PR#37](https://github.com/zscaler/terraform-provider-zia/pull/37)

## Enhacements

- Static IP: Added ``ForceNew`` option to ``ip_address`` in the schema, so the resource will be destroyed and recreated [PR#40](https://github.com/zscaler/terraform-provider-zia/pull/40)

- VPN Credentials: Added ``ForceNew`` option to ``type`` in the schema, so the resource will be destroyed and recreated if the type of the VPN resource needs to be changed from ``IP`` to ``UFQDN`` and vice-versa [PR#41](https://github.com/zscaler/terraform-provider-zia/pull/41)
