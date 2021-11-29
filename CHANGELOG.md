## 1.0.1 (November 29, 2021)

## Bug Fixes

- VPN Credentials: Fixed issue where when creating a VPN credential and `type` was set to `IP`, the field `ip_address` was being returned as a non-expected argument. The issue was addressed on [PR#36](https://github.com/willguibr/terraform-provider-zia/pull/36)

- VPN Credentials: Fixed issue where when creating VPN credential and `type` was set to `UFQDN`, the parameter was not being validated if it was empty. The issue was addressed on [PR#36](https://github.com/willguibr/terraform-provider-zia/pull/36)

- VPN Credentials: Remoted unsupported VPN Credential types `CN` and `XAUTH`. The issue was addressed on [PR#36](https://github.com/willguibr/terraform-provider-zia/pull/36)

- Location Management: Fixed issue where when creating a sub-location and the `ip_addresses` field was empty or the value was not a valid IPv4 address r IPv4 range, the provider pushed partial configuration and then exited with failure. The new validation function, will check if the `parent_id` has been set to a value greater than `0` and if the `ip_addresses` parameter has been fullfilled. The issue was addressed on [PR#37](https://github.com/willguibr/terraform-provider-zia/pull/37)
