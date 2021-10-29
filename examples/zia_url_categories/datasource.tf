data "zia_url_categories" "example"{
    id = "CUSTOM_08"
}

output "zia_url_categories"{
    value = data.zia_url_categories.example
}
