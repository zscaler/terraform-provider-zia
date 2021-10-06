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
  primary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_0.list[0].id
  }
  secondary_dest_vip {
    id = data.zia_gre_virtual_ip_address_list.yvr1_1.list[1].id
  }
  ip_unnumbered = true
}

output "zia_traffic_forwarding_gre_tunnel" {
  value = zia_traffic_forwarding_gre_tunnel.example
}

// Static IP needs to be configured before creating the GRE tunnel.
//"code": "RESOURCE_NOT_FOUND",
// "message": "Static IP 50.98.112.169 has not been configured."
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.169"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = false
}


data "zia_gre_virtual_ip_address_list" "yvr1_0"{
    source_ip = "50.98.112.169"
}

data "zia_gre_virtual_ip_address_list" "yvr1_1"{
    source_ip = "50.98.112.169"
}

resource "zia_activation_status" "example1"{
    status = "ACTIVE"
}

output "zia_activation_status_example1"{
    value = zia_activation_status.example1
}