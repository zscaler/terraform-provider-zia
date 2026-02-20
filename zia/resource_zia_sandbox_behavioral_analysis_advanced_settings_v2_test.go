package zia

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceSandboxSettingsV2_basic(t *testing.T) {
	resourceName := "zia_sandbox_behavioral_analysis_v2.test"
	dataSourceName := "data.zia_sandbox_behavioral_analysis_v2.list_all"

	initialHashes := []struct {
		URL        string
		URLComment string
		Type       string
	}{
		{"42914d6d213a20a2684064be5c80ffa9", "hash_allow_1", "CUSTOM_FILEHASH_ALLOW"},
		{"c0202cf6aeab8437c638533d14563d35", "hash_allow_2", "CUSTOM_FILEHASH_ALLOW"},
		{"1ca31319721740ecb79f4b9ee74cd9b0", "hash_deny_1", "CUSTOM_FILEHASH_DENY"},
		{"2c373a7e86d0f3469849971e053bcc82", "hash_deny_2", "CUSTOM_FILEHASH_DENY"},
	}

	updatedHashes := []struct {
		URL        string
		URLComment string
		Type       string
	}{
		{"42914d6d213a20a2684064be5c80ffa9", "hash_allow_1", "CUSTOM_FILEHASH_ALLOW"},
		{"9578c2be6437dcc8517e78a5de1fa975", "hash_deny_3", "CUSTOM_FILEHASH_DENY"},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSandboxSettingsV2Destroy,
		Steps: []resource.TestStep{
			// Step 1: Create with initial hashes
			{
				Config: testAccSandboxSettingsV2Config(initialHashes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxSettingsV2Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "md5_hash_value_list.#", fmt.Sprintf("%d", len(initialHashes))),
				),
			},
			// Step 2: Verify data source reads the hashes
			{
				Config: testAccSandboxSettingsV2WithDataSource(initialHashes),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSandboxSettingsV2Exists(resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "md5_hash_value_list.#", fmt.Sprintf("%d", len(initialHashes))),
				),
			},
			// Step 3: Update to a different set of hashes
			{
				Config: testAccSandboxSettingsV2Config(updatedHashes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxSettingsV2Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "md5_hash_value_list.#", fmt.Sprintf("%d", len(updatedHashes))),
				),
			},
			// Step 4: Clear all hashes (empty list)
			{
				Config: testAccSandboxSettingsV2EmptyConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxSettingsV2Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "md5_hash_value_list.#", "0"),
				),
			},
			// Step 5: Verify data source returns empty after clearing
			{
				Config: testAccSandboxSettingsV2EmptyWithDataSource(),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "md5_hash_value_list.#", "0"),
				),
			},
			// Step 6: Import
			{
				ResourceName:      resourceName,
				ImportStateId:     "sandbox_settings",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSandboxSettingsV2Config(hashes []struct {
	URL        string
	URLComment string
	Type       string
}) string {
	var sb strings.Builder
	sb.WriteString("resource \"zia_sandbox_behavioral_analysis_v2\" \"test\" {\n")
	for _, h := range hashes {
		sb.WriteString(fmt.Sprintf("  md5_hash_value_list {\n    url         = %q\n    url_comment = %q\n    type        = %q\n  }\n", h.URL, h.URLComment, h.Type))
	}
	sb.WriteString("}\n")
	return sb.String()
}

func testAccSandboxSettingsV2EmptyConfig() string {
	return `resource "zia_sandbox_behavioral_analysis_v2" "test" {}`
}

func testAccSandboxSettingsV2WithDataSource(hashes []struct {
	URL        string
	URLComment string
	Type       string
}) string {
	return testAccSandboxSettingsV2Config(hashes) + "\n" + `data "zia_sandbox_behavioral_analysis_v2" "list_all" {}`
}

func testAccSandboxSettingsV2EmptyWithDataSource() string {
	return testAccSandboxSettingsV2EmptyConfig() + "\n" + `data "zia_sandbox_behavioral_analysis_v2" "list_all" {}`
}

func testAccCheckSandboxSettingsV2Destroy(s *terraform.State) error {
	return nil
}

func testAccCheckSandboxSettingsV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Sandbox Setting ID is set")
		}
		return nil
	}
}
