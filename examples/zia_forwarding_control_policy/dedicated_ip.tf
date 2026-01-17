data "zia_forwarding_control_dedicated_ip_gateway" "this" {
  name = "Default"
}

resource "zia_forwarding_control_rule" "dedicated_ip" {
  name           = "DEDIP_FORWARDING_RULE"
  description    = "DEDIP_FORWARDING_RULE"
  type           = "FORWARDING"
  state          = "ENABLED"
  forward_method = "ENATDEDIP"
  order          = 1
  rank           = 7
  dest_addresses = ["example.com"]
  dedicated_ip_gateway {
    id   = data.zia_forwarding_control_dedicated_ip_gateway.this.id
    name = data.zia_forwarding_control_dedicated_ip_gateway.this.name
  }
}
