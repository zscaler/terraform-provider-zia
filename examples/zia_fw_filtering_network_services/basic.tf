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
    end = 5005
  }
  dest_tcp_ports {
    start = 5000
  }
    dest_tcp_ports {
    start = 5001
  }
  dest_tcp_ports {
    start = 5003
    end = 5005
  }
  type = "CUSTOM"
}