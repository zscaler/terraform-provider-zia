package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
)

func TestAccDataSourceSubscriptionAlerts_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.SubscriptionAlerts)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubscriptionAlertsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSubscriptionAlertsConfigure(resourceTypeAndName, generatedName, variable.AlertDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "email", resourceTypeAndName, "email"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "pt0_severities", resourceTypeAndName, "pt0_severities"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "secure_severities", resourceTypeAndName, "secure_severities"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "manage_severities", resourceTypeAndName, "manage_severities"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "comply_severities", resourceTypeAndName, "comply_severities"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "system_severities", resourceTypeAndName, "system_severities"),
				),
			},
		},
	})
}
