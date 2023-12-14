// Add Hashes to the Sandbox
resource "zia_sandbox_behavioral_analysis" "this" {
  file_hashes_to_be_blocked = [
        "42914d6d213a20a2684064be5c80ffa9",
        "c0202cf6aeab8437c638533d14563d35",
  ]
}

// Erases All Hashes
resource "zia_sandbox_behavioral_analysis" "this" {}
