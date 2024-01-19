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
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlpdictionaries"
)

func TestAccResourceDLPDictionariesBasic(t *testing.T) {
	var dictionary dlpdictionaries.DlpDictionary
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPDictionaries)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPDictionariesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDLPDictionariesConfigure(resourceTypeAndName, initialName, variable.DLPDictionaryDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPDictionariesExists(resourceTypeAndName, &dictionary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPDictionaryDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "phrases.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "patterns.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckDLPDictionariesConfigure(resourceTypeAndName, updatedName, variable.DLPDictionaryDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPDictionariesExists(resourceTypeAndName, &dictionary),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.DLPDictionaryDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "phrases.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "patterns.#", "2"),
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

func testAccCheckDLPDictionariesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPDictionaries {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.dlpdictionaries.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("dlp dictionaries with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDLPDictionariesExists(resource string, dictionary *dlpdictionaries.DlpDictionary) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.dlpdictionaries.Get(id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*dictionary = *receivedRule

		return nil
	}
}

func testAccCheckDLPDictionariesConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name

	return fmt.Sprintf(`
resource "%s" "%s" {
    name = "%s"
	description = "%s"
    phrases {
        action = "PHRASE_COUNT_TYPE_ALL"
        phrase = "Test1"
    }
    phrases {
        action = "PHRASE_COUNT_TYPE_ALL"
        phrase = "Test2"
    }
    custom_phrase_match_type = "MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY"
    patterns {
        action = "PATTERN_COUNT_TYPE_UNIQUE"
        pattern = "Test1"
    }
    patterns {
        action = "PATTERN_COUNT_TYPE_UNIQUE"
        pattern = "Test2"
    }
    dictionary_type = "PATTERNS_AND_PHRASES"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.DLPDictionaries,
		resourceName,
		generatedName,
		description,

		// data source variables
		resourcetype.DLPDictionaries,
		resourceName,
		resourcetype.DLPDictionaries,
		resourceName,
	)
}
