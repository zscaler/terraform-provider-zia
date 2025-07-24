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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/bandwidth_control/bandwidth_classes"
)

func TestAccResourceBandwdithClasses_Basic(t *testing.T) {
	var labels bandwidth_classes.BandwidthClasses
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.BandwdithClasses)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBandwdithClassesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBandwdithClassesConfigure(resourceTypeAndName, initialName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwdithClassesExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "web_applications.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "urls.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "url_categories.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckBandwdithClassesConfigure(resourceTypeAndName, updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBandwdithClassesExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "web_applications.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "urls.#", "3"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "url_categories.#", "2"),
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

func testAccCheckBandwdithClassesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.BandwdithClasses {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		class, err := bandwidth_classes.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if class != nil {
			return fmt.Errorf("class with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckBandwdithClassesExists(resource string, class *bandwidth_classes.BandwidthClasses) resource.TestCheckFunc {
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

		receivedClass, err := bandwidth_classes.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*class = *receivedClass

		return nil
	}
}

func testAccCheckBandwdithClassesConfigure(resourceTypeAndName, generatedName string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
  web_applications = [
    "ACADEMICGPT",
    "AD_CREATIVES",
    "AGENTGPT"
  ]
  urls = [
    "chatgpt.com",
    "chatgpt1.com",
    "openai.com"
  ]
  url_categories = [
    "AI_ML_APPS",
    "GENERAL_AI_ML",
  ]
}

data "%s" "%s" {
  id = %s.%s.id
}
`,
		// resource
		resourcetype.BandwdithClasses,
		resourceName,
		generatedName,

		// data source
		resourcetype.BandwdithClasses,
		resourceName,
		resourcetype.BandwdithClasses, resourceName,
	)
}
