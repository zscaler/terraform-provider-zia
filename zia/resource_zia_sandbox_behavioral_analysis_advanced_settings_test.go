package zia

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestResourceSandboxSettings_basic(t *testing.T) {
	resourceName := "zia_sandbox_behavioral_analysis.test"
	dataSourceName := "data.zia_sandbox_behavioral_analysis.list_all"

	// Define initial and updated hash sets
	initialHashes := []string{
		"42914d6d213a20a2684064be5c80ffa9",
		"c0202cf6aeab8437c638533d14563d35",
		"1ca31319721740ecb79f4b9ee74cd9b0",
		"2c373a7e86d0f3469849971e053bcc82",
		"40858748e03a544f6b562a687777397a",
		"465e89654a72256e7d1fb066388cc2a3",
		"47e7b297f020d53f7de7dc0f450e262d",
		"53d9af8829a9c7f6f177178885901c01",
		"9578c2be6437dcc8517e78a5de1fa975",
		"dfb689196faa945217a8929131f1d670",
		"8f9b7c1c2b84b8c71318b6776d31c9af",
		"a24bb61df75034769ffdda61c7a25926",
		"e5aea3b998644e394f506ac1f0f2f107",
		"1727de1b3d5636f1817d68ba0208fb50",
		"383498f810f0a992b964c19fc21ca398",
		"64990a45cf6b1b900c6b284bb54a1402",
		"97835760aa696d8ab7acbb5a78a5b013",
		"a8ab5aca96d260e649026e7fc05837bf",
		"c63a7c559870873133a84f0eb6ca54cd",
		"cc89100f20002801fa401b77dab0c512",
		"f8c110929606dca4c08ecaa9f9baf140",
		"f3dcf80b6251cfba1cd754006f693a73",
		"2c50efc0fef1601ce1b96b1b7cf991fb",
	}
	resource.ParallelTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSandboxSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAndDataSourceSandboxSettingsConfig(initialHashes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "file_hashes_to_be_blocked.#", fmt.Sprintf("%d", len(initialHashes))),
				),
			},
			{
				// Introduce an explicit delay before checking the data source
				Config: testAccResourceAndDataSourceSandboxSettingsConfig(initialHashes),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSandboxSettingsExists(resourceName),
					resource.TestCheckResourceAttr(dataSourceName, "file_hashes_to_be_blocked.#", fmt.Sprintf("%d", len(initialHashes))),
				),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
			},
			{
				Config: testAccResourceAndDataSourceSandboxSettingsConfig([]string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSandboxSettingsExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "file_hashes_to_be_blocked.#", "0"),
				),
			},
			{
				// Again, introduce a delay before checking the empty state in the data source
				Config: testAccResourceAndDataSourceSandboxSettingsConfig([]string{}),
				Check:  resource.TestCheckResourceAttr(dataSourceName, "file_hashes_to_be_blocked.#", "0"),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
			},
		},
	})
}

func testAccResourceAndDataSourceSandboxSettingsConfig(hashes []string) string {
	resourceConfig := testAccResourceSandboxSettingsConfig(hashes)
	dataSourceConfig := `data "zia_sandbox_behavioral_analysis" "list_all" {}`
	return resourceConfig + "\n" + dataSourceConfig
}

func testAccResourceSandboxSettingsConfig(hashes []string) string {
	var sb strings.Builder
	sb.WriteString(`resource "zia_sandbox_behavioral_analysis" "test" {`)
	if len(hashes) > 0 {
		sb.WriteString(`  file_hashes_to_be_blocked = [`)
		for i, hash := range hashes {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, hash))
		}
		sb.WriteString(`]`)
	} else {
		sb.WriteString(`  file_hashes_to_be_blocked = []`)
	}
	sb.WriteString(`}`)
	return sb.String()
}

func testAccCheckSandboxSettingsDestroy(s *terraform.State) error {
	// Implement checks if there are any resources that need to be cleaned up
	return nil
}

func testAccCheckSandboxSettingsExists(n string) resource.TestCheckFunc {
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
