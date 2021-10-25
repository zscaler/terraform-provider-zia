resource "zia_firewall_filtering_ip_source_groups" "example1"{
    name = "example1"
    description = "example1"
    ip_addresses = ["192.168.1.1", "192.168.1.2", "192.168.1.3"]
}