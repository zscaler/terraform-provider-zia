package zia

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudappcontrol"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_engines"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_notification_templates"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlpdictionaries"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewalldnscontrolpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallipscontrolpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkapplicationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationmanagement"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/rule_labels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/sandbox/sandbox_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/vpncredentials"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlcategories"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/users"
)

var (
	sweeperLogger   hclog.Logger
	sweeperLogLevel hclog.Level
)

func init() {
	sweeperLogLevel = hclog.Warn
	if os.Getenv("TF_LOG") != "" {
		sweeperLogLevel = hclog.LevelFromString(os.Getenv("TF_LOG"))
	}
	sweeperLogger = hclog.New(&hclog.LoggerOptions{
		Level:      sweeperLogLevel,
		TimeFormat: "2006/01/02 03:04:05",
	})
}

func logSweptResource(kind, id, nameOrLabel string) {
	sweeperLogger.Warn(fmt.Sprintf("sweeper found dangling %q %q %q", kind, id, nameOrLabel))
}

type testClient struct {
	sdkV3Client *zscaler.Client
}

var testResourcePrefix = "tf-acc-test-"
var updateResourcePrefix = "tf-updated-"

func TestRunForcedSweeper(t *testing.T) {
	if os.Getenv("ZIA_VCR_TF_ACC") != "" {
		t.Skip("forced sweeper is live and will never be run within VCR")
		return
	}
	if os.Getenv("ZIA_ACC_TEST_FORCE_SWEEPERS") == "" || os.Getenv("TF_ACC") == "" {
		t.Skipf("ENV vars %q and %q must not be blank to force running of the sweepers", "ZIA_ACC_TEST_FORCE_SWEEPERS", "TF_ACC")
		return
	}

	provider := ZIAProvider()
	c := terraform.NewResourceConfigRaw(nil)
	diag := provider.Configure(context.Background(), c)
	if diag.HasError() {
		t.Skipf("sweeper's provider configuration failed: %v", diag)
		return
	}

	sdkClient, err := sdkV3ClientForTest()
	if err != nil {
		t.Fatalf("Failed to get SDK client: %s", err)
	}

	testClient := &testClient{
		sdkV3Client: sdkClient,
	}

	sweepTestRuleLabels(testClient)
	sweepTestSourceIPGroup(testClient)
	sweepTestDestinationIPGroup(testClient)
	sweepTestNetworkServices(testClient)
	sweepTestNetworkServicesGroup(testClient)
	sweepTestNetworkAppGroups(testClient)
	sweepTestLocationManagement(testClient)
	sweepTestGRETunnels(testClient)
	sweepTestForwardingControlRule(testClient)
	sweepTestFirewallFilteringRule(testClient)
	sweepTestFirewallIPSRule(testClient)
	sweepTestFirewallDNSRule(testClient)
	sweepTestFileTypeControlRule(testClient)
	sweepTestSandboxRule(testClient)
	sweepTestURLFilteringRule(testClient)
	sweepTestDLPWebRule(testClient)
	sweepTestDLPEngines(testClient)
	sweepTestDLPDictionary(testClient)
	sweepTestDLPTemplates(testClient)
	sweepTestStaticIP(testClient)
	sweepTestURLCategories(testClient)
	sweepTestVPNCredentials(testClient)
	sweepTestAdminUser(testClient)
	sweepTestUsers(testClient)
}

// Sets up sweeper to clean up dangling resources
func setupSweeper(resourceType string, del func(*testClient) error) {
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			// Retrieve the client and handle the error
			sdkClient, err := sdkV3ClientForTest()
			if err != nil {
				return fmt.Errorf("failed to initialize SDK V3 client for sweeper: %w", err)
			}

			// Pass the client to the deleter function
			return del(&testClient{sdkV3Client: sdkClient})
		},
	})
}

