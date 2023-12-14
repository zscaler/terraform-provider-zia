package zia

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZiaSandboxFileSubmission_basic(t *testing.T) {
	baseURL := "https://github.com/SecurityGeekIO/malware-samples/raw/main/"
	fileNames := []string{
		"2a961d4e5a2100570c942ed20a29735b.bin",
		"327bd8a60fb54aaaba8718c890dda09d.bin",
		"7665f6ee9017276dd817d15212e99ca7.bin",
		"cefb4323ba4deb9dea94dcbe3faa139f.bin",
		"8356bd54e47b000c5fdcf8dc5f6a69fa.apk",
		"841abdc66ea1e208f63d717ebd11a5e9.apk",
		"test-pe-file.exe",
	}

	for _, fileName := range fileNames {
		fileURL := baseURL + fileName
		localFilePath, err := downloadTestFile(fileURL)
		if err != nil {
			t.Fatalf("Error downloading test file from %s: %v", fileURL, err)
		}

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckSandboxSubmissionDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCheckZiaSandboxFileSubmissionConfig(localFilePath, "submit", true),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckZiaSandboxFileSubmissionExists("zia_sandbox_file_submission.this"),
					),
				},
				{
					Config: testAccCheckZiaSandboxFileSubmissionConfig(localFilePath, "discan", false),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckZiaSandboxFileSubmissionExists("zia_sandbox_file_submission.this"),
					),
				},
			},
		})
	}
}

func testAccCheckZiaSandboxFileSubmissionConfig(filePath, submissionMethod string, force bool) string {
	return fmt.Sprintf(`
resource "zia_sandbox_file_submission" "this" {
    file_path          = "%s"
    submission_method  = "%s"
    force              = %t
}
`, filePath, submissionMethod, force)
}

func testAccCheckZiaSandboxFileSubmissionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Since there's no read operation, just check the resource exists in the state
		if _, ok := s.RootModule().Resources[n]; !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		return nil
	}
}

func downloadTestFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpDir, err := os.MkdirTemp("", "testFile")
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(tmpDir, "downloadedFile")
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func testAccCheckSandboxSubmissionDestroy(s *terraform.State) error {
	// Implement if there's anything to check upon resource destruction
	return nil
}
