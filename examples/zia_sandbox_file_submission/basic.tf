// Submit raw or archive files
resource "zia_sandbox_file_submission" "this" {
  file_path     = "test-pe-file.exe"
  submission_method = "submit"
  force = true
}

// Submits raw or archive for out-of-band file inspection
resource "zia_sandbox_file_submission" "this" {
  file_path     = "test-pe-file.exe"
  submission_method = "submit"
  force = true
}