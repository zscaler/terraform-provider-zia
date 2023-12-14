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