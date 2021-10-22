terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}

resource "zia_firewall_filtering_network_service" "example" {
  name        = "example"
  description = "example"
  src_tcp_ports {
    start = 5000
  }
  src_tcp_ports {
    start = 5001
  }
  src_tcp_ports {
    start = 5002
  }
  src_tcp_ports {
    start = 5003
    end = 5005
  }
  dest_tcp_ports {
    start = 5000
  }
    dest_tcp_ports {
    start = 5001
  }
  dest_tcp_ports {
    start = 5002
  }
    dest_tcp_ports {
    start = 5003
    end = 5005
  }
  // src_udp_ports {
    // start = 123
    // end   = 125
  // }
  // src_udp_ports {
    // start = 126
    // end   = 127
  // }

  // dest_udp_ports {
    // start = 123
    // end   = 125
  // }
  // dest_udp_ports {
    // start = 126
    // end   = 127
  // }
  type = "CUSTOM"
}


/*
data "zia_firewall_filtering_network_service" "example" {
  name = zia_firewall_filtering_network_service.example.name
}

output "zia_firewall_filtering_network_service" {
  value = data.zia_firewall_filtering_network_service.example
}
*/