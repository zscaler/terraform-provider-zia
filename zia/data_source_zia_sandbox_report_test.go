package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSandboxReport_Basic(t *testing.T) {
	md5Hash := "F5E282A09B60748513270A7415E3B526" // Example MD5 hash

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceSandboxReportConfig(md5Hash),
				Check: resource.ComposeTestCheckFunc(
					// Check for the 'full' report
					resource.TestCheckResourceAttr("data.zia_sandbox_report.full", "md5_hash", md5Hash),
					resource.TestCheckResourceAttr("data.zia_sandbox_report.full", "details", "full"),

					// Check for the 'summary' report
					resource.TestCheckResourceAttr("data.zia_sandbox_report.summary", "md5_hash", md5Hash),
					resource.TestCheckResourceAttr("data.zia_sandbox_report.summary", "details", "summary"),
				),
			},
		},
	})
}

func testAccCheckDataSourceSandboxReportConfig(md5Hash string) string {
	return fmt.Sprintf(`
data "zia_sandbox_report" "full" {
	md5_hash = "%s"
	details  = "full"
}

data "zia_sandbox_report" "summary" {
	md5_hash = "%s"
	details  = "summary"
}
`, md5Hash, md5Hash)
}
*/
