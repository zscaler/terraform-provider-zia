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
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_engines"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_notification_templates"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlpdictionaries"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkapplicationgroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservices"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/location/locationmanagement"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/rule_labels"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/vpncredentials"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlcategories"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlfilteringpolicies"
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
	labels, err := rule_labels.GetAll(client.sdkClient.rule_labels)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(labels)))
	for _, b := range labels {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := rule_labels.Delete(client.sdkClient.rule_labels, b.ID); err != nil {
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
	ipSourceGroup, err := ipsourcegroups.GetAll(client.sdkClient.ipsourcegroups)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipSourceGroup)))
	for _, b := range ipSourceGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := ipsourcegroups.Delete(client.sdkClient.ipsourcegroups, b.ID); err != nil {
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
	ipDestGroup, err := ipdestinationgroups.GetAll(client.sdkClient.ipdestinationgroups)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(ipDestGroup)))
	for _, b := range ipDestGroup {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := ipdestinationgroups.Delete(client.sdkClient.ipdestinationgroups, b.ID); err != nil {
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
	services, err := networkservices.GetAllNetworkServices(client.sdkClient.networkservices)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(services)))
	for _, b := range services {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkservices.Delete(client.sdkClient.networkservices, b.ID); err != nil {
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
	groups, err := networkservicegroups.GetAllNetworkServiceGroups(client.sdkClient.networkservicegroups)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkservicegroups.DeleteNetworkServiceGroups(client.sdkClient.networkservicegroups, b.ID); err != nil {
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
	groups, err := networkapplicationgroups.GetAllNetworkApplicationGroups(client.sdkClient.networkapplicationgroups)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(groups)))
	for _, b := range groups {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := networkapplicationgroups.Delete(client.sdkClient.networkapplicationgroups, b.ID); err != nil {
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
	locationManagement, err := locationmanagement.GetAll(client.sdkClient.locationmanagement)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(locationManagement)))
	for _, b := range locationManagement {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := locationmanagement.Delete(client.sdkClient.locationmanagement, b.ID); err != nil {
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
	greTunnel, err := gretunnels.GetAll(client.sdkClient.gretunnels)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(greTunnel)))
	for _, b := range greTunnel {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comment, testResourcePrefix) || strings.HasPrefix(b.Comment, updateResourcePrefix) {
			if _, err := gretunnels.DeleteGreTunnels(client.sdkClient.gretunnels, b.ID); err != nil {
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
	rule, err := filteringrules.GetAll(client.sdkClient.filteringrules)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := filteringrules.Delete(client.sdkClient.filteringrules, b.ID); err != nil {
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

func sweepTestForwardingControlRule(client *testClient) error {
	var errorList []error
	rule, err := forwarding_rules.GetAll(client.sdkClient.forwarding_rules)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := forwarding_rules.Delete(client.sdkClient.forwarding_rules, b.ID); err != nil {
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

func sweepTestURLFilteringRule(client *testClient) error {
	var errorList []error

	rule, err := urlfilteringpolicies.GetAll(client.sdkClient.urlfilteringpolicies)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := urlfilteringpolicies.Delete(client.sdkClient.urlfilteringpolicies, b.ID); err != nil {
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
	rule, err := dlp_web_rules.GetAll(client.sdkClient.dlp_web_rules)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlp_web_rules.Delete(client.sdkClient.dlp_web_rules, b.ID); err != nil {
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
	rule, err := dlpdictionaries.GetAll(client.sdkClient.dlpdictionaries)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlpdictionaries.DeleteDlpDictionary(client.sdkClient.dlpdictionaries, b.ID); err != nil {
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

	rule, err := dlp_engines.GetAll(client.sdkClient.dlp_engines)
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))

	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlp_engines.Delete(client.sdkClient.dlp_engines, b.ID); err != nil {
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
	rule, err := dlp_notification_templates.GetAll(client.sdkClient.dlp_notification_templates)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Name, testResourcePrefix) || strings.HasPrefix(b.Name, updateResourcePrefix) {
			if _, err := dlp_notification_templates.Delete(client.sdkClient.dlp_notification_templates, b.ID); err != nil {
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

	rule, err := staticips.GetAll(client.sdkClient.staticips)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comment, testResourcePrefix) || strings.HasPrefix(b.Comment, updateResourcePrefix) {
			if _, err := staticips.Delete(client.sdkClient.staticips, b.ID); err != nil {
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

	rule, err := vpncredentials.GetAll(client.sdkClient.vpncredentials)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.Comments, testResourcePrefix) || strings.HasPrefix(b.Comments, updateResourcePrefix) {
			if err := vpncredentials.Delete(client.sdkClient.vpncredentials, b.ID); err != nil {
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

	rule, err := urlcategories.GetAll(client.sdkClient.urlcategories)
	if err != nil {
		return err
	}

	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))

	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.ConfiguredName, testResourcePrefix) || strings.HasPrefix(b.ConfiguredName, updateResourcePrefix) {
			if _, err := urlcategories.DeleteURLCategories(client.sdkClient.urlcategories, b.ID); err != nil {
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

	rule, err := admins.GetAllAdminUsers(client.sdkClient.admins)
	if err != nil {
		return err
	}
	// Logging the number of identified resources before the deletion loop
	sweeperLogger.Warn(fmt.Sprintf("Found %d resources to sweep", len(rule)))
	for _, b := range rule {
		// Check if the resource name has the required prefix before deleting it
		if strings.HasPrefix(b.UserName, testResourcePrefix) || strings.HasPrefix(b.UserName, updateResourcePrefix) {
			if _, err := admins.DeleteAdminUser(client.sdkClient.admins, b.ID); err != nil {
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
