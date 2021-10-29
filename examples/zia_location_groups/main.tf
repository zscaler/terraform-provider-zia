data "zia_location_groups" "example1"{
    name = "Corporate User Traffic Group"
}

output "zia_location_groups_corporate"{
    value = data.zia_location_groups.example1
}

data "zia_location_groups" "example2"{
    name = "Guest Wifi Group"
}

output "zia_location_groups_wifi"{
    value = data.zia_location_groups.example2
}

data "zia_location_groups" "example3"{
    name = "IoT Traffic Group"
}

output "zia_location_groups_iot"{
    value = data.zia_location_groups.example3
}

data "zia_location_groups" "example4"{
    name = "Server Traffic Group"
}

output "zia_location_groups_server"{
    value = data.zia_location_groups.example4
}
