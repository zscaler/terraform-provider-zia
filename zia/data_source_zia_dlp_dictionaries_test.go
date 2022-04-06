package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDLPDictionaries_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	rDesc := acctest.RandString(15)
	resourceName := "data.zia_dlp_dictionaries.test-dlp-dict"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDLPDictionariesBasic(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDLPDictionariesRule(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-dlp-dict-"+rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test-dlp-dict-"+rDesc),
					resource.TestCheckResourceAttr(resourceName, "dictionary_type", "PATTERNS_AND_PHRASES"),
				),
			},
		},
	})
}

func testAccDataSourceDLPDictionariesBasic(rName, rDesc string) string {
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

data "zia_dlp_dictionaries" "test-dlp-dict" {
	name = zia_dlp_dictionaries.test-dlp-dict.name
}
	`, rName, rDesc)
}

func testAccDataSourceDLPDictionariesRule(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}

		return nil
	}
}
