package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
)

func TestAccResourceTrafficForwardingVPNCredentialsBasic(t *testing.T) {
	var credentials vpncredentials.VPNCredentials
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringVPNCredentials)
	rEmail := acctest.RandomWithPrefix("tf-acc-test-")
	rSharedKey := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingVPNCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				// creation vpn credential type ufqdn
				Config: testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, rEmail, rSharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingVPNCredentialsExists(resourceTypeAndName, &credentials),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", "UFQDN"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "fqdn", rEmail+"@securitygeek.io"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pre_shared_key", rSharedKey),
				),
			},

			// update pre-shared-key and comments vpn credential type ufqdn
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, rEmail, rSharedKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingVPNCredentialsExists(resourceTypeAndName, &credentials),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pre_shared_key", rSharedKey),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingVPNCredentialsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficFilteringVPNCredentials {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.vpncredentials.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("vpn credentials with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckTrafficForwardingVPNCredentialsExists(resource string, rule *vpncredentials.VPNCredentials) resource.TestCheckFunc {
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
		receivedRule, err := apiClient.vpncredentials.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, rEmail, rSharedKey string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	comments = "tf-acc-test-%s"
    type = "UFQDN"
    fqdn = "%s@securitygeek.io"
    pre_shared_key = "%s"
}

data "%s" "%s" {
	id = "${%s.id}"
  }
`,
		// resource variables
		resourcetype.TrafficFilteringVPNCredentials,
		generatedName,
		generatedName,
		rEmail,
		rSharedKey,

		// data source variables
		resourcetype.TrafficFilteringVPNCredentials,
		generatedName,
		resourceTypeAndName,
	)
}
