---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZIA: activation_status"
description: |-
  Triggers activation of ZIA pending configuration changes as part of the Terraform graph.
---

# zia_activation_status (Resource)

* [Official documentation](https://help.zscaler.com/zia/saving-and-activating-changes-zia-admin-portal)
* [API documentation](https://help.zscaler.com/zia/activation#/status-get)

The **zia_activation_status** resource triggers activation of ZIA pending configuration changes. Use it when you want activation to run as part of your Terraform run, after specific resources are applied — for example, by listing those resources in a `depends_on` block so that activation runs only after they are created or updated.

Activation timing is controlled by the ZIA platform, not by Terraform. The provider cannot override ZIA’s native behavior (including auto-activation after inactivity or logout). This resource is one of three ways to activate changes with the provider. For the full picture and recommended options, see the [Activation Overview](https://registry.terraform.io/providers/zscaler/zia/latest/docs/guides/zia-activator-overview) guide.

~> **NOTE** You can also activate in-flight during apply by setting the `ZIA_ACTIVATION` environment variable, or use the recommended out-of-band method with the **ziaActivator** CLI after `terraform apply`. See the [Activation Overview](https://registry.terraform.io/providers/zscaler/zia/latest/docs/guides/zia-activator-overview) guide.

## Example Usage

Trigger activation after specific resources are applied using `depends_on`:

```hcl
resource "zia_url_filtering_rule" "example" {
  name = "Example URL Rule"
  # ... other arguments
}

resource "zia_activation_status" "activation" {
  status = "ACTIVE"

  depends_on = [zia_url_filtering_rule.example]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `status` - (Required) Activates configuration changes.
  * ``0`` = ``ACTIVE``

## Attributes Reference

N/A

## Import

Activation is not an importable resource.

