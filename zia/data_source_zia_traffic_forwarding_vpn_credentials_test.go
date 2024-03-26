package zia

/*
func TestAccDataSourceTrafficForwardingVPNCredentials_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingVPNCredentials)
	rEmail := acctest.RandomWithPrefix("tf-acc-test-")
	rSharedKey := acctest.RandString(20)

	rIP, _ := acctest.RandIpAddress("121.234.54.0/25")
	staticIPTypeAndName, _, staticIPGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficForwardingStaticIP)
	staticIPResourceHCL := testAccCheckTrafficForwardingStaticIPConfigure(staticIPTypeAndName, staticIPGeneratedName, rIP, variable.StaticRoutableIP, variable.StaticGeoOverride)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingVPNCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsUFQDNConfigure(resourceTypeAndName, generatedName, rEmail, rSharedKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comments", resourceTypeAndName, "comments"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "fqdn", resourceTypeAndName, "fqdn"),
				),
			},
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsIPConfigure(resourceTypeAndName, generatedName, staticIPResourceHCL, staticIPTypeAndName, rSharedKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comments", resourceTypeAndName, "comments"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "ip_address", resourceTypeAndName, "ip_address"),
				),
			},
		},
	})
}
*/
