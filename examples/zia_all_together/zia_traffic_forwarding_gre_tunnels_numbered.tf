// Create GRE Tunnel
resource "zia_traffic_forwarding_gre_tunnel" "telus_home_internet_01_gre01" {
  source_ip      = zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre01.ip_address
  comment        = "GRE Tunnel Created with Terraform"
  within_country = true
  country_code   = "CA"
  ip_unnumbered  = false
  depends_on     = [zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre01]
}

resource "zia_traffic_forwarding_gre_tunnel" "au_sydney_branch_gre01" {
  source_ip      = zia_traffic_forwarding_static_ip.au_sydney_branch_gre01.ip_address
  comment        = "GRE Tunnel Created with Terraform"
  within_country = true
  country_code   = "AU"
  ip_unnumbered  = false
  depends_on     = [zia_traffic_forwarding_static_ip.au_sydney_branch_gre01]
}

resource "zia_traffic_forwarding_gre_tunnel" "br_sao_paulo_branch_gre01" {
  source_ip      = zia_traffic_forwarding_static_ip.br_sao_paulo_branch_gre01.ip_address
  comment        = "GRE Tunnel Created with Terraform"
  within_country = true
  country_code   = "BR"
  ip_unnumbered  = false
  depends_on     = [zia_traffic_forwarding_static_ip.br_sao_paulo_branch_gre01]
}