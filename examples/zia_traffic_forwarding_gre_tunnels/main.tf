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
  source_ip = zia_traffic_forwarding_static_ip.example.id
  comment   = "comment test"
  primary_dest_vip {
    // id = data.zia_gre_virtual_ip_address_list.yvr1_0.list[0].id
    //virtual_ip = data.zia_gre_virtual_ip_address_list.yvr1_0.list[0].virtual_ip
    //virtual_ip = "165.225.210.32"
    id = 64199
  }
  secondary_dest_vip {
    //id = data.zia_gre_virtual_ip_address_list.yvr1_1.list[0].id
    //virtual_ip = data.zia_gre_virtual_ip_address_list.yvr1_1.list[0].virtual_ip
  // virtual_ip = "165.225.210.33"
  id = 64197
  }
  ip_unnumbered = true
  within_country = true
}

output "zia_traffic_forwarding_gre_tunnel" {
  value = zia_traffic_forwarding_gre_tunnel.example
}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.169"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}

/*
data "zia_gre_virtual_ip_address_list" "yvr1_0"{
  source_ip = "50.98.112.169"
}

data "zia_gre_virtual_ip_address_list" "yvr1_1"{
  source_ip = "50.98.112.169"
}
*/