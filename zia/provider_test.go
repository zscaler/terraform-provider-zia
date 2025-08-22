package zia

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	// Save and clear environment variables
	envKeys := []string{
		"ZSCALER_CLIENT_ID",
		"ZSCALER_CLIENT_SECRET",
		"ZSCALER_VANITY_DOMAIN",
		"ZSCALER_CLOUD",
	}
	envVals := make(map[string]string)

	for _, key := range envKeys {
		val := os.Getenv(key)
		if val != "" {
			envVals[key] = val
			os.Unsetenv(key)
		}
	}

	// Define test cases using actual env values for valid config
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		vanityDomain string
		cloud        string
		expectError  bool
	}{
		{
			name:         "valid client_id + client_secret",
			clientID:     envVals["ZSCALER_CLIENT_ID"],
			clientSecret: envVals["ZSCALER_CLIENT_SECRET"],
			vanityDomain: envVals["ZSCALER_VANITY_DOMAIN"],
			cloud:        envVals["ZSCALER_CLOUD"],
			expectError:  false,
		},
		{
			name:         "missing client_id",
			clientID:     "",
			clientSecret: envVals["ZSCALER_CLIENT_SECRET"],
			vanityDomain: envVals["ZSCALER_VANITY_DOMAIN"],
			cloud:        envVals["ZSCALER_CLOUD"],
			expectError:  true,
		},
		{
			name:         "missing clientSecret",
			clientID:     envVals["ZSCALER_CLIENT_ID"],
			clientSecret: "",
			vanityDomain: envVals["ZSCALER_VANITY_DOMAIN"],
			cloud:        envVals["ZSCALER_CLOUD"],
			expectError:  true,
		},
		{
			name:         "missing vanity domain",
			clientID:     envVals["ZSCALER_CLIENT_ID"],
			clientSecret: envVals["ZSCALER_CLIENT_SECRET"],
			vanityDomain: "",
			cloud:        envVals["ZSCALER_CLOUD"],
			expectError:  true,
		},
		{
			name:         "valid client_id + client_secret without zscaler_cloud",
			clientID:     envVals["ZSCALER_CLIENT_ID"],
			clientSecret: envVals["ZSCALER_CLIENT_SECRET"],
			vanityDomain: envVals["ZSCALER_VANITY_DOMAIN"],
			cloud:        "",
			expectError:  false,
		},
	}

	// Execute each test case
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
				resourceConfig["zscaler_cloud"] = test.cloud
			}

			provider := ZIAProvider()
			rawData := schema.TestResourceDataRaw(t, provider.Schema, resourceConfig)

			_, diags := provider.ConfigureContextFunc(context.Background(), rawData)

			if test.expectError && !diags.HasError() {
				t.Errorf("expected error but received none")
			}
			if !test.expectError && diags.HasError() {
				t.Errorf("did not expect error but received: %+v", diags)
			}
		})
	}

	// Restore original env vars
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
