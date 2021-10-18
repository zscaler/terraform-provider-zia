resource "zia_traffic_forwarding_gre_tunnel" "nz_auckland_branch_gre01" {
  source_ip = zia_traffic_forwarding_static_ip.nz_auckland_branch_gre01.ip_address
  comment   = "GRE Tunnel Created with Terraform"
  within_country = true
  country_code = "NZ"
  ip_unnumbered = true
  depends_on = [ zia_traffic_forwarding_static_ip.nz_auckland_branch_gre01 ]
}