package zia

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/resourcetype"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/method"
	"github.com/zscaler/terraform-provider-zia/v2/zia/common/testing/variable"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_notification_templates"
)

func TestAccResourceDLPNotificationTemplatesBasic(t *testing.T) {
	var dlpTemplates dlp_notification_templates.DlpNotificationTemplates
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.DLPNotificationTemplates)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDLPNotificationTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDLPNotificationTemplateConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPNotificationTemplateExists(resourceTypeAndName, &dlpTemplates),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "plain_text_message"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "html_message"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "attach_content", strconv.FormatBool(variable.DLPNoticationTemplateAttachContent)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tls_enabled", strconv.FormatBool(variable.DLPNoticationTemplateTLSEnabled)),
				),
			},

			// Update test
			{
				Config: testAccCheckDLPNotificationTemplateConfigure(resourceTypeAndName, generatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDLPNotificationTemplateExists(resourceTypeAndName, &dlpTemplates),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", "tf-acc-test-"+generatedName),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "plain_text_message"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "html_message"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "attach_content", strconv.FormatBool(variable.DLPNoticationTemplateAttachContent)),
					resource.TestCheckResourceAttr(resourceTypeAndName, "tls_enabled", strconv.FormatBool(variable.DLPNoticationTemplateTLSEnabled)),
				),
			},
		},
	})
}

func testAccCheckDLPNotificationTemplateDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.DLPNotificationTemplates {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		rule, err := apiClient.dlp_notification_templates.Get(id)

		if err == nil {
			return fmt.Errorf("id %d already exists", id)
		}

		if rule != nil {
			return fmt.Errorf("dlp dictionaries with id %d exists and wasn't destroyed", id)
		}
	}

	return nil
}

func testAccCheckDLPNotificationTemplateExists(resource string, dlpTemplate *dlp_notification_templates.DlpNotificationTemplates) resource.TestCheckFunc {
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
		receivedTemplate, err := apiClient.dlp_notification_templates.Get(id)

		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}
		*dlpTemplate = *receivedTemplate

		return nil
	}
}

func testAccCheckDLPNotificationTemplateConfigure(resourceTypeAndName, generatedName string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	name               = "tf-acc-test-%s"
	attach_content     = true
	tls_enabled        = true
	subject            = local.subject
	plain_text_message = local.msg_plain_text
	html_message       = local.msg_html
	}

locals {

	subject = <<SUBJECT
	DLP Violation: $${TRANSACTION_ID} $${RULENAME}
	SUBJECT

	msg_plain_text = <<MSGPLAINTEXT
		The attached content triggered a Web DLP rule for your organization.

		Transaction ID: $${TRANSACTION_ID}
		User Accessing the URL: $${USER}
		URL Accessed: $${URL}
		Posting Type: $${TYPE}
		DLP MD5: $${DLPMD5}
		Triggered DLP Violation Engines (assigned to the hit rule): $${ENGINES_IN_RULE}
		Triggered DLP Violation Dictionaries (assigned to the hit rule): $${DICTIONARIES}

		No action is required on your part.
	MSGPLAINTEXT

	msg_html = <<MSGHTML
	<!DOCTYPE html>
	<html>
		<head>
			<style>
				.user {color: rgb(1, 81, 152);}
				.url {color: rgb(1, 81, 152);}
				.postingtype {color: rgb(1, 81, 152);}
				.engines {color: rgb(1, 81, 152);}
				.dictionaries {color: rgb(1, 81, 152);}
			</style>
		</head>
		<body>
			The attached content triggered a Web DLP rule for your organization.
			<br/><br/>
			Transaction ID: <span class="transaction_id">$${TRANSACTION_ID}</span>
			<br/>
			User Accessing the URL: <span class="user">$${USER}</span>
			<br/>
			URL Accessed: <span class="url">$${URL}</span>
			<br/>
			Posting Type: <span class="postingtype">$${TYPE}</span>
			<br/>
			DLP MD5: <span class="dlpmd5">$${DLPMD5}</span>
			<br/>
			Triggered DLP Violation Engines (assigned to the hit rule): <span class="engines">$${ENGINES_IN_RULE}</span>
			<br/>
			Triggered DLP Violation Dictionaries (assigned to the hit rule): <span class="dictionaries">$${DICTIONARIES}</span>
			<br/><br/>
			No action is required on your part.
			<br/><br/>
		</body>
	</html>
	MSGHTML
	}

data "%s" "%s" {
	id = "${%s.id}"
}

`,
		// resource variables
		resourcetype.DLPNotificationTemplates,
		generatedName,
		generatedName,

		// data source variables
		resourcetype.DLPNotificationTemplates,
		generatedName,
		resourceTypeAndName,
	)
}
