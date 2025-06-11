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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxies"
)

func TestAccResourceForwardingControlProxies_Basic(t *testing.T) {
	var proxies proxies.Proxies
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ForwardingControlProxies)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingControlProxiesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingControlProxiesConfigure(resourceTypeAndName, initialName, variable.ProxyDescription, variable.ProxyType, variable.ProxyAddress, variable.ProxyPort, variable.ProxyInsertXauHeader, variable.ProxyBase64EncodeXauHeader),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingControlProxiesExists(resourceTypeAndName, &proxies),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ProxyDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.ProxyType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "address", variable.ProxyAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "port", strconv.Itoa(variable.ProxyPort)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "insert_xau_header", strconv.FormatBool(variable.ProxyInsertXauHeader)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "base64_encode_xau_header", strconv.FormatBool(variable.ProxyBase64EncodeXauHeader)),
				),
			},

			// Update test
			{
				Config: testAccCheckForwardingControlProxiesConfigure(resourceTypeAndName, updatedName, variable.ProxyDescription, variable.ProxyType, variable.ProxyAddress, variable.ProxyPort, variable.ProxyInsertXauHeader, variable.ProxyBase64EncodeXauHeader),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckForwardingControlProxiesExists(resourceTypeAndName, &proxies),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.ProxyDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.ProxyType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "address", variable.ProxyAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "port", strconv.Itoa(variable.ProxyPort)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "insert_xau_header", strconv.FormatBool(variable.ProxyInsertXauHeader)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "base64_encode_xau_header", strconv.FormatBool(variable.ProxyBase64EncodeXauHeader)),
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

func testAccCheckForwardingControlProxiesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ForwardingControlProxies {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		proxy, err := proxies.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if proxy != nil {
			return fmt.Errorf("proxy with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckForwardingControlProxiesExists(resource string, proxy *proxies.Proxies) resource.TestCheckFunc {
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

		receivedProxy, err := proxies.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*proxy = *receivedProxy

		return nil
	}
}

func testAccCheckForwardingControlProxiesConfigure(resourceTypeAndName, generatedName, description, proxyType, address string, port int, xauHeader, encodeXauHeader bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	return fmt.Sprintf(`

resource "%s" "%s" {
    name = "%s"
    description = "%s"
	type = "%s"
	address = "%s"
	port = "%d"
	insert_xau_header = "%s"
	base64_encode_xau_header = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.ForwardingControlProxies,
		resourceName,
		generatedName,
		description,
		proxyType,
		address,
		port,
		strconv.FormatBool(xauHeader),
		strconv.FormatBool(encodeXauHeader),

		// data source variables
		resourcetype.ForwardingControlProxies,
		resourceName,
		// Reference to the resource
		resourcetype.ForwardingControlProxies, resourceName,
	)
}
