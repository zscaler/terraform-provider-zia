# Note: In order to retrieve the Application Segment ID you must use the ZPA Terraform Provider
data "zpa_application_segment" "this" {
    name = "ZPA_App_Segment_IP_Source_Anchoring"
}

# Note: In order to retrieve the Server Group ID you must use the ZPA Terraform Provider
data "zpa_server_group" "this" {
    name = "ZPA_Server_Group_IP_Source_Anchoring"
}

resource "zia_forwarding_control_zpa_gateway" "this" {
  name = "ZPA_GW01"
  description = "ZPA_GW01"
  type = "ZPA"
  zpa_server_group {
    external_id = data.zpa_server_group.this.id
    name = data.zpa_server_group.this.name
  }
    zpa_app_segments {
        external_id        = data.zpa_application_segment.this.id
        name                = data.zpa_application_segment.this.name
    }
}