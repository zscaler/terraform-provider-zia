# Use the ZPA Terraform Provider To Create or Retrieve the Server Group and Application Segment

data "zpa_server_group" "this" {
  name = "Server_Group_IP_Source_Anchoring"
}

data "zpa_application_segment" "this1" {
  name = "App_Segment_IP_Source_Anchoring"
}

data "zpa_application_segment" "this2" {
  name = "App_Segment_IP_Source_Anchoring2"
}

resource "zia_forwarding_control_zpa_gateway" "this" {
    name = "ZPA_GW01"
    description = "ZPA_GW01"
    type = "ZPA"
    zpa_server_group {
      external_id = data.zpa_server_group.this.id
      name = data.zpa_server_group.this.id
    }
    zpa_app_segments {
        external_id = data.zpa_application_segment.this1.id
        name = data.zpa_application_segment.this1.name
    }
    zpa_app_segments {
        external_id = data.zpa_application_segment.this2.id
        name = data.zpa_application_segment.this2.name
    }
}

