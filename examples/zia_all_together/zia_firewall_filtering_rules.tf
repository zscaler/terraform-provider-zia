resource "zia_url_filtering_rules" "zscaler_proxy_traffic" {
    name = "Zscaler Proxy Traffic"
    description = "Zscaler Proxy Traffic"
    action = "ALLOW"
    state = "ENABLED"
    nw_services {
        id = [ zia_firewall_filtering_network_service.example.id ]
    }

}