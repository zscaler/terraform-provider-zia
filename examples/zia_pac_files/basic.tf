resource "zia_pac_files" "this" {
  name               = "example_pac_file"
  description        = "Example hosted PAC file"
  domain             = "acme.com"
  pac_commit_message = "Initial version"
  pac_version_status = "DEPLOYED"
  pac_content        = <<-EOT
function FindProxyForURL(url, host) {
    return "PROXY $${GATEWAY_FX}:80; PROXY $${SECONDARY_GATEWAY_FX}:80; DIRECT";
}
EOT
}
