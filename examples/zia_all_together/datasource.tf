data "zia_traffic_forwarding_gre_tunnel_info" "vancouver_telus_home_internet_gre01" {
    ip_address =  "50.98.112.169"
    gre_enabled = true
}

output "zia_traffic_forwarding_gre_tunnel_info_vancouver_telus_home_internet_gre01" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.vancouver_telus_home_internet_gre01
}

data "zia_traffic_forwarding_gre_tunnel_info" "vancouver_telus_home_internet_gre02" {
    ip_address =  "50.98.112.170"
    gre_enabled = true
}

output "zia_traffic_forwarding_gre_tunnel_info_vancouver_telus_home_internet_gre02" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.vancouver_telus_home_internet_gre02
}

data "zia_traffic_forwarding_gre_tunnel_info" "nz_auckland_branch_gre01" {
    ip_address =  "101.110.112.100"
    gre_enabled = true
}

output "zia_traffic_forwarding_gre_tunnel_info_nz_auckland_branch_gre01" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.nz_auckland_branch_gre01
}

data "zia_traffic_forwarding_gre_tunnel_info" "au_sydney_branch_gre01" {
    ip_address =  "61.68.118.237"
    gre_enabled = true
}

output "zia_traffic_forwarding_gre_tunnel_info_au_sydney_branch_gre01" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.au_sydney_branch_gre01
}

data "zia_traffic_forwarding_gre_tunnel_info" "br_sao_paulo_branch_gre01" {
    ip_address =  "187.22.113.134"
    gre_enabled = true
}

output "zia_traffic_forwarding_gre_tunnel_info_br_sao_paulo_branch_gre01" {
  value = data.zia_traffic_forwarding_gre_tunnel_info.br_sao_paulo_branch_gre01
}