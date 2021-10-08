terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}
provider "zia" {}


resource "zia_traffic_forwarding_gre_tunnel" "example1" {
  source_ip = zia_traffic_forwarding_static_ip.example_1.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  depends_on = [ zia_traffic_forwarding_static_ip.example_1 ]
  within_country = true
  country_code = "CA"
  primary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_0.list[0].id
  }
  secondary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_1.list[1].id
  }
  ip_unnumbered = false
}

output "zia_traffic_forwarding_gre_tunnel_example1" {
  value = zia_traffic_forwarding_gre_tunnel.example1
}

resource "zia_traffic_forwarding_gre_tunnel" "example2" {
  source_ip = zia_traffic_forwarding_static_ip.example_2.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  depends_on = [ zia_traffic_forwarding_static_ip.example_2 ]
  primary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_0.list[0].id
  }
  secondary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_1.list[1].id
  }
  ip_unnumbered = false
}

output "zia_traffic_forwarding_gre_tunnel_example2" {
  value = zia_traffic_forwarding_gre_tunnel.example2
}

// Create static ip addresses
resource "zia_traffic_forwarding_static_ip" "example_1"{
    ip_address =  "50.98.112.170"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}

resource "zia_traffic_forwarding_static_ip" "example_2"{
    ip_address =  "50.98.112.171"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}

data "zia_gre_virtual_ip_address_list" "yvr1_0"{
    source_ip = "50.98.112.170"
}

data "zia_gre_virtual_ip_address_list" "yvr1_1"{
    source_ip = "50.98.112.171"
}

resource "zia_activation_status" "example1"{
    status = "ACTIVE"
    depends_on = [ zia_traffic_forwarding_gre_tunnel.example1,
                   zia_traffic_forwarding_gre_tunnel.example2,
                   zia_traffic_forwarding_static_ip.example_1,
                   zia_traffic_forwarding_static_ip.example_2
                  ]
}

output "zia_activation_status_example1"{
    value = zia_activation_status.example1
}