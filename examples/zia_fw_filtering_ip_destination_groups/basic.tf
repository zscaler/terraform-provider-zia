resource "zia_firewall_filtering_destination_groups" "example"{
    name = "example"
    description = "example"
    type = "DSTN_IP"
    addresses = ["1.2.3.4", "1.2.3.5", "1.2.3.6"]
    countries = ["COUNTRY_CA"]
}