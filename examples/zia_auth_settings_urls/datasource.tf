data "zia_auth_settings_urls" "example" {}

output "zia_auth_settings_urls" {
  value = data.zia_auth_settings_urls.example
}