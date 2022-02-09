package zia

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceTrafficForwardingVPNCredentials_ByIdAndType(t *testing.T) {
	rComments := acctest.RandString(15)
	resourceName := "data.zia_traffic_forwarding_vpn_credentials.by_id"
	resourceName2 := "data.zia_traffic_forwarding_vpn_credentials.by_type"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTrafficForwardingVPNCredentialsById(rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingStaticIP(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comments", rComments),
					// resource.TestCheckResourceAttr(resourceName, "pre_shared_key", "Password@123"),
				),
				PreventPostDestroyRefresh: true,
			},
			{

				Config: testAccDataSourceTrafficForwardingVPNCredentialsByType(rComments),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceTrafficForwardingVPNCredentials(resourceName2),
					resource.TestCheckResourceAttr(resourceName2, "comments", rComments),
					// resource.TestCheckResourceAttr(resourceName2, "pre_shared_key", "Password@123"),
				),
			},
		},
	})
}

func testAccDataSourceTrafficForwardingVPNCredentialsById(rComments string) string {
	return fmt.Sprintf(`
resource "zia_traffic_forwarding_vpn_credentials" "testAcc"{
	type = "UFQDN"
	fqdn = "testAcc@securitygeek.io"
	comments = "%s"
	pre_shared_key = "Password@123"
}
data "zia_traffic_forwarding_vpn_credentials" "by_id" {
	id = zia_traffic_forwarding_vpn_credentials.testAcc.id
}
	`, rComments)
}

func testAccDataSourceTrafficForwardingVPNCredentialsByType(rComments string) string {
	return fmt.Sprintf(`
resource "zia_traffic_forwarding_vpn_credentials" "testAcc"{
	type = "UFQDN"
	fqdn = "testAcc@securitygeek.io"
	comments = "%s"
	pre_shared_key = "Password@123"
}
data "zia_traffic_forwarding_vpn_credentials" "by_type" {
	type = zia_traffic_forwarding_vpn_credentials.testAcc.type
}
	`, rComments)
}

func testAccDataSourceTrafficForwardingVPNCredentials(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no data source called %s", name)
		}
		return nil
	}
}
*/
