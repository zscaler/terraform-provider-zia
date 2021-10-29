resource "zia_location_management" "usa_sjc37_office_branch01"{
    name = "USA_SJC37_Office-Branch01"
    description = "Created with Terraform"
    country = "UNITED_STATES"
    tz = "UNITED_STATES_AMERICA_LOS_ANGELES"
    profile = "CORPORATE"
    parent_id = zia_location_management.usa_sjc37.id
    depends_on = [ zia_traffic_forwarding_static_ip.usa_sjc37, zia_traffic_forwarding_vpn_credentials.usa_sjc37, zia_location_management.usa_sjc37 ]
    auth_required = true
    idle_time_in_minutes = 720
    display_time_unit = "HOUR"
    surrogate_ip = true
    ofw_enabled = true
    ip_addresses = [ "10.5.0.0-10.5.255.255" ]
    up_bandwidth = 10000
    dn_bandwidth = 10000
}