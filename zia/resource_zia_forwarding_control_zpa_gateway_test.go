package zia

/*
import (
	"fmt"
	"log"
	"strconv"
	"testing"

	zpa "github.com/zscaler/terraform-provider-zpa"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/zpa_gateways"
)

func TestAccResourceForwardingControlZPAGatewayBasic(t *testing.T) {
	var groups zpa_gateways.ZPAGateways
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ForwardingControlZPAGateway)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			// "zia": zia.Provider(), // Assuming 'zia' is the alias for your primary provider
			"zpa": zpa.Provider(), // This sets up the secondary provider
		},
		CheckDestroy: testAccCheckForwardingControlZPAGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingControlZPAGatewayConfigure(resourceTypeAndName, generatedName, variable.FowardingControlZPAGWDescription, variable.FowardingControlZPAGWType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingControlZPAGatewayExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FowardingControlZPAGWDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FowardingControlZPAGWType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "zpa_server_group.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "zpa_app_segments.#", "1"),
				),
			},

			// Update test
			{
				Config: testAccCheckForwardingControlZPAGatewayConfigure(resourceTypeAndName, generatedName, variable.FowardingControlZPAGWDescription, variable.FowardingControlZPAGWType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingControlZPAGatewayExists(resourceTypeAndName, &groups),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.FowardingControlZPAGWDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.FowardingControlZPAGWType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "zpa_server_group.#", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "zpa_app_segments.#", "1"),
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

func testAccCheckForwardingControlZPAGatewayDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ForwardingControlZPAGateway {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.zpa_gateways.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("forwarding control zpa gateway with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckForwardingControlZPAGatewayExists(resource string, gw *zpa_gateways.ZPAGateways) resource.TestCheckFunc {
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
		receivedGw, err := apiClient.zpa_gateways.Get(id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*gw = *receivedGw

		return nil
	}
}

func testAccCheckForwardingControlZPAGatewayConfigure(resourceTypeAndName, generatedName, description, gwType string) string {
	return fmt.Sprintf(`

	resource "%s" "%s" {
	name        = "tf-acc-test-%s"
	description = "%s"
	type        = "%s"
    zpa_server_group {
        id = [ 216196257331370183 ]
    }

	zpa_application_segments {
        id = [ 216196257331370184 ]
    }
}

data "%s" "%s" {
	id = "${%s.id}"
  }
`,
		// resource variables
		resourcetype.ForwardingControlZPAGateway,
		generatedName,
		generatedName,
		description,
		gwType,

		// data source variables
		resourcetype.ForwardingControlZPAGateway,
		generatedName,
		resourceTypeAndName,
	)
}
*/
