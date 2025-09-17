package zia

/*
func TestAccResourceVZENNodeBasic(t *testing.T) {
	var labels vzen_nodes.VZENNodes
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ServiceEdgeNode)

	initialName := "tf-acc-test-" + generatedName
	updatedName := "tf-updated-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVZENNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVZENNodeConfigure(resourceTypeAndName, initialName, variable.VzenStatus, variable.VzenNodeType, variable.VzenNodeIPAddress, variable.VzenNodeSubnetMask, variable.VzenNodeDefaultGateway, variable.VzenNodeInProduction, variable.VzenNodeDeploymentMode, variable.VZenSKUType, variable.VzenOnDemandSupportTunnel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVZENNodeExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.VzenStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.VzenNodeType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", variable.VzenNodeIPAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "subnet_mask", variable.VzenNodeSubnetMask),
					resource.TestCheckResourceAttr(resourceTypeAndName, "default_gateway", variable.VzenNodeDefaultGateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "deployment_mode", variable.VzenNodeDeploymentMode),
					resource.TestCheckResourceAttr(resourceTypeAndName, "vzen_sku_type", variable.VZenSKUType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "in_production", strconv.FormatBool(variable.VzenNodeInProduction)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "on_demand_support_tunnel_enabled", strconv.FormatBool(variable.VzenOnDemandSupportTunnel)),
				),
			},

			// Update test
			{
				Config: testAccCheckVZENNodeConfigure(resourceTypeAndName, updatedName, variable.VzenStatus, variable.VzenNodeType, variable.VzenNodeIPAddress, variable.VzenNodeSubnetMask, variable.VzenNodeDefaultGateway, variable.VzenNodeInProduction, variable.VzenNodeDeploymentMode, variable.VZenSKUType, variable.VzenOnDemandSupportTunnel),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVZENNodeExists(resourceTypeAndName, &labels),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "status", variable.VzenStatus),
					resource.TestCheckResourceAttr(resourceTypeAndName, "type", variable.VzenNodeType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "ip_address", variable.VzenNodeIPAddress),
					resource.TestCheckResourceAttr(resourceTypeAndName, "subnet_mask", variable.VzenNodeSubnetMask),
					resource.TestCheckResourceAttr(resourceTypeAndName, "default_gateway", variable.VzenNodeDefaultGateway),
					resource.TestCheckResourceAttr(resourceTypeAndName, "deployment_mode", variable.VzenNodeDeploymentMode),
					resource.TestCheckResourceAttr(resourceTypeAndName, "vzen_sku_type", variable.VZenSKUType),
					resource.TestCheckResourceAttr(resourceTypeAndName, "in_production", strconv.FormatBool(variable.VzenNodeInProduction)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "on_demand_support_tunnel_enabled", strconv.FormatBool(variable.VzenOnDemandSupportTunnel)),
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

func testAccCheckVZENNodeDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.ServiceEdgeNode {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := vzen_nodes.Get(context.Background(), service, id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("rule with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckVZENNodeExists(resource string, node *vzen_nodes.VZENNodes) resource.TestCheckFunc {
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

		receivedNode, err := vzen_nodes.Get(context.Background(), service, id)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*node = *receivedNode

		return nil
	}
}

func testAccCheckVZENNodeConfigure(resourceTypeAndName, generatedName, status, nodeType, nodeIp, nodeMask, nodeGateway string, inProduction bool, nodeDeploymentMode, vzenSKUType string, onDemandSupport bool) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1] // Extract the resource name
	return fmt.Sprintf(`

resource "%s" "%s" {
  name                              = "%s"
  status                            = "%s"
  type                              = "%s"
  ip_address                        = "%s"
  subnet_mask                       = "%s"
  default_gateway                   = "%s"
  in_production                     = %t
  deployment_mode                   = "%s"
  vzen_sku_type                     = "%s"
  on_demand_support_tunnel_enabled  = %t

}

data "%s" "%s" {
	id = "${%s.%s.id}"
  }
`,
		// resource variables
		resourcetype.ServiceEdgeNode,
		resourceName,
		generatedName,
		status,
		nodeType,
		nodeIp,
		nodeMask,
		nodeGateway,
		inProduction,
		nodeDeploymentMode,
		vzenSKUType,
		onDemandSupport,

		// data source variables
		resourcetype.ServiceEdgeNode,
		resourceName,
		// Reference to the resource
		resourcetype.ServiceEdgeNode, resourceName,
	)
}
*/
