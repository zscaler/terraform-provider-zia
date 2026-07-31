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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/security_ueba_alerts/alert_definitions"
)

func TestAccResourceUEBAAlertDefinitionsBasic(t *testing.T) {
	var def alert_definitions.AlertDefinitions
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.UEBAAlertDefinitions)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUEBAAlertDefinitionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUEBAAlertDefinitionsConfigure(resourceTypeAndName, generatedName, variable.UEBAAlertSeverity, variable.UEBAAlertComments),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUEBAAlertDefinitionsExists(resourceTypeAndName, &def),
					resource.TestCheckResourceAttr(resourceTypeAndName, "alert_name", variable.UEBAAlertName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.UEBAAlertStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "scope", variable.UEBAAlertScope),
					resource.TestCheckResourceAttr(resourceTypeAndName, "severity", variable.UEBAAlertSeverity),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", variable.UEBAAlertComments),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "alert_name", resourceTypeAndName, "alert_name"),
				),
			},
			{
				Config: testAccCheckUEBAAlertDefinitionsConfigure(resourceTypeAndName, generatedName, variable.UEBAAlertSeverityUpd, variable.UEBAAlertCommentsUpd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUEBAAlertDefinitionsExists(resourceTypeAndName, &def),
					resource.TestCheckResourceAttr(resourceTypeAndName, "severity", variable.UEBAAlertSeverityUpd),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", variable.UEBAAlertCommentsUpd),
				),
			},
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckUEBAAlertDefinitionsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.UEBAAlertDefinitions {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		def, err := alert_definitions.Get(context.Background(), service, id)
		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}
		if def != nil {
			return fmt.Errorf("UEBA alert definition with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckUEBAAlertDefinitionsExists(resource string, def *alert_definitions.AlertDefinitions) resource.TestCheckFunc {
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

		received, err := alert_definitions.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Received error: %s", resource, err)
		}
		*def = *received

		return nil
	}
}

func testAccCheckUEBAAlertDefinitionsConfigure(resourceTypeAndName, generatedName, severity, comments string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	return fmt.Sprintf(`
resource "%s" "%s" {
    alert_name = "%s"
    status     = "%s"
    scope      = "%s"
    occurrence = "%s"
    interval   = "%s"
    severity   = "%s"
    comments   = "%s"
}

data "%s" "%s" {
    id = "${%s.%s.id}"
}
`,
		resourcetype.UEBAAlertDefinitions,
		resourceName,
		variable.UEBAAlertName,
		variable.UEBAAlertStatus,
		variable.UEBAAlertScope,
		variable.UEBAAlertOccurrence,
		variable.UEBAAlertInterval,
		severity,
		comments,

		resourcetype.UEBAAlertDefinitions,
		resourceName,
		resourcetype.UEBAAlertDefinitions, resourceName,
	)
}
