---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZIA: activation_status"
description: |-
  "Activates configuration changes".
---

# Resource: zia_activation_status

The **zia_activation_status** resource allows the activation of ZIA pending configurations. This resource must always be executed after the resource creation for successfully policy/configuration activation to occur.

~> As of right now, Terraform does not provide native support for commits or post-activation configuration, so configuration and policy activations are handled out-of-band. In order to handle the activation as part of the provider, a separate source code have been developed to generate a CLI binary.

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
