# IP Destination Group of Type DSTN_FQDN
resource "zia_firewall_filtering_destination_groups" "dstn_fqdn" {
  name        = "Example Destination FQDN"
  description = "Example Destination FQDN"
  type        = "DSTN_FQDN"
  addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
}

# IP Destination Group of Type DSTN_IP
resource "zia_firewall_filtering_destination_groups" "example_ip_ranges" {
  name        = "Example - IP Ranges"
  description = "Example - IP Ranges"
  type        = "DSTN_IP"
  addresses = ["3.217.228.0-3.217.231.255",
    "3.235.112.0-3.235.119.255",
    "52.23.61.0-52.23.62.25",
    "35.80.88.0-35.80.95.255"]
}

# IP Destination Group of Type DSTN_DOMAIN
resource "zia_firewall_filtering_destination_groups" "example_dstn_domain" {
  name          = "Example Destination Domain"
  description   = "Example Destination Domain"
  type          = "DSTN_DOMAIN"
  addresses     = ["acme.com", "acme1.com"]
}

# IP Destination Group of Type DSTN_OTHER
resource "zia_firewall_filtering_destination_groups" "example_dstn_other" {
  name          = "Example Destination Other"
  description   = "Example Destination Other"
  type          = "DSTN_OTHER"
  countries     = ["COUNTRY_CA"]
}