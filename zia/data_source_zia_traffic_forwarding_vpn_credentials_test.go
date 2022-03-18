package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/willguibr/terraform-provider-zia/zia/common/resourcetype"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/method"
	"github.com/willguibr/terraform-provider-zia/zia/common/testing/variable"
)

func TestAccDataSourceTrafficForwardingVPNCredentials_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.TrafficFilteringVPNCredentials)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrafficForwardingVPNCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrafficForwardingVPNCredentialsConfigure(resourceTypeAndName, generatedName, variable.VPNCredentialComments, variable.VPNCredentialTypeUFQDN, variable.VPNCredentialFQDN),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comments", resourceTypeAndName, "comments"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "fqdn", resourceTypeAndName, "fqdn"),
				),
			},
		},
	})
}
