package zia

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/advancedthreatsettings"
)

func dataSourceAdvancedThreatSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedThreatSettingsRead,
		Schema: map[string]*schema.Schema{
			"risk_tolerance": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Page Risk tolerance index set between 0 and 100 (100 being the highest risk).",
			},
			"risk_tolerance_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspicious web pages",
			},
			"cmd_ctl_server_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether connections to known Command & Control (C2) Servers are allowed or blocked",
			},
			"cmd_ctl_server_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for connections to known C2 servers",
			},
			"cmd_ctl_traffic_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether botnets are allowed or blocked from sending or receiving commands to unknown servers",
			},
			"cmd_ctl_traffic_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for botnets",
			},
			"malware_sites_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether known malicious sites and content are allowed or blocked",
			},
			"malware_sites_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for malicious sites",
			},
			"activex_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether sites are allowed or blocked from accessing vulnerable ActiveX controls that are known to have been exploited.",
			},
			"activex_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for ActiveX controls",
			},
			"browser_exploits_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether known web browser vulnerabilities prone to exploitation are allowed or blocked.",
			},
			"browser_exploits_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for browser exploits",
			},
			"file_format_vunerabilites_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether known file format vulnerabilities and suspicious or malicious content in Microsoft Office or PDF documents are allowed or blocked",
			},
			"file_format_vunerabilites_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for file format vulnerabilities",
			},
			"known_phishing_sites_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether known phishing sites are allowed or blocked",
			},
			"known_phishing_sites_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for known phishing sites",
			},
			"suspected_phishing_sites_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block suspected phishing sites identified through heuristic detection.",
			},
			"suspected_phishing_sites_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspected phishing sites",
			},
			"suspect_adware_spyware_sites_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block any detections of communication and callback traffic associated with spyware agents and data transmission",
			},
			"suspect_adware_spyware_sites_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for suspected adware and spyware sites",
			},
			"web_spam_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block web pages that pretend to contain useful information, to get higher ranking in search engine results or drive traffic to phishing, adware, or spyware distribution sites.",
			},
			"web_spam_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for web spam",
			},
			"irc_tunnelling_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for web spam",
			},
			"irc_tunnelling_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for IRC tunnels",
			},
			"anonymizer_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block applications and methods used to obscure the destination and the content accessed by the user, therefore blocking traffic to anonymizing web proxies",
			},
			"anonymizer_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for anonymizers",
			},
			"cookie_stealing_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block third-party websites that gather cookie information",
			},
			"cookie_stealing_pcap_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for cookie stealing",
			},
			"potential_malicious_requests_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block this type of cross-site scripting (XSS)",
			},
			"potential_malicious_requests_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for (XSS) attacks",
			},
			"blocked_countries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"block_countries_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for blocked countries",
			},
			"bit_torrent_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for blocked countries",
			},
			"bit_torrent_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for BitTorrent",
			},
			"tor_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block the usage of Tor, a popular P2P anonymizer protocol with support for encryption.",
			},
			"tor_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for Tor",
			},
			"google_talk_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block access to Google Hangouts, a popular P2P VoIP application.",
			},
			"google_talk_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for Google Hangouts",
			},
			"ssh_tunnelling_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block SSH traffic being tunneled over HTTP/Ss",
			},
			"ssh_tunnelling_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for SSH tunnels",
			},
			"crypto_mining_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block cryptocurrency mining network traffic and script",
			},
			"crypto_mining_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for cryptomining",
			},
			"ad_spyware_sites_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block websites known to contain adware or spyware that displays malicious advertisements that can collect users' information without their knowledge",
			},
			"ad_spyware_sites_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for adware and spyware sites",
			},
			"dga_domains_blocked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to allow or block domains that are suspected to be generated using domain generation algorithms (DGA)",
			},
			"alert_for_unknown_suspicious_c2_traffic": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether to send alerts upon detecting unknown or suspicious C2 traffic",
			},
			"dga_domains_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for DGA domains",
			},
			"malicious_urls_capture": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "A Boolean value specifying whether packet capture (PCAP) is enabled or not for malicious URLs",
			},
		},
	}
}

func dataSourceAdvancedThreatSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zClient := meta.(*Client)
	service := zClient.Service

	resp, err := advancedthreatsettings.GetAdvancedThreatSettings(ctx, service)
	if err != nil {
		return nil
	}

	if resp != nil {
		d.SetId("advanced_threat_settings")
		_ = d.Set("risk_tolerance", resp.RiskTolerance)
		_ = d.Set("risk_tolerance_capture", resp.RiskToleranceCapture)
		_ = d.Set("cmd_ctl_server_blocked", resp.CmdCtlServerBlocked)
		_ = d.Set("cmd_ctl_server_capture", resp.CmdCtlServerCapture)
		_ = d.Set("cmd_ctl_traffic_blocked", resp.CmdCtlTrafficBlocked)
		_ = d.Set("cmd_ctl_traffic_capture", resp.CmdCtlTrafficCapture)
		_ = d.Set("malware_sites_blocked", resp.MalwareSitesBlocked)
		_ = d.Set("malware_sites_capture", resp.MalwareSitesCapture)
		_ = d.Set("activex_blocked", resp.ActiveXBlocked)
		_ = d.Set("activex_capture", resp.ActiveXCapture)
		_ = d.Set("browser_exploits_blocked", resp.BrowserExploitsBlocked)
		_ = d.Set("browser_exploits_capture", resp.BrowserExploitsCapture)
		_ = d.Set("file_format_vunerabilites_blocked", resp.FileFormatVulnerabilitiesBlocked)
		_ = d.Set("file_format_vunerabilites_capture", resp.FileFormatVulnerabilitiesCapture)
		_ = d.Set("known_phishing_sites_blocked", resp.KnownPhishingSitesBlocked)
		_ = d.Set("known_phishing_sites_capture", resp.KnownPhishingSitesCapture)
		_ = d.Set("suspected_phishing_sites_blocked", resp.SuspectedPhishingSitesBlocked)
		_ = d.Set("suspected_phishing_sites_capture", resp.SuspectedPhishingSitesCapture)
		_ = d.Set("suspect_adware_spyware_sites_blocked", resp.SuspectAdwareSpywareSitesBlocked)
		_ = d.Set("suspect_adware_spyware_sites_capture", resp.SuspectAdwareSpywareSitesCapture)
		_ = d.Set("web_spam_blocked", resp.WebspamBlocked)
		_ = d.Set("web_spam_capture", resp.WebspamCapture)
		_ = d.Set("irc_tunnelling_blocked", resp.IrcTunnellingBlocked)
		_ = d.Set("irc_tunnelling_capture", resp.IrcTunnellingCapture)
		_ = d.Set("anonymizer_blocked", resp.AnonymizerBlocked)
		_ = d.Set("anonymizer_capture", resp.AnonymizerCapture)
		_ = d.Set("cookie_stealing_blocked", resp.CookieStealingBlocked)
		_ = d.Set("cookie_stealing_pcap_enabled", resp.CookieStealingPCAPEnabled)
		_ = d.Set("potential_malicious_requests_blocked", resp.PotentialMaliciousRequestsBlocked)
		_ = d.Set("potential_malicious_requests_capture", resp.PotentialMaliciousRequestsCapture)
		_ = d.Set("blocked_countries", resp.BlockedCountries)
		_ = d.Set("block_countries_capture", resp.BlockCountriesCapture)
		_ = d.Set("bit_torrent_blocked", resp.BitTorrentBlocked)
		_ = d.Set("bit_torrent_capture", resp.BitTorrentCapture)
		_ = d.Set("tor_blocked", resp.TorBlocked)
		_ = d.Set("tor_capture", resp.TorCapture)
		_ = d.Set("google_talk_blocked", resp.GoogleTalkBlocked)
		_ = d.Set("google_talk_capture", resp.GoogleTalkCapture)
		_ = d.Set("ssh_tunnelling_blocked", resp.SshTunnellingBlocked)
		_ = d.Set("ssh_tunnelling_capture", resp.SshTunnellingCapture)
		_ = d.Set("crypto_mining_blocked", resp.CryptoMiningBlocked)
		_ = d.Set("crypto_mining_capture", resp.CryptoMiningCapture)
		_ = d.Set("ad_spyware_sites_blocked", resp.AdSpywareSitesBlocked)
		_ = d.Set("ad_spyware_sites_capture", resp.AdSpywareSitesCapture)
		_ = d.Set("dga_domains_blocked", resp.DgaDomainsBlocked)
		_ = d.Set("dga_domains_capture", resp.DgaDomainsCapture)
		_ = d.Set("alert_for_unknown_suspicious_c2_traffic", resp.AlertForUnknownOrSuspiciousC2Traffic)
		_ = d.Set("malicious_urls_capture", resp.MaliciousUrlsCapture)

	} else {
		return diag.FromErr(fmt.Errorf("couldn't read advanced threat protection setting options"))
	}

	return nil
}
