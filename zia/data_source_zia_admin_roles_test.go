package zia

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
)

func TestAccDataSourceAdminRoles_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.AdminRoles)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAdminRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAdminRolesConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "rank", resourceTypeAndName, "rank"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "alerting_access", resourceTypeAndName, "alerting_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "dashboard_access", resourceTypeAndName, "dashboard_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "report_access", resourceTypeAndName, "report_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "analysis_access", resourceTypeAndName, "analysis_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "username_access", resourceTypeAndName, "username_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "device_info_access", resourceTypeAndName, "device_info_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "admin_acct_access", resourceTypeAndName, "admin_acct_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "policy_access", resourceTypeAndName, "policy_access"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "role_type", resourceTypeAndName, "role_type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "logs_limit", resourceTypeAndName, "logs_limit"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "report_time_duration", resourceTypeAndName, "report_time_duration"),
				),
			},
		},
	})
}
