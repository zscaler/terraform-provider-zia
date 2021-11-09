data "zia_traffic_forwarding_public_node_vips" "yvr1"{
    datacenter = "YVR1"
}

output "zia_traffic_forwarding_public_node_vips_yvr1"{
    value = data.zia_traffic_forwarding_public_node_vips.yvr1
}

data "zia_traffic_forwarding_public_node_vips" "sea1"{
    datacenter = "SEA1"
}

output "zia_traffic_forwarding_public_node_vips_sea1"{
    value = data.zia_traffic_forwarding_public_node_vips.sea1
}