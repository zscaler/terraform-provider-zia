data "zia_url_filtering_rules" "example"{
    name = "Example"
}

output "zia_url_filtering_rules" {
    value = data.zia_url_filtering_rules.example
}
