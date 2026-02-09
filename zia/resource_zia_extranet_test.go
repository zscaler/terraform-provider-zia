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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/extranet"
)

func TestAccResourceExtranet(t *testing.T) {
	var extranets extranet.Extranet
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.Extranet)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExtranetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckExtranetConfigure(resourceTypeAndName, initialName, variable.ExtranetDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExtranetExists(resourceTypeAndName, &extranets),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ExtranetDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "extranet_dns_list.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "extranet_ip_pool_list.#", "2"),
				),
			},

			// Update test
			{
				Config: testAccCheckExtranetConfigure(resourceTypeAndName, updatedName, variable.ExtranetDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExtranetExists(resourceTypeAndName, &extranets),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ExtranetDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "extranet_dns_list.#", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "extranet_ip_pool_list.#", "2"),
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

func testAccCheckExtranetDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.Extranet {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := extranet.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckExtranetExists(resource string, rule *extranet.Extranet) resource.TestCheckFunc {
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

		receivedRule, err := extranet.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckExtranetConfigure(resourceTypeAndName, generatedName, description string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name = "%s"
    description = "%s"
    extranet_dns_list {
        name                 = "DNS01"
        primary_dns_server   = "8.8.8.8"
        secondary_dns_server = "4.4.4.4"
        use_as_default       = true
    }
    extranet_dns_list {
        name                 = "DNS02"
        primary_dns_server   = "192.168.1.1"
        secondary_dns_server = "192.168.1.2"
        use_as_default       = false
    }
    extranet_ip_pool_list {
        name           = "TFS01"
        ip_start       = "10.0.0.1"
        ip_end         = "10.0.0.21"
        use_as_default = true
    }
    extranet_ip_pool_list {
        name           = "TFS02"
        ip_start       = "10.0.0.22"
        ip_end         = "10.0.0.43"
        use_as_default = false
    }
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.Extranet,
		resourceName,
		generatedName,
		description,

		// data source variables
		resourcetype.Extranet,
		resourceName,
		// Reference to the resource
		resourcetype.Extranet, resourceName,
	)
}
