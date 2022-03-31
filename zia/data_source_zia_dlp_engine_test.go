package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDLPEngines_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDLPEnginesConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.zia_dlp_engines.credit_cards", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_dlp_engines.ssn", "name"),
					resource.TestCheckResourceAttrSet(
						"data.zia_dlp_engines.cyber_bully", "name"),
				),
			},
		},
	})
}

const testAccCheckDataSourceDLPEnginesConfig_basic = `
data "zia_dlp_engines" "credit_cards"{
    name = "Credit Cards"
}
data "zia_dlp_engines" "ssn"{
    name = "Social Security Numbers"
}
data "zia_dlp_engines" "cyber_bully"{
    name = "CYBER_BULLY_ENG"
}
`
*/
