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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_clusters"
)

func TestAccResourceVZENCluster_Basic(t *testing.T) {
	var clusters vzen_clusters.VZENClusters
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ServiceEdgeCluster)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVZENClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVZENClusterConfigure(
					resourceTypeAndName,
					initialName,
					variable.VzenStatus,
					variable.VzenType,
					variable.VzenIPAddress,
					variable.VzenSubnetMask,
					variable.VzenDefaultGateway,
					variable.VzenIpSecEnabled,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVZENClusterExists(resourceTypeAndName, &clusters),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.VzenStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.VzenType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", variable.VzenIPAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "subnet_mask", variable.VzenSubnetMask),
					resource.TestCheckResourceAttr(resourceTypeAndName, "default_gateway", variable.VzenDefaultGateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_sec_enabled", strconv.FormatBool(variable.VzenIpSecEnabled)),
				),
			},
			{
				Config: testAccCheckVZENClusterConfigure(
					resourceTypeAndName,
					updatedName,
					variable.VzenStatus,
					variable.VzenType,
					variable.VzenIPAddress,
					variable.VzenSubnetMask,
					variable.VzenDefaultGateway,
					variable.VzenIpSecEnabled,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVZENClusterExists(resourceTypeAndName, &clusters),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.VzenStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.VzenType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", variable.VzenIPAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "subnet_mask", variable.VzenSubnetMask),
					resource.TestCheckResourceAttr(resourceTypeAndName, "default_gateway", variable.VzenDefaultGateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_sec_enabled", strconv.FormatBool(variable.VzenIpSecEnabled)),
				),
			},
			{
				ResourceName:      resourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVZENClusterDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ServiceEdgeCluster {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := vzen_clusters.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckVZENClusterExists(resource string, rule *vzen_clusters.VZENClusters) resource.TestCheckFunc {
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

		receivedRule, err := vzen_clusters.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*rule = *receivedRule

		return nil
	}
}

func testAccCheckVZENClusterConfigure(resourceTypeAndName, generatedName, status, vzenType, address, mask, gateway string, ipsec bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
    name = "%s"
	status = "%s"
	type = "%s"
	ip_address = "%s"
	subnet_mask = "%s"
	default_gateway = "%s"
	ip_sec_enabled = "%s"
}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.ServiceEdgeCluster,
		resourceName,
		generatedName,
		status,
		vzenType,
		address,
		mask,
		gateway,
		strconv.FormatBool(ipsec),

		// data source variables
		resourcetype.ServiceEdgeCluster,
		resourceName,
		// Reference to the resource
		resourcetype.ServiceEdgeCluster, resourceName,
	)
}
