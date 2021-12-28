resource "zia_firewall_filtering_destination_groups" "example" {
  name        = "Example"
  description = "Example"
  type        = "DSTN_FQDN"
  addresses = [ "test1.acme.com", "test2.acme.com", "test3.acme.com" ]
}