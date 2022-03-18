package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccResourceTrafficForwardingVPNCredentialsBasic(t *testing.T) {
	var credentials vpncredentials.VPNCredentials
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringVPNCredentials)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingVPNCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, variable.VPNCredentialComments, variable.VPNCredentialTypeUFQDN, variable.VPNCredentialFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingVPNCredentialsExists(resourceTypeAndName, &credentials),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", variable.VPNCredentialComments),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.VPNCredentialTypeUFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "fqdn", variable.VPNCredentialFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pre_shared_key", variable.VPNCredentialPreSharedKey),
				),
			},

			// Update test
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, variable.VPNCredentialComments, variable.VPNCredentialTypeUFQDN, variable.VPNCredentialFQDN),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrafficForwardingVPNCredentialsExists(resourceTypeAndName, &credentials),
					resource.TestCheckResourceAttr(resourceTypeAndName, "comments", variable.VPNCredentialComments),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.VPNCredentialTypeUFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "fqdn", variable.VPNCredentialFQDN),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pre_shared_key", variable.VPNCredentialPreSharedKey),
				),
			},
		},
	})
}

func testAccCheckTrafficForwardingVPNCredentialsDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.TrafficFilteringStaticIP {
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

func testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, comments, typeUFQDN, fqdn string) string {
	return fmt.Sprintf(`
// network application group resource
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// resource variables
		TrafficForwardingVPNCredentialsResourceHCL(generatedName, comments, typeUFQDN, fqdn),

		// data source variables
		resourcetype.TrafficFilteringVPNCredentials,
		generatedName,
		resourceTypeAndName,
	)
}

func TrafficForwardingVPNCredentialsResourceHCL(generatedName, comments, typeUFQDN, fqdn string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	comments = "%s"
    type = "%s"
    fqdn = "%s"
    pre_shared_key = "Password@123!"
}
`,
		// resource variables
		resourcetype.TrafficFilteringVPNCredentials,
		generatedName,
		variable.VPNCredentialComments,
		typeUFQDN,
		fqdn,
	)
}
