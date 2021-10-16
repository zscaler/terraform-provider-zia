/*
// Retrieving Zscaler VIPs for GRE Tunnel creation
data "zia_gre_virtual_ip_address_list" "yvr1_0"{
    source_ip = zia_traffic_forwarding_static_ip.vancouver_shaw_business_internet.ip_address
    required_count = 2
}

data "zia_gre_virtual_ip_address_list" "yvr1_1"{
    source_ip = zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_01.ip_address
    required_count = 2
}

data "zia_gre_virtual_ip_address_list" "yvr1_2"{
    source_ip = zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_02.ip_address
    required_count = 2
}

data "zia_gre_virtual_ip_address_list" "akl1_0"{
    source_ip = zia_traffic_forwarding_static_ip.nz_auckland_branch.ip_address
    required_count = 2
}

data "zia_gre_virtual_ip_address_list" "syd1_0"{
    source_ip = zia_traffic_forwarding_static_ip.au_sydney_branch.ip_address
    required_count = 2
}

data "zia_gre_internal_ip_range_list" "gre_list"{
}
*/
