---
subcategory: "Advanced Threat Protection"
layout: "zscaler"
page_title: "ZIA: advanced_threat_settings"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-advanced-threat-protection-policy
  API documentation https://help.zscaler.com/zia/advanced-threat-protection-policy#/cyberThreatProtection/advancedThreatSettings-put
  Retrieves the advanced threat configuration settings in the ZIA Admin Portal
---

# zia_advanced_threat_settings (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-advanced-threat-protection-policy)
* [API documentation](https://help.zscaler.com/zia/advanced-threat-protection-policy#/)

Use the **zia_advanced_threat_settings** data source to retrieve the advanced threat configuration settings in the ZIA Admin Portal. To learn more see [Advanced Threat Protection](https://help.zscaler.com/unified/configuring-security-exceptions-advanced-threat-protection-policy)

## Example Usage

```hcl
data "zia_advanced_threat_settings" "this" {}
```

## Argument Reference

### Read-Only

* `risk_tolerance` - (int) The Page Risk tolerance index set between 0 and 100 (100 being the highest risk).
* `risk_tolerance_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspicious web pages.
* `cmd_ctl_server_blocked` - (bool) A Boolean value specifying whether connections to known Command & Control (C2) Servers are allowed or blocked.
* `cmd_ctl_server_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for connections to known C2 servers.
* `cmd_ctl_traffic_blocked` - (bool) A Boolean value specifying whether botnets are allowed or blocked from sending or receiving commands to unknown servers.
* `cmd_ctl_traffic_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for botnets.
* `malware_sites_blocked` - (bool) A Boolean value specifying whether known malicious sites and content are allowed or blocked.
* `malware_sites_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for malicious sites.
* `activex_blocked` - (bool) A Boolean value specifying whether sites are allowed or blocked from accessing vulnerable ActiveX controls that are known to have been exploited.
* `activex_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for ActiveX controls.
* `browser_exploits_blocked` - (bool) A Boolean value specifying whether known web browser vulnerabilities prone to exploitation are allowed or blocked.
* `browser_exploits_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for browser exploits.
* `file_format_vunerabilites_blocked` - (bool) A Boolean value specifying whether known file format vulnerabilities and suspicious or malicious content in Microsoft Office or PDF documents are allowed or blocked.
* `file_format_vunerabilites_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for file format vulnerabilities.
* `known_phishing_sites_blocked` - (bool) A Boolean value specifying whether known phishing sites are allowed or blocked.
* `known_phishing_sites_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for known phishing sites.
* `suspected_phishing_sites_blocked` - (bool) A Boolean value specifying whether to allow or block suspected phishing sites identified through heuristic detection.
* `suspected_phishing_sites_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspected phishing sites.
* `suspect_adware_spyware_sites_blocked` - (bool) A Boolean value specifying whether to allow or block any detections of communication and callback traffic associated with spyware agents and data transmission.
* `suspect_adware_spyware_sites_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspected adware and spyware sites.
* `web_spam_blocked` - (bool) A Boolean value specifying whether to allow or block web pages that pretend to contain useful information to get higher ranking in search engine results or drive traffic to phishing, adware, or spyware distribution sites.
* `web_spam_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for web spam.
* `irc_tunnelling_blocked` - (bool) A Boolean value specifying whether IRC tunneling is blocked.
* `irc_tunnelling_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for IRC tunnels.
* `anonymizer_blocked` - (bool) A Boolean value specifying whether to allow or block applications and methods used to obscure the destination and the content accessed by the user, therefore blocking traffic to anonymizing web proxies.
* `anonymizer_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for anonymizers.
* `cookie_stealing_blocked` - (bool) A Boolean value specifying whether to allow or block third-party websites that gather cookie information.
* `cookie_stealing_pcap_enabled` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for cookie stealing.
* `potential_malicious_requests_blocked` - (bool) A Boolean value specifying whether to allow or block this type of cross-site scripting (XSS).
* `potential_malicious_requests_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for XSS attacks.
* `blocked_countries` - (list of strings) List of ISO country codes specifying countries to block.

    **NOTE**: Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

* `block_countries_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for blocked countries.
* `bit_torrent_blocked` - (bool) A Boolean value specifying whether BitTorrent traffic is blocked.
* `bit_torrent_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for BitTorrent traffic.
* `tor_blocked` - (bool) A Boolean value specifying whether to allow or block the usage of Tor, a popular P2P anonymizer protocol with support for encryption.
* `tor_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for Tor.
* `google_talk_blocked` - (bool) A Boolean value specifying whether to allow or block access to Google Hangouts, a popular P2P VoIP application.
* `google_talk_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for Google Hangouts.
* `ssh_tunnelling_blocked` - (bool) A Boolean value specifying whether to allow or block SSH traffic being tunneled over HTTP/S.
* `ssh_tunnelling_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for SSH tunnels.
* `crypto_mining_blocked` - (bool) A Boolean value specifying whether to allow or block cryptocurrency mining network traffic and scripts.
* `crypto_mining_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for cryptomining.
* `ad_spyware_sites_blocked` - (bool) A Boolean value specifying whether to allow or block websites known to contain adware or spyware that displays malicious advertisements.
* `ad_spyware_sites_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for adware and spyware sites.
* `dga_domains_blocked` - (bool) A Boolean value specifying whether to allow or block domains that are suspected to be generated using domain generation algorithms (DGA).
* `alert_for_unknown_suspicious_c2_traffic` - (bool) A Boolean value specifying whether to send alerts upon detecting unknown or suspicious C2 traffic.
* `dga_domains_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for DGA domains.
* `malicious_urls_capture` - (bool) A Boolean value specifying whether packet capture (PCAP) is enabled or not for malicious URLs.

## Attribute Reference

N/A
