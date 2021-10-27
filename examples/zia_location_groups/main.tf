data "zia_location_groups" "example"{
    name = "Corporate User Traffic Group"
}

output "zia_location_groups"{
    value = data.zia_location_groups.example
}

