---
layout: "zscaler"
page_title: "Provider: Zscaler Internet Access (ZIA)"
description: |-
    The Zscaler Internet Access provider is used to interact with Zscaler Internet Access (ZIA) API
---

# Zscaler Internet Access (ZIA) Provider

The Zscaler Internet Access provider is used to interact with ZIA API, to automate the provisioning of new locations, IPSec and GRE tunnels, URL filtering policies, Cloud Firewall Policies, DLP Dictionaries, Local Accounts etc. The provider is intended to save time and reducing configuration errors. With this ZIA provider, DevOps teams can automate their security and transform it into DevSecOps workflows. To use this  provider, you must create ZIA API credentials.

Use the navigation on the left to read about the available resources.

Support Disclaimer
-------
!> **Disclaimer:** This Terraform provider is community supported. Although this provider is supported by Zscaler employees, it is **NOT** supported by Zscaler support. Please open all enhancement requests and issues on [Github Issues](https://github.com/zscaler/terraform-provider-zpa/issues) for support.

## Examples Usage

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zia = {
            version = "2.2.3"
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

### Parallelism

Terraform uses goroutines to speed up deployment, but the number of parallel
operations is launches exceeds
[what is recommended](https://help.zscaler.com/zia/about-rate-limiting):

⚠️ **WARNING:** Due to API limitations, we recommend to limit the number of requests ONE, when configuring the following resources rules:
    - ``zia_dlp_web_rules``
    - ``zia_url_filtering_rules``
    - ``zia_firewall_filtering_rule`

  This will allow the API to settle these resources in the correct order. Pushing large batches of security rules at once, may incur in Terraform to Timeout after 20 mins, as it will try to place the rules in the incorrect order. This issue will be addressed in future versions.

In order to accomplish this, make sure you set the
[parallelism](https://www.terraform.io/cli/commands/apply#parallelism-n) value at or
below this limit to prevent performance impacts.

## Support

This template/solution are released under an as-is, best effort, support
policy. These scripts should be seen as community supported and Zscaler
Business Development Team will contribute our expertise as and when possible.
We do not provide technical support or help in using or troubleshooting the components
of the project through our normal support options such as Zscaler support teams,
or ASC (Authorized Support Centers) partners and backline
support options. The underlying product used (Zscaler Internet Access API) by the
scripts or templates are still supported, but the support is only for the
product functionality and not for help in deploying or using the template or
script itself. Unless explicitly tagged, all projects or work posted in our
GitHub repository at (<https://github.com/zscaler>) or sites other
than our official Downloads page on <https://support.zscaler.com>
are provided under the best effort policy.
