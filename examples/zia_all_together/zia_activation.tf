/*
resource "zia_activation_status" "activation"{
    status = "ACTIVE"
     depends_on = [ zia_location_management.ca_vancouver_ipsec,
                    zia_traffic_forwarding_static_ip.ca_vancouver_shaw_business_internet,
                    zia_traffic_forwarding_vpn_credentials.vancouver_ipsec_tunnel,
                    zia_location_management.vancouver_telus_home_internet_gre01,
                    zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre01,
                    zia_traffic_forwarding_gre_tunnel.telus_home_internet_01_gre01,
                    zia_location_management.vancouver_telus_home_internet_gre02,
                    zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre02,
                    zia_traffic_forwarding_gre_tunnel.vancouver_unnumbered_gre02, 
                    zia_location_management.nz_auckland_branch_gre01,
                    zia_traffic_forwarding_static_ip.nz_auckland_branch_gre01,
                    zia_traffic_forwarding_gre_tunnel.nz_auckland_branch_gre01,
                    zia_location_management.au_sydney_branch_gre01,
                    zia_traffic_forwarding_static_ip.au_sydney_branch_gre01,
                   zia_traffic_forwarding_gre_tunnel.au_sydney_branch_gre01
                ]
}

output "zia_activation_status"{
    value = zia_activation_status.activation
}
*/