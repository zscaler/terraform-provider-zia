resource "zia_nat_control_rules" "this" {
    name = "DNAT_01"
    description = "DNAT_01"
    order=1
    rank=7
    state = "ENABLED"
    redirect_port="2000"
    redirect_ip="1.1.1.1"
    src_ips=["192.168.100.0/24", "192.168.200.1"]
    dest_addresses=["3.217.228.0-3.217.231.255", "3.235.112.0-3.235.119.255", "35.80.88.0-35.80.95.255", "server1.acme.com", "*.acme.com"]
    dest_countries=["BR", "CA", "GB"]
    departments {
        id = [8061246]
    }
    dest_ip_groups {
        id = [-4]
    }
    dest_ipv6_groups {
        id = [-5]
    }
    src_ip_groups {
        id = [18448894]
    }
    src_ipv6_groups {
        id = [-3]
    }
    time_windows {
        id = [485]
    }
    nw_services {
        id = [462370, 17472664]
    }
    locations {
        id = [256000852, -3]
    }
    location_groups {
        id = [8061257, 8061256]
    }
    labels {
        id = [1416803]
    }
}
