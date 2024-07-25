package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDLPDictionaryPredefinedIdentifiers_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDLPDictionaryPredefinedIdentifiersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceDLPDictionaryPredefinedIdentifiersCheck("data.zia_dlp_dictionary_predefined_identifiers.identifier1"),
					testAccDataSourceDLPDictionaryPredefinedIdentifiersCheck("data.zia_dlp_dictionary_predefined_identifiers.identifier2"),
					testAccDataSourceDLPDictionaryPredefinedIdentifiersCheck("data.zia_dlp_dictionary_predefined_identifiers.identifier3"),
					testAccDataSourceDLPDictionaryPredefinedIdentifiersCheck("data.zia_dlp_dictionary_predefined_identifiers.identifier4"),
					testAccDataSourceDLPDictionaryPredefinedIdentifiersCheck("data.zia_dlp_dictionary_predefined_identifiers.identifier5"),
				),
			},
		},
	})
}

func testAccDataSourceDLPDictionaryPredefinedIdentifiersCheck(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttrSet(name, "name"),
	)
}

var testAccCheckDataSourceDLPDictionaryPredefinedIdentifiersConfig_basic = `
data "zia_dlp_dictionary_predefined_identifiers" "identifier1"{
    name = "ASPP_LEAKAGE"
}

data "zia_dlp_dictionary_predefined_identifiers" "identifier2"{
    name = "CRED_LEAKAGE"
}

data "zia_dlp_dictionary_predefined_identifiers" "identifier3"{
    name = "EUIBAN_LEAKAGE"
}

data "zia_dlp_dictionary_predefined_identifiers" "identifier4"{
    name = "PPEU_LEAKAGE"
}

data "zia_dlp_dictionary_predefined_identifiers" "identifier5"{
    name = "USDL_LEAKAGE"
}
`
