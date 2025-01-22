package zia

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/advancedthreatsettings"
)

func TestAccResourceAdvancedThreatSettings_basic(t *testing.T) {
	resourceName := "zia_advanced_threat_settings.test"

	testAccPreCheck(t)

	// Initial configuration: risk_tolerance = 50, all booleans = false
	initialConfig := testAccResourceAdvancedThreatSettingsConfig(
		50,
		false, false, false, false, false,
		false, false, false, false, false,
		false, false, false, false, false,
		false, false, false, false, false,
		false, false,
	)

	// Updated configuration: risk_tolerance = 75, all booleans = true
	updatedConfig := testAccResourceAdvancedThreatSettingsConfig(
		75,
		true, true, true, true, true,
		true, true, true, true, true,
		true, true, true, true, true,
		true, true, true, true, true,
		true, false,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAdvancedThreatSettingsDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create resource with initial config
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedThreatSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "risk_tolerance", "50"),
					resource.TestCheckResourceAttr(resourceName, "cmd_ctl_server_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "cmd_ctl_traffic_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "malware_sites_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "activex_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "browser_exploits_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "file_format_vunerabilites_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "known_phishing_sites_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "suspected_phishing_sites_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "suspect_adware_spyware_sites_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "web_spam_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "irc_tunnelling_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "anonymizer_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "cookie_stealing_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "potential_malicious_requests_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "bit_torrent_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "tor_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "google_talk_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssh_tunnelling_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "crypto_mining_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "ad_spyware_sites_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "dga_domains_blocked", "false"),
					resource.TestCheckResourceAttr(resourceName, "alert_for_unknown_suspicious_c2_traffic", "false"),
				),
			},
			// Step 2: Update resource with updated config
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdvancedThreatSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "risk_tolerance", "75"),
					resource.TestCheckResourceAttr(resourceName, "cmd_ctl_server_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "cmd_ctl_traffic_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "malware_sites_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "activex_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "browser_exploits_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "file_format_vunerabilites_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "known_phishing_sites_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "suspected_phishing_sites_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "suspect_adware_spyware_sites_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "web_spam_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "irc_tunnelling_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "anonymizer_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "cookie_stealing_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "potential_malicious_requests_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "bit_torrent_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "tor_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "google_talk_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssh_tunnelling_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "crypto_mining_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "ad_spyware_sites_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "dga_domains_blocked", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_for_unknown_suspicious_c2_traffic", "false"),
				),
			},
			// Step 3: Import the resource and verify
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAdvancedThreatSettingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No advanced_threat_settings ID is set")
		}

		client := testAccProvider.Meta().(*Client)
		service := client.Service

		resp, err := advancedthreatsettings.GetAdvancedThreatSettings(context.Background(), service)
		if err != nil {
			return fmt.Errorf("error getting advanced threat settings: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("advanced threat settings response is nil")
		}

		log.Printf("[INFO] Advanced Threat Settings found: %s", rs.Primary.ID)
		return nil
	}
}

func testAccCheckAdvancedThreatSettingsDestroy(s *terraform.State) error {
	// For global settings, typically there's no "destroy", so just return nil.
	return nil
}

// Updated function signature and body to accept all arguments and set all attributes.
func testAccResourceAdvancedThreatSettingsConfig(
	riskTolerance int,
	cmdCtlServerBlocked bool,
	cmdCtlTrafficBlocked bool,
	malwareSitesBlocked bool,
	activexBlocked bool,
	browserExploitsBlocked bool,
	fileFormatVunerabilitesBlocked bool,
	knownPhishingSitesBlocked bool,
	suspectedPhishingSitesBlocked bool,
	suspectAdwareSpywareSitesBlocked bool,
	webSpamBlocked bool,
	ircTunnellingBlocked bool,
	anonymizerBlocked bool,
	cookieStealingBlocked bool,
	potentialMaliciousRequestsBlocked bool,
	bitTorrentBlocked bool,
	torBlocked bool,
	googleTalkBlocked bool,
	sshTunnellingBlocked bool,
	cryptoMiningBlocked bool,
	adSpywareSitesBlocked bool,
	dgaDomainsBlocked bool,
	alertForUnknownSuspiciousC2Traffic bool,
) string {
	return fmt.Sprintf(`
resource "zia_advanced_threat_settings" "test" {
  risk_tolerance                        = %d
  cmd_ctl_server_blocked                = %t
  cmd_ctl_traffic_blocked               = %t
  malware_sites_blocked                 = %t
  activex_blocked                       = %t
  browser_exploits_blocked              = %t
  file_format_vunerabilites_blocked     = %t
  known_phishing_sites_blocked          = %t
  suspected_phishing_sites_blocked      = %t
  suspect_adware_spyware_sites_blocked  = %t
  web_spam_blocked                      = %t
  irc_tunnelling_blocked                = %t
  anonymizer_blocked                    = %t
  cookie_stealing_blocked               = %t
  potential_malicious_requests_blocked  = %t
  bit_torrent_blocked                   = %t
  tor_blocked                           = %t
  google_talk_blocked                   = %t
  ssh_tunnelling_blocked                = %t
  crypto_mining_blocked                 = %t
  ad_spyware_sites_blocked              = %t
  dga_domains_blocked                   = %t
  alert_for_unknown_suspicious_c2_traffic = %t
}
`,
		riskTolerance,
		cmdCtlServerBlocked,
		cmdCtlTrafficBlocked,
		malwareSitesBlocked,
		activexBlocked,
		browserExploitsBlocked,
		fileFormatVunerabilitesBlocked,
		knownPhishingSitesBlocked,
		suspectedPhishingSitesBlocked,
		suspectAdwareSpywareSitesBlocked,
		webSpamBlocked,
		ircTunnellingBlocked,
		anonymizerBlocked,
		cookieStealingBlocked,
		potentialMaliciousRequestsBlocked,
		bitTorrentBlocked,
		torBlocked,
		googleTalkBlocked,
		sshTunnellingBlocked,
		cryptoMiningBlocked,
		adSpywareSitesBlocked,
		dgaDomainsBlocked,
		alertForUnknownSuspiciousC2Traffic,
	)
}
