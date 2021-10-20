// Creating location management resource
resource "zia_location_management" "ca_vancouver_ipsec" {
  name                 = "CA-Vancouver-IPSec"
  description          = "Created with Terraform"
  country              = "CANADA"
  tz                   = "CANADA_AMERICA_VANCOUVER"
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  xff_forward_enabled  = true
  ofw_enabled          = true
  ips_control          = true
  ip_addresses         = [zia_traffic_forwarding_static_ip.ca_vancouver_shaw_business_internet.ip_address]
  depends_on           = [zia_traffic_forwarding_static_ip.ca_vancouver_shaw_business_internet, zia_traffic_forwarding_vpn_credentials.vancouver_ipsec_tunnel]
  vpn_credentials {
    id   = zia_traffic_forwarding_vpn_credentials.vancouver_ipsec_tunnel.vpn_credental_id
    type = zia_traffic_forwarding_vpn_credentials.vancouver_ipsec_tunnel.type
  }
}

resource "zia_location_management" "vancouver_telus_home_internet_gre01" {
  name                 = "CA - Vancouver - Branch01"
  description          = "Created with Terraform"
  country              = "CANADA"
  tz                   = "CANADA_AMERICA_VANCOUVER"
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  xff_forward_enabled  = true
  ofw_enabled          = true
  ips_control          = true
  ip_addresses         = [zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre01.ip_address]
  depends_on           = [zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre01, zia_traffic_forwarding_gre_tunnel.telus_home_internet_01_gre01]
}

resource "zia_location_management" "vancouver_telus_home_internet_gre02" {
  name                 = "CA - Vancouver - Branch02"
  description          = "Created with Terraform"
  country              = "CANADA"
  tz                   = "CANADA_AMERICA_VANCOUVER"
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  xff_forward_enabled  = true
  ofw_enabled          = true
  ips_control          = true
  ip_addresses         = [zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre02.ip_address]
  depends_on           = [zia_traffic_forwarding_static_ip.vancouver_telus_home_internet_gre02, zia_traffic_forwarding_gre_tunnel.vancouver_unnumbered_gre02]
}

resource "zia_location_management" "nz_auckland_branch_gre01" {
  name                 = "NZ - Auckland - Branch01"
  description          = "Created with Terraform"
  country              = "NEW_ZEALAND"
  tz                   = "NEW_ZEALAND_PACIFIC_AUCKLAND"
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  xff_forward_enabled  = true
  ofw_enabled          = true
  ips_control          = true
  ip_addresses         = [zia_traffic_forwarding_static_ip.nz_auckland_branch_gre01.ip_address]
  depends_on           = [zia_traffic_forwarding_static_ip.nz_auckland_branch_gre01, zia_traffic_forwarding_gre_tunnel.nz_auckland_branch_gre01]
}

resource "zia_location_management" "au_sydney_branch_gre01" {
  name                 = "AU - Sydney - Branch01"
  description          = "Created with Terraform"
  country              = "AUSTRALIA"
  tz                   = "AUSTRALIA_SYDNEY"
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  xff_forward_enabled  = true
  ofw_enabled          = true
  ips_control          = true
  ip_addresses         = [zia_traffic_forwarding_static_ip.au_sydney_branch_gre01.ip_address]
  depends_on           = [zia_traffic_forwarding_static_ip.au_sydney_branch_gre01, zia_traffic_forwarding_gre_tunnel.au_sydney_branch_gre01]
}

resource "zia_location_management" "br_sao_paulo_branch_gre01" {
  name                 = "BR - Sao Paulo - Branch01"
  description          = "Created with Terraform"
  country              = "BRAZIL"
  tz                   = "BRAZIL_AMERICA_SAO_PAULO"
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  xff_forward_enabled  = true
  ofw_enabled          = true
  ips_control          = true
  ip_addresses         = [zia_traffic_forwarding_static_ip.br_sao_paulo_branch_gre01.ip_address]
  depends_on           = [zia_traffic_forwarding_static_ip.br_sao_paulo_branch_gre01, zia_traffic_forwarding_gre_tunnel.br_sao_paulo_branch_gre01]
}

resource "zia_location_management" "br_sao_paulo_branch01_guest_wifi" {
  name         = "Guest Wi-Fi - Branch01"
  description  = "Created with Terraform"
  country      = "BRAZIL"
  tz           = "BRAZIL_AMERICA_SAO_PAULO"
  profile      = "GUESTWIFI"
  parent_id    = zia_location_management.br_sao_paulo_branch_gre01.id
  ofw_enabled  = true
  ip_addresses = ["10.131.2.128-10.131.3.255"]
  up_bandwidth = 2000
  dn_bandwidth = 2000
}

resource "zia_location_management" "br_sao_paulo_branch02_manufacturing_plant" {
  name         = "Manufacturing Plant - Branch02"
  description  = "Created with Terraform"
  country      = "BRAZIL"
  tz           = "BRAZIL_AMERICA_SAO_PAULO"
  parent_id    = zia_location_management.br_sao_paulo_branch_gre01.id
  ofw_enabled  = true
  ip_addresses = ["10.1.0.0-10.1.255.255"]
  up_bandwidth = 10000
  dn_bandwidth = 10000
}

resource "zia_location_management" "br_sao_paulo_branch03_office" {
  name                 = "Office - Branch03"
  description          = "Created with Terraform"
  country              = "BRAZIL"
  tz                   = "BRAZIL_AMERICA_SAO_PAULO"
  profile              = "CORPORATE"
  parent_id            = zia_location_management.br_sao_paulo_branch_gre01.id
  auth_required        = true
  idle_time_in_minutes = 720
  display_time_unit    = "HOUR"
  surrogate_ip         = true
  ofw_enabled          = true
  ip_addresses         = ["10.5.0.0-10.5.255.255"]
  up_bandwidth         = 10000
  dn_bandwidth         = 10000
}