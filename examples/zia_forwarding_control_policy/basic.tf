# Note: In order to retrieve the Application Segment ID you must use the ZPA Terraform Provider
data "zpa_application_segment" "app01" {
    name = "App01"
}

data "zpa_application_segment" "app02" {
    name = "App02"
}

resource "zia_forwarding_control_zpa_gateway" "this" {
    name = "ZPA_GW01"
}

resource "zia_forwarding_control_rule" "example" {
    name = "Example"
    description = "Example"
    type = "FORWARDING"
    state = "ENABLED"
    forward_method = "ZPA"
    order = 1
    zpa_gateway {
        id   = data.zia_forwarding_control_zpa_gateway.id
        name = data.zia_forwarding_control_zpa_gateway.name
    }
    zpa_app_segments {
        name = data.zpa_application_segment.app01.name
        external_id = data.zpa_application_segment.app01.id
    }
    zpa_app_segments {
        name = data.zpa_application_segment.app02.name
        external_id = data.zpa_application_segment.app02.id
    }
}


data "zia_dedicated_ip_proxy" "this" {
  name = "GW01"
}

resource "zia_forwarding_control_rule" "this" {
  name           = "FC_ENATDEDIP_RULE"
  description    = "FC_ENATDEDIP_RULE"
  order          = 1
  rank           = 7
  state          = "ENABLED"
  type           = "FORWARDING"
  forward_method = "ENATDEDIP"
  src_ips            = ["192.168.200.200"]
  dest_addresses     = ["192.168.255.1"]
  dest_ip_categories = ["ZSPROXY_IPS", "CUSTOM_01"]
  dest_countries     = ["CA", "US"]
  dedicated_ip_gateway {
    id   = data.zia_dedicated_ip_proxy.this.id
    name = data.zia_dedicated_ip_proxy.this.name
  }
}
