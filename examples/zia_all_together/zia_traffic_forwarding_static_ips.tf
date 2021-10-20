// Create Static IP address
resource "zia_traffic_forwarding_static_ip" "ca_vancouver_shaw_business_internet" {
  ip_address   = "96.53.93.170"
  routable_ip  = true
  comment      = "CA - Vancouver - Shaw Business Internet"
  geo_override = false
}

resource "zia_traffic_forwarding_static_ip" "vancouver_telus_home_internet_gre01" {
  ip_address   = "50.98.112.169"
  routable_ip  = true
  comment      = "CA - Vancouver - Branch01"
  geo_override = false
}

resource "zia_traffic_forwarding_static_ip" "vancouver_telus_home_internet_gre02" {
  ip_address   = "50.98.112.170"
  routable_ip  = true
  comment      = "CA - Vancouver - Branch02"
  geo_override = false
}

resource "zia_traffic_forwarding_static_ip" "nz_auckland_branch_gre01" {
  ip_address   = "101.110.112.100"
  routable_ip  = true
  comment      = "Auckland - Branch01"
  geo_override = false
}

resource "zia_traffic_forwarding_static_ip" "au_sydney_branch_gre01" {
  ip_address   = "61.68.118.237"
  routable_ip  = true
  comment      = "Sydney - Branch02"
  geo_override = false
}

resource "zia_traffic_forwarding_static_ip" "br_sao_paulo_branch_gre01" {
  ip_address   = "187.22.113.134"
  routable_ip  = true
  comment      = "Sao Paulo - Branch01"
  geo_override = false
}



