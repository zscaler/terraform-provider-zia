data "zia_firewall_filtering_network_application" "apns"{
    id = "APNS"
    locale="en-US"
}

output "zia_firewall_filtering_network_application_apns"{
    value = data.zia_firewall_filtering_network_application.apns
}

data "zia_firewall_filtering_network_application" "dict"{
    id = "DICT"
}

output "zia_firewall_filtering_network_application_dict"{
    value = data.zia_firewall_filtering_network_application.dict
}