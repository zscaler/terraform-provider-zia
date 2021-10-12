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
  country_code = "US"
  ip_unnumbered = true
}

output "zia_traffic_forwarding_gre_tunnel" {
  value = zia_traffic_forwarding_gre_tunnel.example
}

resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address =  "50.98.112.170"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = 49.0
    longitude = -123.0
}


data "zia_gre_virtual_ip_address_list" "yvr1_0"{
    source_ip = zia_traffic_forwarding_static_ip.example.ip_address
    required_count = 2
}


data "zia_gre_internal_ip_range_list" "example"{
}

output "zia_gre_internal_ip_range_example"{
    value = data.zia_gre_internal_ip_range_list.example
}

output "zia_gre_internal_ip_range_example_first"{
    value = data.zia_gre_internal_ip_range_list.example.list[0].start_ip_address
}