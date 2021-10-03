terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}
provider "zia" {}


resource "zia_traffic_forwarding_gre_tunnel" "example" {
  source_ip = "96.53.93.171" 
  comment   = "comment test"
  primary_dest_vip {
    id = 64199
  }
  secondary_dest_vip {
    id = 95619
  }
  ip_unnumbered = true
}
data "zia_traffic_forwarding_gre_tunnel" "example" {
  id = zia_traffic_forwarding_gre_tunnel.example.tunnel_id
}

output "zia_traffic_forwarding_gre_tunnel" {
  value = data.zia_traffic_forwarding_gre_tunnel.example
}
