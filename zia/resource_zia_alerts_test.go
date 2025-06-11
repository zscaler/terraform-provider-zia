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
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/alerts"
)

func TestAccResourceSubscriptionAlertsBasic(t *testing.T) {
	var alerts alerts.AlertSubscriptions
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.SubscriptionAlerts)

	initialEmail := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubscriptionAlertsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSubscriptionAlertsConfigure(resourceTypeAndName, initialEmail, variable.AlertDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubscriptionAlertsExists(resourceTypeAndName, &alerts),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", initialEmail+"@acme.com"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AlertDescription),
				),
			},

			// Update test
			{
				Config: testAccCheckSubscriptionAlertsConfigure(resourceTypeAndName, updatedName, variable.AlertDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubscriptionAlertsExists(resourceTypeAndName, &alerts),
					resource.TestCheckResourceAttr(resourceTypeAndName, "email", updatedName+"@acme.com"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.AlertDescription),
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

func testAccCheckSubscriptionAlertsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.SubscriptionAlerts {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := alerts.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckSubscriptionAlertsExists(resource string, alert *alerts.AlertSubscriptions) resource.TestCheckFunc {
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

		receivedRule, err := alerts.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*alert = *receivedRule

		return nil
	}
}

func testAccCheckSubscriptionAlertsConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	email := generatedName + "@acme.com"

	return fmt.Sprintf(`
resource "%s" "%s" {
    email = "%s"
    description = "%s"
	pt0_severities = ["CRITICAL"]
	secure_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
	manage_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
	comply_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
	system_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
}

data "%s" "%s" {
	id = "${%s.%s.id}"
}
`,
		resourcetype.SubscriptionAlerts,
		resourceName,
		email,
		description,

		resourcetype.SubscriptionAlerts,
		resourceName,
		resourcetype.SubscriptionAlerts, resourceName,
	)
}
