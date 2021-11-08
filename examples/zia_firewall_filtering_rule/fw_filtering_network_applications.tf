resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    nw_applications = [ data.zia_firewall_filtering_network_application_groups.apns.id,
                        data.zia_firewall_filtering_network_application_groups.dict.id
                    ]
}

data "zia_firewall_filtering_network_application" "apns"{
    id = "APNS"
    locale="en-US"
}

data "zia_firewall_filtering_network_application" "dict"{
    id = "DICT"
}