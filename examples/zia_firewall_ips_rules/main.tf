data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
    name = "ZSCALER_PROXY_NW_SERVICES"
}

data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_group_management" "normal_internet" {
    name = "Normal_Internet"
}

data "zia_firewall_filtering_time_window" "work_hours" {
    name = "Work hours"
}

resource "zia_firewall_ips_rule" "example" {
    name = "Example_IPS_Rule01"
    description = "Example_IPS_Rule01"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    threat_categories {
        id = [ 66 ]
    }
    nw_services {
        id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
    }
    departments {
        id = [ data.zia_department_management.engineering.id ]
    }
    groups {
        id = [ data.zia_group_management.normal_internet.id ]
    }
    time_windows {
        id = [ data.zia_firewall_filtering_time_window.work_hours.id ]
    }
}