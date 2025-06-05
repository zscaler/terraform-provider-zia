package zia

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v4/zia/common/testing/variable"
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

					// Custom check for plain_text_message
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[dataSourceTypeAndName]
						got := rs.Primary.Attributes["plain_text_message"]
						if !strings.Contains(got, "Transaction ID: ${TRANSACTION_ID}") {
							return fmt.Errorf("plain_text_message does not contain expected placeholder. Got: %q", got)
						}
						if !strings.Contains(got, "User Accessing the URL: ${USER}") {
							return fmt.Errorf("plain_text_message missing expected user line. Got: %q", got)
						}
						return nil
					},

					// Custom check for html_message
					func(s *terraform.State) error {
						rs := s.RootModule().Resources[dataSourceTypeAndName]
						got := rs.Primary.Attributes["html_message"]
						if !strings.Contains(got, `<span class="transaction_id">${TRANSACTION_ID}</span>`) {
							return fmt.Errorf("html_message does not contain expected HTML transaction ID span. Got: %q", got)
						}
						if !strings.Contains(got, `<span class="url">${URL}</span>`) {
							return fmt.Errorf("html_message missing expected URL span. Got: %q", got)
						}
						return nil
					},

					resource.TestCheckResourceAttr(dataSourceTypeAndName, "attach_content", strconv.FormatBool(variable.DLPNoticationTemplateAttachContent)),
					resource.TestCheckResourceAttr(dataSourceTypeAndName, "tls_enabled", strconv.FormatBool(variable.DLPNoticationTemplateTLSEnabled)),
				),
			},
		},
	})
}
