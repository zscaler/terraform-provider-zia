package zia

import (
	"errors"
	"fmt"
	"log"

	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

var (
	testSdkV3Client          *zscaler.Client
	testAccProvider          *schema.Provider
	testAccProviders         map[string]*schema.Provider
	testAccProviderFactories map[string]func() (*schema.Provider, error)
)

func init() {
	testAccProvider = ZIAProvider()
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
	os.Setenv("TF_VAR_hostname", fmt.Sprintf("%s.%s.%s", os.Getenv("ZSCALER_CLIENT_ID"), os.Getenv("ZSCALER_CLIENT_SECRET"), os.Getenv("ZSCALER_CLOUD")))

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
		setupSweeper(resourcetype.ForwardingControlRule, sweepTestForwardingControlRule)
		setupSweeper(resourcetype.CloudAppControlRule, sweepTestCloudAppControlRule)
		setupSweeper(resourcetype.FirewallFilteringRules, sweepTestFirewallFilteringRule)
		setupSweeper(resourcetype.FileTypeControlRules, sweepTestFileTypeControlRule)
		setupSweeper(resourcetype.FirewallIPSRules, sweepTestFirewallIPSRule)
		setupSweeper(resourcetype.FirewallDNSRules, sweepTestFirewallDNSRule)
		setupSweeper(resourcetype.SandboxRules, sweepTestSandboxRule)
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
	if err := ZIAProvider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
func TestProvider_impl(t *testing.T) {
	_ = ZIAProvider()
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
	// Check for mandatory environment variables for client_id + client_secret authentication
	if v := os.Getenv("ZSCALER_CLIENT_ID"); v == "" {
		return errors.New("ZSCALER_CLIENT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ZSCALER_CLIENT_SECRET"); v == "" {
		return errors.New("ZSCALER_CLIENT_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("ZSCALER_VANITY_DOMAIN"); v == "" {
		return errors.New("ZSCALER_VANITY_DOMAIN must be set for acceptance tests")
	}

	// Optional cloud configuration
	if v := os.Getenv("ZSCALER_CLOUD"); v == "" {
		log.Println("[INFO] ZSCALER_CLOUD is not set. Defaulting to production cloud.")
	}

	return nil
}

func TestProviderValidate(t *testing.T) {
	// Define environment variables related to ZIA authentication
	envKeys := []string{
		"ZSCALER_CLIENT_ID",
		"ZSCALER_CLIENT_SECRET",
		"ZSCALER_VANITY_DOMAIN",
		"ZSCALER_CLOUD",
	}
	envVals := make(map[string]string)

	// Save and clear ZIA env vars to test configuration cleanly
	for _, key := range envKeys {
		val := os.Getenv(key)
		if val != "" {
			envVals[key] = val
			os.Unsetenv(key)
		}
	}

	// Define test cases for the various configurations
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		vanityDomain string
		cloud        string // Optional field
		expectError  bool
	}{
		{"valid client_id + client_secret", "clientID", "clientSecret", "vanityDomain", "cloud", false},
		{"missing client_id", "", "clientSecret", "vanityDomain", "cloud", true},
		{"missing clientSecret", "clientID", "", "vanityDomain", "cloud", true},
		{"missing vanity domain", "clientID", "clientSecret", "", "cloud", true},
		{"valid client_id + client_secret without cloud", "clientID", "clientSecret", "vanityDomain", "", false}, // Ensures cloud is optional
	}

	// Execute each test case
	for _, test := range tests {
		resourceConfig := map[string]interface{}{
			"vanity_domain": test.vanityDomain,
		}
		if test.clientID != "" {
			resourceConfig["client_id"] = test.clientID
		}
		if test.clientSecret != "" {
			resourceConfig["client_secret"] = test.clientSecret
		}
		if test.cloud != "" {
			resourceConfig["cloud"] = test.cloud
		}

		config := terraform.NewResourceConfigRaw(resourceConfig)
		provider := ZIAProvider()
		err := provider.Validate(config)

		// Check expectations based on each test case setup
		if test.expectError && err == nil {
			t.Errorf("test %q: expected error but received none", test.name)
		}
		if !test.expectError && err != nil {
			t.Errorf("test %q: did not expect error but received error: %+v", test.name, err)
		}
	}

	// Restore environment variables after the tests
	for key, val := range envVals {
		os.Setenv(key, val)
	}
}

func sdkV3ClientForTest() (*zscaler.Client, error) {
	if testSdkV3Client != nil {
		return testSdkV3Client, nil
	}

	// Initialize the SDK V3 Client
	client, err := zscalerSDKV3Client(&Config{
		clientID:     os.Getenv("ZSCALER_CLIENT_ID"),
		clientSecret: os.Getenv("ZSCALER_CLIENT_SECRET"),
		vanityDomain: os.Getenv("ZSCALER_VANITY_DOMAIN"),
		cloud:        os.Getenv("ZSCALER_CLOUD"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SDK V3 client: %w", err)
	}

	testSdkV3Client = client
	return testSdkV3Client, nil
}
