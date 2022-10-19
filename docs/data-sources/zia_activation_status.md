---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZIA: activation_status"
description: |-
  "Get Activation Status changes".
---

# Data Source: zia_activation_status

The **zia_activation_status** data source allows to get information about the activation status of ZIA configurations.

~> As of right now, Terraform does not provide native support for commits or post-activation configuration, so configuration and policy activations are handled out-of-band. In order to handle the activation as part of the provider, a separate source code have been developed to generate a CLI binary.

## Example Usage

```hcl
data "zia_activation_status" "activation" {}

```

## Argument Reference

The following arguments are supported:

### Required

There is no required parameter, and the data source will return the one of the current activation statuses below

* `status` - (Required) Activates configuration changes.
  * ``0`` = ``ACTIVE``
  * ``1`` = ``PENDING``
  * ``2`` = ``INPROGRESS``

## Attributes Reference

N/A

## Import

Activation is not an importable resource.
