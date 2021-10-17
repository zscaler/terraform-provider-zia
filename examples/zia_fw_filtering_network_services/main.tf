terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

resource "zia_firewall_filtering_network_service" "example"{
    name = "example"
    description = "example"
    src_tcp_ports {
         start = 123
         end = 125
    }
    dest_tcp_ports {
         start = 123
         end = 125
    }
    src_udp_ports {
         start = 123
         end = 125
         start = 126
         end = 127
    }
    dest_udp_ports {
         start = 123
         end = 125
         start = 126
         end = 127
    }
    type = "CUSTOM"
}

