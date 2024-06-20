package zia

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_engines"
)

func TestAccResourceDLPEnginesBasic(t *testing.T) {
	var engine dlp_engines.DLPEngines
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPEngines)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPEnginesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDLPEnginesConfigure(resourceTypeAndName, initialName, generatedName, variable.DLPCustomEngine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPEnginesExists(resourceTypeAndName, &engine),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_dlp_engine", strconv.FormatBool(variable.DLPCustomEngine)),
				),
			},

			// Update test
			{
				Config: testAccCheckDLPEnginesConfigure(resourceTypeAndName, updatedName, generatedName, variable.DLPCustomEngine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPEnginesExists(resourceTypeAndName, &engine),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "custom_dlp_engine", strconv.FormatBool(variable.DLPCustomEngine)),
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

func testAccCheckDLPEnginesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.dlp_engines

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPEngines {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		engine, err := dlp_engines.Get(service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if engine != nil {
			return fmt.Errorf("dlp engines with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDLPEnginesExists(resource string, engine *dlp_engines.DLPEngines) resource.TestCheckFunc {
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
		service := apiClient.dlp_engines

		receivedEngine, err := dlp_engines.Get(service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*engine = *receivedEngine

		return nil
	}
}

func testAccCheckDLPEnginesConfigure(resourceTypeAndName, generatedName, description string, customEngine bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
	description = "tf-acc-test-%s"
	engine_expression = "((D63.S > 1))"
	custom_dlp_engine = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.DLPEngines,
		resourceName,
		generatedName,
		description,
		strconv.FormatBool(customEngine),

		// data source variables
		resourcetype.DLPEngines,
		resourceName,
		resourcetype.DLPEngines,
		resourceName,
	)
}
