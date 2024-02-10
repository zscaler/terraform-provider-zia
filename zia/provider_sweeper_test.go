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
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
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
	sdkClient *Client
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

	provider := Provider()
	c := terraform.NewResourceConfigRaw(nil)
	diag := provider.Configure(context.TODO(), c)
	if diag.HasError() {
		t.Skipf("sweeper's provider configuration failed: %v", diag)
		return
	}

	sdkClient, err := sdkClientForTest()
	if err != nil {
		t.Fatalf("Failed to get SDK client: %s", err)
	}

	testClient := &testClient{
		sdkClient: sdkClient,
	}

	sweepTestRuleLabels(testClient)
	sweepTestSourceIPGroup(testClient)
	sweepTestDestinationIPGroup(testClient)
	sweepTestNetworkServices(testClient)
	sweepTestNetworkServicesGroup(testClient)
	sweepTestNetworkAppGroups(testClient)
	sweepTestLocationManagement(testClient)
	sweepTestGRETunnels(testClient)
	// sweepTestForwardingControlRule(testClient)
	sweepTestFirewallFilteringRule(testClient)
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
	sdkClient, err := sdkClientForTest()
	if err != nil {
		// You might decide how to handle the error here. Using a panic for simplicity.
		panic(fmt.Sprintf("Failed to get SDK client: %s", err))
	}
	resource.AddTestSweepers(resourceType, &resource.Sweeper{
		Name: resourceType,
		F: func(_ string) error {
			return del(&testClient{sdkClient: sdkClient})
		},
	})
}

func sweepTestRuleLabels(client *testClient) error {
	var errorList []error
	labels, err := client.sdkClient.rule_labels.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(labels)))
	for _, b := range labels {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.rule_labels.Delete(b.ID); err != nil {
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
	ipSourceGroup, err := client.sdkClient.ipsourcegroups.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipSourceGroup)))
	for _, b := range ipSourceGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.ipsourcegroups.Delete(b.ID); err != nil {
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
	ipDestGroup, err := client.sdkClient.ipdestinationgroups.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipDestGroup)))
	for _, b := range ipDestGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.ipdestinationgroups.Delete(b.ID); err != nil {
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
	services, err := client.sdkClient.networkservices.GetAllNetworkServices()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(services)))
	for _, b := range services {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.networkservices.Delete(b.ID); err != nil {
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
	groups, err := client.sdkClient.networkservicegroups.GetAllNetworkServiceGroups()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.networkservicegroups.DeleteNetworkServiceGroups(b.ID); err != nil {
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
	groups, err := client.sdkClient.networkapplicationgroups.GetAllNetworkApplicationGroups()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.networkapplicationgroups.Delete(b.ID); err != nil {
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
	locationManagement, err := client.sdkClient.locationmanagement.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(locationManagement)))
	for _, b := range locationManagement {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.locationmanagement.Delete(b.ID); err != nil {
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
	greTunnel, err := client.sdkClient.gretunnels.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(greTunnel)))
	for _, b := range greTunnel {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comment, testResourcePrefix) || strings.HasPrefix(b.Comment, updateResourcePrefix) {
			if _, err := client.sdkClient.gretunnels.DeleteGreTunnels(b.ID); err != nil {
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
	rule, err := client.sdkClient.filteringrules.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.filteringrules.Delete(b.ID); err != nil {
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

/*
func sweepTestForwardingControlRule(client *testClient) error {
	var errorList []error
	rule, err := client.sdkClient.forwarding_rules.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.forwarding_rules.Delete(b.ID); err != nil {
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
*/

func sweepTestURLFilteringRule(client *testClient) error {
	var errorList []error

	rule, err := client.sdkClient.urlfilteringpolicies.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.urlfilteringpolicies.Delete(b.ID); err != nil {
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
	rule, err := client.sdkClient.dlp_web_rules.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.dlp_web_rules.Delete(b.ID); err != nil {
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
	rule, err := client.sdkClient.dlpdictionaries.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.dlpdictionaries.DeleteDlpDictionary(b.ID); err != nil {
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

	rule, err := client.sdkClient.dlp_engines.GetAll()
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))

	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.dlp_engines.Delete(b.ID); err != nil {
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
	rule, err := client.sdkClient.dlp_notification_templates.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.dlp_notification_templates.Delete(b.ID); err != nil {
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

	rule, err := client.sdkClient.staticips.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comment, testResourcePrefix) || strings.HasPrefix(b.Comment, updateResourcePrefix) {
			if _, err := client.sdkClient.staticips.Delete(b.ID); err != nil {
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

	rule, err := client.sdkClient.vpncredentials.GetAll()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comments, testResourcePrefix) || strings.HasPrefix(b.Comments, updateResourcePrefix) {
			if err := client.sdkClient.vpncredentials.Delete(b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.TrafficForwardingVPNCredentials, fmt.Sprintf("%d", b.ID), b.Comments)
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

func sweepTestURLCategories(client *testClient) error {
	var errorList []error

	rule, err := client.sdkClient.urlcategories.GetAll()
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))

	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.ConfiguredName, testResourcePrefix) || strings.HasPrefix(b.ConfiguredName, updateResourcePrefix) {
			if _, err := client.sdkClient.urlcategories.DeleteURLCategories(b.ID); err != nil {
				errorList = append(errorList, err)
				continue
			}
			logSweptResource(resourcetype.URLCategories, fmt.Sprintf(b.ID), b.ConfiguredName)
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

	rule, err := client.sdkClient.admins.GetAllAdminUsers()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.UserName, testResourcePrefix) || strings.HasPrefix(b.UserName, updateResourcePrefix) {
			if _, err := client.sdkClient.admins.DeleteAdminUser(b.ID); err != nil {
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

	rule, err := client.sdkClient.users.GetAllUsers()
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := client.sdkClient.users.Delete(b.ID); err != nil {
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
