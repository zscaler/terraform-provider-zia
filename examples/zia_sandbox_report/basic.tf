// Summary Report for a MD5 Hash of a file that was analyzed by the Sandbox
data "zia_sandbox_report" "this" {
  md5_hash = "0D5CBE46860E2BD6521D06BED7242FC0"
  details = "summary"
}

// Detailed Report for a MD5 Hash of a file that was analyzed by the Sandbox
data "zia_sandbox_report" "this" {
  md5_hash = "0D5CBE46860E2BD6521D06BED7242FC0"
  details = "full"
}