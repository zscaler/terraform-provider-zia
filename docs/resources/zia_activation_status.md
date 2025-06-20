---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZIA: activation_status"
description: |-
  Official documentation https://help.zscaler.com/zia/saving-and-activating-changes-zia-admin-portal
  API documentation https://help.zscaler.com/zia/activation#/status-get
  Activates configuration changes
---

# zia_activation_status (Resource)

* [Official documentation](https://help.zscaler.com/zia/saving-and-activating-changes-zia-admin-portal)
* [API documentation](https://help.zscaler.com/zia/activation#/status-get)

The **zia_activation_status** resource allows the activation of ZIA pending configurations. This resource must always be executed after the resource creation for successfully policy/configuration activation to occur.

~> **NOTE** As of right now, Terraform does not provide native support for commits or post-activation configuration, so configuration and policy activations are handled out-of-band. In order to handle the activation as part of the provider, a separate source code have been developed to generate a CLI binary.

~> **NOTE** As of version [v2.8.0](https://github.com/zscaler/terraform-provider-zia/releases/tag/v2.8.0) the activation is performed as part of the `terraform apply` during the creation or update of a resource or during the `terraform destroy` during the deletion of a resource. With this improvement the objective is to deprecate the dedicated `zia_activation_status` resource.

## Example Usage

```hcl
data "zia_activation_status" "activation" {}

resource "zia_activation_status" "activation" {
  status                      = "ACTIVE"
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
