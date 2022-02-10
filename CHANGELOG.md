# Changelog

## 2.0.0 (February 9, 2022)

## New Resources and DataSources

The ZIA cloud service API  now includes new endpoints in order to fully support Data Loss Prevention (DLP) rule creation and updates. The following Terraform resources and data source have been added:

DATA SOURCES:

- ``data_source_zia_device_group`` [PR#50](https://github.com/willguibr/terraform-provider-zpa/pull/50) :rocket:
- ``data_source_zia_dlp_notification_templates``.[PR#53](https://github.com/willguibr/terraform-provider-zpa/pull/53) :rocket:
- ``data_source_zia_dlp_web_rules``.[PR#53](https://github.com/willguibr/terraform-provider-zpa/pull/53) :rocket:
- ``data_source_zia_dlp_engines``.[PR#53](https://github.com/willguibr/terraform-provider-zpa/pull/53) :rocket:

RESOURCES:

- ``resource_zia_dlp_notification_templates``.[PR#53](https://github.com/willguibr/terraform-provider-zpa/pull/53):rocket:
- ``resource_zia_dlp_web_rules``.[PR#53](https://github.com/willguibr/terraform-provider-zpa/pull/53) :rocket:
- ``resource_zia_dlp_engines``.[PR#53](https://github.com/willguibr/terraform-provider-zpa/pull/53) :rocket:

UPDATES:

- Added ``zia_device_groups`` to ``resource_zia_url_filtering_rules``.[PR#51](https://github.com/willguibr/terraform-provider-zpa/pull/51) :rocket:

## New Acceptance Tests

- Added multiple acceptance tests to easily and routinely verify that Terraform Plugins produce the expected outcome. [PR#51](https://github.com/willguibr/terraform-provider-zpa/pull/51)
- Added GoRelease workflow to GitHub Actions CI/CD for automatic software release.

## 1.0.3 (December 28, 2021)

## Bug Fixes

- Fixed issue where Terraform showed that resources had been modified even though nothing had been changed in the upstream resources. [PR#45](https://github.com/willguibr/terraform-provider-zia/pull/45) 🔧

## Enhacements

- Added multiple validators across several resources for better API abstraction and mistake prevention during `terraform apply` [PR#46](https://github.com/willguibr/terraform-provider-zia/pull/46) :rocket:

- The provider now supports the ability to import resources via its `name` and/or `id` property to support easier migration of existing ZIA resources via `terraform import` command.
The  following resources are supported:
    - resource_zia_admin_users - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47)] :rocket:
    - resource_zia_dlp_dictionaries - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_firewall_filtering_rules - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_fw_filtering_ip_destination_groups - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_fw_filtering_ip_source_groups - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_fw_filtering_network_application_groups - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_fw_filtering_network_services_groups - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_fw_filtering_network_services - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_location_management - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_url_categories - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_url_filtering_rules - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:
    - resource_zia_user_management_users - [PR#47](https://github.com/willguibr/terraform-provider-zia/pull/47) :rocket:

## 1.0.2 (November 29, 2021)

## Bug Fixes

- VPN Credentials: Fixed issue where when creating a VPN credential and `type` was set to `IP`, the field `ip_address` was being returned as a non-expected argument. The issue was addressed on [PR#36](https://github.com/willguibr/terraform-provider-zia/pull/36)

- VPN Credentials: Fixed issue where when creating VPN credential and `type` was set to `UFQDN`, the parameter was not being validated if it was empty. The issue was addressed on [PR#36](https://github.com/willguibr/terraform-provider-zia/pull/36)

- VPN Credentials: Removed unsupported VPN Credential types `CN` and `XAUTH`. The issue was addressed on [PR#36](https://github.com/willguibr/terraform-provider-zia/pull/36)

- Location Management: Fixed issue where when creating a sub-location and the `ip_addresses` field was empty or the value was not a valid IPv4 address r IPv4 range, the provider pushed partial configuration and then exited with failure. The new validation function, will check if the `parent_id` has been set to a value greater than `0` and if the `ip_addresses` parameter has been fullfilled. The issue was addressed on [PR#37](https://github.com/willguibr/terraform-provider-zia/pull/37)

## Enhacements

- Static IP: Added ``ForceNew`` option to ``ip_address`` in the schema, so the resource will be destroyed and recreated [PR#40](https://github.com/willguibr/terraform-provider-zia/pull/40)

- VPN Credentials: Added ``ForceNew`` option to ``type`` in the schema, so the resource will be destroyed and recreated if the type of the VPN resource needs to be changed from ``IP`` to ``UFQDN`` and vice-versa [PR#41](https://github.com/willguibr/terraform-provider-zia/pull/41)
