---
subcategory: "Advanced Settings"
layout: "zscaler"
page_title: "ZIA: advanced_settings"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-advanced-settings
  API documentation https://help.zscaler.com/zia/advanced-settings#/advancedSettings-get
  Retrieves information about the advanced settings configured in the ZIA Admin Portal
---

# zia_advanced_settings (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-advanced-settings)
* [API documentation](https://help.zscaler.com/zia/advanced-settings#/advancedSettings-get)

The **zia_advanced_settings** Retrieves information about the advanced settings configured in the ZIA Admin Portal. To learn more see [Configuring Advanced Settings](https://help.zscaler.com/zia/configuring-advanced-settings)

## Example Usage

```hcl
data "zia_advanced_settings" "this" {}
```

## Argument Reference

The following arguments are supported:

### Read-Only

* `auth_bypass_apps` - (Set of String) Cloud applications that are exempted from cookie authentication.
* `auth_bypass_urls` - (Set of String) Custom URLs that are exempted from cookie authentication for users.
* `basic_bypass_apps` - (Set of String) Cloud applications that are exempted from Basic authentication.
* `digest_auth_bypass_apps` - (Set of String) Cloud applications that are exempted from Digest authentication.
* `dns_resolution_on_transparent_proxy_exempt_apps` - (Set of String) Cloud applications that are excluded from DNS optimization on transparent proxy mode.
* `dns_resolution_on_transparent_proxy_ipv6_exempt_apps` - (Set of String) Cloud applications that are excluded from DNS optimization for IPv6 addresses on transparent proxy mode.
* `dns_resolution_on_transparent_proxy_apps` - (Set of String) Cloud applications to which DNS optimization on transparent proxy mode applies.
* `dns_resolution_on_transparent_proxy_ipv6_apps` - (Set of String) Cloud applications to which DNS optimization for IPv6 addresses on transparent proxy mode applies.
* `block_domain_fronting_apps` - (Set of String) Applications which are subjected to Domain Fronting.
* `prefer_sni_over_conn_host_apps` - (Set of String) Applications that are exempted from the preferSniOverConnHost setting (i.e., prefer SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections).
* `dns_resolution_on_transparent_proxy_exempt_url_categories` - (Set of String) URL categories that are excluded from DNS optimization on transparent proxy mode.
* `dns_resolution_on_transparent_proxy_ipv6_exempt_url_categories` - (Set of String) URL categories that are excluded from DNS optimization for IPv6 addresses on transparent proxy mode.
* `dns_resolution_on_transparent_proxy_url_categories` - (Set of String) URL categories to which DNS optimization on transparent proxy mode applies.
* `dns_resolution_on_transparent_proxy_ipv6_url_categories` - (Set of String) IPv6 URL categories to which DNS optimization on transparent proxy mode applies.
* `auth_bypass_url_categories` - (Set of String) URL categories that are exempted from cookie authentication.
* `domain_fronting_bypass_url_categories` - (Set of String) URL categories that are exempted from domain fronting.
* `kerberos_bypass_url_categories` - (Set of String) URL categories that are exempted from Kerberos authentication.
* `basic_bypass_url_categories` - (Set of String) URL categories that are exempted from Basic authentication.
* `http_range_header_remove_url_categories` - (Set of String) URL categories for which HTTP range headers must be removed.
* `digest_auth_bypass_url_categories` - (Set of String) URL categories that are exempted from Digest authentication.
* `sni_dns_optimization_bypass_url_categories` - (Set of String) URL categories that are excluded from the preferSniOverConnHost setting (i.e., prefer SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections).
* `kerberos_bypass_urls` - (Set of String) Custom URLs that are exempted from Kerberos authentication.
* `kerberos_bypass_apps` - (Set of String) Cloud applications that are exempted from Kerberos authentication.
* `digest_auth_bypass_urls` - (Set of String) Custom URLs that are exempted from Digest authentication.
* `dns_resolution_on_transparent_proxy_exempt_urls` - (Set of String) URLs that are excluded from DNS optimization on transparent proxy mode.
* `dns_resolution_on_transparent_proxy_urls` - (Set of String) URLs to which DNS optimization on transparent proxy mode applies.
* `enable_dns_resolution_on_transparent_proxy` - (Boolean) Value indicating whether DNS optimization is enabled or disabled for Z-Tunnel 2.0 and transparent proxy mode traffic (e.g., traffic via GRE or IPSec tunnels without a PAC file).
* `enable_ipv6_dns_resolution_on_transparent_proxy` - (Boolean) Value indicating whether DNS optimization is enabled or disabled for IPv6 connections to dual-stack or IPv6-only destinations sent via Z-Tunnel 2.0 and transparent proxy proxy mode (e.g., traffic via GRE or IPSec tunnels without a PAC file).
* `enable_ipv6_dns_optimization_on_all_transparent_proxy` - (Boolean) Enable/Disable DNS optimization for all IPv6 transparent proxy traffic.
* `enable_evaluate_policy_on_global_ssl_bypass` - (Boolean) Enable/Disable DNS optimization for all IPv6 transparent proxy traffic.
* `enable_office365` - (Boolean) Value indicating whether Microsoft Office 365 One Click Configuration is enabled or not.
* `log_internal_ip` - (Boolean) Value indicating whether to log internal IP address present in X-Forwarded-For (XFF) proxy header or not.
* `enforce_surrogate_ip_for_windows_app` - (Boolean) Enforce Surrogate IP authentication for Windows app traffic.
* `track_http_tunnel_on_http_ports` - (Boolean) Value indicating whether to apply configured policies on tunneled HTTP traffic sent via a CONNECT method request on port 80.
* `block_http_tunnel_on_non_http_ports` - (Boolean) Value indicating whether HTTP CONNECT method requests to non-standard ports are allowed or not (i.e., requests directed to ports other than the standard HTTP/S ports 80 and 443).
* `block_domain_fronting_on_host_header` - (Boolean) Value indicating whether to block or allow HTTP/S transactions in which the FQDN of the request URL is different than the FQDN of the request's host header.
* `zscaler_client_connector_1_and_pac_road_warrior_in_firewall` - (Boolean) Value indicating whether to apply the Firewall rules configured without a specified location criteria (or with the Road Warrior location) to remote user traffic forwarded via Z-Tunnel 1.0 or PAC files.
* `cascade_url_filtering` - (Boolean) Value indicating whether to apply the URL Filtering policy even when the Cloud App Control policy already allows a transaction explicitly.
* `enable_policy_for_unauthenticated_traffic` - (Boolean) Value indicating whether policies that include user and department criteria can be configured and applied for unauthenticated traffic.
* `block_non_compliant_http_request_on_http_ports` - (Boolean) Value indicating whether to allow or block traffic that is not compliant with RFC HTTP protocol standards.
* `enable_admin_rank_access` - (Boolean) Value indicating whether ranks are enabled for admins to allow admin ranks in policy configuration and management.
* `http2_nonbrowser_traffic_enabled` - (Boolean) Value indicating whether or not HTTP/2 should be the default web protocol for accessing various applications at your organizational level.
* `ecs_for_all_enabled` - (Boolean) Value indicating whether or not to include the ECS option in all DNS queries, originating from all locations and remote users.
* `dynamic_user_risk_enabled` - (Boolean) Value indicating whether to dynamically update user risk score by tracking risky user activities in real time.
* `block_connect_host_sni_mismatch` - (Boolean) Value indicating whether CONNECT host and SNI mismatch (i.e., CONNECT host doesn't match the SSL/TLS client hello SNI) is blocked or not.
* `prefer_sni_over_conn_host` - (Boolean) Value indicating whether or not to use the SSL/TLS client hello SNI for DNS resolution instead of the CONNECT host for forward proxy connections.
* `sipa_xff_header_enabled` - (Boolean) Value indicating whether or not to insert XFF header to all traffic forwarded from ZIA to ZPA, including source IP-anchored and ZIA-inspected ZPA application traffic.
* `block_non_http_on_http_port_enabled` - (Boolean) Value indicating whether non-HTTP Traffic on HTTP/S ports are allowed or blocked.
* `ui_session_timeout` - (Integer) Specifies the login session timeout for admins accessing the ZIA Admin Portal. Minimum value is 300 seconds.
