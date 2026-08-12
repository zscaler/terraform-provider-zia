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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/pacfiles"
)

const testAccPacContentV1 = `function FindProxyForURL(url, host) { return "DIRECT"; }`
const testAccPacContentV2 = `function FindProxyForURL(url, host) { /* updated */ return "DIRECT"; }`

func TestAccResourcePacFilesBasic(t *testing.T) {
	var pacFile pacfiles.PACFileConfig
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.PacFiles)

	initialName := "tf-acc-test-" + generatedName

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacFilesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacFilesConfigure(resourceTypeAndName, initialName, variable.PacFileDescription, testAccPacContentV1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacFilesExists(resourceTypeAndName, &pacFile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "description", variable.PacFileDescription),
					resource.TestCheckResourceAttr(resourceTypeAndName, "domain", variable.PacFileDomain),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pac_version", "1"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pac_version_status", "DEPLOYED"),
					resource.TestCheckResourceAttrSet(resourceTypeAndName, "pac_url"),
				),
			},

			// Content update creates a new version of the same PAC file.
			{
				Config: testAccCheckPacFilesConfigure(resourceTypeAndName, initialName, variable.PacFileDescription, testAccPacContentV2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacFilesExists(resourceTypeAndName, &pacFile),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", initialName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pac_version", "2"),
					resource.TestCheckResourceAttr(resourceTypeAndName, "pac_version_status", "DEPLOYED"),
				),
			},
			// Import test
			{
				ResourceName:            resourceTypeAndName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_version"},
			},
		},
	})
}

func testAccCheckPacFilesDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*Client)
	service := apiClient.Service

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourcetype.PacFiles {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			log.Println("Failed in conversion with error:", err)
			return err
		}

		allFiles, err := pacfiles.GetPacFiles(context.Background(), service, "pac_content")
		if err != nil {
			return err
		}
		for _, pacFile := range allFiles {
			if pacFile.ID == id {
				return fmt.Errorf("pac file with id %d exists and wasn't destroyed", id)
			}
		}
	}

	return nil
}

func testAccCheckPacFilesExists(resource string, pacFile *pacfiles.PACFileConfig) resource.TestCheckFunc {
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

		allFiles, err := pacfiles.GetPacFiles(context.Background(), service, "")
		if err != nil {
			return fmt.Errorf("failed fetching pac files: %s", err)
		}
		for _, f := range allFiles {
			if f.ID == id {
				*pacFile = f
				return nil
			}
		}
		return fmt.Errorf("pac file %d not found", id)
	}
}

func testAccCheckPacFilesConfigure(resourceTypeAndName, generatedName, description, pacContent string) string {
	resourceName := strings.Split(resourceTypeAndName, ".")[1]
	return fmt.Sprintf(`

resource "%s" "%s" {
    name               = "%s"
    description        = "%s"
    domain             = "%s"
    pac_commit_message = "%s"
    pac_version_status = "DEPLOYED"
    pac_content        = <<-EOT
%s
EOT
}

data "%s" "%s" {
	id = "${%s.%s.id}"
}
`,
		// resource variables
		resourcetype.PacFiles,
		resourceName,
		generatedName,
		description,
		variable.PacFileDomain,
		variable.PacFileCommitMessage,
		pacContent,

		// data source variables
		resourcetype.PacFiles,
		resourceName,
		// Reference to the resource
		resourcetype.PacFiles, resourceName,
	)
}
