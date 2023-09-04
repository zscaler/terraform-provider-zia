package zia

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
)

var (
	testSdkClient            *Client
	testAccProvider          *schema.Provider
	testAccProviders         map[string]*schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"zia": testAccProvider,
	}

	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"zia": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

// TestMain overridden main testing function. Package level BeforeAll and AfterAll.
// It also delineates between acceptance tests and unit tests
func TestMain(m *testing.M) {
	// TF_VAR_hostname allows the real hostname to be scripted into the config tests
	// see examples/okta_resource_set/basic.tf
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s.%s.%s", os.Getenv("ZIA_USERNAME"), os.Getenv("ZIA_PASSWORD"), os.Getenv("ZIA_API_KEY"), os.Getenv("ZIA_CLOUD")))

	// NOTE: Acceptance test sweepers are necessary to prevent dangling
	// resources.
	// NOTE: Don't run sweepers if we are playing back VCR as nothing should be
	// going over the wire
	if os.Getenv("ZIA_VCR_TF_ACC") != "play" {
		setupSweeper(resourcetype.RuleLabels, sweepTestRuleLabels)
		setupSweeper(resourcetype.TrafficForwardingLocManagement, sweepTestLocationManagement)
		setupSweeper(resourcetype.TrafficForwardingGRETunnel, sweepTestGRETunnels)
		setupSweeper(resourcetype.TrafficForwardingStaticIP, sweepTestStaticIP)
		setupSweeper(resourcetype.TrafficForwardingVPNCredentials, sweepTestVPNCredentials)
		setupSweeper(resourcetype.FirewallFilteringRules, sweepTestFirewallFilteringRule)
		setupSweeper(resourcetype.FWFilteringSourceGroup, sweepTestSourceIPGroup)
		setupSweeper(resourcetype.FWFilteringDestinationGroup, sweepTestDestinationIPGroup)
		setupSweeper(resourcetype.FWFilteringNetworkServices, sweepTestNetworkServices)
		setupSweeper(resourcetype.FWFilteringNetworkServiceGroups, sweepTestNetworkServicesGroup)
		setupSweeper(resourcetype.FWFilteringNetworkAppGroups, sweepTestNetworkAppGroups)
		setupSweeper(resourcetype.URLFilteringRules, sweepTestURLFilteringRule)
		setupSweeper(resourcetype.DLPWebRules, sweepTestDLPWebRule)
		setupSweeper(resourcetype.DLPDictionaries, sweepTestDLPDictionary)
		setupSweeper(resourcetype.DLPEngines, sweepTestDLPEngines)
		setupSweeper(resourcetype.DLPNotificationTemplates, sweepTestDLPTemplates)
		setupSweeper(resourcetype.URLCategories, sweepTestURLCategories)
		setupSweeper(resourcetype.AdminUsers, sweepTestAdminUser)
		setupSweeper(resourcetype.Users, sweepTestUsers)
	}
	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	_ = Provider()
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		err := accPreCheck()
		if err != nil {
			t.Fatalf("%v", err)
		}
	}
}

// accPreCheck checks if the necessary environment variables for acceptance tests are set.
func accPreCheck() error {
	username := os.Getenv("ZIA_USERNAME")
	password := os.Getenv("ZIA_PASSWORD")
	apiKey := os.Getenv("ZIA_API_KEY")
	ziaCloud := os.Getenv("ZIA_CLOUD")

	// Check for the presence of necessary environment variables.
	if username == "" {
		return errors.New("ZIA_USERNAME must be set for acceptance tests")
	}

	if password == "" {
		return errors.New("ZIA_PASSWORD must be set for acceptance tests")
	}

	if apiKey == "" {
		return errors.New("ZIA_API_KEY must be set for acceptance tests")
	}

	if ziaCloud == "" {
		return errors.New("ZIA_CLOUD must be set for acceptance tests")
	}

	return nil
}

func sdkClientForTest() (*Client, error) {
	if testSdkClient == nil {
		sweeperLogger.Warn("testSdkClient is not initialized. Initializing now...")

		config := &Config{
			Username:   os.Getenv("ZIA_USERNAME"),
			Password:   os.Getenv("ZIA_PASSWORD"),
			APIKey:     os.Getenv("ZIA_API_KEY"),
			ZIABaseURL: os.Getenv("ZIA_CLOUD"),
			UserAgent:  "terraform-provider-zia",
		}

		var err error
		testSdkClient, err = config.Client()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize testSdkClient: %v", err)
		}
	}
	return testSdkClient, nil
}
