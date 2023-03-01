data "zia_traffic_forwarding_gre_vip_recommended_list" "this"{
    source_ip = zia_traffic_forwarding_static_ip.this.ip_address
    required_count = 2
}

data "zia_gre_internal_ip_range_list" "this"{
    required_count = 10
}

resource "zia_traffic_forwarding_static_ip" "this"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "Created with Terraform"
    geo_override = true
    latitude = 49.0526
    longitude = -122.8291
}

resource "zia_traffic_forwarding_gre_tunnel" "this" {
  source_ip      = zia_traffic_forwarding_static_ip.this.ip_address
  comment        = "GRE Tunnel Created with Terraform"
  internal_ip_range = data.zia_gre_internal_ip_range_list.this.list[0].start_ip_address
  within_country = false
  country_code   = "CA"
  ip_unnumbered  = false
  depends_on     = [zia_traffic_forwarding_static_ip.this]
  primary_dest_vip {
    datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].datacenter
    id = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].id
    virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[0].virtual_ip
  }
  secondary_dest_vip {
    datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].datacenter
    id = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].id
    virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.this.list[1].virtual_ip
  }
}

# Note: The attributes ``primary_dest_vip`` and ``secondary_dest_vip`` are considered optional
# The provider automatically selects the VIPs if none is indicated, using the longitude and latitude coordinates of the ``source_ip``.
# The ``source_ip`` attribute can obtain its IP address using the ``zia_traffic_forwarding_static_ip`` resource or datasource.