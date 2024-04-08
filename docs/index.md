---
layout: "zscaler"
page_title: "Provider: Zscaler Internet Access (ZIA)"
description: |-
    The Zscaler Internet Access provider is used to interact with Zscaler Internet Access (ZIA) API
---

# Zscaler Internet Access (ZIA) Provider

The Zscaler Internet Access provider is used to interact with ZIA API, to automate the provisioning of new locations, IPSec and GRE tunnels, URL filtering policies, Cloud Firewall Policies, DLP Dictionaries, Local Accounts etc. The provider is intended to save time and reducing configuration errors. With this ZIA provider, DevOps teams can automate their security and transform it into DevSecOps workflows. To use this  provider, you must create ZIA API credentials.

Use the navigation on the left to read about the available resources.

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.

## Examples Usage

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zia = {
            version = "2.7.0"
            source = "zscaler/zia"
        }
    }
}

provider "zia" {}
```

```hcl
# Create a location management resource
resource "zia_location_management" "testAcc_location"{
    # ...
}
```

## Authentication

The ZIA provider offers various means of providing credentials for authentication. The following methods are supported:

* Static credentials
* Environment variables

### Static credentials

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

Static credentials can be provided by specifying the `username`, `password`, `api_key`, `zia_cloud` arguments in-line in the ZIA provider block:

**Usage:**

```hcl
provider "zia" {
  username      = 'xxxxxxxxxxxxxxxx'
  password      = 'xxxxxxxxxxxxxxxx'
  api_key       = 'xxxxxxxxxxxxxxxx'
  zia_cloud     = '<zscaler_cloud_name>'
}
```

### Support Zscaler Internet Access Clouds

The ZIA Terraform Provider supports the following environments:

* zscaler
* zscalerone
* zscalertwo
* zscalerthree
* zscloud
* zscalerbeta
* zscalergov
* zscalerten
* zspreview

### Environment variables

You can provide credentials via the `ZIA_USERNAME`, `ZIA_PASSWORD`, `ZIA_API_KEY`, `ZIA_CLOUD` environment variables, representing your ZIA username, password, API Key credentials and tenant base URL, respectively.

```hcl
provider "zia" {}
```

**Usage:**

```sh
export ZIA_USERNAME = "xxxxxxxxxxxxxxxx"
export ZIA_PASSWORD = "xxxxxxxxxxxxxxxx"
export ZIA_API_KEY  = "xxxxxxxxxxxxxxxx"
export ZIA_CLOUD    = "xxxxxxxxxxxxxxxx"
terraform plan
```

If you are on Windows, use PowerShell to set the environmenr variables using the following commands:

```sh
$env:username   = 'xxxxxxxxxxxxxxxx'
$env:password   = 'xxxxxxxxxxxxxxxx'
$env:api_key    = 'xxxxxxxxxxxxxxxx'
$env:zia_cloud  = '<zscaler_cloud_name>'
```

For details about how to retrieve your tenant Base URL and API key/token refer to the Zscaler help portal. <https://help.zscaler.com/zia/getting-started-zia-api>

### ZIA Configuration Activation

The ZIA platform requires every configuration to be activated. As of version [v2.8.0](https://github.com/zscaler/terraform-provider-zia/releases/tag/v2.8.0) the provider supports implicit activation. In order to make this process more flexible, we have implemented a dedicated environment variable `ZIA_ACTIVATION`, which when set to `true` will implicitly activate the changes as resources are configured.
If the environment variable `ZIA_ACTIVATION` is not set, you must then use the out of band activation method described here [zia activator](guides/zia-activator-overview.md) or leverage the dedicated activation resource `zia_activation_status`.

### Zscaler Sandbox Authentication

The ZIA Terraform provider requires both the `ZIA_CLOUD` and `ZIA_SANDBOX_TOKEN` in order to authenticate to the Zscaler Cloud Sandbox environment. For details on how obtain the API Token visit the Zscaler help portal [About Sandbox API Token](https://help.zscaler.com/zia/about-sandbox-api-token)

### Parallelism

Terraform uses goroutines to speed up deployment, but the number of parallel
operations is launches exceeds
[what is recommended](https://help.zscaler.com/zia/about-rate-limiting):

⚠️ **WARNING:** Due to API limitations, we recommend to limit the number of requests to ONE, when configuring the following resources rules:
    - ``zia_dlp_web_rules``
    - ``zia_url_filtering_rules``
    - ``zia_firewall_filtering_rule``
    - ``zia_forwarding_control_rule``

In order to accomplish this, make sure you set the [parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at or below this limit to prevent performance impacts.

## General Support Statement

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.
