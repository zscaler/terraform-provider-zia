---
subcategory: "Forwarding Control Policy"
layout: "zscaler"
page_title: "ZIA: forwarding_control_zpa_gateway"
description: |-
    Creates and manages ZIA forwarding control zpa gateway used in IP Source Anchoring.

---
# Resource: forwarding_control_zpa_gateway

Use the **forwarding_control_zpa_gateway** resource allows the creation and management of ZIA forwarding control ZPA Gateway used in IP Source Anchoring integration between Zscaler Internet Access and Zscaler Private Access. This resource can then be associated with a ZIA Forwarding Control Rule.

⚠️ **IMPORTANT:**: To configure a ZPA Gateway you **MUST** use the ZPA Terraform Provider to configure a Server Group and Application Segment with the Source IP Anchoring feature enabled at the Application Segment resource. Please refer to the ZPA Terraform Provider documentation [here](https://registry.terraform.io/providers/zscaler/zpa/latest/docs)

## Example Usage

```hcl
# ZIA Forwarding Control - ZPA Gateway
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
```

## Argument Reference

The following arguments are supported:

* `name` (Required) The name of the forwarding control ZPA Gateway to be exported.
* `zpa_server_group` (Required) - The ZPA Server Group that is configured for Source IP Anchoring
  * `external_id` (Required) - An external identifier used for an entity that is managed outside of ZIA. Examples include zpaServerGroup and zpaAppSegments. This field is not applicable to ZIA-managed entities.
  * `name` (Required) - The configured name of the entity
* `zpa_app_segments` - (Required) The ZPA Server Group that is configured for Source IP Anchoring
  * `external_id` (Required) - An external identifier used for an entity that is managed outside of ZIA. Examples include zpaServerGroup and zpaAppSegments. This field is not applicable to ZIA-managed entities.
  * `name` (Required) - The configured name of the entity

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - Additional details about the ZPA gateway
* `type` - (string) - Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA). Supported values: ``ZPA`` and ``ECZPA``

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**forwarding_control_zpa_gateway** can be imported by using `<GATEWAY_ID>` or `<GATEWAY_NAME>` as the import ID.

For example:

```shell
terraform import forwarding_control_zpa_gateway.example <gateway_id>
```

or

```shell
terraform import forwarding_control_zpa_gateway.example <gateway_name>
```