func sweepTestRuleLabels(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	labels, err := rule_labels.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(labels)))
	for _, b := range labels {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := rule_labels.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.RuleLabels, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestSourceIPGroup(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	ipSourceGroup, err := ipsourcegroups.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipSourceGroup)))
	for _, b := range ipSourceGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := ipsourcegroups.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FWFilteringSourceGroup, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestDestinationIPGroup(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	ipDestGroup, err := ipdestinationgroups.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipDestGroup)))
	for _, b := range ipDestGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := ipdestinationgroups.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FWFilteringDestinationGroup, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestNetworkServices(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	services, err := networkservices.GetAllNetworkServices(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(services)))
	for _, b := range services {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkservices.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FWFilteringNetworkServices, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestNetworkServicesGroup(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	groups, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkservicegroups.DeleteNetworkServiceGroups(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FWFilteringNetworkServiceGroups, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestNetworkAppGroups(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	groups, err := networkapplicationgroups.GetAllNetworkApplicationGroups(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkapplicationgroups.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FWFilteringNetworkAppGroups, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestLocationManagement(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	locationManagement, err := locationmanagement.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(locationManagement)))
	for _, b := range locationManagement {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := locationmanagement.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.TrafficForwardingLocManagement, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestGRETunnels(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	greTunnel, err := gretunnels.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(greTunnel)))
	for _, b := range greTunnel {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comment, testResourcePrefix) || strings.HasPrefix(b.Comment, updateResourcePrefix) {
			if _, err := gretunnels.DeleteGreTunnels(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.TrafficForwardingGRETunnel, fmt.Sprintf("%d", b.ID), b.Comment)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestFirewallFilteringRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := filteringrules.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := filteringrules.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FirewallFilteringRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestFileTypeControlRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := firewalldnscontrolpolicies.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := firewalldnscontrolpolicies.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FirewallDNSRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestFirewallDNSRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := firewalldnscontrolpolicies.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := firewalldnscontrolpolicies.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FirewallDNSRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestFirewallIPSRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := firewallipscontrolpolicies.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := firewallipscontrolpolicies.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.FirewallIPSRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestSandboxRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := sandbox_rules.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := sandbox_rules.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.SandboxRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestForwardingControlRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := forwarding_rules.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := forwarding_rules.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.ForwardingControlRule, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestCloudAppControlRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	ruleTypes := []string{"STREAMING_MEDIA"}

	for _, ruleType := range ruleTypes {
		rules, err := cloudappcontrol.GetByRuleType(context.Background(), service, ruleType)
		if err != nil {
			return err
		}

		// Logging the number of identified resources before the deletion loop
		sweeperLogger.Warn(fmt.Sprintf("Found %d resources of type %s to sweep", len(rules), ruleType))
		for _, b := range rules {
			// Check if the resource name has the required prefix before deleting it
			if strings.HasPrefix(b.Name, testResourcePrefix) {
				if _, err := cloudappcontrol.Delete(context.Background(), service, ruleType, b.ID); err != nil {
					errorList = append(errorList, err)
					continue
				}
				logSweptResource(resourcetype.CloudAppControlRule, fmt.Sprintf("%d", b.ID), b.Name)
			}
		}
	}

	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}

	return condenseError(errorList)
}

func sweepTestURLFilteringRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := urlfilteringpolicies.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := urlfilteringpolicies.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.URLFilteringRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestDLPWebRule(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := dlp_web_rules.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlp_web_rules.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.DLPWebRules, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestDLPDictionary(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := dlpdictionaries.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlpdictionaries.DeleteDlpDictionary(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.DLPDictionaries, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestDLPEngines(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := dlp_engines.GetAll(context.Background(), service)
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))

	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlp_engines.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.DLPEngines, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}

	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}

	return condenseError(errorList)
}

func sweepTestDLPTemplates(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := dlp_notification_templates.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlp_notification_templates.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.DLPNotificationTemplates, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestStaticIP(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := staticips.GetAll(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comment, testResourcePrefix) || strings.HasPrefix(b.Comment, updateResourcePrefix) {
			if _, err := staticips.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.TrafficForwardingStaticIP, fmt.Sprintf("%d", b.ID), b.Comment)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}

	return condenseError(errorList)
}

func sweepTestVPNCredentials(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for VPN credentials
	service := &zscaler.Service{
		Client: client.sdkV3Client,
	}

	// Fetch all VPN credentials
	vpnCredentials, err := vpncredentials.GetAll(context.Background(), service)
	if err != nil {
		return fmt.Errorf("failed to fetch VPN credentials: %v", err)
	}

	// Filter IDs for deletion based on prefix
	var idsToDelete []int
	for _, credential := range vpnCredentials {
		if strings.HasPrefix(credential.Comments, testResourcePrefix) || strings.HasPrefix(credential.Comments, updateResourcePrefix) {
			idsToDelete = append(idsToDelete, credential.ID)
		}
	}

	// Log the number of resources identified for deletion
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(idsToDelete)))

	// Perform bulk deletions in batches
	for len(idsToDelete) > 0 {
		batchSize := 100 // Assuming 100 is the max batch size for bulk delete
		if len(idsToDelete) < batchSize {
			batchSize = len(idsToDelete)
		}

		batch := idsToDelete[:batchSize]
		idsToDelete = idsToDelete[batchSize:]

		// Perform the bulk delete
		_, err := vpncredentials.BulkDelete(context.Background(), service, batch)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to delete batch: %v", err))
			continue
		}

		// Log successful deletion of each resource in the batch
		for _, id := range batch {
			logSweptResource(resourcetype.TrafficForwardingVPNCredentials, fmt.Sprintf("%d", id), "Bulk deletion")
		}
	}

	// Log and return any errors encountered
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
		return condenseError(errorList)
	}

	return nil
}

func sweepTestURLCategories(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := urlcategories.GetAll(context.Background(), service, true, false)
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))

	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.ConfiguredName, testResourcePrefix) || strings.HasPrefix(b.ConfiguredName, updateResourcePrefix) {
			if _, err := urlcategories.DeleteURLCategories(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.URLCategories, (b.ID), b.ConfiguredName)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestAdminUser(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := admins.GetAllAdminUsers(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.UserName, testResourcePrefix) || strings.HasPrefix(b.UserName, updateResourcePrefix) {
			if _, err := admins.DeleteAdminUser(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.AdminUsers, fmt.Sprintf("%d", b.ID), b.UserName)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}

func sweepTestUsers(client *testClient) error {
	var errorList []error

	// Instantiate the specific service for App Connector Group
	service := &zscaler.Service{
		Client: client.sdkV3Client, // Use the existing SDK client
	}

	rule, err := users.GetAllUsers(context.Background(), service)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := users.Delete(context.Background(), service, b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.Users, fmt.Sprintf("%d", b.ID), b.Name)
		}
	}
	// Log errors encountered during the deletion process
	if len(errorList) > 0 {
		for _, err := range errorList {
			sweeperLogger.Error(err.Error())
		}
	}
	return condenseError(errorList)
}
