data "zia_rule_labels" "example" {
	name = "Example"
}

output "zia_rule_labels" {
	value = data.zia_rule_labels.example
}