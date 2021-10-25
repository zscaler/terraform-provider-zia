data "zia_firewall_filtering_rule" "example" {
    name = "Office 365 One Click Rule"
}

output "zia_firewall_filtering_rule" {
  value = data.zia_firewall_filtering_rule.example
}