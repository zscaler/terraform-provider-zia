package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/roles"
)

func TestAccResourceAdminRolesBasic(t *testing.T) {
	var roles roles.AdminRoles
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminRoles)

	initialName := "tf-acc-test-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminRolesConfigure(resourceTypeAndName, initialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAdminRolesExists(resourceTypeAndName, &roles),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "alerting_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dashboard_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "report_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "analysis_access", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username_access", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "device_info_access", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "admin_acct_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "policy_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "role_type", "EXEC_INSIGHT_AND_ORG_ADMIN"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "logs_limit", "UNRESTRICTED"),
				),
			},

			// Update test
			{
				Config: testAccCheckAdminRolesConfigure(resourceTypeAndName, initialName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "alerting_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "dashboard_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "report_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "analysis_access", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "username_access", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "device_info_access", "READ_ONLY"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "admin_acct_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "policy_access", "READ_WRITE"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "role_type", "EXEC_INSIGHT_AND_ORG_ADMIN"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "logs_limit", "UNRESTRICTED"),
				),
			},
			// Import test
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAdminRolesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.AdminRoles {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := roles.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckAdminRolesExists(resource string, rule *roles.AdminRoles) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		apiClient := testAccProvider.Meta().(*Client)
		service := apiClient.Service

		receivedRole, err := roles.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRole

		return nil
	}
}

func testAccCheckAdminRolesConfigure(resourceTypeAndName, generatedName string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name = "%s"
	rank = 7
	alerting_access    = "READ_WRITE"
	dashboard_access   = "READ_WRITE"
	report_access      = "READ_WRITE"
	analysis_access    = "READ_ONLY"
	username_access    = "READ_ONLY"
	device_info_access = "READ_ONLY"
	admin_acct_access  = "READ_WRITE"
	policy_access      = "READ_WRITE"
	role_type          = "EXEC_INSIGHT_AND_ORG_ADMIN"
  permissions = [
    "ADVANCED_SETTINGS",
    "COMPLY",
    "FIREWALL_DNS",
    "NSS_CONFIGURATION",
    "SECURE",
    "SSL_POLICY",
    "VZEN_CONFIGURATION",
    "PARTNER_INTEGRATION",
    "REMOTE_ASSISTANCE_MANAGEMENT",
    "LOCATIONS",
    "VPN_CREDENTIALS",
    "HOSTED_PAC_FILES",
    "EZ_AGENT_CONFIGURATIONS",
    "SECURE_AGENT_NOTIFICATIONS",
    "PROXY_GATEWAY",
    "STATIC_IPS",
    "GRE_TUNNELS",
    "SUBCLOUDS",
    "AUTHENTICATION_SETTINGS",
    "USER_MANAGEMENT",
    "IDENTITY_PROXY_SETTINGS",
    "APIKEY_MANAGEMENT",
    "POLICY_RESOURCE_MANAGEMENT",
    "CLIENT_CONNECTOR_PORTAL",
    "CUSTOM_URL_CAT",
    "OVERRIDE_EXISTING_CAT",
    "TENANT_PROFILE_MANAGEMENT"
  ]
  ext_feature_permissions = {
    INCIDENT_WORKFLOW = "FULL"
  }

  logs_limit           = "UNRESTRICTED"
  report_time_duration = -1
}

data "%s" "%s" {
    name = "${%s.%s.name}"
}
`,
		// resource variables
		resourcetype.AdminRoles,
		resourceName,
		generatedName,

		// data source variables
		resourcetype.AdminRoles,
		resourceName,
		// Reference to the resource
		resourcetype.AdminRoles,
		resourceName,
	)
}
