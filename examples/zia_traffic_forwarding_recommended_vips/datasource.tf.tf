data "zia_traffic_forwarding_gre_vip_recommended_list" "yvr1_0"{
    source_ip = "1.1.1.1"
    required_count = 2
}

output "zia_traffic_forwarding_gre_vip_recommended_list"{
    value = data.zia_gre_virtual_ip_address_list.example
}