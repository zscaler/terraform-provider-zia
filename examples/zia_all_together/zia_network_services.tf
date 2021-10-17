resource "zia_firewall_filtering_network_service" "example"{
    name = "example"
    description = "example"
    src_tcp_ports {
         start = 80
    }
    dest_tcp_ports {
        start = 123
    }
    src_udp_ports {
        start = 123
    }
    dest_udp_ports {
        start = 123
    }
    type = "CUSTOM"
}

