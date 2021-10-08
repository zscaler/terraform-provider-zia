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
  source_ip = zia_traffic_forwarding_static_ip.example.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  within_country = true
  primary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_0.list[0].id
  }
  secondary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_0.list[1].id
  }
  ip_unnumbered = true
}

output "zia_traffic_forwarding_gre_tunnel" {
  value = zia_traffic_forwarding_gre_tunnel.example
}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.170"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}


data "zia_gre_virtual_ip_address_list" "yvr1_0"{
<<<<<<< HEAD
    source_ip = "50.98.112.170"
}

data "zia_gre_virtual_ip_address_list" "yvr1_1"{
    source_ip = "50.98.112.170"
}

resource "zia_activation_status" "example"{
=======
    source_ip = zia_traffic_forwarding_static_ip.example.ip_address
    required_count = 2
}

/*
resource "zia_activation_status" "example1"{
>>>>>>> master
    status = "ACTIVE"
    depends_on = [ zia_traffic_forwarding_gre_tunnel.example, zia_traffic_forwarding_static_ip.example ]
}

<<<<<<< HEAD
output "zia_activation_status_example"{
    value = zia_activation_status.example
=======
output "zia_activation_status_example1"{
    value = zia_activation_status.example1
}
*/


data "zia_gre_internal_ip_range_list" "example"{
}

output "zia_gre_internal_ip_range_example"{
    value = data.zia_gre_internal_ip_range_list.example
}

output "zia_gre_internal_ip_range_example_first"{
    value = data.zia_gre_internal_ip_range_list.example.list[0].start_ip_address
>>>>>>> master
}