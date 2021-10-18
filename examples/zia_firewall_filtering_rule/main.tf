terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}
resource "zia_firewall_filtering_rule" "zscaler_proxy_traffic" {
    name = "Zscaler Proxy Traffic"
    description = "Zscaler Proxy Traffic"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    rank = 1
    enable_full_logging = true
    nw_services {
        id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
    }
}


data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
    name = "ZSCALER_PROXY_NW_SERVICES"
}

output "zia_firewall_filtering_rule" {
  value = zia_firewall_filtering_rule.zscaler_proxy_traffic
}