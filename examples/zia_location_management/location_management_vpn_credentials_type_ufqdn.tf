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
    depends_on = [ zia_traffic_forwarding_vpn_credentials.usa_sjc37 ]
    vpn_credentials {
       id = zia_traffic_forwarding_vpn_credentials.usa_sjc37.id
       type = zia_traffic_forwarding_vpn_credentials.usa_sjc37.type
    }
}

######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########
resource "zia_traffic_forwarding_vpn_credentials" "usa_sjc37"{
    type = "UFQDN"
    fqdn = "usa_sjc37@acme.com"
    comments = "USA - San Jose IPSec Tunnel"
    pre_shared_key = "*************"
}
