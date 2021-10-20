// VPN Credential resource for IPSec Tunnels
resource "zia_traffic_forwarding_vpn_credentials" "vancouver_ipsec_tunnel" {
  type           = "UFQDN"
  fqdn           = "vpn@securitygeek.io"
  comments       = "Vancouver IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}


resource "zia_traffic_forwarding_vpn_credentials" "br_sao_paulo" {
  type           = "UFQDN"
  fqdn           = "saopaulo@securitygeek.io"
  comments       = "Sao Paulo - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "nl_amsterdam" {
  type           = "UFQDN"
  fqdn           = "amsterdam@securitygeek.io"
  comments       = "Amsterdam - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "ge_berlin" {
  type           = "UFQDN"
  fqdn           = "berlin@securitygeek.io"
  comments       = "Berlin - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "hu_budapest" {
  type           = "UFQDN"
  fqdn           = "budapest@securitygeek.io"
  comments       = "Hungary - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "us_chicago" {
  type           = "UFQDN"
  fqdn           = "chicago@securitygeek.io"
  comments       = "Chicago - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "us_newyork" {
  type           = "UFQDN"
  fqdn           = "newyork@securitygeek.io"
  comments       = "New York - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "fr_paris" {
  type           = "UFQDN"
  fqdn           = "paris@securitygeek.io"
  comments       = "Paris - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "us_san_francisco" {
  type           = "UFQDN"
  fqdn           = "sanfrancisco@securitygeek.io"
  comments       = "San Francisco - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "au_sydney" {
  type           = "UFQDN"
  fqdn           = "sydney@securitygeek.io"
  comments       = "Sydney - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}

resource "zia_traffic_forwarding_vpn_credentials" "jp_tokyo" {
  type           = "UFQDN"
  fqdn           = "tokyo@securitygeek.io"
  comments       = "Tokyo - IPSec Tunnel"
  pre_shared_key = "TBt*kJ8623!hL97uPsvdcDq_7AEf!7k@Kfc"
}