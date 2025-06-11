package zia

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/nss_servers"
)

func TestAccResourceNSSServer_Basic(t *testing.T) {
	var nssServers nss_servers.NSSServers
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.NSSServer)

	initialName := "tf-acc-test-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNSSServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNSSServerConfigure(resourceTypeAndName, initialName, variable.NSSStatus, variable.NSSType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSSServerExists(resourceTypeAndName, &nssServers),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.NSSStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.NSSType),
				),
			},

			// Update test
			{
				Config: testAccCheckNSSServerConfigure(resourceTypeAndName, initialName, variable.NSSStatus, variable.NSSType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNSSServerExists(resourceTypeAndName, &nssServers),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.NSSStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.NSSType),
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

func testAccCheckNSSServerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.NSSServer {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := nss_servers.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckNSSServerExists(resource string, rule *nss_servers.NSSServers) resource.TestCheckFunc {
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
		service := apiClient.Service

		receivedServer, err := nss_servers.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedServer

		return nil
	}
}

func testAccCheckNSSServerConfigure(resourceTypeAndName, generatedName, status, nssType string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name   = "%s"
    status = "%s"
	type   = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.NSSServer,
		resourceName,
		generatedName,
		status,
		nssType,

		// data source variables
		resourcetype.NSSServer,
		resourceName,
		// Reference to the resource
		resourcetype.NSSServer, resourceName,
	)
}
