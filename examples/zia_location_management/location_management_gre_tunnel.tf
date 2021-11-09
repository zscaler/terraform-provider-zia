resource "zia_location_management" "usa_sjc37"{
    name = "USA_SJC_37"
    description = "Created with Terraform"
    country = "UNITED_STATES"
    tz = "UNITED_STATES_AMERICA_LOS_ANGELES"
    auth_required = true
    idle_time_in_minutes = 720
    display_time_unit = "HOUR"
    surrogate_ip = true
    xff_forward_enabled = true
    ofw_enabled = true
    ips_control = true
    ip_addresses = [ zia_traffic_forwarding_static_ip.usa_sjc37.ip_address ]
    depends_on = [ zia_traffic_forwarding_static_ip.usa_sjc37]
}

resource "zia_traffic_forwarding_static_ip" "usa_sjc37"{
    ip_address =  "1.1.1.1"
    routable_ip = true
    comment = "SJC37 - GRE Tunnel"
    geo_override = false
}