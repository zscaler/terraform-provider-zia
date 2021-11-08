resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    nw_services {
        id = [ data.zia_firewall_filtering_network_service.zscaler_proxy_nw_services.id ]
    }
}

data "zia_firewall_filtering_network_service" "zscaler_proxy_nw_services" {
    name = "ZSCALER_PROXY_NW_SERVICES"
}