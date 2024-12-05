# Example Usage - Create Firewall DNS Rules - Redirect Action

resource "zia_firewall_dns_rule" "this" {
    name = "Example_DNS_Rule01"
    description = "Example_DNS_Rule01"
    action = "REDIR_REQ"
    state = "ENABLED"
    order = 10
    rank = 7
    redirect_ip = "8.8.8.8"
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
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

data "zia_department_management" "engineering" {
 name = "Engineering"
}

data "zia_group_management" "normal_internet" {
    name = "Normal_Internet"
}

data "zia_firewall_filtering_time_window" "work_hours" {
    name = "Work hours"
}

# Example Usage - Create Firewall DNS Rules - Redirect Request DOH

resource "zia_firewall_dns_rule" "this2" {
    name = "Example_DNS_Rule02"
    description = "Example_DNS_Rule02"
    action = "REDIR_REQ_DOH"
    state = "ENABLED"
    order = 12
    rank = 7
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
    dns_gateway {
      id = 18207342
      name = "DNS_GW01"
    }
}

# Example Usage - Create Firewall DNS Rules - Redirect Request DOH

resource "zia_firewall_dns_rule" "this2" {
    name = "Example_DNS_Rule02"
    description = "Example_DNS_Rule02"
    action = "REDIR_REQ_DOH"
    state = "ENABLED"
    order = 12
    rank = 7
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
    dns_gateway {
      id = 18207342
      name = "DNS_GW01"
    }
}

# Example Usage - Create Firewall DNS Rules - Redirect TCP Request

resource "zia_firewall_dns_rule" "this3" {
    name = "Example_DNS_Rule03"
    description = "Example_DNS_Rule03"
    action = "REDIR_REQ_TCP"
    state = "ENABLED"
    order = 13
    rank = 7
    dest_countries = ["CA", "US"]
    source_countries = ["CA", "US"]
    protocols = ["ANY_RULE"]
    dns_gateway {
      id = 18207342
      name = "DNS_GW01"
    }
}