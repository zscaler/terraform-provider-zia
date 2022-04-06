package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlpdictionaries"
)

func TestAccResourceDLPDictionaries_basic(t *testing.T) {
	var dictionary dlpdictionaries.DlpDictionary
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(20)
	resourceName := "zia_dlp_dictionaries.test-dlp-dict"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPDictionariesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDLPDictionariesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPDictionariesExists("zia_dlp_dictionaries.test-dlp-dict", &dictionary),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dlp-dict-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-dlp-dict-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "dictionary_type", "PATTERNS_AND_PHRASES"),
				),
			},
		},
	})
}

func testAccDLPDictionariesBasic(rName, rDesc string) string {
	return fmt.Sprintf(`

resource "zia_dlp_dictionaries" "test-dlp-dict"{
	name = "test-dlp-dict-%s"
	description = "test-dlp-dict-%s"
	phrases {
		action = "PHRASE_COUNT_TYPE_ALL"
		phrase = "Test1"
	}
	phrases {
		action = "PHRASE_COUNT_TYPE_UNIQUE"
		phrase = "Test2"
	}
	custom_phrase_match_type = "MATCH_ALL_CUSTOM_PHRASE_PATTERN_DICTIONARY"
	patterns {
		action = "PATTERN_COUNT_TYPE_ALL"
		pattern = "Test1"
	}
	patterns {
		action = "PATTERN_COUNT_TYPE_UNIQUE"
		pattern = "Test2"
	}
	dictionary_type = "PATTERNS_AND_PHRASES"
}
	`, rName, rDesc)
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
		receivedDictionary, err := apiClient.dlpdictionaries.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*dictionary = *receivedDictionary

		return nil
	}
}

func testAccCheckDLPDictionariesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zia_dlp_dictionaries" {
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
