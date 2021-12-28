resource "zia_firewall_filtering_destination_groups" "example"{
    name = "example"
    description = "example"
    type = "DSTN_IP"
    addresses = ["1.2.3.4", "1.2.3.5", "1.2.3.6"]
    countries = ["COUNTRY_CA"]
}


resource "zia_firewall_filtering_destination_groups" "example_ip_ranges" {
  name        = "Example - IP Ranges"
  description = "Example - IP Ranges"
  type        = "DSTN_IP"
  addresses = ["3.217.228.0-3.217.231.255",
    "3.235.112.0-3.235.119.255",
    "52.23.61.0-52.23.62.25",
    "35.80.88.0-35.80.95.255",