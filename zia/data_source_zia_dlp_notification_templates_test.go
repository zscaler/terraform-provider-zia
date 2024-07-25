package zia

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v3/zia/common/testing/variable"
)

func TestAccDataSourceDLPNotificationTemplates_Basic(t *testing.T) {
	resourceTypeAndName, dataSourceTypeAndName, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPNotificationTemplates)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPNotificationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDLPNotificationTemplateConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "id", resourceTypeAndName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "name", resourceTypeAndName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "plain_text_message", resourceTypeAndName, "plain_text_message"),
					resource.TestCheckResourceAttrPair(dataSourceTypeAndName, "html_message", resourceTypeAndName, "html_message"),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "attach_content", strconv.FormatBool(variable.DLPNoticationTemplateAttachContent)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "tls_enabled", strconv.FormatBool(variable.DLPNoticationTemplateTLSEnabled)),
				),
			},
		},
	})
}
