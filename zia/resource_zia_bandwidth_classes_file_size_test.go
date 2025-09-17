package zia

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceBandwdithClassesFileSize_Basic(t *testing.T) {
	resourceName := "zia_bandwidth_classes_file_size.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceBandwdithClassesFileSizeDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create the resource with an initial file size
			{
				Config: testAccResourceBandwdithClassesFileSizeConfig("FILE_250MB"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "file_size", "FILE_250MB"),
					resource.TestCheckResourceAttr(resourceName, "name", "LARGE_FILE"),
					resource.TestCheckResourceAttr(resourceName, "type", "BANDWIDTH_CAT_LARGE_FILE"),
				),
			},
			// Step 2: Update the resource to a different file size
			{
				Config: testAccResourceBandwdithClassesFileSizeConfig("FILE_1GB"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "file_size", "FILE_1GB"),
				),
			},
			// Step 3: Import test
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckResourceBandwdithClassesFileSizeDestroy(s *terraform.State) error {
	// Since this resource updates an existing static object, no destroy check is needed
	return nil
}

func testAccResourceBandwdithClassesFileSizeConfig(fileSize string) string {
	return fmt.Sprintf(`
resource "zia_bandwidth_classes_file_size" "test" {
  file_size = "%s"
}
`, fileSize)
}
