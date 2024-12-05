package zia

/*
import (
	"testing"

	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: NEEDS FIXING BY ENGINEERING: "{"code":"RBA_LIMITED","message":"Functional scope restriction requires PROXY_GATEWAY"}"
// ONEAPI-915 - ZIA API Tests – Results (RBA_LIMITED) and Other Errors


func TestAccDataSourceForwardingControlRule_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.ForwardingControlRule)

	// Generate Rule Label HCL Resource
	ruleLabelTypeAndName, _, ruleLabelGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.RuleLabels)
	ruleLabelHCL := testAccCheckRuleLabelsConfigure(ruleLabelTypeAndName, ruleLabelGeneratedName, variable.RuleLabelDescription)

	// Generate Source IP Group HCL Resource
	sourceIPGroupTypeAndName, _, sourceIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringSourceGroup)
	sourceIPGroupHCL := testAccCheckFWIPSourceGroupsConfigure(sourceIPGroupTypeAndName, sourceIPGroupGeneratedName, variable.FWSRCGroupDescription)

	// Generate Destination IP Group HCL Resource
	dstIPGroupTypeAndName, _, dstIPGroupGeneratedName := method.GenerateRandomSourcesTypeAndName(resourcetype.FWFilteringDestinationGroup)
	dstIPGroupHCL := testAccCheckFWIPDestinationGroupsConfigure(dstIPGroupTypeAndName, dstIPGroupGeneratedName, variable.FWDSTGroupDescription, variable.FWDSTGroupTypeDSTNFQDN)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckForwardingControlRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckForwardingControlRuleConfigure(resourceTypeAndName, generatedName, generatedName, variable.FowardingControlDescription, variable.FowardingControlType, variable.FowardingControlState, ruleLabelTypeAndName, ruleLabelHCL, sourceIPGroupTypeAndName, sourceIPGroupHCL, dstIPGroupTypeAndName, dstIPGroupHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "description", resourceTypeAndName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "type", resourceTypeAndName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "state", resourceTypeAndName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "order", resourceTypeAndName, "order"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "nw_services.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "departments.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "groups.#", "2"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "labels.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "src_ip_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "dest_ip_groups.#", "1"),
				),
			},
		},
	})
}
*/
